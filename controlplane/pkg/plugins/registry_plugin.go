package plugins

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/plugins"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/registry"
)

// RegistryPluginManager transmits each method call to all registered registry plugins
type RegistryPluginManager interface {
	PluginManager
	plugins.RegistryPluginServer
}

type registryPluginManager struct {
	sync.RWMutex
	pluginClients map[string]plugins.RegistryPluginClient
}

func createRegistryPluginManager() RegistryPluginManager {
	return &registryPluginManager{
		pluginClients: make(map[string]plugins.RegistryPluginClient),
	}
}

func (rpm *registryPluginManager) Register(name string, conn *grpc.ClientConn) {
	client := plugins.NewRegistryPluginClient(conn)
	rpm.addClient(name, client)
}

func (rpm *registryPluginManager) addClient(name string, client plugins.RegistryPluginClient) {
	rpm.Lock()
	defer rpm.Unlock()

	rpm.pluginClients[name] = client
}

func (rpm *registryPluginManager) getClients() map[string]plugins.RegistryPluginClient {
	rpm.RLock()
	defer rpm.RUnlock()

	return rpm.pluginClients
}

func (rpm *registryPluginManager) RegisterNSM(ctx context.Context, nsm *registry.NetworkServiceManager) (*registry.NetworkServiceManager, error) {
	for name, plugin := range rpm.getClients() {
		pluginCtx, cancel := context.WithTimeout(ctx, pluginCallTimeout)

		var err error
		nsm, err = plugin.RegisterNSM(pluginCtx, nsm)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("'%s' registry plugin returned an error: %v", name, err)
		}
	}
	return nsm, nil
}

func (rpm *registryPluginManager) RegisterNSE(ctx context.Context, registration *registry.NSERegistration) (*registry.NSERegistration, error) {
	for name, plugin := range rpm.getClients() {
		pluginCtx, cancel := context.WithTimeout(ctx, pluginCallTimeout)

		var err error
		registration, err = plugin.RegisterNSE(pluginCtx, registration)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("'%s' registry plugin returned an error: %v", name, err)
		}
	}
	return registration, nil
}

func (rpm *registryPluginManager) RemoveNSE(ctx context.Context, request *registry.RemoveNSERequest) (*empty.Empty, error) {
	for name, plugin := range rpm.getClients() {
		pluginCtx, cancel := context.WithTimeout(ctx, pluginCallTimeout)

		_, err := plugin.RemoveNSE(pluginCtx, request)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("'%s' registry plugin returned an error: %v", name, err)
		}
	}
	return &empty.Empty{}, nil
}

func (rpm *registryPluginManager) GetNSEs(ctx context.Context, empty *empty.Empty) (*plugins.NSEList, error) {
	endpoints := []*registry.NetworkServiceEndpoint{}
	for name, plugin := range rpm.getClients() {
		pluginCtx, cancel := context.WithTimeout(ctx, pluginCallTimeout)

		response, err := plugin.GetNSEs(pluginCtx, empty)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("'%s' registry plugin returned an error: %v", name, err)
		}

		endpoints = append(endpoints, response.GetNetworkServiceEndpoints()...)
	}
	return &plugins.NSEList{NetworkServiceEndpoints: endpoints}, nil
}