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

package handler

const (
	NetworkServiceResource         = "networkservice"
	NetworkServiceChannelResource  = "networkservicechannel"
	NetworkServiceEndpointResource = "networkserviceendpoint"
	HandlerCreate                  = "create"
	HandlerDelete                  = "delete"
	HandlerUpdate                  = "update"
)

type NsmEvent struct {
	Key          string
	EventType    string
	ResourceType string
}

// API is the interface to a CRD handler plugin
type API interface {
	ObjectCreated(obj interface{}, event NsmEvent)
	ObjectDeleted(obj interface{}, event NsmEvent)
	ObjectUpdated(objOld, objNew interface{}, event NsmEvent)
}
