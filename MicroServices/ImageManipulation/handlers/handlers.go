package handlers

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
	"image_manip/jobs"
	"time"
)

var (
	jobId = 0
)

func getNewJobId() uint8 {
	jobId++
	return uint8(jobId + 1)
}

func ProcessImage(job *jobs.Job, jobRequests chan *jobs.JobRequest, ctx context.Context) ([]byte, error) {
	jobRequest := jobs.NewJobRequest(job)
	resultChan := jobRequest.Result
	errChan := jobRequest.Error
	jobRequests <- jobRequest

	//defer cancel()
	defer close(resultChan)
	defer close(errChan)

	select {
	case result := <-resultChan:
		{
			buf, err := gocv.IMEncode(".png", *result)
			if err != nil {
				log.Err(err).Msg("failed to encode result into png format")
				return nil, err
			}
			return buf.GetBytes(), nil
		}
	case err := <-errChan:
		{
			return nil, err
		}

	case <-ctx.Done():
		return nil, errors.New("job cancelled due to timeout")

	}

}

func InvertImage(jobRequests chan *jobs.JobRequest, image *gocv.Mat) ([]byte, error) {
	job := jobs.NewJob(getNewJobId(), jobs.NewInvert(), image)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ProcessImage(job, jobRequests, ctx)

}

func SaturateImage(jobRequests chan *jobs.JobRequest, image *gocv.Mat, value float32) ([]byte, error) {
	job := jobs.NewJob(getNewJobId(), jobs.NewSaturate(value), image)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ProcessImage(job, jobRequests, ctx)
}

func DetectEdges(jobRequests chan *jobs.JobRequest, image *gocv.Mat, tLower, tHigher float32) ([]byte, error) {
	job := jobs.NewJob(getNewJobId(), jobs.NewEdgeDetection(tLower, tHigher), image)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ProcessImage(job, jobRequests, ctx)

}
