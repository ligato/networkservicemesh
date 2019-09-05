package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/crossconnect"
	local_connection "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	local_networkservice "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/networkservice"
	nsm2 "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/nsm/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/nsmdapi"
	pluginsapi "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/plugins"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/registry"
	remote_networkservice "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsmd"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/plugins"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/prefix_pool"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/serviceregistry"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/vni"
	"github.com/networkservicemesh/networkservicemesh/dataplane/pkg/apis/dataplane"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
)

const (
	Master = "master"
	Worker = "worker"
)

type sharedStorage struct {
	sync.RWMutex

	services  map[string]*registry.NetworkService
	managers  map[string]*registry.NetworkServiceManager
	endpoints map[string]*registry.NetworkServiceEndpoint
}

func newSharedStorage() *sharedStorage {
	return &sharedStorage{
		services:  make(map[string]*registry.NetworkService),
		managers:  make(map[string]*registry.NetworkServiceManager),
		endpoints: make(map[string]*registry.NetworkServiceEndpoint),
	}
}

type nsmdTestServiceDiscovery struct {
	apiRegistry *testApiRegistry
	storage     *sharedStorage
	nsmCounter  int
	nsmgrName   string
}

func (impl *nsmdTestServiceDiscovery) RegisterNSE(ctx context.Context, in *registry.NSERegistration, opts ...grpc.CallOption) (*registry.NSERegistration, error) {
	logrus.Infof("Register NSE: %v", in)

	if in.GetNetworkService() != nil {
		impl.storage.services[in.GetNetworkService().GetName()] = in.GetNetworkService()
	}
	if in.GetNetworkServiceManager() != nil {
		in.NetworkServiceManager.Name = impl.nsmgrName
		impl.nsmCounter++
	}
	if in.GetNetworkServiceEndpoint() != nil {
		impl.storage.Lock()
		impl.storage.endpoints[in.GetNetworkServiceEndpoint().GetName()] = in.GetNetworkServiceEndpoint()
		impl.storage.Unlock()
	}
	in.NetworkServiceManager = impl.storage.managers[impl.nsmgrName]
	return in, nil
}

func (impl *nsmdTestServiceDiscovery) RemoveNSE(ctx context.Context, in *registry.RemoveNSERequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	impl.storage.Lock()
	delete(impl.storage.endpoints, in.GetNetworkServiceEndpointName())
	impl.storage.Unlock()
	return nil, nil
}

func newNSMDTestServiceDiscovery(testAPI *testApiRegistry, nsmgrName string, storage *sharedStorage) *nsmdTestServiceDiscovery {
	return &nsmdTestServiceDiscovery{
		storage:     storage,
		apiRegistry: testAPI,
		nsmCounter:  0,
		nsmgrName:   nsmgrName,
	}
}

func (impl *nsmdTestServiceDiscovery) FindNetworkService(ctx context.Context, in *registry.FindNetworkServiceRequest, opts ...grpc.CallOption) (*registry.FindNetworkServiceResponse, error) {
	impl.storage.RLock()
	storageEndpoints := impl.storage.endpoints
	impl.storage.RUnlock()

	endpoints := []*registry.NetworkServiceEndpoint{}
	managers := map[string]*registry.NetworkServiceManager{}

	for _, ep := range storageEndpoints {
		if ep.NetworkServiceName == in.NetworkServiceName {
			endpoints = append(endpoints, ep)

			mgr := impl.storage.managers[ep.NetworkServiceManagerName]
			if mgr != nil {
				managers[mgr.Name] = mgr
			}
		}
	}

	return &registry.FindNetworkServiceResponse{
		NetworkService:          impl.storage.services[in.NetworkServiceName],
		NetworkServiceEndpoints: endpoints,
		NetworkServiceManagers:  managers,
	}, nil
}

func (impl *nsmdTestServiceDiscovery) RegisterNSM(ctx context.Context, in *registry.NetworkServiceManager, opts ...grpc.CallOption) (*registry.NetworkServiceManager, error) {
	logrus.Infof("Register NSM: %v", in)
	in.Name = impl.nsmgrName
	impl.nsmCounter++
	impl.storage.managers[impl.nsmgrName] = in
	return in, nil
}

