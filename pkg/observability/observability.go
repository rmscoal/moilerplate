package observability

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	_defaultTraceEndpoint   string = "localhost:4137"
	_defaultMetricsEndpoint string = "localhost:9496"

	_defaultServiceName       string = "moilerplate-app"
	_defaultServiceVersion    string = "v0.1.0"
	_defaultServiceInstanceID string = "moilerplate-app-1"

	once sync.Once
	obsv *observability
)

type observability struct {
	// Endpoint to the trace exporter
	traceExporterEndpoint string
	// Endpoint to the metrics exporter
	metricsExporterEndpoint string

	// Services attributes

	serviceName       string
	serviceVersion    string
	serviceInstanceID string

	// Shutdownfuncs for all the otel providers
	shutdownFuncs []func(context.Context) error
}

func Init(ctx context.Context, opts ...Option) {
	if obsv == nil {
		once.Do(func() {
			obsv = &observability{
				traceExporterEndpoint:   _defaultTraceEndpoint,
				metricsExporterEndpoint: _defaultMetricsEndpoint,
				serviceName:             _defaultServiceName,
				serviceVersion:          _defaultServiceVersion,
				serviceInstanceID:       _defaultServiceInstanceID,
			}
		})
	}

	for _, opt := range opts {
		opt(obsv)
	}

	resource, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(obsv.serviceName), semconv.ServiceVersion(obsv.serviceVersion),
			semconv.ServiceInstanceID(obsv.serviceInstanceID),
		),
	)
	if err != nil {
		obsv.handleErr(ctx, err)
	}

	tpShutdown, err := obsv.newTraceProvider(ctx, resource)
	if err != nil {
		obsv.handleErr(ctx, err)
	}
	obsv.shutdownFuncs = append(obsv.shutdownFuncs, tpShutdown)

	mpShutdown, err := obsv.newMeterProvider(ctx, resource)
	if err != nil {
		obsv.handleErr(ctx, err)
	}
	obsv.shutdownFuncs = append(obsv.shutdownFuncs, mpShutdown)
}

// Trace Provider -.
func (ob observability) newTraceProvider(ctx context.Context, resource *resource.Resource) (func(context.Context) error, error) {
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(ob.traceExporterEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)

	return traceExporter.Shutdown, nil
}

// Metrics Provider -.
func (ob observability) newMeterProvider(ctx context.Context, resource *resource.Resource) (func(context.Context) error, error) {
	metricsExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(ob.metricsExporterEndpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(metricsExporter, sdkmetric.WithInterval(5*time.Second)),
		),
	)
	otel.SetMeterProvider(mp)

	return metricsExporter.Shutdown, nil
}

// handleErr handlers error by joining errors and safely shutting down otel's
// providers and then panics.
func (ob observability) handleErr(ctx context.Context, err error) {
	for _, f := range ob.shutdownFuncs {
		err = errors.Join(err, f(ctx))
	}

	if err != nil {
		log.Fatalf("FATAL - error while initiating opentelemetry: %s\n", err)
	}
}
