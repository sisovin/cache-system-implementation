package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer      trace.Tracer
	cacheHits   metric.Int64Counter
	cacheMisses metric.Int64Counter
)

func init() {
	tracer = otel.Tracer("cache-system")
	meter := otel.Meter("cache-system")
	cacheHits = meter.NewInt64Counter("cache_hits", metric.WithDescription("Number of cache hits"))
	cacheMisses = meter.NewInt64Counter("cache_misses", metric.WithDescription("Number of cache misses"))
}

func InjectSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, operationName)
	span.SetAttributes(attribute.String("component", "cache"))
	return ctx, span
}

func EndSpan(span trace.Span) {
	span.End()
}

func RecordCacheHit(ctx context.Context) {
	cacheHits.Add(ctx, 1)
}

func RecordCacheMiss(ctx context.Context) {
	cacheMisses.Add(ctx, 1)
}

func RecordCacheOperation(ctx context.Context, operationName string, operation func() error) error {
	ctx, span := InjectSpan(ctx, operationName)
	defer EndSpan(span)

	err := operation()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetStatus(trace.StatusCodeOk, "Success")
	}

	return err
}
