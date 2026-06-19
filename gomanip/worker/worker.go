package worker

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"goManip/jobs"
	"io"
	"log/slog"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrWorkerPanic = errors.New("worker panicked")
	ErrOpenCV      = errors.New("opencv error")
	ErrUnknown     = errors.New("unknown error")
)

type ImageStore interface {
	// SaveImage saves the processed image bytes with the associated job ID.
	SaveImage(ctx context.Context, jobId string, image []byte) error
	// WriteError reports an error that occurred during job processing, associated with the job ID and
	// trace ID for the failed job.
	WriteError(ctx context.Context, jobId string, err error, traceId string) error
	// GetImage retrieves the processed image bytes for a given job ID.
	// If the job failed, returns an error and the trace ID associated with the failed job if the error wraps ErrJobFailed.
	GetImage(ctx context.Context, jobId string) ([]byte, string, error)
}

type Worker struct {
	shutdownctx context.Context
	id          int
	jobRequests <-chan *jobs.JobRequest
	wg          *sync.WaitGroup
	store       ImageStore
	logger      *slog.Logger
}

func NewWorker(ctx context.Context, id int, jobRequests <-chan *jobs.JobRequest, wg *sync.WaitGroup, store ImageStore, logger *slog.Logger) *Worker {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &Worker{
		shutdownctx: ctx,
		id:          id,
		jobRequests: jobRequests,
		wg:          wg,
		store:       store,
		logger:      logger,
	}
}

// Work drains the job request channel until it is closed, processing each request.
func (w *Worker) Work() {
	defer w.wg.Done()

	for jobRequest := range w.jobRequests {
		w.process(jobRequest)
	}

	w.logger.InfoContext(w.shutdownctx, "Worker shutting down", "worker", w.id)
}

func (w *Worker) process(jobRequest *jobs.JobRequest) {
	job := jobRequest.Job
	ctx := jobRequest.Ctx
	jobId := job.GetJobId()

	requestSpanCtx := trace.SpanContextFromContext(ctx)
	tracer := otel.Tracer("goManip-worker")
	spanCtx, span := tracer.Start(ctx, "worker.work", trace.WithLinks(
		trace.Link{SpanContext: requestSpanCtx}))
	span.SetAttributes(attribute.Int("workerId", w.id), attribute.String("jobId", jobId))
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("%w: %v", ErrWorkerPanic, r)
			w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d panicked while processing job %s", w.id, jobId), "error", panicErr)
			span.RecordError(panicErr)
			span.SetStatus(codes.Error, "worker panicked")
			if writeErr := w.store.WriteError(spanCtx, jobId, panicErr, span.SpanContext().TraceID().String()); writeErr != nil {
				w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: failed to write panic error for job %s", w.id, jobId), "error", writeErr)
			}
		}
	}()

	processedImage, err := job.Process(spanCtx)

	switch {
	case errors.Is(err, jobs.ErrOpenCV):
		w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: error during opencv op: %v", w.id, err), "error", ErrOpenCV)
	case err != nil:
		w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: job %s failed for an unknown reason: %v", w.id, jobId, err), "error", ErrUnknown)
	}
	if err != nil {
		if writeErr := w.store.WriteError(spanCtx, jobId, err, span.SpanContext().TraceID().String()); writeErr != nil {
			w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: failed to write error for job %s", w.id, jobId), "error", writeErr)
		}
		return
	}
	defer processedImage.Close()

	span.AddEvent("worker.complete", trace.WithAttributes(
		attribute.Int("workerId", w.id),
		attribute.String("jobId", jobId),
		attribute.Int64("jobStart", job.GetStartTime().UnixMilli()),
		attribute.Int64("jobEnd", job.GetEndTime().UnixMilli()),
	))

	imageBytes := processedImage.ToBytes()

	span.AddEvent("worker.compress", trace.WithAttributes(attribute.Int("workerId", w.id), attribute.String("jobId", jobId)))
	var buf bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: failed to set gzip compression level for job %s", w.id, jobId), "error", err)
		return
	}

	if _, err := gzipWriter.Write(imageBytes); err != nil {
		w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: failed to compress image for job %s", w.id, jobId), "error", err)
		return
	}

	if err := gzipWriter.Close(); err != nil {
		w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: failed to finalize compression for job %s", w.id, jobId), "error", err)
		return
	}
	compressedImage := buf.Bytes()

	span.AddEvent("worker.compress.finished", trace.WithAttributes(
		attribute.Int("workerId", w.id),
		attribute.String("jobId", jobId),
		attribute.Int64("compressedSize", int64(len(compressedImage))),
		attribute.Int64("originalSize", int64(len(imageBytes))),
	))

	if err := w.store.SaveImage(spanCtx, jobId, compressedImage); err != nil {
		w.logger.ErrorContext(spanCtx, fmt.Sprintf("Worker %d: failed to save image for job %s", w.id, jobId), "error", err)
		return
	}

	w.logger.InfoContext(spanCtx, fmt.Sprintf("Job %s processed successfully", jobId), "worker", w.id, "jobId", jobId)
}
