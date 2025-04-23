package handlers_test

import (
	"go.uber.org/goleak"
	"sync"

	"context"
	"errors"
	"gocv.io/x/gocv"
	"image_manip/handlers"
	"image_manip/jobs"
	"image_manip/worker"
	"testing"
	"time"
)

type MockOperationSuccess struct {
}
type MockOperationErr struct {
}

type MockOperationTimeOut struct{}

func (m MockOperationTimeOut) Run(input *gocv.Mat) (*gocv.Mat, error) {

	time.Sleep(time.Millisecond * 3)
	return input, nil
}

func (m MockOperationErr) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return nil, errors.New("error processing job")
}

func (m MockOperationSuccess) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return input, nil
}

func TestImageHandler(t *testing.T) {
	defer goleak.VerifyNone(t)
	requests := make(chan *jobs.JobRequest)
	testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV32F)
	testCtx, timeOutCancel := context.WithTimeout(context.Background(), time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer testImage.Close()
	defer cancel()
	defer timeOutCancel()

	tests := []struct {
		name        string
		wantErr     bool
		wantTimeout bool
		job         *jobs.Job
		ctx         context.Context
	}{
		{
			name:        "Test success",
			wantErr:     false,
			wantTimeout: false,
			job:         jobs.NewJob(0, MockOperationSuccess{}, &testImage),
		},
		{
			name:        "Test Failure",
			wantErr:     true,
			wantTimeout: false,
			job:         jobs.NewJob(0, MockOperationErr{}, &testImage),
		},
		{
			name:        "Test Timeout",
			wantErr:     true,
			wantTimeout: false,
			job:         jobs.NewJob(0, MockOperationTimeOut{}, &testImage),
		},
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go worker.Worker(ctx, 0, requests, wg)

	timeOutError := errors.New("job cancelled due to timeout")

	for _, tt := range tests {
		_, err := handlers.ProcessImage(tt.job, requests, testCtx)

		if tt.wantTimeout && !errors.Is(err, timeOutError) {
			t.Errorf("TestInvertImageHandler() error = %v, wantErr %v, expected %s", err, tt.wantErr, timeOutError.Error())
		}

		if (err != nil) != tt.wantErr {
			t.Errorf("TestInvertImageHandler() Test %s, error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}

	}
}
