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
	"context"
	"net"

	"github.com/networkservicemesh/networkservicemesh/dataplane/pkg/common"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/crossconnect"
	local "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	remote "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/connection"
	monitor_crossconnect "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/monitor/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/dataplane/pkg/apis/dataplane"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
)

// Kernel forwarding plane related constants
const (
	DataplaneNameKey                              = "DATAPLANE_NAME"
	DataplaneNameDefault                          = "kernel-forwarder"
	DataplaneSocketKey                            = "DATAPLANE_SOCKET"
	DataplaneSocketDefault                        = "/var/lib/networkservicemesh/nsm-vppagent.dataplane.sock"
	DataplaneSocketTypeKey                        = "DATAPLANE_SOCKET_TYPE"
	DataplaneSocketTypeDefault                    = "unix"
	SrcIPEnvKey                                   = "NSM_DATAPLANE_SRC_IP"
	cCONNECT = true
	cDISCONNECT = false
)

type KernelForwarder struct {
	common               *common.DataplaneConfigBase
	mechanisms           *Mechanisms
	updateCh             chan *Mechanisms
	srcIP                net.IP
	egressInterface      common.EgressInterface
	monitor              monitor_crossconnect.MonitorServer
}

func CreateKernelForwarder() *KernelForwarder {
	return &KernelForwarder{}
}

// Mechanisms is a message used to communicate any changes in operational parameters and constraints
type Mechanisms struct {
	remoteMechanisms []*remote.Mechanism
	localMechanisms  []*local.Mechanism
}

func (v *KernelForwarder) MonitorMechanisms(empty *empty.Empty, updateSrv dataplane.Dataplane_MonitorMechanismsServer) error {
	logrus.Infof("MonitorMechanisms was called")
	initialUpdate := &dataplane.MechanismUpdate{
		RemoteMechanisms: v.mechanisms.remoteMechanisms,
		LocalMechanisms:  v.mechanisms.localMechanisms,
	}
	logrus.Infof("Sending MonitorMechanisms update: %v", initialUpdate)
	if err := updateSrv.Send(initialUpdate); err != nil {
		logrus.Errorf("Kernel forwarding plane server: Detected error %s, grpc code: %+v on grpc channel", err.Error(), status.Convert(err).Code())
		return nil
	}
	for {
		select {
		// Waiting for any updates which might occur during a life of dataplane module and communicating
		// them back to NSM.
		case update := <-v.updateCh:
			v.mechanisms = update
			logrus.Infof("Sending MonitorMechanisms update: %v", update)
			if err := updateSrv.Send(&dataplane.MechanismUpdate{
				RemoteMechanisms: update.remoteMechanisms,
				LocalMechanisms:  update.localMechanisms,
			}); err != nil {
				logrus.Errorf("Kernel forwarding plane server: Detected error %s, grpc code: %+v on grpc channel", err.Error(), status.Convert(err).Code())
				return nil
			}
		}
	}
}

func (v *KernelForwarder) Request(ctx context.Context, crossConnect *crossconnect.CrossConnect) (*crossconnect.CrossConnect, error) {
	logrus.Infof("Request() called with %v", crossConnect)
	xcon, err := v.connectOrDisconnect(ctx, crossConnect, cCONNECT)
	if err != nil {
		return nil, err
	}
	v.monitor.Update(xcon)
	logrus.Infof("Request() called with %v returning: %v", crossConnect, xcon)
	return xcon, err
}

func (v *KernelForwarder) connectOrDisconnect(ctx context.Context, crossConnect *crossconnect.CrossConnect, connect bool) (*crossconnect.CrossConnect, error) {
	/* 1. Handle local connection */
	if crossConnect.GetLocalSource().GetMechanism().GetType() == local.MechanismType_KERNEL_INTERFACE &&
		crossConnect.GetLocalDestination().GetMechanism().GetType() == local.MechanismType_KERNEL_INTERFACE {
		return handleKernelConnectionLocal(crossConnect, connect)
	}
	/* 2. Handle remote connection */
	return handleKernelConnectionRemote(ctx, crossConnect, connect)
}

func (v *KernelForwarder) Close(ctx context.Context, crossConnect *crossconnect.CrossConnect) (*empty.Empty, error) {
	logrus.Infof("Close() called with %#v", crossConnect)
	xcon, err := v.connectOrDisconnect(ctx, crossConnect, cDISCONNECT)
	if err != nil {
		logrus.Warn(err)
	}
	v.monitor.Delete(xcon)
	return &empty.Empty{}, nil
}

// Init makes setup for the KernelForwarder
func (v *KernelForwarder) Init(common *common.DataplaneConfigBase, monitor monitor_crossconnect.MonitorServer) error {
	tracer, closer := tools.InitJaeger("kernel-forwarder")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	setDataplaneConfigBase(v, common)
	setDataplaneConfigKernelForwarder(v, monitor)
	return nil
}