func (impl *nsmdTestServiceDiscovery) GetEndpoints(ctx context.Context, empty *empty.Empty, opts ...grpc.CallOption) (*registry.NetworkServiceEndpointList, error) {
	return &registry.NetworkServiceEndpointList{
		NetworkServiceEndpoints: []*registry.NetworkServiceEndpoint{
			{
				Name:                      "ep1",
				NetworkServiceManagerName: "nsm1",
			},
		},
	}, nil
}

type testPluginRegistry struct {
	connectionPluginManager *testConnectionPluginManager
}

type testConnectionPluginManager struct {
	plugins.PluginManager
	plugins []pluginsapi.ConnectionPluginServer
}

func newTestPluginRegistry() *testPluginRegistry {
	return &testPluginRegistry{
		connectionPluginManager: &testConnectionPluginManager{},
	}
}

func (pr *testPluginRegistry) Start(ctx context.Context) error {
	return nil
}

func (pr *testPluginRegistry) Stop() error {
	return nil
}

func (pr *testPluginRegistry) GetConnectionPluginManager() plugins.ConnectionPluginManager {
	return pr.connectionPluginManager
}

func (cpm *testConnectionPluginManager) addPlugin(plugin pluginsapi.ConnectionPluginServer) {
	cpm.plugins = append(cpm.plugins, plugin)
}

func (cpm *testConnectionPluginManager) UpdateConnection(ctx context.Context, wrapper *pluginsapi.ConnectionWrapper) (*pluginsapi.ConnectionWrapper, error) {
	for _, plugin := range cpm.plugins {
		var err error
		wrapper, err = plugin.UpdateConnection(ctx, wrapper)
		if err != nil {
			return nil, err
		}
	}
	return wrapper, nil
}

func (cpm *testConnectionPluginManager) ValidateConnection(ctx context.Context, wrapper *pluginsapi.ConnectionWrapper) (*pluginsapi.ConnectionValidationResult, error) {
	for _, plugin := range cpm.plugins {
		result, err := plugin.ValidateConnection(ctx, wrapper)
		if err != nil {
			return nil, err
		}
		if result.GetStatus() != pluginsapi.ConnectionValidationStatus_SUCCESS {
			return result, nil
		}
	}
	return &pluginsapi.ConnectionValidationResult{Status: pluginsapi.ConnectionValidationStatus_SUCCESS}, nil
}

type nsmdTestServiceRegistry struct {
	nseRegistry             *nsmdTestServiceDiscovery
	apiRegistry             *testApiRegistry
	testDataplaneConnection *testDataplaneConnection
	localTestNSE            local_networkservice.NetworkServiceClient
	vniAllocator            vni.VniAllocator
	rootDir                 string
}

func (impl *nsmdTestServiceRegistry) VniAllocator() vni.VniAllocator {
	return impl.vniAllocator
}

func (impl *nsmdTestServiceRegistry) NewWorkspaceProvider() serviceregistry.WorkspaceLocationProvider {
	return nsmd.NewWorkspaceProvider(impl.rootDir)
}

func (impl *nsmdTestServiceRegistry) WaitForDataplaneAvailable(ctx context.Context, model model.Model, timeout time.Duration) error {
	return nsmd.NewServiceRegistry().WaitForDataplaneAvailable(ctx, model, timeout)
}

func (impl *nsmdTestServiceRegistry) WorkspaceName(endpoint *registry.NSERegistration) string {
	return ""
}

func (impl *nsmdTestServiceRegistry) RemoteNetworkServiceClient(ctx context.Context, nsm *registry.NetworkServiceManager) (remote_networkservice.NetworkServiceClient, *grpc.ClientConn, error) {
	err := tools.WaitForPortAvailable(context.Background(), "tcp", nsm.Url, 100*time.Millisecond)
	if err != nil {
		return nil, nil, err
	}

	logrus.Println("Remote Network Service is available, attempting to connect...")
	conn, err := tools.DialTCP(nsm.GetUrl())
	if err != nil {
		logrus.Errorf("Failed to dial Network Service Registry at %s: %s", nsm.Url, err)
		return nil, nil, err
	}
	client := remote_networkservice.NewNetworkServiceClient(conn)
	return client, conn, nil
}

