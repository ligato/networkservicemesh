package tests

import (
	"context"

	"github.com/ligato/vpp-agent/api/configurator"
	"github.com/pkg/errors"

	"github.com/networkservicemesh/networkservicemesh/sdk/vppagent"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/networkservicemesh/api/pkg/api/networkservice"

	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
)

// TestCommit is a VPP Agent TestCommit composite
type TestCommit struct {
	VppConfig *configurator.Config
}

// Request implements the request handler
// Provides/Consumes from ctx context.Context:
//     VppAgentConfig
//	   Next
func (c *TestCommit) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	ctx = vppagent.WithConfig(ctx) // Guarantees we will retrieve a non-nil VppAgentConfig from context.Context
	vppAgentConfig := vppagent.Config(ctx)
	if vppAgentConfig == nil {
		return nil, errors.New("received empty VppAgentConfig")
	}

	endpoint.Log(ctx).Infof("Sending VppAgentConfig to VPP Agent: %v", vppAgentConfig)
	c.VppConfig = proto.Clone(vppAgentConfig).(*configurator.Config)

	if endpoint.Next(ctx) != nil {
		return endpoint.Next(ctx).Request(ctx, request)
	}
	return request.GetConnection(), nil
}

// Close implements the close handler
// Provides/Consumes from ctx context.Context:
//     VppAgentConfig
//	   Next
func (c *TestCommit) Close(ctx context.Context, connection *networkservice.Connection) (*empty.Empty, error) {
	ctx = vppagent.WithConfig(ctx) // Guarantees we will retrieve a non-nil VppAgentConfig from context.Context
	vppAgentConfig := vppagent.Config(ctx)

	if vppAgentConfig == nil {
		return nil, errors.New("received empty vppAgentConfig")
	}

	endpoint.Log(ctx).Infof("Sending vppAgentConfig to VPP Agent: %v", vppAgentConfig)
	c.VppConfig = proto.Clone(vppAgentConfig).(*configurator.Config)

	if endpoint.Next(ctx) != nil {
		return endpoint.Next(ctx).Close(ctx, connection)
	}
	return &empty.Empty{}, nil
}

// NewTestCommit creates a new TestCommit endpoint.
func NewTestCommit() *TestCommit {
	return &TestCommit{}
}

// Init will reset the vpp shouldResetVpp is true
func (c *TestCommit) Init(*endpoint.InitContext) error {
	return nil
}
