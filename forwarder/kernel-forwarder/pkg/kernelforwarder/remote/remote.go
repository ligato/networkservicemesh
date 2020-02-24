// Copyright (c) 2020 VMware, Inc.
//
// Copyright (c) 2020 Doc.ai and/or its affiliates.
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

// Package remote - controlling remote mechanisms interfaces
package remote

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	sdk "github.com/networkservicemesh/networkservicemesh/forwarder/sdk/wireguard"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/vxlan"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/wireguard"
)

// INCOMING, OUTGOING - packet direction constants
const (
	INCOMING = iota
	OUTGOING = iota
)

// Connect - struct with remote mechanism interfaces creation and deletion methods
type Connect struct {
	*sdk.DeviceManager
}

// NewConnect - creates instance of remote Connect
func NewConnect() *Connect {
	return &Connect{
		DeviceManager: sdk.NewWireguardDeviceManager(false),
	}
}

// CreateInterface - creates interface to remote connection
func (c *Connect) CreateInterface(ifaceName string, remoteConnection *connection.Connection, direction uint8) error {
	logrus.Infof("Create interface for %v", remoteConnection)
	switch remoteConnection.GetMechanism().GetType() {
	case vxlan.MECHANISM:
		return c.createVXLANInterface(ifaceName, remoteConnection, direction)
	case wireguard.MECHANISM:
		return c.CreateWireguardDevice(ifaceName, remoteConnection, direction == INCOMING)
	}
	return errors.Errorf("unknown remote mechanism - %v", remoteConnection.GetMechanism().GetType())
}

// DeleteInterface - deletes interface to remote connection
func (c *Connect) DeleteInterface(ifaceName string, remoteConnection *connection.Connection) error {
	switch remoteConnection.GetMechanism().GetType() {
	case vxlan.MECHANISM:
		return c.deleteVXLANInterface(ifaceName)
	case wireguard.MECHANISM:
		c.DeleteWireguardDevice(ifaceName)
		return nil
	}
	return errors.Errorf("unknown remote mechanism - %v", remoteConnection.GetMechanism().GetType())
}
