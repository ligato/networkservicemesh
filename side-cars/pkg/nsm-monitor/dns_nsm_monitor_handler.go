package nsmmonitor

import (
	"context"

	"github.com/networkservicemesh/networkservicemesh/utils"

	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/k8s/api/nsm-coredns/update"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
)

const (
	//UpdateAPIClientSock means path to client socket for dns context update server
	UpdateAPIClientSock = utils.EnvVar("UPDATE_API_CLIENT_SOCKET")
)

//nsmDNSMonitorHandler implements Handler interface for handling dnsConfigs
type nsmDNSMonitorHandler struct {
	EmptyNSMMonitorHandler
	dnsConfigUpdateClient update.DNSConfigServiceClient
}

func (h *nsmDNSMonitorHandler) Updated(old, new *networkservice.Connection) {
	logrus.Infof("Deleting config with id %v", old.Id)
	_, _ = h.dnsConfigUpdateClient.RemoveDNSContext(context.Background(), &update.RemoveDNSContextMessage{ConnectionID: old.Id})
	logrus.Infof("Adding config with id %v", new.Id)
	_, _ = h.dnsConfigUpdateClient.AddDNSContext(context.Background(), &update.AddDNSContextMessage{ConnectionID: new.Id, Context: new.Context.DnsContext})
}

//NewNsmDNSMonitorHandler creates new DNS monitor handler
func NewNsmDNSMonitorHandler() Handler {
	clientSock := UpdateAPIClientSock.StringValue()
	if clientSock == "" {
		logrus.Fatalf("unable to create Handler instance. Expect %v is not empty", UpdateAPIClientSock.Name())
	}
	conn, err := tools.DialUnix(clientSock)
	if err != nil {
		logrus.Fatalf("An error during dial unix socket by path %v, error: %v", clientSock, err.Error())
	}
	return &nsmDNSMonitorHandler{
		dnsConfigUpdateClient: update.NewDNSConfigServiceClient(conn),
	}
}

func (h *nsmDNSMonitorHandler) Connected(conns map[string]*networkservice.Connection) {
	for _, conn := range conns {
		if conn.Context == nil || conn.Context.DnsContext == nil {
			continue
		}
		logrus.Info(conn.Context.DnsContext)
		_, _ = h.dnsConfigUpdateClient.AddDNSContext(context.Background(), &update.AddDNSContextMessage{ConnectionID: conn.Id, Context: conn.Context.DnsContext})
	}
}

func (h *nsmDNSMonitorHandler) Closed(conn *networkservice.Connection) {
	logrus.Infof("Deleting config with id %v", conn.Id)
	_, _ = h.dnsConfigUpdateClient.RemoveDNSContext(context.Background(), &update.RemoveDNSContextMessage{ConnectionID: conn.Id})
}
