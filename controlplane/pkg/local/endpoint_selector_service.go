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
package local

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/networkservice"
	unified_nsm "github.com/networkservicemesh/networkservicemesh/controlplane/api/nsm"
	unified_connection "github.com/networkservicemesh/networkservicemesh/controlplane/api/nsm/connection"
	plugin_api "github.com/networkservicemesh/networkservicemesh/controlplane/api/plugins"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/common"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/plugins"
)

// ConnectionService makes basic Mechanism selection for the incoming connection
type endpointSelectorService struct {
	nseManager     unified_nsm.NetworkServiceEndpointManager
	pluginRegistry plugins.PluginRegistry
}

func (cce *endpointSelectorService) updateConnection(ctx context.Context, conn *connection.Connection) (*connection.Connection, error) {
	if conn.GetContext() == nil {
		c := &connectioncontext.ConnectionContext{}
		conn.SetContext(c)
	}

	wrapper := plugin_api.NewConnectionWrapper(conn)
	wrapper, err := cce.pluginRegistry.GetConnectionPluginManager().UpdateConnection(ctx, wrapper)
	if err != nil {
		return conn, err
	}

	return wrapper.GetConnection().(*connection.Connection), nil
}

func (cce *endpointSelectorService) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	logger := common.Log(ctx)
	clientConnection := common.ModelConnection(ctx)
	dp := common.Dataplane(ctx)

	if clientConnection == nil {
		return nil, fmt.Errorf("client connection need to be passed")
	}

	// 4. Check if Heal/Update if we need to ask remote NSM or it is a just local mechanism change requested.
	// true if we detect we need to request NSE to upgrade/update connection.
	// 4.1 New Network service is requested, we need to close current connection and do re-request of NSE.
	requestNSEOnUpdate := cce.checkNSEUpdateIsRequired(ctx, clientConnection, request, logger, dp)

	// 7. do a Request() on NSE and select it.
	if clientConnection.ConnectionState == model.ClientConnectionHealing && !requestNSEOnUpdate {
		return cce.checkUpdateConnectionContext(ctx, request, clientConnection)
	}

	// 7.1 try find NSE and do a Request to it.
	var lastError error
	ignoreEndpoints := common.IgnoredEndpoints(ctx)
	for {
		if ctx.Err() != nil {
			logger.Errorf("NSM:(7.1.0) Context timeout, during find/call NSE... %v", ctx.Err())
			return nil, ctx.Err()
		}
		// 7.1.1 Clone Connection to support iteration via EndPoints
		newRequest, endpoint, err := cce.prepareRequest(ctx, request, clientConnection, ignoreEndpoints)
		if err != nil {
			if lastError != nil {
				return nil, fmt.Errorf("NSM:(7.1.5) %v. Last NSE Error: %v", err, lastError)
			}
			return nil, err
		}
		ctx = common.WithEndpoint(ctx, endpoint)

		// Perform passing execution to next chain element.
		conn, err := ProcessNext(ctx, newRequest)

		// 7.1.8 in case of error we put NSE into ignored list to check another one.
		if err != nil {
			logger.Errorf("NSM:(7.1.8) NSE respond with error: %v ", err)
			lastError = err
			ignoreEndpoints[endpoint.GetNetworkServiceEndpoint().GetName()] = endpoint
			continue
		}

		// We could put endpoint to clientConnection.
		clientConnection.Endpoint = endpoint

		if !cce.nseManager.IsLocalEndpoint(endpoint) {
			clientConnection.RemoteNsm = endpoint.GetNetworkServiceManager()
		}

		// 7.1.9 We are fine with NSE connection and could continue.
		return conn, nil
	}
}

func (cce *endpointSelectorService) selectEndpoint(ctx context.Context, clientConnection *model.ClientConnection, ignoreEndpoints map[string]*registry.NSERegistration, nseConn unified_connection.Connection) (*registry.NSERegistration, error) {
	var endpoint *registry.NSERegistration
	if clientConnection.ConnectionState == model.ClientConnectionHealing {
		// 7.1.2 Check previous endpoint, and it we will be able to contact it, it should be fine.
		endpointName := clientConnection.Endpoint.GetNetworkServiceEndpoint().GetName()
		if clientConnection.Endpoint != nil && ignoreEndpoints[endpointName] == nil {
			endpoint = clientConnection.Endpoint
		} else {
			// Ignored, we need to update DSTid.
			clientConnection.GetConnectionDestination().SetID("-")
		}
		//TODO: Add check if endpoint are in registry or not.
	}
	// 7.1.3 Check if endpoint is not ignored yet
	if endpoint == nil {
		// 7.1.4 Choose a new endpoint
		return cce.nseManager.GetEndpoint(ctx, nseConn, ignoreEndpoints)
	}
	return endpoint, nil
}

func (cce *endpointSelectorService) checkNSEUpdateIsRequired(ctx context.Context, clientConnection *model.ClientConnection, request *networkservice.NetworkServiceRequest, logger logrus.FieldLogger, dp *model.Dataplane) bool {
	requestNSEOnUpdate := false
	if clientConnection.ConnectionState == model.ClientConnectionHealing {
		if request.Connection.GetNetworkService() != clientConnection.GetNetworkService() {
			requestNSEOnUpdate = true

			// Just close, since client connection already passed with context.
			_, err := ProcessClose(ctx, request.GetConnection())
			// Network service is closing, we need to close remote NSM and re-program local one.
			if err != nil {
				logger.Errorf("NSM:(4.1) Error during close of NSE during Request.Upgrade %v Existing connection: %v error %v", request, clientConnection, err)
			}

		} else {
			// 4.2 Check if NSE is still required, if some more context requests are different.
			requestNSEOnUpdate = cce.checkNeedNSERequest(logger, request.Connection, clientConnection, dp)
		}
	}
	return requestNSEOnUpdate
}

