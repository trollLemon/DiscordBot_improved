package jobs_test

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"

	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"gocv.io/x/gocv"

	"goManip/jobs"
	"goManip/store"
	"goManip/worker"
)

func TestDispatchJob(t *testing.T) {
	defer goleak.VerifyNone(t)

	testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3)
	defer testImage.Close()

	requests := make(chan *jobs.JobRequest, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobDispatcher := jobs.NewJobDispatcher(requests)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	store := store.NewRedisStore(mr.Addr(), "", "", 0, time.Hour)
	defer store.Close()

	worker := worker.NewWorker(ctx, 0, requests, wg, store, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	go worker.Work()

	job := jobs.NewJob("0", jobs.NewInvert(), &testImage)
	jobId, err := jobDispatcher.DispatchJob(ctx, job)
	assert.NoError(t, err)
	assert.Equal(t, "0", jobId)

	cancel()
	jobDispatcher.Close()
	wg.Wait()
}

func TestDispatchWithActualOps(t *testing.T) {
	testctx := context.Background()
	tests := []struct {
		name string
		fn   func(*jobs.JobDispatcher, *gocv.Mat) (string, error)
	}{
		{
			name: "Test Invert Image",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueInvertImage(testctx, jobDispatcher, image)
			},
		},
		{
			name: "Test Saturate Image",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueSaturateImage(testctx, jobDispatcher, image, 1.3)
			},
		},
		{
			name: "Test Edge Detection",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueDetectEdges(testctx, jobDispatcher, image, 100, 200)
			},
		},
		{
			name: "Test Edge Detection (zero threshold)",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueDetectEdges(testctx, jobDispatcher, image, 0, 200)
			},
		},
		{
			name: "Test Morphology (Dilation)",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueMorphImage(testctx, jobDispatcher, image, jobs.Dilate, 3, 3)
			},
		},
		{
			name: "Test Morphology (Erosion)",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueMorphImage(testctx, jobDispatcher, image, jobs.Erode, 3, 3)
			},
		},

		{
			name: "Test Image Reduction",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueReduceImage(testctx, jobDispatcher, image, 0.5)
			},
		},
		{
			name: "Test Add Text",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueAddText(testctx, jobDispatcher, image, "I love golang", 1.0, 0.5, 0.5)
			},
		},
		{
			name: "Test Random Filter",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueRandomFilter(testctx, jobDispatcher, image, -1, 1, 3, true)
			},
		},

		{
			name: "Test Shuffle",
			fn: func(jobDispatcher *jobs.JobDispatcher, image *gocv.Mat) (string, error) {
				return jobs.EnqueueShuffle(testctx, jobDispatcher, image, 64)
			},
		},
	}

	defer goleak.VerifyNone(t)

	numWorkers := 10

	requestChan := make(chan *jobs.JobRequest, numWorkers)
	jobDispatcher := jobs.NewJobDispatcher(requestChan)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	store := store.NewRedisStore(mr.Addr(), "", "", 0, time.Hour)
	defer store.Close()

	for idx := range numWorkers {
		worker := worker.NewWorker(ctx, idx, requestChan, wg, store, nil)
		wg.Add(1)

		go worker.Work()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3)
			jobId, err := tt.fn(jobDispatcher, &testImage)
			assert.NoError(t, err)
			assert.NotEmpty(t, jobId)
		})
	}

	jobDispatcher.Close()
	wg.Wait()
}
