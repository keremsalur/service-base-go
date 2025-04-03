package otel

import (
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/zipkin"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer trace.Tracer

func InitTracer(zipkinURL string) (*sdktrace.TracerProvider, error) {
	exporter, err := zipkin.New(
		zipkinURL,
	)
	if err != nil {
		return nil, err
	}

	// Trace provider'ı ayarla
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-app"),
		)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Global tracer provider'ı set et
	otel.SetTracerProvider(tp)
	tracer = otel.Tracer("my-app")
	return tp, nil
}

func GetTracer() trace.Tracer {
	return tracer
}
