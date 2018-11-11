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

package main

import (
	"context"
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/ligato/networkservicemesh/controlplane/pkg/model/networkservice"
	"github.com/ligato/networkservicemesh/controlplane/pkg/model/registry"

	"github.com/ligato/networkservicemesh/pkg/tools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	// networkServiceName defines Network Service Name the NSE is serving for
	networkServiceName = "icmp-responder"
	// SocketBaseDir defines the location of NSM Endpoints listen socket
	SocketBaseDir = "/var/lib/networkservicemesh"
	// RegistrySocketFile defines the name of NSM Endpoints operations socket
	RegistrySocketFile = "nsm.client.io.sock"
)

func main() {
	var wg sync.WaitGroup

	// For NSE to program container's dataplane, container's linux namespace must be sent to NSM
	linuxNS, err := tools.GetCurrentNS()
	if err != nil {
		logrus.Fatalf("nse: failed to get a linux namespace with error: %+v, exiting...", err)
		os.Exit(1)
	}
	logrus.Infof("Starting NSE, linux namespace: %s", linuxNS)

	// NSM socket path will be used to drop NSE socket for NSM's Connection request
	connectionServerSocket := path.Join(SocketBaseDir, linuxNS+".nse.io.sock")
	if err := tools.SocketCleanup(connectionServerSocket); err != nil {
		logrus.Fatalf("nse: failure to cleanup stale socket %s with error: %+v", connectionServerSocket, err)
	}

	logrus.Infof("nse: listening socket %s", connectionServerSocket)
	connectionServer, err := net.Listen("unix", connectionServerSocket)
	if err != nil {
		logrus.Fatalf("nse: failure to listen on a socket %s with error: %+v", connectionServerSocket, err)
	}
	grpcServer := grpc.NewServer()

	// Registering NSE API, it will listen for Connection requests from NSM and return information
	// needed for NSE's dataplane programming.
	nseConn := New()

	networkservice.RegisterNetworkServiceServer(grpcServer, nseConn)

	go func() {
		wg.Add(1)
		if err := grpcServer.Serve(connectionServer); err != nil {
			logrus.Fatalf("nse: failed to start grpc server on socket %s with error: %+v ", connectionServerSocket, err)
		}
	}()
	// Check if the socket of Endpoint Connection Server is operation
	testSocket, err := tools.SocketOperationCheck(connectionServerSocket)
	if err != nil {
		logrus.Fatalf("nse: failure to communicate with the connectionServerSocket %s with error: %+v", connectionServerSocket, err)
	}
	testSocket.Close()

	// NSE connection server is ready and now endpoints can be advertised to NSM
	registrySocket := path.Join(SocketBaseDir, RegistrySocketFile)

	if _, err := os.Stat(registrySocket); err != nil {
		logrus.Errorf("nse: failure to access nsm socket at %s with error: %+v, exiting...", registrySocket, err)
		os.Exit(1)
	}

	conn, err := tools.SocketOperationCheck(registrySocket)
	if err != nil {
		logrus.Fatalf("nse: failure to communicate with the registrySocket %s with error: %+v", registrySocket, err)
	}
	defer conn.Close()
	logrus.Infof("nsm: connection to nsm server on socket: %s succeeded.", registrySocket)

	registryConnection := registry.NewNetworkServiceRegistryClient(conn)

	nseid := rand.Uint64()

	nse := &registry.NetworkServiceEndpoint{
		NetworkServiceName: networkServiceName,
		EndpointName:       networkServiceName + "-" + strconv.FormatUint(nseid, 36),
		Payload:            "IP",
		Labels:             make(map[string]string),
		SocketLocation:     connectionServerSocket,
	}

	registeredNSE, err := registryConnection.RegisterNSE(context.Background(), nse)
	if err != nil {
		logrus.Fatalln("unable to register endpoint", err)
	}
	logrus.Infoln("NSE registered: " + registeredNSE.EndpointName)

	logrus.Infof("nse: channel has been successfully advertised, waiting for connection from NSM...")
	// Now block on WaitGroup
	wg.Wait()
}
