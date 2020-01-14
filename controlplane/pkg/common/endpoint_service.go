// Copyright (c) 2019 Cisco and/or its affiliates.
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

package common

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	mechanismCommon "github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/common"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/kernel"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/api/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/properties"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools/spanhelper"
)

type endpointService struct {
	nseManager     nsm.NetworkServiceEndpointManager
	props          *properties.Properties
	model          model.Model
	requestBuilder RequestBuilder
}

func (cce *endpointService) closeEndpoint(ctx context.Context, cc *model.ClientConnection) error {
	span := spanhelper.FromContext(ctx, "closeEndpoint")
	defer span.Finish()
	ctx = span.Context()
	logger := span.Logger()

	if cc.Endpoint == nil {
		logger.Infof("No need to close, since NSE is we know is dead at this point.")
		return nil
	}
	closeCtx, closeCancel := context.WithTimeout(ctx, cce.props.CloseTimeout)
	defer closeCancel()

	client, nseClientError := cce.nseManager.CreateNSEClient(closeCtx, cc.Endpoint)

	if client != nil {
		if ld := cc.Xcon.GetDestination(); ld != nil {
			return client.Close(ctx, ld)
		}
		err := client.Cleanup()
		span.LogError(err)
	} else {
		span.LogError(nseClientError)
	}
	return nseClientError
}

func (cce *endpointService) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	logger := Log(ctx)
	clientConnection := ModelConnection(ctx)
	endpoint := Endpoint(ctx)

	if clientConnection == nil {
		return nil, errors.Errorf("client connection need to be passed")
	}
	client, err := cce.nseManager.CreateNSEClient(ctx, endpoint)
	if err != nil {
		return nil, errors.Errorf("NSM: endpointService failed to create NSE Client. %v", err)
	}
	defer func() {
		if cleanupErr := client.Cleanup(); cleanupErr != nil {
			logger.Errorf("NSM: endpointService: error during Cleanup: %v", cleanupErr)
		}
	}()

	var conn *connection.Connection
	var connID string

	// Try Heal only if endpoint are same as for existing connection.
	if clientConnection.ConnectionState == model.ClientConnectionHealing && endpoint == clientConnection.Endpoint {
		conn = clientConnection.Xcon.GetDestination()
		connID = conn.Id
	} else {
		conn = request.Connection
		connID = ""
	}

	message := cce.requestBuilder.Build(ctx, connID, endpoint, conn)

	logger.Infof("NSM: endpointService: requesting NSE with request %v", message)

	span := spanhelper.FromContext(ctx, "nse.request")
	ctx = span.Context()
	defer span.Finish()
	span.LogObject("nse.request", message)

	nseConn, e := client.Request(ctx, message)
	span.LogObject("nse.response", nseConn)
	if e != nil {
		e = errors.Errorf("NSM: endpointService: error requesting networkservice from %+v with message %#v error: %s", endpoint, message, e)
		span.LogError(e)
		return nil, e
	}
	if err = cce.updateConnectionContext(ctx, request.GetConnection(), nseConn); err != nil {
		err = errors.Errorf("NSM: endpointService: failure Validating NSE Connection: %s", err)
		span.LogError(err)
		return nil, err
	}

	// Update connection parameters, add workspace if local nse
	cce.updateConnectionParameters(nseConn, endpoint)

	ctx = WithEndpointConnection(ctx, nseConn)

	return ProcessNext(ctx, request)
}

func (cce *endpointService) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	clientConnection := ModelConnection(ctx)
	if clientConnection != nil {
		if err := cce.closeEndpoint(ctx, clientConnection); err != nil {
			return &empty.Empty{}, err
		}
	}

	return ProcessClose(ctx, connection)
}

func (cce *endpointService) validateConnection(ctx context.Context, conn *connection.Connection) error {
	if err := conn.IsComplete(); err != nil {
		return err
	}

	return nil
}

func (cce *endpointService) updateConnectionContext(ctx context.Context, source, destination *connection.Connection) error {
	if err := cce.validateConnection(ctx, destination); err != nil {
		return err
	}

	if err := source.UpdateContext(destination.GetContext()); err != nil {
		return err
	}

	return nil
}

func (cce *endpointService) updateConnectionParameters(nseConn *connection.Connection, endpoint *registry.NSERegistration) {
	if cce.nseManager.IsLocalEndpoint(endpoint) {
		modelEp := cce.model.GetEndpoint(endpoint.GetNetworkServiceEndpoint().GetName())
		if modelEp != nil { // In case of tests this could be empty
			nseConn.GetMechanism().GetParameters()[mechanismCommon.Workspace] = modelEp.Workspace
			nseConn.GetMechanism().GetParameters()[kernel.WorkspaceNSEName] = modelEp.Endpoint.GetNetworkServiceEndpoint().GetName()
		}
		logrus.Infof("NSM: endpointService: update Local NSE connection parameters: %v", nseConn.Mechanism)
	}
}

// NewEndpointService -  creates a service to connect to endpoint
func NewEndpointService(nseManager nsm.NetworkServiceEndpointManager, properties *properties.Properties, mdl model.Model, builder RequestBuilder) networkservice.NetworkServiceServer {
	return &endpointService{
		nseManager:     nseManager,
		props:          properties,
		model:          mdl,
		requestBuilder: builder,
	}
}
