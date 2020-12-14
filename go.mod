module github.com/networkservicemesh/sdk

go 1.15

require (
	cloud.google.com/go v0.46.3 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/RoaringBitmap/roaring v0.4.23
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/edwarnicke/exechelper v1.0.2
	github.com/edwarnicke/grpcfd v0.0.0-20200920223154-d5b6e1f19bd0
	github.com/edwarnicke/serialize v1.0.7
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/nats-io/nats-streaming-server v0.17.0
	github.com/nats-io/stan.go v0.6.0
	github.com/networkservicemesh/api v0.0.0-20210202152048-ec956057eb3a
	github.com/open-policy-agent/opa v0.16.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spiffe/go-spiffe/v2 v2.0.0-alpha.4.0.20200528145730-dc11d0c74e85
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.16.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/otlp v0.17.0
	go.opentelemetry.io/otel/metric v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	go.opentelemetry.io/otel/sdk/metric v0.17.0
	go.opentelemetry.io/otel/trace v0.17.0
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/goleak v1.1.10
	gonum.org/v1/gonum v0.6.2
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.16.0 => ../../../../../xor/opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc
