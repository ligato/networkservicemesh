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

import (
	context "context"

	"github.com/pkg/errors"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor"
)

type event struct {
	monitor.BaseEvent
}

func (e *event) Message() (interface{}, error) {
	eventType, err := eventTypeToConnectionEventType(e.EventType())
	if err != nil {
		return nil, err
	}

	connections, err := connectionsFromEntities(e.Entities())
	if err != nil {
		return nil, err
	}

	return &connection.ConnectionEvent{
		Type:        eventType,
		Connections: connections,
	}, nil
}

type eventFactory struct {
	factoryName string
}

func (m *eventFactory) FactoryName() string {
	return m.factoryName
}
func (m *eventFactory) NewEvent(ctx context.Context, eventType monitor.EventType, entities map[string]monitor.Entity) monitor.Event {
	return &event{
		BaseEvent: monitor.NewBaseEvent(ctx, eventType, entities),
	}
}

func (m *eventFactory) EventFromMessage(ctx context.Context, message interface{}) (monitor.Event, error) {
	connectionEvent, ok := message.(*connection.ConnectionEvent)
	if !ok {
		return nil, errors.Errorf("unable to cast %v to local.ConnectionEvent", message)
	}

	eventType, err := connectionEventTypeToEventType(connectionEvent.GetType())
	if err != nil {
		return nil, err
	}

	entities := entitiesFromConnections(connectionEvent.Connections)

	return &event{
		BaseEvent: monitor.NewBaseEvent(ctx, eventType, entities),
	}, nil
}

func eventTypeToConnectionEventType(eventType monitor.EventType) (connection.ConnectionEventType, error) {
	switch eventType {
	case monitor.EventTypeInitialStateTransfer:
		return connection.ConnectionEventType_INITIAL_STATE_TRANSFER, nil
	case monitor.EventTypeUpdate:
		return connection.ConnectionEventType_UPDATE, nil
	case monitor.EventTypeDelete:
		return connection.ConnectionEventType_DELETE, nil
	default:
		return 0, errors.Errorf("unable to cast %v to local.ConnectionEventType", eventType)
	}
}

func connectionEventTypeToEventType(connectionEventType connection.ConnectionEventType) (monitor.EventType, error) {
	switch connectionEventType {
	case connection.ConnectionEventType_INITIAL_STATE_TRANSFER:
		return monitor.EventTypeInitialStateTransfer, nil
	case connection.ConnectionEventType_UPDATE:
		return monitor.EventTypeUpdate, nil
	case connection.ConnectionEventType_DELETE:
		return monitor.EventTypeDelete, nil
	default:
		return "", errors.Errorf("unable to cast %v to monitor.EventType", connectionEventType)
	}
}

func connectionsFromEntities(entities map[string]monitor.Entity) (map[string]*connection.Connection, error) {
	connections := map[string]*connection.Connection{}

	for k, v := range entities {
		if conn, ok := v.(*connection.Connection); ok {
			connections[k] = conn
		} else {
			return nil, errors.New("unable to cast Entity to connection.Connection")
		}
	}

	return connections, nil
}

func entitiesFromConnections(connections map[string]*connection.Connection) map[string]monitor.Entity {
	entities := map[string]monitor.Entity{}

	for k, v := range connections {
		entities[k] = v
	}

	return entities
}