type localTestNSENetworkServiceClient struct {
	req        *local_networkservice.NetworkServiceRequest
	prefixPool prefix_pool.PrefixPool
}

func (impl *localTestNSENetworkServiceClient) Request(ctx context.Context, in *local_networkservice.NetworkServiceRequest, opts ...grpc.CallOption) (*local_connection.Connection, error) {
	impl.req = in
	netns, _ := tools.GetCurrentNS()
	if netns == "" {
		netns = "12"
	}
	mechanism := &local_connection.Mechanism{
		Type: local_connection.MechanismType_KERNEL_INTERFACE,
		Parameters: map[string]string{
			local_connection.NetNsInodeKey: netns,
			// TODO: Fix this terrible hack using xid for getting a unique interface name
			local_connection.InterfaceNameKey: "nsm" + in.GetConnection().GetId(),
		},
	}

	// TODO take into consideration LocalMechnism preferences sent in request
	srcIP, dstIP, requested, err := impl.prefixPool.Extract(in.Connection.Id, connectioncontext.IpFamily_IPV4, in.Connection.GetContext().GetIpContext().ExtraPrefixRequest...)
	if err != nil {
		return nil, err
	}
	conn := &local_connection.Connection{
		Id:             in.GetConnection().GetId(),
		NetworkService: in.GetConnection().GetNetworkService(),
		Mechanism:      mechanism,
		Context: &connectioncontext.ConnectionContext{
			IpContext: &connectioncontext.IPContext{
				SrcIpAddr:     srcIP.String(),
				DstIpAddr:     dstIP.String(),
				ExtraPrefixes: requested,
			},
		},
	}
	err = conn.IsComplete()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return conn, nil
}

func (impl *localTestNSENetworkServiceClient) Close(ctx context.Context, in *local_connection.Connection, opts ...grpc.CallOption) (*empty.Empty, error) {
	//panic("implement me")
	return nil, nil
}

func (impl *nsmdTestServiceRegistry) EndpointConnection(ctx context.Context, endpoint *model.Endpoint) (local_networkservice.NetworkServiceClient, *grpc.ClientConn, error) {
	return impl.localTestNSE, nil, nil
}

type testDataplaneConnection struct {
	connections []*crossconnect.CrossConnect
}

func (impl *testDataplaneConnection) Request(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*crossconnect.CrossConnect, error) {
	impl.connections = append(impl.connections, in)

	if source := in.GetLocalSource(); source != nil && source.Labels != nil {
		if source.Labels != nil {
			if val, ok := source.Labels["dataplane_sleep"]; ok {
				delay, err := strconv.Atoi(val)
				if err == nil {
					logrus.Infof("Delaying Dataplane Request: %v", delay)
					<-time.After(time.Duration(delay) * time.Second)
				}
			}
		}
	}
	return in, nil
}

