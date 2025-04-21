package worker_test

import (
	"errors"
	"gocv.io/x/gocv"
	"image_manip/jobs"
	"image_manip/worker"
	"testing"
	"time"
)

type MockOperationSuccess struct {
}
type MockOperationErr struct {
}

func (m MockOperationErr) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return nil, errors.New("error processing job")
}

func (m MockOperationSuccess) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return input, nil
}

func TestWorker(t *testing.T) {
	timeLimit := time.Second * 5
	testImage := gocv.NewMat()
	tests := []struct {
		name              string
		wantErr           bool
		jobRequests       []*jobs.JobRequest
		jobRequestChannel chan *jobs.JobRequest
	}{
		{
			name:              "Test success",
			wantErr:           false,
			jobRequests:       []*jobs.JobRequest{jobs.NewJobRequest(jobs.NewJob(0, MockOperationSuccess{}, &testImage))},
			jobRequestChannel: make(chan *jobs.JobRequest, 1),
		},
		{
			name:              "Test Failure",
			wantErr:           true,
			jobRequests:       []*jobs.JobRequest{jobs.NewJobRequest(jobs.NewJob(0, MockOperationErr{}, &testImage))},
			jobRequestChannel: make(chan *jobs.JobRequest, 1),
		},
		{
			name:    "Test Success multiple jobs",
			wantErr: false,
			jobRequests: []*jobs.JobRequest{
				jobs.NewJobRequest(jobs.NewJob(0, MockOperationSuccess{}, &testImage)),
				jobs.NewJobRequest(jobs.NewJob(1, MockOperationSuccess{}, &testImage)),
				jobs.NewJobRequest(jobs.NewJob(2, MockOperationSuccess{}, &testImage)),
				jobs.NewJobRequest(jobs.NewJob(3, MockOperationSuccess{}, &testImage)),
			},
			jobRequestChannel: make(chan *jobs.JobRequest, 1),
		},
	}

	for _, tt := range tests {
		go worker.Worker(0, tt.jobRequestChannel)

		for _, jobRequest := range tt.jobRequests {
			tt.jobRequestChannel <- jobRequest

			select {
			case err := <-jobRequest.Error:
				if (err != nil) != tt.wantErr {
					t.Errorf("Worker() error = %v, wantErr %v", err, tt.wantErr)
				}
			case result := <-jobRequest.Result:
				if result == nil {
					t.Errorf("Worker() result is nil, expected result")
				}
			case <-time.After(timeLimit):
				t.Errorf("worker() timeout, did not send to err or result channels")
			}
			close(jobRequest.Error)
			close(jobRequest.Result)
		}

		close(tt.jobRequestChannel)
	}

}
