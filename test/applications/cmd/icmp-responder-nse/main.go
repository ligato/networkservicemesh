// Copyright 2018 VMware, Inc.
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

package main

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/monitor/local"

	"github.com/networkservicemesh/networkservicemesh/test/applications/cmd/icmp-responder-nse/flags"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/networkservice"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
)

var version string

func main() {
	logrus.Info("Starting icmp-responder-nse...")
	logrus.Infof("Version: %v", version)

	flags := flags.ParseFlags()
	// Capture signals to cleanup before exiting
	c := tools.NewOSSignalChannel()

	endpoints := []networkservice.NetworkServiceServer{
		endpoint.NewMonitorEndpoint(nil),
		endpoint.NewConnectionEndpoint(nil),
	}

	if flags.Neighbors {
		logrus.Infof("Adding neighbors endpoint to chain")
		endpoints = append(endpoints,
			endpoint.NewCustomFuncEndpoint("neighbor", ipNeighborMutator))
	}

	ipamEndpoint := endpoint.NewIpamEndpoint(nil)
	endpoints = append(endpoints, ipamEndpoint)

	routeAddr := endpoint.CreateRouteMutator([]string{"8.8.8.8/30"})
	if common.IsIPv6(ipamEndpoint.PrefixPool.GetPrefixes()[0]) {
		routeAddr = endpoint.CreateRouteMutator([]string{"2001:4860:4860::8888/126"})
	}
	if flags.Routes {
		logrus.Infof("Adding routes endpoint to chain")
		endpoints = append(endpoints, endpoint.NewCustomFuncEndpoint("route", routeAddr))
	}

	if flags.DNS {
		logrus.Info("Adding dns endpoint to chain")
		endpoints = append(endpoints, endpoint.NewCustomFuncEndpoint("dns", dnsMutator))
	}

	if flags.Update {
		logrus.Infof("Adding updating endpoint to chain")
		endpoints = append(endpoints,
			endpoint.NewCustomFuncEndpoint("update", func(ctx context.Context, conn *connection.Connection) error {
				monitorServer := endpoint.MonitorServer(ctx)
				logrus.Infof("Delaying 5 seconds before send update event.")
				go func() {
					<-time.After(5 * time.Second)
					updateConnections(monitorServer)
					logrus.Infof("Update event sended.")
				}()
				return nil
			}))
	}

	composite := endpoint.NewCompositeEndpoint(endpoints...)

	nsmEndpoint, err := endpoint.NewNSMEndpoint(context.Background(), nil, composite)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	_ = nsmEndpoint.Start()
	if !flags.Dirty {
		defer func() { _ = nsmEndpoint.Delete() }()
	}

	// Capture signals to cleanup before exiting
	<-c
}

func dnsMutator(ctc context.Context, c *connection.Connection) error {
	defaultIP := strings.Split(c.Context.IpContext.DstIpAddr, "/")[0]
	c.Context.DnsContext = &connectioncontext.DNSContext{
		Configs: []*connectioncontext.DNSConfig{
			{
				DnsServerIps:  ServerIPsEnv.GetStringListValueOrDefault(defaultIP),
				SearchDomains: SearchDomainsEnv.GetStringListValueOrDefault("icmp.app"),
			},
		},
	}
	return nil
}

func ipNeighborMutator(ctc context.Context, c *connection.Connection) error {
	addrs, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range addrs {
		adrs, err := iface.Addrs()
		if err != nil {
			logrus.Error(err)
			continue
		}

		for _, a := range adrs {
			addr, _, _ := net.ParseCIDR(a.String())
			if !addr.IsLoopback() {
				c.GetContext().GetIpContext().IpNeighbors = append(c.GetContext().GetIpContext().GetIpNeighbors(),
					&connectioncontext.IpNeighbor{
						Ip:              addr.String(),
						HardwareAddress: iface.HardwareAddr.String(),
					},
				)
			}
		}
	}
	return nil
}

func updateConnections(monitorServer local.MonitorServer) {
	for _, entity := range monitorServer.Entities() {
		localConnection := proto.Clone(entity.(*connection.Connection)).(*connection.Connection)
		localConnection.GetContext().GetIpContext().ExcludedPrefixes =
			append(localConnection.GetContext().GetIpContext().GetExcludedPrefixes(), "255.255.255.255/32")

		monitorServer.Update(localConnection)
	}
}
