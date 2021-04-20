// Copyright (c) 2020-2021 Doc.ai and/or its affiliates.
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

// Package memory provides registry chain based on memory chain elements
package memory

import (
	"context"
	"net/url"
	"time"

	"github.com/networkservicemesh/api/pkg/api/registry"
	"google.golang.org/grpc"

	registryserver "github.com/networkservicemesh/sdk/pkg/registry"
	"github.com/networkservicemesh/sdk/pkg/registry/common/checkid"
	"github.com/networkservicemesh/sdk/pkg/registry/common/connect"
	"github.com/networkservicemesh/sdk/pkg/registry/common/expire"
	"github.com/networkservicemesh/sdk/pkg/registry/common/memory"
	"github.com/networkservicemesh/sdk/pkg/registry/common/proxy"
	"github.com/networkservicemesh/sdk/pkg/registry/common/serialize"
	"github.com/networkservicemesh/sdk/pkg/registry/core/adapters"
	"github.com/networkservicemesh/sdk/pkg/registry/core/chain"
)

type serverOptions struct {
	expiryDuration time.Duration
	dialOptions    []grpc.DialOption
}

// Option modifies option values
type Option func(o *serverOptions)

// WithExpiryDuration sets expire duration time for the server
func WithExpiryDuration(duration time.Duration) Option {
	return func(o *serverOptions) {
		o.expiryDuration = duration
	}
}

// WithDialOptions sets gRPC Dial Options for the server
func WithDialOptions(options ...grpc.DialOption) Option {
	return func(o *serverOptions) {
		o.dialOptions = options
	}
}

// NewServer creates new registry server based on memory storage
func NewServer(ctx context.Context, proxyRegistryURL *url.URL, options ...Option) registryserver.Registry {
	opts := &serverOptions{
		expiryDuration: time.Minute,
	}
	for _, opt := range options {
		opt(opts)
	}

	nseChain := chain.NewNetworkServiceEndpointRegistryServer(
		serialize.NewNetworkServiceEndpointRegistryServer(),
		expire.NewNetworkServiceEndpointRegistryServer(ctx, opts.expiryDuration),
		checkid.NewNetworkServiceEndpointRegistryServer(),
		memory.NewNetworkServiceEndpointRegistryServer(),
		proxy.NewNetworkServiceEndpointRegistryServer(proxyRegistryURL),
		connect.NewNetworkServiceEndpointRegistryServer(ctx, func(ctx context.Context, cc grpc.ClientConnInterface) registry.NetworkServiceEndpointRegistryClient {
			return chain.NewNetworkServiceEndpointRegistryClient(
				registry.NewNetworkServiceEndpointRegistryClient(cc),
			)
		}, connect.WithClientDialOptions(opts.dialOptions...)),
	)
	nsChain := chain.NewNetworkServiceRegistryServer(
		serialize.NewNetworkServiceRegistryServer(),
		expire.NewNetworkServiceServer(ctx, adapters.NetworkServiceEndpointServerToClient(nseChain)),
		memory.NewNetworkServiceRegistryServer(),
		proxy.NewNetworkServiceRegistryServer(proxyRegistryURL),
		connect.NewNetworkServiceRegistryServer(ctx, func(ctx context.Context, cc grpc.ClientConnInterface) registry.NetworkServiceRegistryClient {
			return chain.NewNetworkServiceRegistryClient(
				registry.NewNetworkServiceRegistryClient(cc),
			)
		}, connect.WithClientDialOptions(opts.dialOptions...)),
	)

	return registryserver.NewServer(nsChain, nseChain)
}
