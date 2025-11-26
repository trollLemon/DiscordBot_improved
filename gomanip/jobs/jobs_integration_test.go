package jobs_test

import (
	"github.com/stretchr/testify/assert"
	"goManip/jobs"
	"gocv.io/x/gocv"
	"testing"
)

func createJobs(operation jobs.Operation, images []*gocv.Mat) []*jobs.Job {

	var jobsArr []*jobs.Job

	for i, image := range images {
		jobsArr = append(jobsArr, jobs.NewJob(uint32(i), operation, image))
	}
	return jobsArr

}

func TestIntegration(t *testing.T) {

	testImages := generateTestImages() // from operations_test.go

	tests := []struct {
		name      string
		wantError bool
		jobs      []*jobs.Job
	}{
		{
			name:      "Test Invert Job",
			wantError: false,
			jobs:      createJobs(jobs.NewInvert(), testImages),
		},
		{
			name:      "Test Shuffle Job",
			wantError: false,
			jobs:      createJobs(jobs.NewShuffle(15), testImages),
		},
		{
			name:      "Test Shuffle Job Error",
			wantError: true,
			jobs:      createJobs(jobs.NewShuffle(0), testImages),
		},
		{
			name:      "Test Edge Detection",
			wantError: false,
			jobs:      createJobs(jobs.NewEdgeDetection(100.0, 200.0), testImages),
		},
		{
			name:      "Test Edge Detection Error",
			wantError: true,
			jobs:      createJobs(jobs.NewEdgeDetection(-1.0, 200.0), testImages),
		},
		{
			name:      "Test Saturation",
			wantError: false,
			jobs:      createJobs(jobs.NewSaturate(1.6), testImages),
		},
		{
			name:      "Test Saturation Error",
			wantError: true,
			jobs:      createJobs(jobs.NewSaturate(-1.6), testImages),
		},
		{
			name:      "Test Dilate",
			wantError: false,
			jobs:      createJobs(jobs.NewMorphology(3, 5, jobs.Dilate), testImages),
		},
		{
			name:      "Test Dilate Error",
			wantError: true,
			jobs:      createJobs(jobs.NewMorphology(0, 5, jobs.Dilate), testImages),
		},
		{
			name:      "Test Erode",
			wantError: false,
			jobs:      createJobs(jobs.NewMorphology(3, 5, jobs.Erode), testImages),
		},
		{
			name:      "Test Erode Error",
			wantError: true,
			jobs:      createJobs(jobs.NewMorphology(0, 5, jobs.Erode), testImages),
		},
		{
			name:      "Test Reduce",
			wantError: false,
			jobs:      createJobs(jobs.NewReduce(0.5), testImages),
		},
		{
			name:      "Test Reduce Error",
			wantError: true,
			jobs:      createJobs(jobs.NewReduce(0.0), testImages),
		},
		{
			name:      "Test Random Filter",
			wantError: false,
			jobs:      createJobs(jobs.NewRandomFilter(3, -1, 1, true), testImages),
		},
		{
			name:      "Test Random Filter Error",
			wantError: true,
			jobs:      createJobs(jobs.NewRandomFilter(0, -1, 1, false), testImages),
		},
		{
			name:      "Test Add Text",
			wantError: false,
			jobs:      createJobs(jobs.NewAddText("text", 1.0, 0.5, 0.5), testImages),
		},
		{
			name:      "Test Add Text Error",
			wantError: false,
			jobs:      createJobs(jobs.NewAddText("text", 1.0, 0.5, 0.5), testImages),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, job := range tt.jobs {
				result, err := job.Process()

				assert.Equal(t, tt.wantError, err != nil)
				if result != nil && job.GetTimeElapsed() == 0 {
					t.Errorf("Test: %s, expected job to take time to complete, StartTime: %d, EndTime: %d ", tt.name, job.GetStartTime().UnixMilli(), job.GetEndTime().UnixMilli())
				}

			}
		})

	}
}
