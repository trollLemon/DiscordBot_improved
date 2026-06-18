package jobs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"

	"goManip/jobs"
)

func createJobs(operation jobs.Operation, images []*gocv.Mat) []*jobs.Job {

	var jobArr []*jobs.Job

	for i, image := range images {
		jobArr = append(jobArr, jobs.NewJob(fmt.Sprintf("%d", i), operation, image))
	}
	return jobArr

}

func TestJob(t *testing.T) {

	tests := []struct {
		name string
		op   func() (jobs.Operation, error)
	}{
		{
			name: "Test Edge Detection",
			op:   func() (jobs.Operation, error) { return jobs.NewEdgeDetection(100.0, 200.0) },
		},
		{
			name: "Test Saturation",
			op:   func() (jobs.Operation, error) { return jobs.NewSaturate(1.6) },
		},
		{
			name: "Test Dilate",
			op:   func() (jobs.Operation, error) { return jobs.NewMorphology(3, 5, jobs.Dilate) },
		},
		{
			name: "Test Erode",
			op:   func() (jobs.Operation, error) { return jobs.NewMorphology(3, 5, jobs.Erode) },
		},
		{
			name: "Test Reduce",
			op:   func() (jobs.Operation, error) { return jobs.NewReduce(0.5) },
		},
		{
			name: "Test Random Filter",
			op:   func() (jobs.Operation, error) { return jobs.NewRandomFilter(3, -1, 1, true) },
		},
		{
			name: "Test Add Text",
			op:   func() (jobs.Operation, error) { return jobs.NewAddText("text", 1.0, 0.5, 0.5) },
		},
		{
			name: "Test Shuffle",
			op:   func() (jobs.Operation, error) { return jobs.NewShuffle(15, 1920, 1080) },
		},
		{
			name: "Test Invert",
			op:   func() (jobs.Operation, error) { return jobs.NewInvert(), nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := tt.op()
			assert.NoError(t, err)
			testImages := generateTestImages() // from operations_test.go

			for _, job := range createJobs(op, testImages) {
				ctx := context.Background()
				result, err := job.Process(ctx)
				assert.Nil(t, err)
				if result != nil {
					if job.GetTimeElapsed() == 0 {
						t.Errorf("Test: %s, expected job to take time to complete, StartTime: %d, EndTime: %d",
							tt.name, job.GetStartTime().UnixMilli(), job.GetEndTime().UnixMilli())
					}
					result.Close()
				}
			}
		})
	}

}
