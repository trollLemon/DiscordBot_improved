package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"

	"goManip/jobs"
)

type mockOperation struct{}

func (m *mockOperation) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return input, nil
}

func TestJobMock(t *testing.T) {

	mockImage := gocv.NewMatWithSize(64, 64, gocv.MatTypeCV16UC3)

	tests := []struct {
		name   string
		job    *jobs.Job
		wantId uint32
	}{
		{
			name:   "TestNewJob",
			job:    jobs.NewJob(1, &mockOperation{}, &mockImage),
			wantId: 1,
		},
		{
			name:   "TestProcess",
			job:    jobs.NewJob(2, &mockOperation{}, &mockImage),
			wantId: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.job.GetJobId(); got != tt.wantId {
				t.Errorf("GetJobId() = %v, want %v", got, tt.wantId)
			}

			_, err := tt.job.Process()

			assert.Nil(t, err)

			if !tt.job.GetEndTime().After(tt.job.GetStartTime()) {
				t.Error("endTime must be after startTime")
			}

			elapsed := tt.job.GetTimeElapsed()
			assert.True(t, tt.job.GetStartTime().UnixMilli() != 0)
			assert.True(t, tt.job.GetEndTime().UnixMilli() != 0)
			assert.True(t, elapsed > 0)
		})
	}

	// Test case where Process hasn't been called
	job := jobs.NewJob(3, &mockOperation{}, &mockImage)
	if job.GetEndTime().IsZero() {
		t.Run("TestJobNotRun", func(t *testing.T) {
			if got := job.GetTimeElapsed(); got != 0 {
				t.Errorf("GetTimeElapsed returned %d when Process wasn't called", got)
			}
		})
	}

}

func createJobs(operation jobs.Operation, images []*gocv.Mat) []*jobs.Job {

	var jobArr []*jobs.Job

	for i, image := range images {
		jobArr = append(jobArr, jobs.NewJob(uint32(i), operation, image))
	}
	return jobArr

}

func TestJob(t *testing.T) {
	testImages := generateTestImages() // from operations_test.go

	tests := []struct {
		name      string
		wantError string
		job       []*jobs.Job
	}{
		{
			name: "Test Invert Job",
			job:  createJobs(jobs.NewInvert(), testImages),
		},
		{
			name: "Test Shuffle Job",
			job:  createJobs(jobs.NewShuffle(15), testImages),
		},
		{
			name:      "Test Shuffle Job Error",
			job:       createJobs(jobs.NewShuffle(0), testImages),
			wantError: "expected partitions to be greater than 1, got 0",
		},
		{
			name: "Test Edge Detection",
			job:  createJobs(jobs.NewEdgeDetection(100.0, 200.0), testImages),
		},
		{
			name:      "Test Edge Detection Error",
			job:       createJobs(jobs.NewEdgeDetection(-1.0, 200.0), testImages),
			wantError: "expected t_lower and t_higher to be greater than or equal to 0, got -1.0 and 200.0",
		},
		{
			name: "Test Saturation",
			job:  createJobs(jobs.NewSaturate(1.6), testImages),
		},
		{
			name:      "Test Saturation Error",
			job:       createJobs(jobs.NewSaturate(-1.6), testImages),
			wantError: "expected saturation value to be greater than 0, got -1.6",
		},
		{
			name: "Test Dilate",
			job:  createJobs(jobs.NewMorphology(3, 5, jobs.Dilate), testImages),
		},
		{
			name:      "Test Dilate Error",
			job:       createJobs(jobs.NewMorphology(0, 5, jobs.Dilate), testImages),
			wantError: "expected kernel size and iterations to be greater than 0, got 0 and 5",
		},
		{
			name: "Test Erode",
			job:  createJobs(jobs.NewMorphology(3, 5, jobs.Erode), testImages),
		},
		{
			name:      "Test Erode Error",
			job:       createJobs(jobs.NewMorphology(0, 5, jobs.Erode), testImages),
			wantError: "expected kernel size and iterations to be greater than 0, got 0 and 5",
		},
		{
			name: "Test Reduce",
			job:  createJobs(jobs.NewReduce(0.5), testImages),
		},
		{
			name:      "Test Reduce Error",
			job:       createJobs(jobs.NewReduce(0.0), testImages),
			wantError: "expected quality to be greater than 0.0, got 0.0",
		},
		{
			name: "Test Random Filter",
			job:  createJobs(jobs.NewRandomFilter(3, -1, 1, true), testImages),
		},
		{
			name:      "Test Random Filter Error",
			job:       createJobs(jobs.NewRandomFilter(0, -1, 1, false), testImages),
			wantError: "expected kernel size to be greater than 0, got 0",
		},
		{
			name: "Test Add Text",
			job:  createJobs(jobs.NewAddText("text", 1.0, 0.5, 0.5), testImages),
		},
		{
			name:      "Test Add Text Error",
			job:       createJobs(jobs.NewAddText("text", -1.0, 0.5, 0.5), testImages),
			wantError: "expected font scale to be greater than 0, got -1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, job := range tt.job {
				result, err := job.Process()

				if tt.wantError != "" {
					assert.ErrorContains(t, err, tt.wantError)
				} else {
					assert.Nil(t, err)
				}

				if result != nil && job.GetTimeElapsed() == 0 {
					t.Errorf("Test: %s, expected job to take time to complete, StartTime: %d, EndTime: %d ", tt.name, job.GetStartTime().UnixMilli(), job.GetEndTime().UnixMilli())
				}

			}
		})

	}
}
