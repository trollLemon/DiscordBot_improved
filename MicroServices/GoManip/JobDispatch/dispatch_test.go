package JobDispatch_test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"sync"

	"context"
	"errors"
	"goManip/JobDispatch"
	"goManip/jobs"
	"goManip/worker"
	"gocv.io/x/gocv"
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

func TestDispatchJob(t *testing.T) {

	testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3)

	dispatchTests := []struct {
		name        string
		wantErr     bool
		wantTimeout bool
		job         *jobs.Job
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
			job:         jobs.NewJob(1, MockOperationErr{}, &testImage),
		},
		{
			name:        "Test Timeout",
			wantErr:     true,
			wantTimeout: true,
			job:         jobs.NewJob(2, MockOperationTimeOut{}, &testImage),
		},
	}

	defer goleak.VerifyNone(t)
	requests := make(chan *jobs.JobRequest)
	ctx, cancel := context.WithCancel(context.Background())

	jobDispatcher := JobDispatch.NewJobDispatcher(requests, time.Millisecond*2)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go worker.Worker(ctx, 0, requests, wg)

	timeOutError := errors.New("job cancelled due to timeout")

	for _, tt := range dispatchTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jobDispatcher.DispatchJob(tt.job)
			if tt.wantTimeout {
				assert.Error(t, err, timeOutError.Error())
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

	cancel()
	jobDispatcher.Close()
	wg.Wait()
	testImage.Close()

}
