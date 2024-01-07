package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

// Tracer provider for running unit test
type TestingTracerProvider struct{ embedded.TracerProvider }

func NewTestingTracerProvider() TestingTracerProvider {
	return TestingTracerProvider{}
}

func (TestingTracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return TestingTracer{}
}

// Tracer object for running unit test
type TestingTracer struct{ embedded.Tracer }

func (t TestingTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	// We return the ctx directly without needing to modify it
	return ctx, trace.SpanFromContext(context.TODO())
}
