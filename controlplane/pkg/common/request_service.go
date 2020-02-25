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

package common

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
)

// requestValidator -
type requestValidator struct {
}

func (cce *requestValidator) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	err := request.IsValid()

	if err != nil {
		Log(ctx).Error(err)
		return nil, err
	}
	return ProcessNext(ctx, request)
}

func (cce *requestValidator) Close(ctx context.Context, connection *networkservice.Connection) (*empty.Empty, error) {
	return ProcessClose(ctx, connection)
}

// NewRequestValidator -  creates a service to verify request
func NewRequestValidator() networkservice.NetworkServiceServer {
	return &requestValidator{}
}
