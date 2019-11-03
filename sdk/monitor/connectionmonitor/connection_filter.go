// Copyright (c) 2019 Cisco Systems, Inc and/or its affiliates.
//
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

package connectionmonitor

import "github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"

type monitorConnectionFilter struct {
	connection.MonitorConnection_MonitorConnectionsServer

	selector *connection.MonitorScopeSelector
}

// NewMonitorConnectionFilter - create a connection montior server filter
func NewMonitorConnectionFilter(selector *connection.MonitorScopeSelector, monitor connection.MonitorConnection_MonitorConnectionsServer) connection.MonitorConnection_MonitorConnectionsServer {
	return &monitorConnectionFilter{
		selector: selector,
		MonitorConnection_MonitorConnectionsServer: monitor,
	}
}

// Send filters event connections and pass it to the next sending layer
func (d *monitorConnectionFilter) Send(in *connection.ConnectionEvent) error {
	out := &connection.ConnectionEvent{
		Type:        in.Type,
		Connections: make(map[string]*connection.Connection),
	}
	for key, value := range in.GetConnections() {
		if len(d.selector.NetworkServiceManagers) > 0 && value.GetSourceNetworkServiceManagerName() == d.selector.NetworkServiceManagers[0] {
			out.Connections[key] = value
		}
		if len(d.selector.NetworkServiceManagers) > 1 && value.GetDestinationNetworkServiceManagerName() == d.selector.NetworkServiceManagers[1] {
			out.Connections[key] = value
		}
	}
	if len(out.Connections) > 0 || out.Type == connection.ConnectionEventType_INITIAL_STATE_TRANSFER {
		return d.MonitorConnection_MonitorConnectionsServer.Send(out)
	}
	return nil
}
