package JobDispatch

import (
	"context"
	"errors"
	"goManip/jobs"
	"gocv.io/x/gocv"
	"sync/atomic"
	"time"
)

type JobDispatcher struct {
	jobId       uint32
	jobRequests chan<- *jobs.JobRequest
	maxTime     time.Duration
}

func NewJobDispatcher(jobRequests chan<- *jobs.JobRequest, maxTime time.Duration) *JobDispatcher {
	return &JobDispatcher{jobId: 0, jobRequests: jobRequests, maxTime: maxTime}
}

func (j *JobDispatcher) awaitResult(jobRequest *jobs.JobRequest, ctx context.Context) (*gocv.NativeByteBuffer, error) {

	select {
	case result := <-jobRequest.Result:

		image, err := result.Image, result.Error

		if err != nil {
			return nil, err
		}

		imageBytes, err := gocv.IMEncode(".png", *image)

		if err != nil {
			return nil, err
		}

		return imageBytes, nil

	case <-ctx.Done():
		return nil, errors.New("job cancelled due to timeout")

	}

}

func (j *JobDispatcher) getNewJobId() uint32 {
	return atomic.AddUint32(&j.jobId, 1)
}

func (j *JobDispatcher) Close() {
	close(j.jobRequests)
}

func (j *JobDispatcher) DispatchJob(job *jobs.Job) (*gocv.NativeByteBuffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), j.maxTime)
	jobRequest := jobs.NewJobRequest(job, ctx)
	j.jobRequests <- jobRequest
	defer cancel()
	return j.awaitResult(jobRequest, ctx)

}

func EnqueueInvertImage(dispatcher *JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewInvert(), image)
	return dispatcher.DispatchJob(job)
}

func EnqueueSaturateImage(dispatcher *JobDispatcher, image *gocv.Mat, value float32) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewSaturate(value), image)
	return dispatcher.DispatchJob(job)
}

func EnqueueDetectEdges(dispatcher *JobDispatcher, image *gocv.Mat, tLower, tHigher float32) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewEdgeDetection(tLower, tHigher), image)
	return dispatcher.DispatchJob(job)

}

func EnqueueMorphImage(dispatcher *JobDispatcher, image *gocv.Mat, choice jobs.Choice, kernelSize, iterations int) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewMorphology(kernelSize, iterations, choice), image)
	return dispatcher.DispatchJob(job)
}

func EnqueueReduceImage(dispatcher *JobDispatcher, image *gocv.Mat, reduceValue float32) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewReduce(reduceValue), image)
	return dispatcher.DispatchJob(job)
}

func EnqueueAddText(dispatcher *JobDispatcher, image *gocv.Mat, text string, fontScale, xPerc, yPerc float64) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewAddText(text, fontScale, xPerc, yPerc), image)
	return dispatcher.DispatchJob(job)

}

func EnqueueRandomFilter(dispatcher *JobDispatcher, image *gocv.Mat, min, max, kernelSize int, normalize bool) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewRandomFilter(kernelSize, min, max, normalize), image)
	return dispatcher.DispatchJob(job)
}

func EnqueueShuffle(dispatcher *JobDispatcher, image *gocv.Mat, partitions int) (*gocv.NativeByteBuffer, error) {
	job := jobs.NewJob(dispatcher.getNewJobId(), jobs.NewShuffle(partitions), image)
	return dispatcher.DispatchJob(job)
}
