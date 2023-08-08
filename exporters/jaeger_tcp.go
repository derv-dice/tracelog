package exporters

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func ExporterJaegerTCP(url, serviceName, environment string, attrs ...attribute.KeyValue) (tp *tracesdk.TracerProvider, err error) {
	attrs = append(attrs, semconv.ServiceName(serviceName), attribute.String(envKey, environment))

	var exp *jaeger.Exporter
	if exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url))); err != nil {
		return
	}

	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attrs...)))

	return
}
