// Copyright (c) 2018 Cisco and/or its affiliates.
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

package objectstore

import (
	"time"

	"github.com/ligato/networkservicemesh/pkg/nsm/apis/netmesh"
	"github.com/ligato/networkservicemesh/plugins/idempotent"
)

const (
	// ObjectStoreReadyInterval defines readiness retry interval
	ObjectStoreReadyInterval = time.Second * 1
	// ObjectStoreReadyTimeout defines readiness timeout
	ObjectStoreReadyTimeout = time.Second * 5
)

// Interface is the interface to a ObjectStore handler plugin
type Interface interface {
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	GetNetworkService(nsName, nsNamespace string) *netmesh.NetworkService
	AddChannelToNetworkService(nsName string, nsNamespace string, ch *netmesh.NetworkServiceChannel) error
	ListNetworkServices() []*netmesh.NetworkService
	ListNetworkServiceChannels() []*netmesh.NetworkServiceChannel
	ListNetworkServiceEndpoints() []*netmesh.NetworkServiceEndpoint
}

// PluginAPI - API for the Plugin
type PluginAPI interface {
	idempotent.PluginAPI
	Interface
}
