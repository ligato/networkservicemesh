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

package remote

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/srv6"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/vxlan"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/common"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/serviceregistry"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools/spanhelper"
	"github.com/networkservicemesh/networkservicemesh/utils"
)

const (
	// ForwarderRetryCount - A number of times to call Forwarder Request, TODO: Remove after DP will be stable.
	ForwarderRetryCount = 10
	// ForwarderRetryDelay - a delay between operations.
	ForwarderRetryDelay = 500 * time.Millisecond
	// ForwarderTimeout - A forwarder timeout
	ForwarderTimeout = 15 * time.Second
	// ErrorCloseTimeout - timeout to close all stuff in case of error
	ErrorCloseTimeout = 15 * time.Second
	// PreferredRemoteMechanism - mechanism name will be chosen by default if supported
	PreferredRemoteMechanism = utils.EnvVar("PREFERRED_REMOTE_MECHANISM")
)

// forwarderService -
type forwarderService struct {
	serviceRegistry serviceregistry.ServiceRegistry
	model           model.Model
}

func (cce *forwarderService) selectForwarder(request *networkservice.NetworkServiceRequest) (*model.Forwarder, error) {
	dp, err := cce.model.SelectForwarder(func(dp *model.Forwarder) bool {
		for _, m := range request.GetRequestMechanismPreferences() {
			if cce.findMechanism(dp.RemoteMechanisms, m.GetType()) != nil {
				return true
			}
		}
		return false
	})
	return dp, err
}
func (cce *forwarderService) findMechanism(mechanismPreferences []*networkservice.Mechanism, mechanismType string) *networkservice.Mechanism {
	for _, m := range mechanismPreferences {
		if m.GetType() == mechanismType {
			return m
		}
	}
	return nil
}

func (cce *forwarderService) selectRemoteMechanism(request *networkservice.NetworkServiceRequest, dp *model.Forwarder) (*networkservice.Mechanism, error) {
	var mechanism *networkservice.Mechanism
	var dpMechanism *networkservice.Mechanism

	if preferredMechanismName := PreferredRemoteMechanism.StringValue(); len(preferredMechanismName) > 0 {
		for _, m := range request.GetRequestMechanismPreferences() {
			if m.GetType() == preferredMechanismName {
				if dpm := cce.findMechanism(dp.RemoteMechanisms, m.GetType()); dpm != nil {
					mechanism = m
					dpMechanism = dpm
					break
				}
			}
		}
	}

	if mechanism == nil {
		for _, m := range request.GetRequestMechanismPreferences() {
			dpm := cce.findMechanism(dp.RemoteMechanisms, m.GetType())
			if dpm != nil {
				mechanism = m
				dpMechanism = dpm
				break
			}
		}
	}

	if mechanism == nil || dpMechanism == nil {
		return nil, errors.Errorf("failed to select mechanism, no matched mechanisms found")
	}

	switch mechanism.GetType() {
	case vxlan.MECHANISM:
		cce.configureVXLANParameters(mechanism.GetParameters(), dpMechanism.GetParameters())

	case srv6.MECHANISM:
		connectionID := request.GetConnection().GetId()
		parameters := mechanism.GetParameters()
		dpParameters := dpMechanism.GetParameters()

		cce.configureSRv6Parameters(connectionID, parameters, dpParameters)
	}

	logrus.Infof("NSM:(5.1) Remote mechanism selected %v", mechanism)
	return mechanism, nil
}

func (cce *forwarderService) configureVXLANParameters(parameters, dpParameters map[string]string) {
	parameters[vxlan.DstIP] = dpParameters[vxlan.SrcIP]

	extSrcIP := parameters[vxlan.SrcIP]
	extDstIP := dpParameters[vxlan.SrcIP]
	srcIP := parameters[vxlan.SrcIP]
	dstIP := dpParameters[vxlan.SrcIP]

	if ip, ok := parameters[vxlan.SrcOriginalIP]; ok {
		srcIP = ip
	}

	if ip, ok := parameters[vxlan.DstExternalIP]; ok {
		extDstIP = ip
	}

	var vni uint32
	if extDstIP != extSrcIP {
		vni = cce.serviceRegistry.VniAllocator().Vni(extDstIP, extSrcIP)
	} else {
		vni = cce.serviceRegistry.VniAllocator().Vni(dstIP, srcIP)
	}

	parameters[vxlan.VNI] = strconv.FormatUint(uint64(vni), 10)
}

