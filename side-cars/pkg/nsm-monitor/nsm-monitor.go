// Copyright (c) 2019 Cisco and/or its affiliates.
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

package nsmmonitor

import (
	"context"
	"fmt"
	"io"
	"time"

	nsminit "github.com/networkservicemesh/networkservicemesh/side-cars/pkg/nsm-init"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/networkservice"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/sdk/client"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor"
	"github.com/networkservicemesh/networkservicemesh/sdk/monitor/local"
)

const (
	nsmMonitorLogFormat          = "NSM Monitor: %v"
	nsmMonitorLogWithParamFormat = "NSM Monitor: %v: %v"

	nsmMonitorRetryDelay = 5 * time.Second
)

// Handler - handler to perform configuration of monitoring app
type Handler interface {
	//Connected occurs when the nsm-monitor connected
	Connected(map[string]*connection.Connection)
	//Healing occurs when the healing started
	Healing(conn *connection.Connection)
	//Closed occurs when the connection closed
	Closed(conn *connection.Connection)
	//GetConfiguration gets custom network service configuration
	GetConfiguration() *common.NSConfiguration
	//ProcessHealing occurs when the restore failed, the error pass as the second parameter
	ProcessHealing(newConn *connection.Connection, e error)
	//Stopped occurs when the invoked App.Stop()
	Stopped()
	//IsEnableJaeger returns is Jaeger needed
	IsEnableJaeger() bool
}

// App - application to perform monitoring.
type App interface {
	nsminit.NSMApp
	// SetHandler - sets a handler instance
	SetHandler(helper Handler)
	Stop()
}

//EmptyNSMMonitorHandler has empty implementation of each method of interface Handler
type EmptyNSMMonitorHandler struct {
}

//Connected occurs when the nsm-monitor connected
func (h *EmptyNSMMonitorHandler) Connected(map[string]*connection.Connection) {}

//Healing occurs when the healing started
func (h *EmptyNSMMonitorHandler) Healing(conn *connection.Connection) {}

//Closed occurs when the connection closed
func (h *EmptyNSMMonitorHandler) Closed(conn *connection.Connection) {}

//GetConfiguration returns nil by default
func (h *EmptyNSMMonitorHandler) GetConfiguration() *common.NSConfiguration { return nil }

//ProcessHealing occurs when the restore failed, the error pass as the second parameter
func (h *EmptyNSMMonitorHandler) ProcessHealing(newConn *connection.Connection, e error) {}

//Stopped occurs when the invoked App.Stop()
func (h *EmptyNSMMonitorHandler) Stopped() {}

//IsEnableJaeger returns false by default
func (h *EmptyNSMMonitorHandler) IsEnableJaeger() bool {
	return tools.IsOpentracingEnabled()
}

type nsmMonitorApp struct {
	connections map[string]*connection.Connection
	helper      Handler
	stop        chan struct{}

	initRecieved  bool
	recovery      bool
	configuration *common.NSConfiguration
}

func (c *nsmMonitorApp) Stop() {
	close(c.stop)
}

func (c *nsmMonitorApp) SetHandler(listener Handler) {
	c.helper = listener
}

func (c *nsmMonitorApp) Run() {
	// Capture signals to cleanup before exiting
	var tracingCloser io.Closer
	var tracer opentracing.Tracer
	if c.helper == nil || c.helper.IsEnableJaeger() {
		tracer, tracingCloser = tools.InitJaeger("nsm-monitor")
		opentracing.SetGlobalTracer(tracer)
	}

	go func() {
		c.beginMonitoring()
		if tracingCloser != nil {
			defer func() {
				if err := tracingCloser.Close(); err != nil {
					logrus.Errorf("error closing opentracing context %v", err)
				}
			}()
		}
	}()
}

// NewNSMMonitorApp - creates a monitoring application.
func NewNSMMonitorApp(configuration *common.NSConfiguration) App {
	return &nsmMonitorApp{
		connections:   map[string]*connection.Connection{},
		stop:          make(chan struct{}),
		configuration: configuration,
	}
}

