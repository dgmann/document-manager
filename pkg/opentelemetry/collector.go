package opentelemetry

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var ErrNoUrl = errors.New("no collector url")

// Provider is an opinionated client to quickly instrument an application via OpenTelemetry.
type Provider struct {
	TraceProvider *sdktrace.TracerProvider
	MeterProvider *sdkmetric.MeterProvider
	ServiceName   string
	conn          *grpc.ClientConn
}

func NewProvider(ctx context.Context, serviceName string, collectorUrl string) (*Provider, error) {
	if len(collectorUrl) == 0 {
		return nil, ErrNoUrl
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorUrl,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceProvider, err := initOtelTraceProvider(ctx, conn, commonOptions(serviceName))
	if err != nil {
		return nil, err
	}
	meterProvider, err := initOtelMeterProvider(ctx, conn, commonOptions(serviceName))
	if err != nil {
		return nil, err
	}

	return &Provider{
		TraceProvider: traceProvider,
		MeterProvider: meterProvider,
		ServiceName:   serviceName,
		conn:          conn,
	}, nil
}

// Shutdown cleans up all open resources.
//
// Can also be called on a nil Provider, in that case it is a NoOp.
func (p *Provider) Shutdown(ctx context.Context) error {
	// If nothing was initialized, we can just return
	if p == nil {
		return nil
	}

	return errors.Join(p.conn.Close(), p.TraceProvider.Shutdown(ctx), p.MeterProvider.Shutdown(ctx))
}

func commonOptions(serviceName string) []resource.Option {
	return []resource.Option{
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
		resource.WithSchemaURL(semconv.SchemaURL),
	}
}

func initOtelTraceProvider(ctx context.Context, conn *grpc.ClientConn, commonOptions []resource.Option) (*sdktrace.TracerProvider, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	res, err := resource.New(ctx, commonOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider, nil
}

func initOtelMeterProvider(ctx context.Context, conn *grpc.ClientConn, commonOptions []resource.Option) (*sdkmetric.MeterProvider, error) {
	exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	res, err := resource.New(ctx, commonOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}