func (impl *testDataplaneConnection) Close(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (impl *testDataplaneConnection) MonitorMechanisms(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (dataplane.Dataplane_MonitorMechanismsClient, error) {
	return nil, nil
}

func (impl *nsmdTestServiceRegistry) DataplaneConnection(ctx context.Context, dataplane *model.Dataplane) (dataplane.DataplaneClient, *grpc.ClientConn, error) {
	return impl.testDataplaneConnection, nil, nil
}

func (impl *nsmdTestServiceRegistry) NSMDApiClient() (nsmdapi.NSMDClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", impl.apiRegistry.nsmdPort)
	logrus.Infof("Connecting to nsmd on socket: %s...", addr)

	// Wait to be sure it is already initialized
	err := tools.WaitForPortAvailable(context.Background(), "tcp", addr, 100*time.Millisecond)
	if err != nil {
		return nil, nil, err
	}
	conn, err := tools.DialTCP(addr)
	if err != nil {
		logrus.Errorf("Failed to dial Network Service Registry at %s: %s", addr, err)
		return nil, nil, err
	}

	logrus.Info("Requesting nsmd for client connection...")
	return nsmdapi.NewNSMDClient(conn), conn, nil
}

func (impl *nsmdTestServiceRegistry) GetPublicAPI() string {
	return fmt.Sprintf("%s:%d", "127.0.0.1", impl.apiRegistry.nsmdPublicPort)
}

func (impl *nsmdTestServiceRegistry) DiscoveryClient() (registry.NetworkServiceDiscoveryClient, error) {
	return impl.nseRegistry, nil
}

func (impl *nsmdTestServiceRegistry) NseRegistryClient() (registry.NetworkServiceRegistryClient, error) {
	return impl.nseRegistry, nil
}

func (impl *nsmdTestServiceRegistry) NsmRegistryClient() (registry.NsmRegistryClient, error) {
	return impl.nseRegistry, nil
}

func (impl *nsmdTestServiceRegistry) Stop() {
	logrus.Printf("Delete temporary workspace root: %s", impl.rootDir)
	os.RemoveAll(impl.rootDir)
}

type testApiRegistry struct {
	nsmdPort       int
	nsmdPublicPort int
}

func (impl *testApiRegistry) NewNSMServerListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	impl.nsmdPort = listener.Addr().(*net.TCPAddr).Port
	return listener, err
}

func (impl *testApiRegistry) NewPublicListener(nsmdAPIAddress string) (net.Listener, error) {
	listener, err := net.Listen("tcp", nsmdAPIAddress)
	impl.nsmdPublicPort = listener.Addr().(*net.TCPAddr).Port
	return listener, err
}

func newTestApiRegistry() *testApiRegistry {
	return &testApiRegistry{
		nsmdPort:       0,
		nsmdPublicPort: 0,
	}
}

func newNetworkServiceClient(nsmServerSocket string) (local_networkservice.NetworkServiceClient, *grpc.ClientConn, error) {
	// Wait till we actually have an nsmd to talk to
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := tools.WaitForPortAvailable(ctx, "unix", nsmServerSocket, 100*time.Millisecond)
	if err != nil {
		return nil, nil, err
	}

	conn, err := tools.DialUnix(nsmServerSocket)
	if err != nil {
		return nil, nil, err
	}
	// Init related activities start here
	nsmConnectionClient := local_networkservice.NewNetworkServiceClient(conn)
	return nsmConnectionClient, conn, nil
}

type nsmdFullServer interface {
	Stop()
}
type nsmdFullServerImpl struct {
	apiRegistry     *testApiRegistry
	nseRegistry     *nsmdTestServiceDiscovery
	pluginRegistry  *testPluginRegistry
	serviceRegistry *nsmdTestServiceRegistry
	testModel       model.Model
	manager         nsm2.NetworkServiceManager
	nsmServer       nsmd.NSMServer
	rootDir         string
}

func (srv *nsmdFullServerImpl) Stop() {
	srv.serviceRegistry.Stop()
	if srv.nsmServer != nil {
		srv.nsmServer.Stop()
	}

	dir, _ := ioutil.ReadDir(srv.rootDir)
	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{srv.rootDir, d.Name()}...))
	}
}

func (srv *nsmdFullServerImpl) StopNoClean() {
	if srv.nsmServer != nil {
		srv.nsmServer.Stop()
	}
}

func (impl *nsmdFullServerImpl) addFakeDataplane(dp_name string, dp_addr string) {
	impl.testModel.AddDataplane(&model.Dataplane{
		RegisteredName: dp_name,
		SocketLocation: dp_addr,
		LocalMechanisms: []connection.Mechanism{
			&local_connection.Mechanism{
				Type: local_connection.MechanismType_KERNEL_INTERFACE,
			},
		},
		MechanismsConfigured: true,
	})
}

func (srv *nsmdFullServerImpl) registerFakeEndpoint(networkServiceName string, payload string, nse_address string) *model.Endpoint {
	return srv.registerFakeEndpointWithName(networkServiceName, payload, nse_address, networkServiceName+"provider")
}
func (srv *nsmdFullServerImpl) registerFakeEndpointWithName(networkServiceName string, payload string, nsmgrName string, endpointname string) *model.Endpoint {
	reg := &registry.NSERegistration{
		NetworkService: &registry.NetworkService{
			Name:    networkServiceName,
			Payload: payload,
		},
		NetworkServiceEndpoint: &registry.NetworkServiceEndpoint{
			Name:                      endpointname,
			Payload:                   payload,
			NetworkServiceManagerName: nsmgrName,
			NetworkServiceName:        networkServiceName,
		},
	}
	regResp, err := srv.nseRegistry.RegisterNSE(context.Background(), reg)
	if err != nil {
		panic(err)
	}
	if regResp.NetworkService.Name != networkServiceName {
		panic(fmt.Errorf("%s is not equal to %s", regResp.NetworkService.Name, networkServiceName))
	}

	return &model.Endpoint{
		Endpoint:       reg,
		Workspace:      "nsm-1",
		SocketLocation: "nsm-1/client",
	}
}

