package jobs

import (
	"testing"

	"gocv.io/x/gocv"
)

func createJobs(operation Operation, images []*gocv.Mat) []*Job {

	var jobs []*Job

	for i, image := range images {
		jobs = append(jobs, NewJob(uint8(i), operation, image))
	}
	return jobs

}

func TestIntegration(t *testing.T) {

	testImages := generateTestImages() // from operations_test.go

	tests := []struct {
		name      string
		wantError bool
		jobs      []*Job
	}{
		{
			name:      "Test Invert Job",
			wantError: false,
			jobs:      createJobs(NewInvert(), testImages),
		},
		{
			name:      "Test Shuffle Job",
			wantError: false,
			jobs:      createJobs(NewShuffle(15), testImages),
		},
		{
			name:      "Test Shuffle Job Error",
			wantError: true,
			jobs:      createJobs(NewShuffle(0), testImages),
		},
		{
			name:      "Test Edge Detection",
			wantError: false,
			jobs:      createJobs(NewEdgeDetection(100.0, 200.0), testImages),
		},
		{
			name:      "Test Edge Detection Error",
			wantError: true,
			jobs:      createJobs(NewEdgeDetection(-1.0, 200.0), testImages),
		},
		{
			name:      "Test Saturation",
			wantError: false,
			jobs:      createJobs(NewSaturate(1.6), testImages),
		},
		{
			name:      "Test Saturation Error",
			wantError: true,
			jobs:      createJobs(NewSaturate(-1.6), testImages),
		},
		{
			name:      "Test Dilate",
			wantError: false,
			jobs:      createJobs(NewMorphology(3, 5, Dilate), testImages),
		},
		{
			name:      "Test Dilate Error",
			wantError: true,
			jobs:      createJobs(NewMorphology(0, 5, Dilate), testImages),
		},
		{
			name:      "Test Erode",
			wantError: false,
			jobs:      createJobs(NewMorphology(3, 5, Erode), testImages),
		},
		{
			name:      "Test Erode Error",
			wantError: true,
			jobs:      createJobs(NewMorphology(0, 5, Erode), testImages),
		},
		{
			name:      "Test Reduce",
			wantError: false,
			jobs:      createJobs(NewReduce(0.5), testImages),
		},
		{
			name:      "Test Reduce Error",
			wantError: true,
			jobs:      createJobs(NewReduce(0.0), testImages),
		},
		{
			name:      "Test Random Filter",
			wantError: false,
			jobs:      createJobs(NewRandomFilter(3, -1, 1, false), testImages),
		},
		{
			name:      "Test Random Filter Error",
			wantError: true,
			jobs:      createJobs(NewRandomFilter(0, -1, 1, false), testImages),
		},
		{
			name:      "Test Add Text",
			wantError: false,
			jobs:      createJobs(NewAddText("text", 1.0, 0.5, 0.5), testImages),
		},
		{
			name:      "Test Add Text Error",
			wantError: false,
			jobs:      createJobs(NewAddText("text", 1.0, 0.5, 0.5), testImages),
		},
	}

	for _, tt := range tests {

		for _, job := range tt.jobs {
			result, err := job.Process()

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}

			if result != nil && job.GetTimeElapsed() == 0 {
				t.Errorf("Test: %s, expected job to take time to complete, StartTime: %d, EndTime: %d ", tt.name, job.GetStartTime(), job.GetEndTime())
			}

		}

	}
}
