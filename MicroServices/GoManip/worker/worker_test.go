package worker_test

import (
	"context"
	"errors"
	"go.uber.org/goleak"
	"goManip/jobs"
	"goManip/worker"
	"gocv.io/x/gocv"
	"sync"
	"testing"
	"time"
)

type MockOperationSuccess struct {
}
type MockOperationErr struct {
}

type MockOperationTimeOut struct{}

func (m MockOperationErr) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return nil, errors.New("error processing job")
}

func (m MockOperationSuccess) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return input, nil
}

func (m MockOperationTimeOut) Run(input *gocv.Mat) (*gocv.Mat, error) {

	time.Sleep(1 * time.Second)

	return input, nil
}

func TestWorker(t *testing.T) {

	var tests = []struct {
		name      string
		operation jobs.Operation
		wantErr   bool
		canceled  bool
	}{
		{
			name:      "Test Success",
			operation: MockOperationSuccess{},
			wantErr:   false,
			canceled:  false,
		},
		{
			name:      "Test Error",
			operation: MockOperationErr{},
			wantErr:   true,
			canceled:  false,
		},
		{
			name:      "Test Canceled by timeout",
			operation: MockOperationTimeOut{},
			wantErr:   true,
			canceled:  true,
		},
	}

	defer goleak.VerifyNone(t)

	timeLimit := time.Second * 5
	testImage := gocv.NewMatWithSize(64, 64, gocv.MatTypeCV8UC3)
	defer testImage.Close()

	shutDownCtx, shutDownCancel := context.WithCancel(context.Background())

	workerWg := &sync.WaitGroup{}

	jobReqs := make(chan *jobs.JobRequest)

	go worker.Worker(shutDownCtx, 0, jobReqs, workerWg)
	workerWg.Add(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeLimit)

			job := jobs.NewJob(0, tt.operation, &testImage)
			jobRequest := jobs.NewJobRequest(job, ctx)

			jobReqs <- jobRequest

			result := <-jobRequest.Result

			image, err := result.Image, result.Error

			if !tt.wantErr && err != nil {
				t.Errorf("TestWorker() %s, expected err, got %v", tt.name, err)
			}

			if tt.canceled && image != nil && err != nil {
				t.Errorf("TestWorker() %s, expected err and nil image due to timeout", tt.name)
			}

			if result.Image != nil {
				result.Image.Close()
			}
			cancel()
		})
	}
	close(jobReqs)
	shutDownCancel()
	workerWg.Wait()

}

func TestWorkerShutdown(t *testing.T) {
	wg := &sync.WaitGroup{}
	jobReqs := make(chan *jobs.JobRequest)
	testImage := gocv.NewMat()
	ctx, cancel := context.WithCancel(context.Background())
	go worker.Worker(ctx, 0, jobReqs, wg)

	wg.Add(1)
	cancel()
	jobReqs <- jobs.NewJobRequest(jobs.NewJob(0, MockOperationSuccess{}, &testImage), context.Background())

	close(jobReqs)
	wg.Wait()
}

func TestWorkerCancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	jobReqs := make(chan *jobs.JobRequest)
	testImage := gocv.NewMat()
	go worker.Worker(context.Background(), 0, jobReqs, wg)
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	jobReqs <- jobs.NewJobRequest(jobs.NewJob(0, MockOperationSuccess{}, &testImage), ctx)

	close(jobReqs)
	wg.Wait()
}
