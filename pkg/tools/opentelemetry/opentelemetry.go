// Copyright (c) 2020 Doc.ai and/or its affiliates.
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

// Package opentelemetry provides a set of utilities for assisting with telemetry data
package opentelemetry

import (
	"context"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	// DefaultAgentAddr - opentelemetry agent default address
	DefaultAgentAddr = "0.0.0.0:55680"

	// InstrumentationName - denotes the library that provides the instrumentation
	InstrumentationName = "NSM"
)


// Init - creates opentelemetry trace and metrics providers
func Init(ctx context.Context, otelAgentAddr, service string) (func(), error) {
	if !log.IsOpentelemetryEnabled() {
		return nil, nil
	}

	metricsDriver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(otelAgentAddr),
	)
	tracesDriver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(otelAgentAddr),
	)
	splitCfg := otlp.SplitConfig{
		ForMetrics: metricsDriver,
		ForTraces:  tracesDriver,
	}
	driver := otlp.NewSplitDriver(splitCfg)
	exp, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(service),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)
	otel.SetTracerProvider(tracerProvider)

	/* Create meter provider*/
	cont := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exp,
		),
		controller.WithPusher(exp),
		controller.WithCollectPeriod(time.Second),
		controller.WithPushTimeout(0),
	)
	global.SetMeterProvider(cont.MeterProvider())
	if err := cont.Start(ctx); err != nil {
		return nil, err
	}

	return func() {
		// pushes any last exports to the receiver
		handleErr(ctx, cont.Stop(ctx), "failed to shutdown pusher")
		handleErr(ctx, tracerProvider.Shutdown(ctx), "failed to shutdown provider")
		handleErr(ctx, exp.Shutdown(ctx), "failed to stop exporter")
	}, nil
}

func handleErr(ctx context.Context, err error, message string) {
	if err != nil {
		log.FromContext(ctx).Errorf("%s: %v", message, err)
	}
}
