package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"gocv.io/x/gocv"

	"goManip/jobs"
	"goManip/store"
	"goManip/util"
	"goManip/worker"
)

var (
	ErrParseParams   = errors.New("failed to parse query parameters")
	ErrJobDispatcher = errors.New("failed to get job dispatcher")
)

func getDispatcher(c echo.Context) *jobs.JobDispatcher {
	return c.Get("jobDispatcher").(*jobs.JobDispatcher)
}

// JobResponse is the response body returned on a successful job enqueue.
type JobResponse struct {
	JobId string `json:"jobId"`
}

func handleImageOperation(
	c echo.Context,
	enqueueFunc func(image *gocv.Mat) (string, error),
) error {
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()
	image, err := util.GetImageFromBody(c)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to read image from request body", "error", err)
		return SendGomanipError(c, err, spanContext.TraceID().String())
	}

	jobId, err := enqueueFunc(image)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to enqueue job", "error", err)
		return SendGomanipError(c, err, spanContext.TraceID().String())
	}

	return c.JSON(http.StatusOK, &JobResponse{JobId: jobId})
}

func InvertEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()

	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	ctx := context.WithoutCancel(c.Request().Context())
	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueInvertImage(ctx, jobDispatcher, image)
	})

}

func SaturateEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()
	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	saturation, err := util.ParseSaturation(c)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse saturation", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())
	}

	ctx := context.WithoutCancel(c.Request().Context())

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueSaturateImage(ctx, jobDispatcher, image, saturation)
	})

}

func EdgeDetectionEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()
	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	tLower, tHigher, err := util.ParseEdgeDetection(c)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse edge detection", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())

	}

	ctx := context.WithoutCancel(c.Request().Context())

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueDetectEdges(ctx, jobDispatcher, image, tLower, tHigher)
	})
}

func MorphologyEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()
	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	morphType, kernelSize, iterations, err := util.ParseMorphology(c)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse morphology", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())

	}

	ctx := context.WithoutCancel(c.Request().Context())

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueMorphImage(ctx, jobDispatcher, image, jobs.Choice(morphType), kernelSize, iterations)
	})
}

func ReduceEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()

	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	quality, err := util.ParseReduce(c)

	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse reduce", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())

	}

	ctx := context.WithoutCancel(c.Request().Context())

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueReduceImage(ctx, jobDispatcher, image, quality)
	})
}

func AddTextEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()

	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	text, fontScale, xPerc, yPerc, err := util.ParseAddText(c)

	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse add text", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())

	}

	ctx := context.WithoutCancel(c.Request().Context())

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueAddText(ctx, jobDispatcher, image, text, fontScale, xPerc, yPerc)

	})
}

func RandomFilterEndpoint(c echo.Context) error {
	spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()
	jobDispatcher := getDispatcher(c)
	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	minVal, maxVal, kernelSize, normalize, err := util.ParseRandomFilter(c)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse random filter", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())

	}

	ctx := context.WithoutCancel(c.Request().Context())

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueRandomFilter(ctx, jobDispatcher, image, minVal, maxVal, kernelSize, normalize)
	})
}

func ShuffleEndpoint(c echo.Context) error {
	jobDispatcher := getDispatcher(c)
	span := trace.SpanFromContext(c.Request().Context())
	spanContext := span.SpanContext()

	if jobDispatcher == nil {
		slog.ErrorContext(c.Request().Context(), "Job dispatcher is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	partitions, err := util.ParseShuffle(c)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Failed to parse shuffle", "error", err)
		return SendGomanipError(c, ErrParseParams, spanContext.TraceID().String())

	}

	ctx := context.WithoutCancel(c.Request().Context())
	ctx = trace.ContextWithRemoteSpanContext(ctx, spanContext)

	return handleImageOperation(c, func(image *gocv.Mat) (string, error) {
		return jobs.EnqueueShuffle(ctx, jobDispatcher, image, partitions)
	})
}

