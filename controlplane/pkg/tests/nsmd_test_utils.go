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

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/monitor/remote"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/services"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor/connectionmonitor"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/common"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/kernel"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/nsmdapi"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
	nsm2 "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/api/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsmd"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/serviceregistry"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/sid"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/vni"
	"github.com/networkservicemesh/networkservicemesh/forwarder/api/forwarder"
	"github.com/networkservicemesh/networkservicemesh/pkg/probes"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools/spanhelper"
	monitor_crossconnect "github.com/networkservicemesh/networkservicemesh/sdk/monitor/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/sdk/prefix_pool"
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

func NewSharedStorage() *sharedStorage {
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
	logrus.Infof("Test Register NSE: %v", in)

	if in.GetNetworkService() != nil {
		impl.storage.Lock()
		impl.storage.services[in.GetNetworkService().GetName()] = in.GetNetworkService()
		impl.storage.Unlock()
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

func (impl *nsmdTestServiceDiscovery) BulkRegisterNSE(ctx context.Context, opts ...grpc.CallOption) (registry.NetworkServiceRegistry_BulkRegisterNSEClient, error) {
	return nil, errors.Errorf("not implemented")
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

type nsmdTestServiceRegistry struct {
	nseRegistry             *nsmdTestServiceDiscovery
	apiRegistry             *testApiRegistry
	testForwarderConnection *testForwarderConnection
	localTestNSE            networkservice.NetworkServiceClient
	vniAllocator            vni.VniAllocator
	sidAllocator            sid.Allocator
	rootDir                 string
}

func (impl *nsmdTestServiceRegistry) SIDAllocator() sid.Allocator {
	return impl.sidAllocator
}

func (impl *nsmdTestServiceRegistry) VniAllocator() vni.VniAllocator {
	return impl.vniAllocator
}

func (impl *nsmdTestServiceRegistry) NewWorkspaceProvider() serviceregistry.WorkspaceLocationProvider {
	return nsmd.NewWorkspaceProvider(impl.rootDir)
}

func (impl *nsmdTestServiceRegistry) WaitForForwarderAvailable(ctx context.Context, model model.Model, timeout time.Duration) error {
	return nsmd.NewServiceRegistry().WaitForForwarderAvailable(ctx, model, timeout)
}

func (impl *nsmdTestServiceRegistry) WorkspaceName(endpoint *registry.NSERegistration) string {
	return ""
}

func (impl *nsmdTestServiceRegistry) RemoteNetworkServiceClient(ctx context.Context, nsm *registry.NetworkServiceManager) (networkservice.NetworkServiceClient, *grpc.ClientConn, error) {
	span := spanhelper.FromContext(ctx, "RemoteNetworkServiceClient")
	defer span.Finish()
	err := tools.WaitForPortAvailable(span.Context(), "tcp", nsm.Url, 100*time.Millisecond)
	if err != nil {
		return nil, nil, err
	}

	span.Logger().Info("Remote Network Service is available, attempting to connect...")
	conn, err := tools.DialContextTCP(span.Context(), nsm.GetUrl())
	span.LogError(err)
	if err != nil {
		span.Logger().Errorf("Failed to dial Network Service Registry at %s: %s", nsm.Url, err)
		return nil, nil, err
	}
	client := networkservice.NewNetworkServiceClient(conn)
	return client, conn, nil
}

type localTestNSENetworkServiceClient struct {
	sync.Mutex
	req                  *networkservice.NetworkServiceRequest
	prefixPool           prefix_pool.PrefixPool
	requestHandleCounter int
}

func (impl *localTestNSENetworkServiceClient) Request(ctx context.Context, in *networkservice.NetworkServiceRequest, opts ...grpc.CallOption) (*networkservice.Connection, error) {
	impl.Lock()
	impl.requestHandleCounter++
	impl.Unlock()
	impl.req = in
	netns, _ := tools.GetCurrentNS()
	if netns == "" {
		netns = "12"
	}
	mechanism := &networkservice.Mechanism{
		Type: kernel.MECHANISM,
		Parameters: map[string]string{
			common.NetNSInodeKey: netns,
			// TODO: Fix this terrible hack using xid for getting a unique interface name
			common.InterfaceNameKey: "nsm" + in.GetConnection().GetId(),
		},
	}

	// TODO take into consideration LocalMechnism preferences sent in request
	srcIP, dstIP, requested, err := impl.prefixPool.Extract(in.Connection.Id, networkservice.IpFamily_IPV4, in.Connection.GetContext().GetIpContext().ExtraPrefixRequest...)
	if err != nil {
		return nil, err
	}
	conn := &networkservice.Connection{
		Id:             in.GetConnection().GetId(),
		NetworkService: in.GetConnection().GetNetworkService(),
		Mechanism:      mechanism,
		Context: &networkservice.ConnectionContext{
			IpContext: &networkservice.IPContext{
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

func (impl *localTestNSENetworkServiceClient) Close(ctx context.Context, in *networkservice.Connection, opts ...grpc.CallOption) (*empty.Empty, error) {
	//panic("implement me")
	return nil, nil
}

func (impl *nsmdTestServiceRegistry) EndpointConnection(ctx context.Context, endpoint *model.Endpoint) (networkservice.NetworkServiceClient, *grpc.ClientConn, error) {
	return impl.localTestNSE, nil, nil
}

type testForwarderConnection struct {
	connections []*crossconnect.CrossConnect
}

func (impl *testForwarderConnection) Request(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*crossconnect.CrossConnect, error) {
	impl.connections = append(impl.connections, in)

	if source := in.Source; source != nil && source.Labels != nil {
		if source.Labels != nil {
			if val, ok := source.Labels["forwarder_sleep"]; ok {
				delay, err := strconv.Atoi(val)
				if err == nil {
					logrus.Infof("Delaying Forwarder Request: %v", delay)
					<-time.After(time.Duration(delay) * time.Second)
				}
			}
		}
	}
	return in, nil
}

func (impl *testForwarderConnection) Close(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*empty.Empty, error) {
	return nil, nil
}

func (impl *testForwarderConnection) MonitorMechanisms(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (forwarder.MechanismsMonitor_MonitorMechanismsClient, error) {
	return nil, nil
}

func (impl *nsmdTestServiceRegistry) ForwarderConnection(ctx context.Context, forwarder *model.Forwarder) (forwarder.ForwarderClient, *grpc.ClientConn, error) {
	return impl.testForwarderConnection, nil, nil
}

func (impl *nsmdTestServiceRegistry) NSMDApiClient(ctx context.Context) (nsmdapi.NSMDClient, *grpc.ClientConn, error) {
	span := spanhelper.FromContext(ctx, "NSMDApiClient")
	defer span.Finish()
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", impl.apiRegistry.nsmdPort)
	span.Logger().Infof("Connecting to nsmd on socket: %s...", addr)

	// Wait to be sure it is already initialized
	err := tools.WaitForPortAvailable(span.Context(), "tcp", addr, 100*time.Millisecond)
	span.LogError(err)
	if err != nil {
		return nil, nil, err
	}
	conn, err := tools.DialContextTCP(span.Context(), addr)
	if err != nil {
		err = errors.Errorf("failed to dial Network Service Registry at %s: %s", addr, err)
		span.LogError(err)
		return nil, nil, err
	}

	span.Logger().Info("Requesting nsmd for client connection...")
	return nsmdapi.NewNSMDClient(conn), conn, nil
}

func (impl *nsmdTestServiceRegistry) GetPublicAPI() string {
	return fmt.Sprintf("%s:%d", "127.0.0.1", impl.apiRegistry.nsmdPublicPort)
}

func (impl *nsmdTestServiceRegistry) DiscoveryClient(context.Context) (registry.NetworkServiceDiscoveryClient, error) {
	return impl.nseRegistry, nil
}

func (impl *nsmdTestServiceRegistry) NseRegistryClient(context.Context) (registry.NetworkServiceRegistryClient, error) {
	return impl.nseRegistry, nil
}

func (impl *nsmdTestServiceRegistry) NsmRegistryClient(context.Context) (registry.NsmRegistryClient, error) {
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

func newNetworkServiceClient(nsmServerSocket string) (networkservice.NetworkServiceClient, *grpc.ClientConn, error) {
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
	nsmConnectionClient := networkservice.NewNetworkServiceClient(conn)
	return nsmConnectionClient, conn, nil
}

type nsmdFullServer interface {
	Stop()
}
type nsmdFullServerImpl struct {
	apiRegistry               *testApiRegistry
	nseRegistry               *nsmdTestServiceDiscovery
	serviceRegistry           *nsmdTestServiceRegistry
	TestModel                 model.Model
	manager                   nsm2.NetworkServiceManager
	nsmServer                 nsmd.NSMServer
	rootDir                   string
	monitorCrossConnectClient *nsmd.NsmMonitorCrossConnectClient
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

func (srv *nsmdFullServerImpl) AddFakeForwarder(dp_name string, dp_addr string) {
	srv.TestModel.AddForwarder(context.Background(), &model.Forwarder{
		RegisteredName: dp_name,
		SocketLocation: dp_addr,
		LocalMechanisms: []*networkservice.Mechanism{
			{
				Type: kernel.MECHANISM,
			},
		},
		MechanismsConfigured: true,
	})
}

func (srv *nsmdFullServerImpl) RegisterFakeEndpoint(networkServiceName string, payload string, nse_address string) *model.Endpoint {
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
		panic(errors.Errorf("%s is not equal to %s", regResp.NetworkService.Name, networkServiceName))
	}

	return &model.Endpoint{
		Endpoint:       reg,
		Workspace:      "nsm-1",
		SocketLocation: "nsm-1/client",
	}
}

func (srv *nsmdFullServerImpl) requestNSMConnection(clientName string) (networkservice.NetworkServiceClient, *grpc.ClientConn) {
	response := srv.RequestNSM(clientName)
	// Now we could try to connect via Client API
	nsmClient, conn := srv.CreateNSClient(response)
	return nsmClient, conn
}

func (srv *nsmdFullServerImpl) CreateNSClient(response *nsmdapi.ClientConnectionReply) (networkservice.NetworkServiceClient, *grpc.ClientConn) {
	nsmClient, conn, err := newNetworkServiceClient(response.HostBasedir + "/" + response.Workspace + "/" + response.NsmServerSocket)
	if err != nil {
		panic(err)
	}
	return nsmClient, conn
}

func (srv *nsmdFullServerImpl) RequestNSM(clientName string) *nsmdapi.ClientConnectionReply {
	client, con, err := srv.serviceRegistry.NSMDApiClient(context.Background())
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
		panic(errors.Errorf("%s is not equal to %s", response.Workspace, clientName))
	}
	return response
}

func NewNSMDFullServer(nsmgrName string, storage *sharedStorage) *nsmdFullServerImpl {
	rootDir, err := ioutil.TempDir("", "nsmd_test")
	if err != nil {
		panic(err)
	}

	return newNSMDFullServerAt(context.Background(), nsmgrName, storage, rootDir, model.NewModel())
}

func NewNSMDFullServerWithModel(nsmgrName string, storage *sharedStorage, testModel model.Model) *nsmdFullServerImpl {
	rootDir, err := ioutil.TempDir("", "nsmd_test")
	if err != nil {
		panic(err)
	}

	return newNSMDFullServerAt(context.Background(), nsmgrName, storage, rootDir, testModel)
}

func newNSMDFullServerAt(ctx context.Context, nsmgrName string, storage *sharedStorage, rootDir string, testModel model.Model) *nsmdFullServerImpl {
	srv := &nsmdFullServerImpl{}
	srv.apiRegistry = newTestApiRegistry()
	srv.nseRegistry = newNSMDTestServiceDiscovery(srv.apiRegistry, nsmgrName, storage)
	srv.rootDir = rootDir

	prefixPool, err := prefix_pool.NewPrefixPool("10.20.1.0/24")
	if err != nil {
		panic(err)
	}
	srv.serviceRegistry = &nsmdTestServiceRegistry{
		nseRegistry:             srv.nseRegistry,
		apiRegistry:             srv.apiRegistry,
		testForwarderConnection: &testForwarderConnection{},
		localTestNSE: &localTestNSENetworkServiceClient{
			prefixPool:           prefixPool,
			requestHandleCounter: 0,
		},
		vniAllocator: vni.NewVniAllocator(),
		rootDir:      rootDir,
	}

	srv.TestModel = testModel
	srv.manager = nsm.NewNetworkServiceManager(ctx, srv.TestModel, srv.serviceRegistry)

	// Choose a public API listener
	sock, err := srv.apiRegistry.NewPublicListener("127.0.0.1:0")
	if err != nil {
		logrus.Errorf("Failed to start Public API server...")
		return nil
	}

	// Lets start NSMD NSE registry service
	nsmServer, err := nsmd.StartNSMServer(ctx, srv.TestModel, srv.manager, srv.apiRegistry)
	srv.nsmServer = nsmServer
	if err != nil {
		panic(err)
	}

	srv.monitorCrossConnectClient = nsmd.NewMonitorCrossConnectClient(srv.TestModel, nsmServer, nsmServer.XconManager(), srv.nsmServer)
	srv.TestModel.AddListener(srv.monitorCrossConnectClient)
	probes := probes.New("Test probes", nil)

	nsmServer.StartAPIServerAt(ctx, sock, probes)

	return srv
}

// CreateRequest - create test request
func CreateRequest() *networkservice.NetworkServiceRequest {
	request := &networkservice.NetworkServiceRequest{
		Connection: &networkservice.Connection{
			NetworkService: "golden_network",
			Context: &networkservice.ConnectionContext{
				IpContext: &networkservice.IPContext{
					DstIpRequired: true,
					SrcIpRequired: true,
				},
			},
			Labels: make(map[string]string),
		},
		MechanismPreferences: []*networkservice.Mechanism{
			{
				Type: kernel.MECHANISM,
				Parameters: map[string]string{
					common.NetNSInodeKey:    "10",
					common.InterfaceNameKey: "icmp-responder1",
				},
			},
		},
	}

	return request
}

type monitorManager struct {
	crossConnectMonitor     monitor_crossconnect.MonitorServer
	remoteConnectionMonitor remote.MonitorServer
	localConnectionMonitors map[string]connectionmonitor.MonitorServer
}

func (m *monitorManager) CrossConnectMonitor() monitor_crossconnect.MonitorServer {
	return m.crossConnectMonitor
}

func (m *monitorManager) RemoteConnectionMonitor() monitor.Server {
	return m.remoteConnectionMonitor
}

func (m *monitorManager) LocalConnectionMonitor(workspace string) connectionmonitor.MonitorServer {
	return m.localConnectionMonitors[workspace]
}

type endpointManager struct {
	model model.Model
}

func (stub *endpointManager) DeleteEndpointWithBrokenConnection(ctx context.Context, endpoint *model.Endpoint) error {
	stub.model.DeleteEndpoint(ctx, endpoint.EndpointName())
	return nil
}

func startAPIServer(model model.Model, nsmdAPIAddress string) (*grpc.Server, monitor_crossconnect.MonitorServer, net.Listener, error) {
	sock, err := net.Listen("tcp", nsmdAPIAddress)
	if err != nil {
		return nil, nil, sock, err
	}
	grpcServer := tools.NewServer(context.Background())
	serviceRegistry := nsmd.NewServiceRegistry()

	xconManager := services.NewClientConnectionManager(model, nil, serviceRegistry)

	monitorManager := &monitorManager{
		crossConnectMonitor:     monitor_crossconnect.NewMonitorServer(),
		remoteConnectionMonitor: remote.NewMonitorServer(xconManager),
		localConnectionMonitors: map[string]connectionmonitor.MonitorServer{},
	}

	crossconnect.RegisterMonitorCrossConnectServer(grpcServer, monitorManager.crossConnectMonitor)
	networkservice.RegisterMonitorConnectionServer(grpcServer, monitorManager.remoteConnectionMonitor)

	monitorClient := nsmd.NewMonitorCrossConnectClient(model, monitorManager, xconManager, &endpointManager{model: model})
	model.AddListener(monitorClient)
	// Add more public API services here.

	go func() {
		if err := grpcServer.Serve(sock); err != nil {
			logrus.Errorf("failed to start gRPC NSMD API server %+v", err)
		}
	}()
	logrus.Infof("NSM gRPC API Server: %s is operational", nsmdAPIAddress)

	return grpcServer, monitorManager.crossConnectMonitor, sock, nil
}

func readNMSDCrossConnectEvents(address string, count int) []*crossconnect.CrossConnectEvent {
	var err error
	conn, err := tools.DialTCP(address)
	if err != nil {
		logrus.Errorf("Failure to communicate with the socket %s with error: %+v", address, err)
		return nil
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			logrus.Errorf("An error during close conn: %v", err)
		}
	}()
	forwarderClient := crossconnect.NewMonitorCrossConnectClient(conn)

	// Looping indefinitely or until grpc returns an error indicating the other end closed connection.
	stream, err := forwarderClient.MonitorCrossConnects(context.Background(), &empty.Empty{})
	if err != nil {
		logrus.Warningf("Error: %+v.", err)
		return nil
	}
	pos := 0
	result := []*crossconnect.CrossConnectEvent{}
	for {
		event, err := stream.Recv()
		logrus.Infof("(test) receive event: %s %v", event.Type, event.CrossConnects)
		if err != nil {
			logrus.Errorf("Error2: %+v.", err)
			return result
		}
		result = append(result, event)
		pos++
		if pos == count {
			return result
		}
	}
}

type eventCollector struct {
	messages []interface{}
	events   chan interface{}
}

func (e *eventCollector) SendMsg(msg interface{}) error {
	e.messages = append(e.messages, msg)
	e.events <- msg
	return nil
}
