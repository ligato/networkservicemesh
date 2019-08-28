// Copyright 2019 VMware, Inc.
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

package kernelforwarder

import (
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/crossconnect"
	local "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	remote "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/connection"

	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"

	"github.com/networkservicemesh/networkservicemesh/dataplane/pkg/common"
)

func (v *KernelForwarder) handleRemoteConnection(egress common.EgressInterfaceType, crossConnect *crossconnect.CrossConnect, connect bool) (*crossconnect.CrossConnect, error) {
	var devices map[string]string
	var xcon *crossconnect.CrossConnect
	var err error
	if crossConnect.GetRemoteSource().GetMechanism().GetType() == remote.MechanismType_VXLAN &&
		crossConnect.GetLocalDestination().GetMechanism().GetType() == local.MechanismType_KERNEL_INTERFACE {
		/* 1. Incoming remote connection */
		logrus.Info("Incoming connection - remote source/local destination")
		xcon, devices, err = handleConnection(egress, crossConnect, connect, cINCOMING)
	} else if crossConnect.GetLocalSource().GetMechanism().GetType() == local.MechanismType_KERNEL_INTERFACE &&
		crossConnect.GetRemoteDestination().GetMechanism().GetType() == remote.MechanismType_VXLAN {
		/* 2. Outgoing remote connection */
		logrus.Info("Outgoing connection - local source/remote destination")
		xcon, devices, err = handleConnection(egress, crossConnect, connect, cOUTGOING)
	} else {
		logrus.Errorf("invalid remote connection type")
		return crossConnect, fmt.Errorf("invalid remote connection type")
	}
	if devices != nil {
		v.updateDeviceList(devices, connect)
	}
	return xcon, err
}

func handleConnection(egress common.EgressInterfaceType, crossConnect *crossconnect.CrossConnect, connect bool, direction uint8) (*crossconnect.CrossConnect, map[string]string, error) {
	var devices map[string]string
	/* 1. Get the connection configuration */
	cfg, err := newConnectionConfig(crossConnect, direction)
	if err != nil {
		logrus.Errorf("failed to get the configuration for remote connection - %v", err)
		return crossConnect, nil, err
	}
	nsPath, name, ifaceIP, vxlanIP, routes := modifyConfiguration(cfg, direction)
	if connect {
		/* 2. Create a connection */
		devices, err = createRemoteConnection(nsPath, name, ifaceIP, egress.SrcIPNet().IP, vxlanIP, cfg.vni, routes, cfg.neighbors)
		if err != nil {
			logrus.Errorf("failed to create remote connection - %v", err)
		}
	} else {
		/* 3. Delete a connection */
		devices, err = deleteRemoteConnection(nsPath, name)
		if err != nil {
			logrus.Errorf("failed to delete remote connection - %v", err)
		}
	}
	return crossConnect, devices, err
}

func createRemoteConnection(nsPath, ifaceName, ifaceIP string, egressIP, remoteIP net.IP, vni int, routes []*connectioncontext.Route, neighbors []*connectioncontext.IpNeighbor) (map[string]string, error) {
	logrus.Info("Creating remote connection...")
	/* 1. Get handler for container namespace */
	containerNs, err := netns.GetFromPath(nsPath)
	defer func() {
		if err = containerNs.Close(); err != nil {
			logrus.Error("error when closing:", err)
		}
	}()
	if err != nil {
		logrus.Errorf("failed to get namespace handler from path - %v", err)
		return nil, err
	}

	/* 2. Prepare interface - VXLAN */
	iface := newVXLAN(ifaceName, egressIP, remoteIP, vni)

	/* 3. Create interface - host namespace */
	if err = netlink.LinkAdd(iface); err != nil {
		logrus.Errorf("failed to create VXLAN interface - %v", err)
	}

	/* 4. Setup interface */
	if err = setupLinkInNs(containerNs, ifaceName, ifaceIP, routes, neighbors, true); err != nil {
		logrus.Errorf("failed to setup container interface %q: %v", ifaceName, err)
		return nil, err
	}
	return map[string]string{ifaceName: nsPath}, nil
}

func deleteRemoteConnection(nsPath, ifaceName string) (map[string]string, error) {
	logrus.Info("Deleting remote connection...")
	/* 1. Get handler for container namespace */
	containerNs, err := netns.GetFromPath(nsPath)
	defer func() {
		if err = containerNs.Close(); err != nil {
			logrus.Error("error when closing:", err)
		}
	}()
	if err != nil {
		logrus.Errorf("failed to get namespace handler from path - %v", err)
		return nil, err
	}

	/* 2. Setup interface */
	if err = setupLinkInNs(containerNs, ifaceName, "", nil, nil, false); err != nil {
		logrus.Errorf("failed to setup container interface %q: %v", ifaceName, err)
		return nil, err
	}

	/* 3. Get a link object for the interface */
	ifaceLink, err := netlink.LinkByName(ifaceName)
	if err != nil {
		logrus.Errorf("failed to get link for %q - %v", ifaceName, err)
		return nil, err
	}

	/* 4. Delete the VXLAN interface - host namespace */
	if err := netlink.LinkDel(ifaceLink); err != nil {
		logrus.Errorf("failed to delete the VXLAN - %v", err)
		return nil, err
	}
	return map[string]string{ifaceName: nsPath}, nil
}

/* modifyConfiguration swaps the values based on the direction of the connection - incoming or outgoing */
func modifyConfiguration(cfg *connectionConfig, direction uint8) (string, string, string, net.IP, []*connectioncontext.Route) {
	if direction == cINCOMING {
		return cfg.dstNsPath, cfg.dstName, cfg.dstIP, cfg.srcIPVXLAN, cfg.dstRoutes
	}
	return cfg.srcNsPath, cfg.srcName, cfg.srcIP, cfg.dstIPVXLAN, cfg.srcRoutes
}

func newVXLAN(ifaceName string, egressIP, remoteIP net.IP, vni int) *netlink.Vxlan {
	/* Populate the VXLAN interface configuration */
	return &netlink.Vxlan{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifaceName,
		},
		VxlanId: vni,
		Group:   remoteIP,
		SrcAddr: egressIP,
	}
}
