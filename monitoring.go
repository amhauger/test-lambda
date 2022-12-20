package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

func CreateTracerAndBatcher() (func(), error) {
	exporter, ctx, err := createExporter()
	if err != nil {
		return nil, err
	}

	tracer := createBatcher(exporter)

	otel.SetTracerProvider(tracer)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return func() { _ = tracer.Shutdown(ctx) }, nil
}

func createExporter() (*otlptrace.Exporter, context.Context, error) {
	ctx := context.Background()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, ctx, err
	}

	return exporter, ctx, nil
}

func createBatcher(exporter *otlptrace.Exporter) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes("", attribute.KeyValue{
			Key:   "service.name",
			Value: attribute.StringValue("test-layers"),
		})),
	)
}