func (cce *forwarderService) configureSRv6Parameters(connectionID string, parameters, dpParameters map[string]string) {
	parameters[srv6.DstHardwareAddress] = dpParameters[srv6.SrcHardwareAddress]
	parameters[srv6.DstHostIP] = dpParameters[srv6.SrcHostIP]
	parameters[srv6.DstHostLocalSID] = dpParameters[srv6.SrcHostLocalSID]
	parameters[srv6.DstBSID] = cce.serviceRegistry.SIDAllocator().SID(connectionID)
	parameters[srv6.DstLocalSID] = cce.serviceRegistry.SIDAllocator().SID(connectionID)
}

func (cce *forwarderService) updateMechanism(request *networkservice.NetworkServiceRequest, dp *model.Forwarder) error {
	conn := request.GetConnection()
	// 5.x
	if m, err := cce.selectRemoteMechanism(request, dp); err == nil {
		conn.Mechanism = m.Clone()
	} else {
		return err
	}

	if conn.GetMechanism() == nil {
		return errors.Errorf("required mechanism are not found... %v ", request.GetRequestMechanismPreferences())
	}

	if conn.GetMechanism().GetParameters() == nil {
		conn.Mechanism.Parameters = map[string]string{}
	}

	return nil
}

func (cce *forwarderService) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	logger := common.Log(ctx)
	span := spanhelper.GetSpanHelper(ctx)

	clientConnection := common.ModelConnection(ctx)
	// 3. get forwarder
	if err := cce.serviceRegistry.WaitForForwarderAvailable(ctx, cce.model, ForwarderTimeout); err != nil {
		logger.Errorf("Error waiting for forwarder: %v", err)
		return nil, err
	}

	// TODO: We could iterate forwarders to match required one, if failed with first one.
	dp, err := cce.selectForwarder(request)
	if err != nil {
		return nil, err
	}

	// 5. Select a local forwarder and put it into conn object
	err = cce.updateMechanism(request, dp)
	if err != nil {
		// 5.1 Close forwarder connection, if had existing one and NSE is closed.
		cce.doFailureClose(ctx)
		return nil, errors.Errorf("NSM:(5.1) %v", err)
	}

	span.LogObject("dataplane", dp)

	ctx = common.WithForwarder(ctx, dp)
	conn, connErr := common.ProcessNext(ctx, request)
	if connErr != nil {
		cce.doFailureClose(ctx)
		return conn, connErr
	}
	// We need to program forwarder.
	return cce.programForwarder(ctx, conn, dp, clientConnection)
}

func (cce *forwarderService) doFailureClose(ctx context.Context) {
	clientConnection := common.ModelConnection(ctx)

	newCtx, cancel := context.WithTimeout(context.Background(), ErrorCloseTimeout)
	defer cancel()

	span := spanhelper.CopySpan(newCtx, spanhelper.GetSpanHelper(ctx), "doForwarderClose")
	defer span.Finish()

	newCtx = span.Context()

	newCtx = common.WithLog(newCtx, span.Logger())
	newCtx = common.WithModelConnection(newCtx, clientConnection)

	closeErr := cce.performClose(newCtx, clientConnection, span.Logger())
	span.LogError(closeErr)
}

func (cce *forwarderService) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	cc := common.ModelConnection(ctx)
	logger := common.Log(ctx)
	empt, err := common.ProcessClose(ctx, conn)
	if closeErr := cce.performClose(ctx, cc, logger); closeErr != nil {
		logger.Errorf("Failed to close: %v", closeErr)
	}
	return empt, err
}

func (cce *forwarderService) performClose(ctx context.Context, cc *model.ClientConnection, logger logrus.FieldLogger) error {
	// Close endpoints, etc
	if cc.ForwarderState != model.ForwarderStateNone {
		logger.Info("NSM.Forwarder: Closing cross connection on forwarder...")
		dp := cce.model.GetForwarder(cc.ForwarderRegisteredName)
		forwarderClient, conn, err := cce.serviceRegistry.ForwarderConnection(ctx, dp)
		if err != nil {
			logger.Error(err)
			return err
		}
		if conn != nil {
			defer func() { _ = conn.Close() }()
		}
		if _, err := forwarderClient.Close(ctx, cc.Xcon); err != nil {
			logger.Error(err)
			return err
		}
		logger.Info("NSM.Forwarder: Cross connection successfully closed on forwarder")
		cc.ForwarderState = model.ForwarderStateNone
	}
	return nil
}

