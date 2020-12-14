package spanlogger

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"

	opentelemetrynsm "github.com/networkservicemesh/sdk/pkg/tools/opentelemetry"
	opentelemetry "go.opentelemetry.io/otel/trace"
)

type Span interface {
	Log(level, format string, v ...interface{})
	LogObject(k, v interface{})
	WithField(k, v interface{}) Span
	Finish()

	ToString() string
}

 // Opentracing span
type otSpan struct {
	span    opentracing.Span
}

func (ot *otSpan) Log(level, format string, v ...interface{}) {
	ot.span.LogFields(
		opentracinglog.String("event", level),
		opentracinglog.String("message", fmt.Sprintf(format, v...)),
	)
}

func (ot *otSpan) LogObject(k, v interface{}) {
	ot.span.LogFields(opentracinglog.Object(k.(string), v))
}

func (ot *otSpan) WithField(k, v interface{}) Span {
	ot.span = ot.span.SetTag(k.(string), v)
	return ot
}

func (ot *otSpan) ToString() string {
	if spanStr := fmt.Sprintf("%v", ot.span); spanStr != "{}" {
		return spanStr
	}
	return ""
}

func (ot *otSpan) Finish() {
	ot.span.Finish()
}

func newOTSpan(ctx context.Context, operationName string, additionalFields map[string]interface{}) (c context.Context, s Span) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	for k, v := range additionalFields  {
		span = span.SetTag(k, v)
	}
	return ctx, &otSpan{span: span}
}

// Opentelemetry span
type otelSpan struct {
	span    opentelemetry.Span
}

func (otel *otelSpan) Log(level, format string, v ...interface{}) {
	otel.span.AddEvent(
		"",
		opentelemetry.WithAttributes([]label.KeyValue{
			label.String("event", level),
			label.String("message", fmt.Sprintf(format, v...)),
		}...),
	)
}

func (otel *otelSpan) LogObject(k, v interface{}) {
	otel.span.AddEvent(
		"",
		opentelemetry.WithAttributes([]label.KeyValue{
			label.String(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v)),
		}...),
	)
}

func (otel *otelSpan) WithField(k, v interface{}) Span {
	otel.span.SetAttributes(label.Any(k.(string), v))
	return otel
}

func (otel *otelSpan) ToString() string {
	if spanID := otel.span.SpanContext().SpanID; spanID.IsValid() {
		return spanID.String()
	}
	return ""
}

func (otel *otelSpan) Finish() {
	otel.span.End()
}

func newOTELSpan(ctx context.Context, operationName string, additionalFields map[string]interface{}) (c context.Context, s Span) {
	var add []label.KeyValue

	for k, v := range additionalFields  {
		add = append(add, label.Any(k, v))
	}

	ctx, span := otel.Tracer(opentelemetrynsm.InstrumentationName).Start(ctx, operationName)
	span.SetAttributes(add...)

	return ctx, &otelSpan{span: span}
}