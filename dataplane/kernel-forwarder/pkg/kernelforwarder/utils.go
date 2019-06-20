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
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/crossconnect"
	local "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	remote "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/connection"

	"github.com/networkservicemesh/networkservicemesh/dataplane/pkg/common"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"runtime"
)

type KernelConnectionConfig struct {
	srcNsPath string
	dstNsPath string
	srcName   string
	dstName   string
	srcIP     string
	dstIP     string
}

func handleKernelConnectionLocal(crossConnect *crossconnect.CrossConnect, connect bool) (*crossconnect.CrossConnect, error) {
	/* Create a connection */
	if connect {
		/* 1. Get the connection configuration */
		cfg, err := getConnectionConfig(crossConnect)
		if err != nil {
			logrus.Errorf("Failed to get the configuration for local connection - %v", err)
			return crossConnect, err
		}
		/* 2. Get namespace handlers from their path - source and destination */
		srcNsHandle, err := netns.GetFromPath(cfg.srcNsPath)
		defer srcNsHandle.Close()
		if err != nil {
			logrus.Errorf("Failed to get source namespace handler from path - %v", err)
			return crossConnect, err
		}
		dstNsHandle, err := netns.GetFromPath(cfg.dstNsPath)
		defer dstNsHandle.Close()
		if err != nil {
			logrus.Errorf("Failed to get destination namespace handler from path - %v", err)
			return crossConnect, err
		}
		/* 3. Create a VETH pair and inject each end in the corresponding namespace */
		if err = createVETH(cfg, srcNsHandle, dstNsHandle); err != nil {
			logrus.Errorf("Failed to create the VETH pair - %v", err)
			return crossConnect, err
		}
		/* 4. Bring up and configure each pair end with its IP address */
		setupVETHEnd(srcNsHandle, cfg.srcName, cfg.srcIP)
		setupVETHEnd(dstNsHandle, cfg.dstName, cfg.dstIP)
	}
	/* Delete a connection */
	return crossConnect, nil
}

func handleKernelConnectionRemote(crossConnect *crossconnect.CrossConnect, connect bool) (*crossconnect.CrossConnect, error) {
	if crossConnect.GetRemoteSource().GetMechanism().GetType() == remote.MechanismType_VXLAN &&
		crossConnect.GetLocalDestination().GetMechanism().GetType() == local.MechanismType_KERNEL_INTERFACE {
		/* 1. Incoming connection */
		logrus.Info("Incoming connection - remote source/local destination")
		return handleKernelConnectionRemoteIncoming(crossConnect, connect)
	} else if crossConnect.GetLocalSource().GetMechanism().GetType() == local.MechanismType_KERNEL_INTERFACE &&
		crossConnect.GetRemoteDestination().GetMechanism().GetType() == remote.MechanismType_VXLAN {
		/* 2. Outgoing connection */
		logrus.Info("Outgoing connection - local source/remote destination")
		return handleKernelConnectionRemoteOutgoing(crossConnect, connect)
	} else {
		logrus.Errorf("Remote connection is not supported yet.")
	}
	return crossConnect, nil
}

func handleKernelConnectionRemoteIncoming(crossConnect *crossconnect.CrossConnect, connect bool) (*crossconnect.CrossConnect, error) {
	/* Create a connection */
	if connect {

	}
	/* Delete a connection */
	return crossConnect, nil
}

func handleKernelConnectionRemoteOutgoing(crossConnect *crossconnect.CrossConnect, connect bool) (*crossconnect.CrossConnect, error) {
	/* Create a connection */
	if connect {

	}
	/* Delete a connection */
	return crossConnect, nil
}

func getConnectionConfig(crossConnect *crossconnect.CrossConnect) (*KernelConnectionConfig, error) {
	srcNsPath, err := crossConnect.GetLocalSource().GetMechanism().NetNsFileName()
	if err != nil {
		logrus.Errorf("Failed to get source namespace path - %v", err)
		return nil, err
	}
	dstNsPath, err := crossConnect.GetLocalDestination().GetMechanism().NetNsFileName()
	if err != nil {
		logrus.Errorf("Failed to get destination namespace path - %v", err)
		return nil, err
	}
	return &KernelConnectionConfig{
		srcNsPath: srcNsPath,
		dstNsPath: dstNsPath,
		srcName:   crossConnect.GetLocalSource().GetMechanism().GetParameters()[local.InterfaceNameKey],
		dstName:   crossConnect.GetLocalDestination().GetMechanism().GetParameters()[local.InterfaceNameKey],
		srcIP:     crossConnect.GetLocalSource().GetContext().GetSrcIpAddr(),
		dstIP:     crossConnect.GetLocalSource().GetContext().GetDstIpAddr(),
	}, nil
}

func setupVETHEnd(nsHandle netns.NsHandle, ifName, addrIP string) {
	/* Lock the OS thread so we don't accidentally switch namespaces */
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	/* Save the current network namespace */
	oNsHandle, _ := netns.Get()
	defer oNsHandle.Close()

	/* Switch to the new namespace */
	netns.Set(nsHandle)
	defer nsHandle.Close()

	/* Get a link for the interface name */
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		logrus.Errorf("Failed to lookup %q, %v", ifName, err)
	}

	/* Setup the interface with an IP address */
	addr, _ := netlink.ParseAddr(addrIP)
	netlink.AddrAdd(link, addr)

	/* Bring the interface UP */
	if err = netlink.LinkSetUp(link); err != nil {
		logrus.Errorf("Failed to set %q up: %v", ifName, err)
	}

	/* Switch back to the original namespace */
	netns.Set(oNsHandle)
}

func createVETH(cfg *KernelConnectionConfig, srcNsHandle, dstNsHandle netns.NsHandle) error {
	/* Initial VETH configuration */
	cfgVETH := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: cfg.srcName,
			MTU:  16000,
		},
		PeerName: cfg.dstName,
	}

	/* Create the VETH pair - host namespace */
	if err := netlink.LinkAdd(cfgVETH); err != nil {
		logrus.Errorf("Failed to create the VETH pair - %v", err)
		return err
	}

	/* Get a link for each VETH pair ends */
	srcLink, err := netlink.LinkByName(cfg.srcName)
	if err != nil {
		logrus.Errorf("Failed to get source link from name - %v", err)
		return err
	}
	dstLink, err := netlink.LinkByName(cfg.dstName)
	if err != nil {
		logrus.Errorf("Failed to get destination link from name - %v", err)
		return err
	}

	/* Inject each end in its corresponding client/endpoint namespace */
	if err = netlink.LinkSetNsFd(srcLink, int(srcNsHandle)); err != nil {
		logrus.Errorf("Failed to inject the VETH end in the source namespace - %v", err)
		return err
	}
	if err = netlink.LinkSetNsFd(dstLink, int(dstNsHandle)); err != nil {
		logrus.Errorf("Failed to inject the VETH end in the destination namespace - %v", err)
		return err
	}
	return nil
}

func (v *KernelForwarder) configureKernelForwarder() {
	v.common.MechanismsUpdateChannel = make(chan *common.Mechanisms, 1)
	v.common.Mechanisms = &common.Mechanisms{
		LocalMechanisms: []*local.Mechanism{
			{
				Type: local.MechanismType_KERNEL_INTERFACE,
			},
		},
		RemoteMechanisms: []*remote.Mechanism{
			{
				Type: remote.MechanismType_VXLAN,
				Parameters: map[string]string{
					remote.VXLANSrcIP: v.common.EgressInterface.SrcIPNet().IP.String(),
				},
			},
		},
	}
}
