// Copyright (c) 2021 Doc.ai and/or its affiliates.
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

// Package metric provides a chain element that sends metrics to collector
package metrics

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"
	"github.com/networkservicemesh/sdk/pkg/tools/opentelemetry/meterhelper"
)

type metricClient struct {
	helpers map[string]meterhelper.MeterHelper
}

// NewClient returns a new metric client chain element
func NewClient() networkservice.NetworkServiceClient {
	return &metricClient{
		helpers: make(map[string]meterhelper.MeterHelper),
	}
}

func (t *metricClient) Request(ctx context.Context, request *networkservice.NetworkServiceRequest, opts ...grpc.CallOption) (*networkservice.Connection, error) {
	conn, err := next.Client(ctx).Request(ctx, request, opts...)
	if err != nil {
		return nil, err
	}

	if conn.Path != nil {
		for _, pathSegment := range conn.GetPath().GetPathSegments() {
			if pathSegment.Metrics == nil {
				continue
			}

			_, ok := t.helpers[pathSegment.Id]
			if !ok {
				t.helpers[pathSegment.Id] = meterhelper.NewMeterHelper(pathSegment.Name, conn.GetPath().GetPathSegments()[0].Id)
			}

			t.helpers[pathSegment.Id].WriteMetrics(ctx, pathSegment.Metrics)
		}
	}

	return conn, nil
}

func (t *metricClient) Close(ctx context.Context, conn *networkservice.Connection, opts ...grpc.CallOption) (*empty.Empty, error) {
	return next.Client(ctx).Close(ctx, conn, opts...)
}
