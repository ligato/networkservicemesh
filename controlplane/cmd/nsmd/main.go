package main

import (
	"os"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsmd"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/plugins"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
)

var version string

// Default values and environment variables of proxy connection
const (
	NsmdAPIAddressEnv      = "NSMD_API_ADDRESS"
	NsmdAPIAddressDefaults = ":5001"
)

func main() {
	logrus.Info("Starting nsmd...")
	logrus.Infof("Version: %v", version)
	start := time.Now()

	// Capture signals to cleanup before exiting
	c := tools.NewOSSignalChannel()

	tracer, closer := tools.InitJaeger("nsmd")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	nsmdProbes := nsmd.NewProbes()
	go nsmdProbes.BeginHealthCheck()

	apiRegistry := nsmd.NewApiRegistry()
	serviceRegistry := nsmd.NewServiceRegistry()
	pluginRegistry := plugins.NewPluginRegistry()

	if err := pluginRegistry.Start(); err != nil {
		logrus.Errorf("Failed to start Plugin Registry: %v", err)
		return
	}
	defer func() {
		if err := pluginRegistry.Stop(); err != nil {
			logrus.Errorf("Failed to stop Plugin Registry: %v", err)
		}
	}()

	model := model.NewModel() // This is TCP gRPC server uri to access this NSMD via network.
	defer serviceRegistry.Stop()
	manager := nsm.NewNetworkServiceManager(model, serviceRegistry, pluginRegistry)

	var server nsmd.NSMServer
	var err error
	// Start NSMD server first, load local NSE/client registry and only then start dataplane/wait for it and recover active connections.
	if server, err = nsmd.StartNSMServer(model, manager, serviceRegistry, pluginRegistry, apiRegistry); err != nil {
		logrus.Errorf("Error starting nsmd service: %+v", err)
		return
	}
	defer server.Stop()
	nsmdProbes.SetNSMServerReady()

	// Register CrossConnect monitorCrossConnectServer client as ModelListener
	monitorCrossConnectClient := nsmd.NewMonitorCrossConnectClient(server, server.XconManager(), server)
	model.AddListener(monitorCrossConnectClient)

	// Starting dataplane
	logrus.Info("Starting Dataplane registration server...")
	if err := server.StartDataplaneRegistratorServer(); err != nil {
		logrus.Errorf("Error starting dataplane service: %+v", err)
		return
	}

	// Wait for dataplane to be connecting to us
	if err := manager.WaitForDataplane(nsmd.DataplaneTimeout); err != nil {
		logrus.Errorf("Error waiting for dataplane..")
		return
	}
	nsmdProbes.SetDPServerReady()

	// Choose a public API listener
	nsmdAPIAddress := os.Getenv(NsmdAPIAddressEnv)
	if strings.TrimSpace(nsmdAPIAddress) == "" {
		nsmdAPIAddress = NsmdAPIAddressDefaults
	}
	sock, err := apiRegistry.NewPublicListener(nsmdAPIAddress)
	if err != nil {
		logrus.Errorf("Failed to start Public API server...")
		return
	}
	nsmdProbes.SetPublicListenerReady()

	server.StartAPIServerAt(sock)
	nsmdProbes.SetAPIServerReady()

	elapsed := time.Since(start)
	logrus.Debugf("Starting NSMD took: %s", elapsed)

	<-c
}