func (c *nsmMonitorApp) beginMonitoring() {
	for {
		nsmClient, err := client.NewNSMClient(context.Background(), c.configuration)
		if err != nil {
			logrus.Errorf(nsmMonitorLogWithParamFormat, "unable to create the NSM client", err)

			c.waitRetry()
			continue
		}

		logrus.Infof(nsmMonitorLogFormat, "connection to NSM established")

		monitorClient, err := local.NewMonitorClient(nsmClient.NsmConnection.GrpcClient)
		if err != nil {
			logrus.Errorf(nsmMonitorLogWithParamFormat, "failed to start monitor client", err)

			c.waitRetry()
			if err := nsmClient.Destroy(context.Background()); err != nil {
				logrus.Errorf("failed to close NSM client connection")
			}
			continue
		}
		defer monitorClient.Close()

		for {
			if c.initRecieved && !c.recovery {
				// Performing recovery if required.
				if c.helper != nil {
					c.helper.Connected(c.connections)
				}
				// Since NSMD will setup public socket only when all connections will be ok, we need to perform request only on ones it loose.
				if c.performRecovery(nsmClient) {
					// since we not recovered, we will continue after delay
					c.waitRetry()
					continue
				} else {
					c.recovery = true
				}
			}
			if !c.readEvents(monitorClient) {
				break // If someting happened we need to retry
			}
		}

		// Close current NSM client connection.
		if err := nsmClient.Destroy(context.Background()); err != nil {
			logrus.Errorf("failed to close NSM client connection")
		}
	}
}

func (c *nsmMonitorApp) readEvents(monitorClient monitor.Client) bool {
	select {
	case err := <-monitorClient.ErrorChannel():
		logrus.Errorf(nsmMonitorLogWithParamFormat, "NSM die, re-connecting", err)
		for _, c := range c.connections {
			c.State = connection.State_DOWN // Mark all as down.
		}
		return false
	case event := <-monitorClient.EventChannel():
		if event.EventType() == monitor.EventTypeInitialStateTransfer {
			logrus.Infof(nsmMonitorLogFormat, "Monitor started")
			c.initRecieved = true
		}

		for _, entity := range event.Entities() {
			switch event.EventType() {
			case monitor.EventTypeInitialStateTransfer, monitor.EventTypeUpdate:
				c.updateConnection(entity)
			case monitor.EventTypeDelete:
				logrus.Infof(nsmMonitorLogFormat, "Connection closed")
				if c.helper != nil {
					conn, ok := entity.(*connection.Connection)
					if ok {
						c.helper.Closed(conn)
					}
				}
			}
		}
		return true
	case <-c.stop:
		if c.helper != nil {
			c.helper.Stopped()
			logrus.Infof(nsmMonitorLogFormat, "Processing stop")
			break
		}
		return false
	}
	return true
}

func (c *nsmMonitorApp) updateConnection(entity monitor.Entity) {
	conn, ok := entity.(*connection.Connection)
	// update connections
	if ok {
		if existingConn, exists := c.connections[conn.Id]; exists {
			logrus.Infof(nsmMonitorLogWithParamFormat, "Connection updated", fmt.Sprintf("%v %v", existingConn, conn))
		} else {
			logrus.Infof(nsmMonitorLogWithParamFormat, "Initial connection accepted", conn)
		}
		c.connections[conn.Id] = conn
	}
}

func (c *nsmMonitorApp) waitRetry() {
	logrus.Errorf(nsmMonitorLogWithParamFormat, "Retry delay %v sec", nsmMonitorRetryDelay/time.Second)
	<-time.After(nsmMonitorRetryDelay)
}

func (c *nsmMonitorApp) performRecovery(nsmClient *client.NsmClient) bool {
	logrus.Infof(nsmMonitorLogFormat, "Performing recovery if needed...")

	needRetry := false
	for _, conn := range c.connections {
		if conn.State == connection.State_UP {
			continue
		}
		cClone := (conn.Clone()).(*connection.Connection)

		ipCtx := cClone.Context.IpContext
		if ipCtx != nil {
			if ipCtx.DstIpAddr != "" {
				ipCtx.DstIpRequired = true
			}
			if ipCtx.SrcIpAddr != "" {
				ipCtx.SrcIpRequired = true
			}
		}

		outgoingRequest := networkservice.NetworkServiceRequest{
			Connection: cClone,
			MechanismPreferences: []*connection.Mechanism{
				conn.Mechanism,
			},
		}
		if c.helper != nil {
			c.helper.Healing(cClone)
		}

		outgoingConnection, err := nsmClient.NsClient.Request(context.Background(), &outgoingRequest)

		if err != nil {
			logrus.Errorf(nsmMonitorLogWithParamFormat, "failed to restore connection. Will retry", err)
			// Let's drop connection id, since we failed one time.
			conn.Id = "-"
			needRetry = true
			continue
		} else {
			logrus.Errorf(nsmMonitorLogWithParamFormat, "connection restored", outgoingConnection)
			delete(c.connections, conn.Id)
			c.connections[outgoingConnection.Id] = outgoingConnection
		}
		if c.helper != nil {
			c.helper.ProcessHealing(outgoingConnection, err)
		}
	}
	return needRetry
}
