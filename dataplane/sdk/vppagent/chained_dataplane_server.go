package forwarder

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/dataplane/api/dataplane"
)

type chainedDataplaneServer struct {
	handlers []dataplane.DataplaneServer
}

func (c *chainedDataplaneServer) Request(ctx context.Context, crossConnect *crossconnect.CrossConnect) (*crossconnect.CrossConnect, error) {

	if len(c.handlers) == 0 {
		logrus.Info("chainedDataplaneServer: has not handlers for next request")
		return crossConnect, nil
	}
	next := &next{handlers: c.handlers, index: 0}
	nextCtx := withNext(ctx, next)
	return next.Request(nextCtx, crossConnect)

}

func (c *chainedDataplaneServer) Close(ctx context.Context, crossConnect *crossconnect.CrossConnect) (*empty.Empty, error) {
	if len(c.handlers) == 0 {
		logrus.Info("chainedDataplaneServer: has not handlers for next close")
		return new(empty.Empty), nil
	}
	next := &next{handlers: c.handlers, index: 0}
	nextCtx := withNext(ctx, next)
	return next.Close(nextCtx, crossConnect)
}

// ChainOf makes chain of dataplane server handlers
func ChainOf(handlers ...dataplane.DataplaneServer) dataplane.DataplaneServer {
	return &chainedDataplaneServer{handlers: handlers}
}
