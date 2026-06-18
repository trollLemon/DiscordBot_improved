package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gocv.io/x/gocv"
)

var (
	ErrImgEmpty          = errors.New("input image is empty")
	ErrInvalidParameters = errors.New("invalid parameters provided for operation")
)

// JobRequest represents a request to process a job, containing the job and its associated context for tracing.
type JobRequest struct {
	Job *Job
	Ctx context.Context
}

type JobDispatcher struct {
	jobRequests chan<- *JobRequest
}

func NewJobDispatcher(jobRequests chan<- *JobRequest) *JobDispatcher {
	return &JobDispatcher{jobRequests: jobRequests}
}

// getNewJobId generates a UUID for a new job.
func (j *JobDispatcher) getNewJobId() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate job ID: %w", err)
	}
	return id.String(), nil
}

func (j *JobDispatcher) Close() {
	close(j.jobRequests)
}

// DispatchJob enqueues the job for processing by a worker and immediately returns the job ID.
func (j *JobDispatcher) DispatchJob(ctx context.Context, job *Job) (string, error) {
	j.jobRequests <- &JobRequest{Job: job, Ctx: ctx}
	return job.GetJobId(), nil
}

func EnqueueInvertImage(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op := NewInvert()

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueSaturateImage(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, value float32) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewSaturate(value)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueDetectEdges(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, tLower, tHigher float32) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewEdgeDetection(tLower, tHigher)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueMorphImage(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, choice Choice, kernelSize, iterations int) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewMorphology(kernelSize, iterations, choice)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueReduceImage(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, reduceValue float32) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewReduce(reduceValue)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueAddText(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, text string, fontScale, xPerc, yPerc float64) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewAddText(text, fontScale, xPerc, yPerc)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueRandomFilter(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, min, max, kernelSize int, normalize bool) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewRandomFilter(kernelSize, min, max, normalize)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}

func EnqueueShuffle(ctx context.Context, dispatcher *JobDispatcher, image *gocv.Mat, partitions int) (string, error) {
	if image.Empty() {
		return "", ErrImgEmpty
	}

	op, err := NewShuffle(partitions, image.Rows(), image.Cols())
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidParameters, err)
	}

	jobId, err := dispatcher.getNewJobId()
	if err != nil {
		return "", err
	}

	return dispatcher.DispatchJob(ctx, NewJob(jobId, op, image))
}
