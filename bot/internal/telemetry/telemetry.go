package telemetry

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const serviceName = "discord-bot"

// InitTracer configures the global OpenTelemetry tracer provider so spans are exported to the OTLP
// endpoint over HTTP. It returns a shutdown function used to flush pending spans.
func InitTracer(endpoint string) (func(context.Context) error, error) {
	ctx := context.Background()
	clientOptions := []otlptracehttp.Option{otlptracehttp.WithInsecure()}
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		clientOptions = append(clientOptions, otlptracehttp.WithEndpointURL(endpoint))
	} else {
		clientOptions = append(clientOptions, otlptracehttp.WithEndpoint(endpoint))
	}

	exporter, err := otlptracehttp.New(ctx, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize otlp trace exporter: %w", err)
	}

	res := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// InitLogger configures the global OpenTelemetry logger provider so that otelslog-backed loggers
// export records to the OTLP endpoint. It returns a shutdown function used to flush pending logs.
func InitLogger(endpoint string) (func(context.Context) error, error) {
	ctx := context.Background()
	clientOptions := []otlploghttp.Option{otlploghttp.WithInsecure()}
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		clientOptions = append(clientOptions, otlploghttp.WithEndpointURL(endpoint))
	} else {
		clientOptions = append(clientOptions, otlploghttp.WithEndpoint(endpoint))
	}

	exporter, err := otlploghttp.New(ctx, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize otlp log exporter: %w", err)
	}

	res := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)

	global.SetLoggerProvider(lp)

	return lp.Shutdown, nil
}
