package remote

import (
	"github.com/networkservicemesh/api/pkg/api/networkservice"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/services"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor/connectionmonitor"
)

// MonitorServer is a monitor.Server for remote/connection GRPC API
type MonitorServer interface {
	monitor.Server
	networkservice.MonitorConnectionServer
}

type monitorServer struct {
	connectionmonitor.MonitorServer
	manager *services.ClientConnectionManager
}

// NewMonitorServer creates a new MonitorServer
func NewMonitorServer(manager *services.ClientConnectionManager) MonitorServer {
	rv := &monitorServer{
		MonitorServer: connectionmonitor.NewMonitorServer("RemoteConnection"),
		manager:       manager,
	}
	return rv
}

// MonitorConnections adds recipient for MonitorServer events
func (s *monitorServer) MonitorConnections(selector *networkservice.MonitorScopeSelector, recipient networkservice.MonitorConnection_MonitorConnectionsServer) error {
	err := s.MonitorServer.MonitorConnections(selector, recipient)
	if s.manager != nil {
		s.manager.UpdateRemoteMonitorDone(selector.GetPathSegments())
	}
	return err
}
