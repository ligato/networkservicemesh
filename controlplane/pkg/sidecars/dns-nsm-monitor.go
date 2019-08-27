package sidecars

import (
	"context"
	"path"
	"time"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/reload"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"

	"github.com/networkservicemesh/networkservicemesh/utils"

	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
)

const (
	//DefaultPathToCorefile contains default path to Corefile
	DefaultPathToCorefile = "/etc/coredns/Corefile"
	//DefaultReloadCorefileTime means time to reload Corefile
	DefaultReloadCorefileTime = 2 * time.Second
	defaultK8sDNSServer       = "10.96.0.10"
)

//NsmDNSMonitorHandler implements NSMMonitorHandler interface for handling dnsConfigs
type NsmDNSMonitorHandler struct {
	EmptyNSMMonitorHandler
	dnsConfigManager *utils.DNSConfigManager
	corefileUpdater  utils.Operation
}

//DefaultDNSNsmMonitor creates default DNS nsm monitor
func DefaultDNSNsmMonitor() NSMApp {
	return NewDNSNsmMonitor(DefaultPathToCorefile, DefaultReloadCorefileTime)
}

//NewDNSNsmMonitor creates new dns nsm monitor with a specific path to corefile and time to reload corefile
func NewDNSNsmMonitor(pathToCorefile string, reloadTime time.Duration) NSMApp {
	dnsConfigManager := utils.NewDNSConfigManager(defaultBasicDNSConfig(), reloadTime)
	logrus.Infof("Created corefile %v", pathToCorefile)
	result := NewNSMMonitorApp()
	pathToSocket := path.Dir(pathToCorefile) + "/client.sock"
	conn, err := tools.DialContextUnix(context.TODO(), pathToSocket)
	if err != nil {
		logrus.Error(err)
	}
	client := reload.NewReloadServiceClient(conn)
	corefileUpdater := utils.NewSingleAsyncOperation(func() {
		file := dnsConfigManager.Caddyfile(pathToCorefile)
		err, _ := client.Reload(context.Background(), &reload.ReloadMessage{Context: file.String()})
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("Corefile updated.")
	})
	result.SetHandler(&NsmDNSMonitorHandler{
		corefileUpdater:  corefileUpdater,
		dnsConfigManager: dnsConfigManager,
	})
	return result
}

//Connected checks connection and l handle all dns configs
func (h *NsmDNSMonitorHandler) Connected(conns map[string]*connection.Connection) {
	logrus.Info("NsmDNSMonitor: connected")
	for _, conn := range conns {
		if conn.Context == nil {
			logrus.Infof("conn %v has empty ConnectionContext", conn)
			continue
		}
		if conn.Context.DnsContext == nil {
			logrus.Infof("conn %v has empty DnsContext", conn)
			continue
		}
		for _, config := range conn.Context.DnsContext.Configs {
			logrus.Infof("Adding dns config with id: %v, value: %v", conn.Id, config)
			h.dnsConfigManager.Store(conn.Id, *config)
		}
	}
	h.corefileUpdater.Run()
}

//Closed removes all dns configs related to connection ID
func (h *NsmDNSMonitorHandler) Closed(conn *connection.Connection) {
	logrus.Infof("Deleting config with id %v", conn.Id)
	h.dnsConfigManager.Delete(conn.Id)
	h.corefileUpdater.Run()
}

func defaultBasicDNSConfig() connectioncontext.DNSConfig {
	return connectioncontext.DNSConfig{
		DnsServerIps:  []string{defaultK8sDNSServer},
		SearchDomains: []string{},
	}
}
