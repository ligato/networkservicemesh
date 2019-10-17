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
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/remote/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/remote/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/common"
	"github.com/networkservicemesh/networkservicemesh/pkg/security"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools/spanhelper"
	"github.com/sirupsen/logrus"
)

type securityService struct {
	provider security.Provider
}

func NewSecurityService(provider security.Provider) *securityService {
	return &securityService{
		provider: provider,
	}
}

func (s *securityService) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	span := spanhelper.GetSpanHelper(ctx)
	sc := security.NewContext()

	nctx := common.WithSecurityContext(ctx, sc)
	nsc := common.SecurityContext(nctx)
	logrus.Infof("remote.New SC = %p", nsc)

	conn, err := ProcessNext(nctx, request)
	if err != nil {
		span.LogError(err)
		return nil, err
	}

	if s.provider == nil {
		logrus.Warn("insecure: provider is not set, return Connection without signature")
		return conn, nil
	}

	logrus.Infof("remote.SecurityService: sc.GetResponseOboToken() = %v", sc.GetResponseOboToken())

	if err := security.SignConnection(conn, sc.GetResponseOboToken(), s.provider); err != nil {
		span.LogError(err)
		return nil, err
	}

	return conn, nil
}

func (s *securityService) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	return ProcessClose(ctx, connection)
}

func getOboToken(ctx context.Context) string {
	nseConn := common.EndpointConnection(ctx)
	if nseConn == nil {
		return ""
	}

	return nseConn.GetSignature()
}
