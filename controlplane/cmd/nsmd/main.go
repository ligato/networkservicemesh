package main

import (
	"context"
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

	span := opentracing.StartSpan("nsmd")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	nsmdProbes := nsmd.NewProbes()
	go nsmdProbes.BeginHealthCheck()

	apiRegistry := nsmd.NewApiRegistry()
	serviceRegistry := nsmd.NewServiceRegistry()
	pluginRegistry := plugins.NewPluginRegistry()

	if err := pluginRegistry.Start(ctx); err != nil {
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
	var srvErr error
	// Start NSMD server first, load local NSE/client registry and only then start dataplane/wait for it and recover active connections.

	if server, srvErr = nsmd.StartNSMServer(ctx, model, manager, serviceRegistry, apiRegistry); srvErr != nil {
		logrus.Errorf("error starting nsmd service: %+v", srvErr)
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
	if err := manager.WaitForDataplane(ctx, nsmd.DataplaneTimeout); err != nil {
		logrus.Errorf("Error waiting for dataplane..")
		return
	}
	nsmdProbes.SetDPServerReady()

	// Choose a public API listener
	nsmdAPIAddress := os.Getenv(NsmdAPIAddressEnv)
	if strings.TrimSpace(nsmdAPIAddress) == "" {
		nsmdAPIAddress = NsmdAPIAddressDefaults
	}
	sock, sockErr := apiRegistry.NewPublicListener(nsmdAPIAddress)
	if sockErr != nil {
		logrus.Errorf("failed to start Public API server %v", sockErr)
		return
	}
	nsmdProbes.SetPublicListenerReady()

	server.StartAPIServerAt(ctx, sock)
	nsmdProbes.SetAPIServerReady()

	elapsed := time.Since(start)
	logrus.Debugf("Starting NSMD took: %s", elapsed)

	<-c
}