func GetImageEndpoint(c echo.Context) error {
	jobId := c.Param("id")

	span := trace.SpanFromContext(c.Request().Context())
	spanContext := span.SpanContext()
	imageStore, ok := c.Get("store").(worker.ImageStore)
	if !ok {
		slog.ErrorContext(c.Request().Context(), "Image store is not present in the context")
		return SendGomanipError(c, ErrJobDispatcher, spanContext.TraceID().String())
	}

	ctx := context.WithoutCancel(c.Request().Context())
	ctx = trace.ContextWithRemoteSpanContext(ctx, spanContext)

	span.AddEvent("server.getImage", trace.WithAttributes(attribute.String("jobId", jobId)))
	imageBytes, traceID, err := imageStore.GetImage(ctx, jobId)

	if errors.Is(err, store.ErrJobFailed) {
		slog.ErrorContext(c.Request().Context(), fmt.Sprintf("Job %s failed during processing", jobId), "error", err, "traceId", traceID)
		return SendGomanipError(c, err, traceID)
	}
	if err != nil {
		slog.ErrorContext(c.Request().Context(), fmt.Sprintf("Failed to get image for job ID %s", jobId), "error", err)
		return SendGomanipError(c, err, spanContext.TraceID().String())
	}
	r, err := gzip.NewReader(bytes.NewReader(imageBytes))
	if err != nil {
		slog.ErrorContext(c.Request().Context(), fmt.Sprintf("Failed to create gzip reader for job ID %s", jobId), "error", err)
		return SendGomanipError(c, err, spanContext.TraceID().String())
	}
	defer r.Close()

	span.AddEvent("server.decompressImage", trace.WithAttributes(attribute.String("jobId", jobId)))
	decompressed, err := io.ReadAll(r)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), fmt.Sprintf("Failed to decompress image for job ID %s", jobId), "error", err)
		return SendGomanipError(c, err, spanContext.TraceID().String())
	}

	return c.Blob(http.StatusOK, "image/png", decompressed)

}

func MetricsEndpoint(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func InitRouting(e *echo.Echo, jobDispatcher *jobs.JobDispatcher, store worker.ImageStore) {
	e.Use(otelecho.Middleware("goManip"))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogURIPath: true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Status < http.StatusBadRequest {
				return nil
			}

			attrs := []any{
				"uri", v.URI,
				"path", v.URIPath,
				"status", v.Status,
				"method", c.Request().Method,
				"latency", v.Latency,
			}

			if v.Status >= http.StatusInternalServerError {
				slog.ErrorContext(c.Request().Context(), "http_request_failed", attrs...)
			} else {
				slog.WarnContext(c.Request().Context(), "http_request_failed", attrs...)
			}

			return nil
		},
	}))

	jobsGroup := e.Group("")
	jobsGroup.Use(JobDispatcherMiddleware(jobDispatcher))
	jobsGroup.Use(FileTypeVerifyMiddleware())

	jobsGroup.POST("/invert/", InvertEndpoint)
	jobsGroup.POST("/saturate/", SaturateEndpoint)
	jobsGroup.POST("/edgeDetection/", EdgeDetectionEndpoint)
	jobsGroup.POST("/morphology/", MorphologyEndpoint)
	jobsGroup.POST("/reduction/", ReduceEndpoint)
	jobsGroup.POST("/text/", AddTextEndpoint)
	jobsGroup.POST("/randomFilter/", RandomFilterEndpoint)
	jobsGroup.POST("/shuffle/", ShuffleEndpoint)

	resultGroup := e.Group("/image")
	resultGroup.Use(StoreMiddleware(store))
	e.GET("/metrics/", MetricsEndpoint)
	resultGroup.GET("/:id", GetImageEndpoint)

}

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
		semconv.ServiceNameKey.String("goManip"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.5))),
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
		semconv.ServiceNameKey.String("goManip"),
	)

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)

	global.SetLoggerProvider(lp)

	return lp.Shutdown, nil
}

func Start(e *echo.Echo, address, port string) {
	e.Logger.Fatal(e.Start(address + ":" + port))
}

func GraceFullShutdown(jobDispatcher *jobs.JobDispatcher, wg *sync.WaitGroup, cancel context.CancelFunc) {

	slog.Info("Closing worker request channels")
	jobDispatcher.Close()
	slog.Info("Waiting for workers to drain remaining jobs")
	wg.Wait()
	cancel()

}
