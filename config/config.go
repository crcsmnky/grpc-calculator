package config

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	serviceName = "grpc-calculator-otel"
	serviceVersion = "1.0.0"
)

// func Init() (*sdktrace.TracerProvider, error) {
// 	ctx := context.Background()

// 	// exporter, err := otlptracegrpc.New(
// 	// 	ctx,
// 	// 	otlptracegrpc.WithInsecure(),
// 	// )

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// var traceExporter sdktrace.SpanExporter

// 	// grpcEndpoint := "localhost:4317"
// 	grpcEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	
// 	log.Printf("setting up gRPC endpoint: %s", grpcEndpoint)

// 	traceExporter, err := otlptracegrpc.New(ctx,
// 		otlptracegrpc.WithInsecure(),
// 		otlptracegrpc.WithEndpoint(grpcEndpoint),
// 		// otlptracegrpc.WithDialOption(grpc.with),
// 	)

// 	if err != nil {
// 		log.Fatalf("%s: %v", "failed to create trace exporter", err)
// 	}

// 	resources := resource.NewWithAttributes(
// 		semconv.SchemaURL,
// 		semconv.ServiceNameKey.String(serviceName),
// 		semconv.ServiceVersionKey.String(serviceVersion),
// 		semconv.ServiceInstanceIDKey.String(uuid.New().String()),
// 	)

// 	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithBatcher(traceExporter),
// 		sdktrace.WithResource(resources),
// 		sdktrace.WithSpanProcessor(bsp),
// 	)

// 	otel.SetTracerProvider(tracerProvider)
// 	otel.SetTextMapPropagator(propagation.TraceContext{})
// 	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

// 	return tracerProvider, err
// }

func NewTraceExporter(ctx context.Context, otlpEndpoint string) (*otlptrace.Exporter, error) {
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
	)

	return traceExporter, err
}

func NewTraceProvider(exporter sdktrace.SpanExporter) *sdktrace.TracerProvider {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			semconv.ServiceInstanceIDKey.String(uuid.New().String()),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
}