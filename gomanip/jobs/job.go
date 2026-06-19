package jobs

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gocv.io/x/gocv"
)

// Operation represents an abstraction for some operation to perform on an image.
// A struct implementing Operation should have any parameters stored within the struct, so the Job struct can call Run and get a
// result image.
type Operation interface {
	// Run executes the image manipulation operation on the given input image, and returns the resulting image and error (if any).
	Run(input *gocv.Mat) (*gocv.Mat, error)

	// ToAttributes returns the operations parameters as a slice of OpenTelemetry attributes.
	// If this is not requird for a given operation, this can return nil or an empty slice.
	ToAttributes() []attribute.KeyValue
}

// NewJob creates a new job given an id, operation, and input image.
func NewJob(id string, operation Operation, image *gocv.Mat) *Job {
	return &Job{jobId: id, operation: operation, inputImage: image}
}

// Job contains the input image, operation object, and a jobId, as well as startTime and endTime which are populated after a worker calls Process().
type Job struct {
	jobId       string
	operation   Operation
	inputImage  *gocv.Mat
	startTime   time.Time
	endTime     time.Time
	elapsedTime time.Duration
}

// Process runs the operation on the stored input image, and returns the result image and error (if any)
func (j *Job) Process(ctx context.Context) (*gocv.Mat, error) {

	requestSpanCtx := trace.SpanContextFromContext(ctx)

	operationAttributes := j.operation.ToAttributes()

	tracer := otel.Tracer("goManip-worker")
	_, span := tracer.Start(ctx, "job.process", trace.WithLinks(
		trace.Link{SpanContext: requestSpanCtx}))
	span.SetAttributes(attribute.String("jobId", j.jobId))
	span.SetAttributes(operationAttributes...)
	defer span.End()

	defer j.inputImage.Close()

	j.startTime = time.Now()

	result, err := j.operation.Run(j.inputImage)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	j.elapsedTime = time.Since(j.startTime)

	j.endTime = time.Now()

	return result, err
}

// GetJobId returns the job id.
func (j *Job) GetJobId() string {
	return j.jobId
}

// GetTimeElapsed returns the total amount of time (in ns) the job took.
func (j *Job) GetTimeElapsed() int {
	return int(j.elapsedTime)
}

// GetStartTime returns the time the job started.
func (j *Job) GetStartTime() time.Time {
	return j.startTime
}

// GetEndTime returns the time the job ended.
func (j *Job) GetEndTime() time.Time {
	return j.endTime
}
