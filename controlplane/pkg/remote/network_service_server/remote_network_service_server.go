package network_service_server

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/nsm"
	remote_connection "github.com/networkservicemesh/networkservicemesh/controlplane/api/remote/connection"
	remote_networkservice "github.com/networkservicemesh/networkservicemesh/controlplane/api/remote/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/serviceregistry"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor"
)

const (
	// nseConnectionTimeout defines a timoute for NSM to succeed connection to NSE (seconds)
	nseConnectionTimeout = 15 * time.Second
)

type remoteNetworkServiceServer struct {
	model           model.Model
	serviceRegistry serviceregistry.ServiceRegistry
	monitor         monitor.Server
	manager         nsm.NetworkServiceManager
}

// NewRemoteNetworkServiceServer creates a new remote.NetworkServiceServer
func NewRemoteNetworkServiceServer(model model.Model, manager nsm.NetworkServiceManager, serviceRegistry serviceregistry.ServiceRegistry, connectionMonitor monitor.Server) remote_networkservice.NetworkServiceServer {
	server := &remoteNetworkServiceServer{
		model:           model,
		serviceRegistry: serviceRegistry,
		monitor:         connectionMonitor,
		manager:         manager,
	}
	return server
}

func (srv *remoteNetworkServiceServer) Request(ctx context.Context, request *remote_networkservice.NetworkServiceRequest) (*remote_connection.Connection, error) {
	logrus.Infof("RemoteNSMD: Received request from client to connect to NetworkService: %v", request)
	conn, err := srv.manager.Request(ctx, request)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	result := conn.(*remote_connection.Connection)
	srv.monitor.Update(ctx, result)

	logrus.Info("RemoteNSMD: Dataplane configuration done...")
	return result, nil
}

func (srv *remoteNetworkServiceServer) Close(ctx context.Context, connection *remote_connection.Connection) (*empty.Empty, error) {
	logrus.Infof("Remote closing connection: %v", *connection)
	clientConnection := srv.model.GetClientConnection(connection.GetId())
	if clientConnection == nil {
		return nil, errors.Errorf("There is no such client connection %v", connection)
	}
	err := srv.manager.Close(ctx, clientConnection)
	if err != nil {
		logrus.Errorf("Error during connection close: %v", err)
	}
	srv.monitor.Delete(ctx, connection)
	return &empty.Empty{}, nil
}