func (cce *forwarderService) programForwarder(ctx context.Context, conn *networkservice.Connection, dp *model.Forwarder, clientConnection *model.ClientConnection) (*networkservice.Connection, error) {
	span := spanhelper.FromContext(ctx, "programForwarder")
	defer span.Finish()
	// We need to program forwarder.
	forwarderClient, forwarderConn, err := cce.serviceRegistry.ForwarderConnection(ctx, dp)
	if err != nil {
		span.Logger().Errorf("Error creating forwarder connection %v. Performing close", err)
		cce.doFailureClose(span.Context())
		return conn, err
	}
	if forwarderConn != nil { // Required for testing
		defer func() {
			if closeErr := forwarderConn.Close(); closeErr != nil {
				span.Logger().Errorf("NSM:(6.1) Error during close Forwarder connection: %v", closeErr)
			}
		}()
	}

	var newXcon *crossconnect.CrossConnect
	// 9. We need to program forwarder with our values.
	// 9.1 Sending updated request to forwarder.
	for dpRetry := 0; dpRetry < ForwarderRetryCount; dpRetry++ {
		if ctx.Err() != nil {
			cce.doFailureClose(ctx)
			return nil, ctx.Err()
		}

		attemptSpan := spanhelper.FromContext(span.Context(), fmt.Sprintf("ProgramAttempt-%v", dpRetry))
		defer attemptSpan.Finish()
		attemptSpan.LogObject("attempt", dpRetry)

		span.Logger().Infof("NSM:(9.1) Sending request to forwarder")
		attemptSpan.LogObject("request", clientConnection.Xcon)

		dpCtx, cancel := context.WithTimeout(attemptSpan.Context(), ForwarderTimeout)
		newXcon, err = forwarderClient.Request(dpCtx, clientConnection.Xcon)
		cancel()
		if err != nil {
			attemptSpan.Logger().Errorf("NSM:(9.1.1) Forwarder request failed: %v retry: %v", err, dpRetry)

			// Let's try again with a short delay
			if dpRetry < ForwarderRetryCount-1 {
				<-time.After(ForwarderRetryDelay)
				continue
			}
			attemptSpan.Logger().Errorf("NSM:(9.1.2) Forwarder request  all retry attempts failed: %v", clientConnection.Xcon)
			// 9.3 If datplane configuration are failed, we need to close remore NSE actually.
			cce.doFailureClose(attemptSpan.Context())
			attemptSpan.Finish()
			return conn, err
		}

		// In case of context deadline, we need to close NSE and forwarder.
		ctxErr := attemptSpan.Context().Err()
		if ctxErr != nil {
			attemptSpan.LogError(ctxErr)
			cce.doFailureClose(attemptSpan.Context())
			attemptSpan.Finish()
			return nil, ctxErr
		}

		clientConnection.Xcon = newXcon

		attemptSpan.Logger().Infof("NSM:(9.2) Forwarder configuration successful")
		attemptSpan.LogObject("crossConnection", clientConnection.Xcon)
		break
	}

	// Update connection context if it updated from forwarder.
	return cce.updateClientConnection(ctx, conn, clientConnection, dp)
}

func (cce *forwarderService) updateClientConnection(ctx context.Context, conn *networkservice.Connection, clientConnection *model.ClientConnection, dp *model.Forwarder) (*networkservice.Connection, error) {
	// Update connection context if it updated from forwarder.
	err := conn.UpdateContext(clientConnection.GetConnectionSource().GetContext())
	if err != nil {
		cce.doFailureClose(ctx)
		return nil, err
	}

	clientConnection.ForwarderRegisteredName = dp.RegisteredName
	clientConnection.ForwarderState = model.ForwarderStateReady
	if clientConnection.GetConnectionDestination() != nil && clientConnection.GetConnectionDestination().GetContext() != nil {
		conn.Context.EthernetContext = clientConnection.GetConnectionDestination().GetContext().EthernetContext
	}
	return conn, nil
}

// NewForwarderService -  creates a service to program forwarder.
func NewForwarderService(model model.Model, serviceRegistry serviceregistry.ServiceRegistry) networkservice.NetworkServiceServer {
	return &forwarderService{
		model:           model,
		serviceRegistry: serviceRegistry,
	}
}
