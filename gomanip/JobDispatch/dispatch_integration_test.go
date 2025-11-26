package JobDispatch_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"goManip/JobDispatch"
	"goManip/jobs"
	"goManip/worker"
	"gocv.io/x/gocv"
	"sync"
	"testing"
	"time"
)

func TestDispatchIntegration(t *testing.T) {

	integrationTests := []struct {
		name    string
		wantErr bool
		fn      func(*JobDispatch.JobDispatcher, *gocv.Mat) (*gocv.NativeByteBuffer, error)
	}{
		{
			name:    "Test Invert Image",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueInvertImage(jobDispatcher, image)
			},
		},
		{
			name:    "Test Invert Image Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueInvertImage(jobDispatcher, nil)
			},
		},
		{
			name:    "Test Saturate Image",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueSaturateImage(jobDispatcher, image, 1.3)
			},
		},
		{
			name:    "Test Saturate Image Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueSaturateImage(jobDispatcher, image, -1.3)
			},
		},
		{
			name:    "Test Edge Detection",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueDetectEdges(jobDispatcher, image, 100, 200)
			},
		},
		{
			name:    "Test Edge Detection Error",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueDetectEdges(jobDispatcher, image, 0, 200)
			},
		},

		{
			name:    "Test Morphology (Dilation)",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Dilate, 3, 3)
			},
		},
		{
			name:    "Test Morphology (Erosion)",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Erode, 3, 3)
			},
		},
		{
			name:    "Test Morphology Error (invalid morph op)",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueMorphImage(jobDispatcher, image, "wrongOP", 3, 3)
			},
		},
		{
			name:    "Test Morphology Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueMorphImage(jobDispatcher, image, jobs.Erode, -3, 3)
			},
		},
		{
			name:    "Test Image Reduction",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueReduceImage(jobDispatcher, image, 0.5)
			},
		},
		{
			name:    "Test Image Reduction Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueReduceImage(jobDispatcher, image, 0.0)
			},
		},
		{
			name:    "Test Add Text",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueAddText(jobDispatcher, image, "I love golang", 1.0, 0.5, 0.5)
			},
		},
		{
			name:    "Test Add Text Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueAddText(jobDispatcher, image, "", 1.0, 0.5, 0.5)
			},
		},
		{
			name:    "Test Random Filter",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueRandomFilter(jobDispatcher, image, -1, 1, 3, true)
			},
		},
		{
			name:    "Test Random Filter Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueRandomFilter(jobDispatcher, image, -1, 1, 0, true)
			},
		},
		{
			name:    "Test Shuffle",
			wantErr: false,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueShuffle(jobDispatcher, image, 64)
			},
		},
		{
			name:    "Test Shuffle Error",
			wantErr: true,
			fn: func(jobDispatcher *JobDispatch.JobDispatcher, image *gocv.Mat) (*gocv.NativeByteBuffer, error) {
				return JobDispatch.EnqueueShuffle(jobDispatcher, image, 0)
			},
		},
	}

	defer goleak.VerifyNone(t)
	numWorkers := 10

	requestChan := make(chan *jobs.JobRequest, numWorkers)
	maxTime := time.Second * 1
	jobDispatcher := JobDispatch.NewJobDispatcher(requestChan, maxTime)
	wg := new(sync.WaitGroup)
	timeOutContext, cancel := context.WithTimeout(context.Background(), maxTime)
	defer cancel()

	testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3)

	for idx := range numWorkers {
		wg.Add(1)
		go worker.Worker(timeOutContext, idx, requestChan, wg)
	}

	for _, tt := range integrationTests {
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