func (srv *nsmdFullServerImpl) requestNSMConnection(clientName string) (local_networkservice.NetworkServiceClient, *grpc.ClientConn) {
	response := srv.requestNSM(clientName)
	// Now we could try to connect via Client API
	nsmClient, conn := srv.createNSClient(response)
	return nsmClient, conn
}

func (srv *nsmdFullServerImpl) createNSClient(response *nsmdapi.ClientConnectionReply) (local_networkservice.NetworkServiceClient, *grpc.ClientConn) {
	nsmClient, conn, err := newNetworkServiceClient(response.HostBasedir + "/" + response.Workspace + "/" + response.NsmServerSocket)
	if err != nil {
		panic(err)
	}
	return nsmClient, conn
}

func (srv *nsmdFullServerImpl) requestNSM(clientName string) *nsmdapi.ClientConnectionReply {
	client, con, err := srv.serviceRegistry.NSMDApiClient()
	if err != nil {
		panic(err)
	}
	defer con.Close()

	response, err := client.RequestClientConnection(context.Background(), &nsmdapi.ClientConnectionRequest{
		Workspace: clientName,
	})

	if err != nil {
		panic(err)
	}

	logrus.Printf("workspace %s", response.Workspace)

	if response.Workspace != clientName {
		panic(fmt.Errorf("%s is not equal to %s", response.Workspace, clientName))
	}
	return response
}

func newNSMDFullServer(nsmgrName string, storage *sharedStorage) *nsmdFullServerImpl {
	rootDir, err := ioutil.TempDir("", "nsmd_test")
	if err != nil {
		panic(err)
	}

	return newNSMDFullServerAt(context.Background(), nsmgrName, storage, rootDir)
}

func newNSMDFullServerAt(ctx context.Context, nsmgrName string, storage *sharedStorage, rootDir string) *nsmdFullServerImpl {
	srv := &nsmdFullServerImpl{}
	srv.apiRegistry = newTestApiRegistry()
	srv.nseRegistry = newNSMDTestServiceDiscovery(srv.apiRegistry, nsmgrName, storage)
	srv.pluginRegistry = newTestPluginRegistry()
	srv.rootDir = rootDir

	prefixPool, err := prefix_pool.NewPrefixPool("10.20.1.0/24")
	if err != nil {
		panic(err)
	}
	srv.serviceRegistry = &nsmdTestServiceRegistry{
		nseRegistry:             srv.nseRegistry,
		apiRegistry:             srv.apiRegistry,
		testDataplaneConnection: &testDataplaneConnection{},
		localTestNSE: &localTestNSENetworkServiceClient{
			prefixPool: prefixPool,
		},
		vniAllocator: vni.NewVniAllocator(),
		rootDir:      rootDir,
	}

	srv.testModel = model.NewModel()
	srv.manager = nsm.NewNetworkServiceManager(srv.testModel, srv.serviceRegistry, srv.pluginRegistry)

	// Choose a public API listener
	sock, err := srv.apiRegistry.NewPublicListener("127.0.0.1:0")
	if err != nil {
		logrus.Errorf("Failed to start Public API server...")
		return nil
	}

	// Lets start NSMD NSE registry service
	nsmServer, err := nsmd.StartNSMServer(ctx, srv.testModel, srv.manager, srv.serviceRegistry, srv.apiRegistry)
	srv.nsmServer = nsmServer
	if err != nil {
		panic(err)
	}

	monitorCrossConnectClient := nsmd.NewMonitorCrossConnectClient(nsmServer, nsmServer.XconManager(), srv.nsmServer)
	srv.testModel.AddListener(monitorCrossConnectClient)

	// Start API Server
	nsmServer.StartAPIServerAt(ctx, sock)

	return srv
}