func (cce *endpointSelectorService) validateConnection(ctx context.Context, conn unified_connection.Connection) error {
	if err := conn.IsComplete(); err != nil {
		return err
	}

	wrapper := plugin_api.NewConnectionWrapper(conn)
	result, err := cce.pluginRegistry.GetConnectionPluginManager().ValidateConnection(ctx, wrapper)
	if err != nil {
		return err
	}

	if result.GetStatus() != plugin_api.ConnectionValidationStatus_SUCCESS {
		return fmt.Errorf(result.GetErrorMessage())
	}

	return nil
}

func (cce *endpointSelectorService) updateConnectionContext(ctx context.Context, source *connection.Connection, destination unified_connection.Connection) error {
	if err := cce.validateConnection(ctx, destination); err != nil {
		return err
	}

	if err := source.UpdateContext(destination.GetContext()); err != nil {
		return err
	}

	return nil
}

/**
check if we need to do a NSE/Remote NSM request in case of our connection Upgrade/Healing procedure.
*/
func (cce *endpointSelectorService) checkNeedNSERequest(logger logrus.FieldLogger, nsmConn *connection.Connection, existingCC *model.ClientConnection, dp *model.Dataplane) bool {
	// 4.2.x
	// 4.2.1 Check if context is changed, if changed we need to
	if !proto.Equal(nsmConn.GetContext(), existingCC.GetConnectionSource().GetContext()) {
		return true
	}
	// We need to check, dp has mechanism changes in our Remote connection selected mechanism.

	if dst := existingCC.GetConnectionDestination(); dst.IsRemote() {
		dstM := dst.GetConnectionMechanism()

		// 4.2.2 Let's check if remote destination is matchs our dataplane destination.
		if dpM := cce.findMechanism(dp.RemoteMechanisms, dstM.GetMechanismType()); dpM != nil {
			// 4.2.3 We need to check if source mechanism type and source parameters are different
			for k, v := range dpM.GetParameters() {
				rmV := dstM.GetParameters()[k]
				if v != rmV {
					logger.Infof("NSM:(4.2.3) Remote mechanism parameter %s was different with previous one : %v  %v", k, rmV, v)
					return true
				}
			}
			if !dpM.Equals(dstM) {
				logger.Infof("NSM:(4.2.4)  Remote mechanism was different with previous selected one : %v  %v", dstM, dpM)
				return true
			}
		} else {
			logger.Infof("NSM:(4.2.5) Remote mechanism previously selected was not found: %v  in dataplane %v", dstM, dp.RemoteMechanisms)
			return true
		}
	}

	return false
}

func (cce *endpointSelectorService) findMechanism(mechanismPreferences []unified_connection.Mechanism, mechanismType unified_connection.MechanismType) unified_connection.Mechanism {
	for _, m := range mechanismPreferences {
		if m.GetMechanismType() == mechanismType {
			return m
		}
	}
	return nil
}

func (cce *endpointSelectorService) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	return ProcessClose(ctx, connection)
}

func (cce *endpointSelectorService) checkUpdateConnectionContext(ctx context.Context, request *networkservice.NetworkServiceRequest, clientConnection *model.ClientConnection) (*connection.Connection, error) {
	logger := common.Log(ctx)
	// We do not need to do request to endpoint and just need to update all stuff.
	// 7.2 We do not need to access NSE, since all parameters are same.
	clientConnection.GetConnectionSource().SetConnectionMechanism(request.Connection.GetConnectionMechanism())
	clientConnection.GetConnectionSource().SetConnectionState(unified_connection.StateUp)

	// 7.3 Destination context probably has been changed, so we need to update source context.
	if err := cce.updateConnectionContext(ctx, request.GetConnection(), clientConnection.GetConnectionDestination()); err != nil {
		err = fmt.Errorf("NSM:(7.3) Failed to update source connection context: %v", err)

		// Just close since client connection is already passed with context
		if _, closeErr := ProcessClose(ctx, request.GetConnection()); closeErr != nil {
			logger.Errorf("Close error: %v", closeErr)
		}
		return nil, err
	}

	if !cce.nseManager.IsLocalEndpoint(clientConnection.Endpoint) {
		clientConnection.RemoteNsm = clientConnection.Endpoint.GetNetworkServiceManager()
	}
	return request.Connection, nil
}

func (cce *endpointSelectorService) prepareRequest(ctx context.Context, request *networkservice.NetworkServiceRequest, clientConnection *model.ClientConnection, ignoreEndpoints map[string]*registry.NSERegistration) (*networkservice.NetworkServiceRequest, *registry.NSERegistration, error) {
	newRequest := request.Clone().(*networkservice.NetworkServiceRequest)
	nseConn := newRequest.Connection

	logger := common.Log(ctx)
	endpoint, err := cce.selectEndpoint(ctx, clientConnection, ignoreEndpoints, nseConn)
	if err != nil {
		return nil, nil, err
	}

	logger.Infof("selected endpoint %v", endpoint)
	// 7.1.6 Update Request with exclude_prefixes, etc
	nseConn, err = cce.updateConnection(ctx, nseConn)
	if err != nil {
		return nil, nil, fmt.Errorf("NSM:(7.1.6) Failed to update connection: %v", err)
	}

	newRequest.Connection = nseConn
	return newRequest, endpoint, nil
}

// NewEndpointSelectorService - creates a service to select endpoint
func NewEndpointSelectorService(nseManager unified_nsm.NetworkServiceEndpointManager, pluginRegistry plugins.PluginRegistry) networkservice.NetworkServiceServer {
	return &endpointSelectorService{
		nseManager:     nseManager,
		pluginRegistry: pluginRegistry,
	}
}
