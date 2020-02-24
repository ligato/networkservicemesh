package vppagent

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.ligato.io/vpp-agent/v3/proto/ligato/configurator"

	"github.com/networkservicemesh/networkservicemesh/forwarder/api/forwarder"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	monitor_crossconnect "github.com/networkservicemesh/networkservicemesh/sdk/monitor/crossconnect"
)

type contextKeyType string

const (
	nextKey       contextKeyType = "nextKey"
	clientKey     contextKeyType = "clientKey"
	dataChangeKey contextKeyType = "dataChangeKey"
	loggerKey     contextKeyType = "loggerKey"
	monitorKey    contextKeyType = "monitorKey"
)

func withNext(ctx context.Context, handler forwarder.ForwarderServer) context.Context {
	return context.WithValue(ctx, nextKey, handler)
}

// Next returns next forwarder server of current chain state. Returns nil if context has not chain.
func Next(ctx context.Context) forwarder.ForwarderServer {
	if v, ok := ctx.Value(nextKey).(forwarder.ForwarderServer); ok {
		return v
	}
	return nil
}

// WithConfiguratorClient adds to context value with configurator client
func WithConfiguratorClient(ctx context.Context, endpoint string) (context.Context, func() error, error) {
	conn, err := tools.DialTCPInsecure(endpoint)
	if err != nil {
		Logger(ctx).Errorf("Can't dial grpc server: %v", err)
		return nil, nil, err
	}
	client := configurator.NewConfiguratorServiceClient(conn)

	return context.WithValue(ctx, clientKey, client), conn.Close, nil
}

//ConfiguratorClient returns configurator client or nill if client not created
func ConfiguratorClient(ctx context.Context) configurator.ConfiguratorServiceClient {
	if client, ok := ctx.Value(clientKey).(configurator.ConfiguratorServiceClient); ok {
		return client
	}
	return nil
}

//WithDataChange puts dataChange config into context
func WithDataChange(ctx context.Context, dataChange *configurator.Config) context.Context {
	return context.WithValue(ctx, dataChangeKey, dataChange)
}

//DataChange gets dataChange config from context
func DataChange(ctx context.Context) *configurator.Config {
	if dataChange, ok := ctx.Value(dataChangeKey).(*configurator.Config); ok {
		return dataChange
	}
	return nil
}

//Logger returns logger from context
func Logger(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value(loggerKey).(logrus.FieldLogger); ok {
		return logger
	}
	return logrus.New()
}

//WithLogger puts logger into context
func WithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

//WithMonitor puts into context cross connect monitor server
func WithMonitor(ctx context.Context, monitor monitor_crossconnect.MonitorServer) context.Context {
	return context.WithValue(ctx, monitorKey, monitor)
}

//MonitorServer gets from context cross connect monitor server
func MonitorServer(ctx context.Context) monitor_crossconnect.MonitorServer {
	if monitor, ok := ctx.Value(monitorKey).(monitor_crossconnect.MonitorServer); ok {
		return monitor
	}
	return nil
}
