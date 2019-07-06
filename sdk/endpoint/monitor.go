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
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/monitor/local"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/sirupsen/logrus"
)

// MonitorEndpoint is a monitoring composite
type MonitorEndpoint struct {
	monitorConnectionServer local.MonitorServer
}

// Request implements the request handler
// Consumes from ctx context.Context:
//     MonitorServer
//	   Next
func (mce *MonitorEndpoint) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	if Next(ctx) != nil {
		incomingConnection, err := Next(ctx).Request(ctx, request)
		if err != nil {
			logrus.Errorf("Next request failed: %v", err)
			return nil, err
		}

		logrus.Infof("Monitor UpdateConnection: %v", incomingConnection)
		MonitorServer(ctx).Update(incomingConnection)

		return incomingConnection, nil
	}
	return nil, fmt.Errorf("MonitorEndpoint.Request - cannot create requested connection")
}

// Close implements the close handler
// Request implements the request handler
// Consumes from ctx context.Context:
//     MonitorServer
//	   Next
func (mce *MonitorEndpoint) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	logrus.Infof("Monitor DeleteConnection: %v", connection)
	if Next(ctx) != nil {
		rv, err := Next(ctx).Close(ctx, connection)
		MonitorServer(ctx).Delete(connection)
		return rv, err
	}
	return nil, fmt.Errorf("Monitor DeleteConnection cannot close connection")
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
	configuration.CompleteNSConfiguration()

	self := &MonitorEndpoint{
		monitorConnectionServer: local.NewMonitorServer(),
	}

	return self
}
