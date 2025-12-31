package jobdispatch_test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"sync"

	"context"
	"errors"
	"goManip/jobdispatch"
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

	jobDispatcher := jobdispatch.NewJobDispatcher(requests, time.Millisecond*2)

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

func TestDispatchWithActualOps(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		fn      func(*jobdispatch.JobDispatcher, *gocv.Mat) (*gocv.NativeByteBuffer, error)
	}{
		{
			name:    "Test Invert Image",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueInvertImage(jobDispatcher, image)
			},
		},
		{
			name:    "Test Invert Image Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueInvertImage(jobDispatcher, nil)
			},
		},
		{
			name:    "Test Saturate Image",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueSaturateImage(jobDispatcher, image, 1.3)
			},
		},
		{
			name:    "Test Saturate Image Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueSaturateImage(jobDispatcher, image, -1.3)
			},
		},
		{
			name:    "Test Edge Detection",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueDetectEdges(jobDispatcher, image, 100, 200)
			},
		},
		{
			name:    "Test Edge Detection Error",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueDetectEdges(jobDispatcher, image, 0, 200)
			},
		},

		{
			name:    "Test Morphology (Dilation)",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Dilate, 3, 3)
			},
		},
		{
			name:    "Test Morphology (Erosion)",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Erode, 3, 3)
			},
		},
		{
			name:    "Test Morphology Error (invalid morph op)",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueMorphImage(jobDispatcher, image, "wrongOP", 3, 3)
			},
		},
		{
			name:    "Test Morphology Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Erode, -3, 3)
			},
		},
		{
			name:    "Test Image Reduction",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueReduceImage(jobDispatcher, image, 0.5)
			},
		},
		{
			name:    "Test Image Reduction Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueReduceImage(jobDispatcher, image, 0.0)
			},
		},
		{
			name:    "Test Add Text",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueAddText(jobDispatcher, image, "I love golang", 1.0, 0.5, 0.5)
			},
		},
		{
			name:    "Test Add Text Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueAddText(jobDispatcher, image, "", 1.0, 0.5, 0.5)
			},
		},
		{
			name:    "Test Random Filter",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueRandomFilter(jobDispatcher, image, -1, 1, 3, true)
			},
		},
		{
			name:    "Test Random Filter Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueRandomFilter(jobDispatcher, image, -1, 1, 0, true)
			},
		},
		{
			name:    "Test Shuffle",
			wantErr: false,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueShuffle(jobDispatcher, image, 64)
			},
		},
		{
			name:    "Test Shuffle Error",
			wantErr: true,
			fn: func(jobDispatcher *jobdispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return jobdispatch.EnqueueShuffle(jobDispatcher, image, 0)
			},
		},
	}

	defer goleak.VerifyNone(t)
	numWorkers := 10

	requestChan := make(chan *jobs.JobRequest, numWorkers)
	maxTime := time.Second * 1
	jobDispatcher := jobdispatch.NewJobDispatcher(requestChan, maxTime)
	wg := new(sync.WaitGroup)
	timeOutContext, cancel := context.WithTimeout(context.Background(), maxTime)
	defer cancel()

	testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3)

	for idx := range numWorkers {
		wg.Add(1)
		go worker.Worker(timeOutContext, idx, requestChan, wg)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bytes, err := tt.fn(jobDispatcher, &testImage)

			assert.Equal(t, tt.wantErr, err != nil)

			if (bytes == nil || bytes.Len() == 0) && !tt.wantErr {
				t.Errorf("TestDispatchIntegration() %s, bytes is empty", tt.name)
			}
			if bytes != nil {
				bytes.Close()

			}
		})
	}

	jobDispatcher.Close()
	wg.Wait()
}
