// Copyright 2018, 2019 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package endpoint

import (
	"context"

	connectionMonitor "github.com/networkservicemesh/networkservicemesh/sdk/monitor/connectionmonitor"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	unified "github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
)

// MonitorEndpoint is a monitoring composite
type MonitorEndpoint struct {
	monitorConnectionServer connectionMonitor.MonitorServer
}

// Init will be called upon NSM Endpoint instantiation with the proper context
func (mce *MonitorEndpoint) Init(context *InitContext) error {
	grpcServer := context.GrpcServer
	unified.RegisterMonitorConnectionServer(grpcServer, mce.monitorConnectionServer)
	return nil
}

// Request implements the request handler
// Consumes from ctx context.Context:
//     ConnectionMonitor
//	   Next
func (mce *MonitorEndpoint) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	if Next(ctx) != nil {
		// Pass monitor server
		ctx = WithMonitorServer(ctx, mce.monitorConnectionServer)

		incomingConnection, err := Next(ctx).Request(ctx, request)
		if err != nil {
			Log(ctx).Errorf("Next request failed: %v", err)
			return nil, err
		}

		Log(ctx).Infof("Monitor UpdateConnection: %v", incomingConnection)
		mce.monitorConnectionServer.Update(ctx, incomingConnection)

		return incomingConnection, nil
	}
	return nil, errors.New("MonitorEndpoint.Request - cannot create requested connection")
}

// Close implements the close handler
// Request implements the request handler
// Consumes from ctx context.Context:
//     ConnectionMonitor
//	   Next
func (mce *MonitorEndpoint) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	Log(ctx).Infof("Monitor DeleteConnection: %v", connection)

	// Pass monitor server
	ctx = WithMonitorServer(ctx, mce.monitorConnectionServer)

	if Next(ctx) != nil {
		rv, err := Next(ctx).Close(ctx, connection)
		mce.monitorConnectionServer.Delete(ctx, connection)
		return rv, err
	}
	return nil, errors.New("monitor DeleteConnection cannot close connection")
}

// Name returns the composite name
func (mce *MonitorEndpoint) Name() string {
	return "monitor"
}

// NewMonitorEndpoint creates a MonitorEndpoint
func NewMonitorEndpoint(configuration *common.NSConfiguration) *MonitorEndpoint {
	// ensure the env variables are processed
	if configuration == nil {
		configuration = &common.NSConfiguration{}
	}

	self := &MonitorEndpoint{
		monitorConnectionServer: connectionMonitor.NewMonitorServer("EndpointConnection"),
	}

	return self
}
