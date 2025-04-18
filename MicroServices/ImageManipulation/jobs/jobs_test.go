package jobs

import (
	"testing"

	"gocv.io/x/gocv"
)

type mockOperation struct{}

func (m *mockOperation) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return input, nil
}

func TestJob(t *testing.T) {

	mockImage := gocv.NewMatWithSize(64, 64, gocv.MatTypeCV16UC3)

	tests := []struct {
		name    string
		job     *Job
		wantId  uint8
		wantErr bool
	}{
		{
			name:    "TestNewJob",
			job:     NewJob(1, &mockOperation{}, &mockImage),
			wantId:  1,
			wantErr: false,
		},
		{
			name:    "TestProcess",
			job:     NewJob(2, &mockOperation{}, &mockImage),
			wantId:  2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.job.GetJobId(); got != tt.wantId {
				t.Errorf("GetJobId() = %v, want %v", got, tt.wantId)
			}

			_, err := tt.job.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Process() returned unexpected error: %v", tt.name, err)
			}

			if !tt.job.endTime.After(tt.job.startTime) {
				t.Error("endTime must be after startTime")
			}

			elapsed := tt.job.GetTimeElapsed()

			if elapsed == 0 {
				t.Errorf("GetTimeElapsed reported no time elapsed")
			}
		})
	}

	// Test case where Process hasn't been called
	job := NewJob(3, &mockOperation{}, &mockImage)
	if job.endTime.IsZero() {
		t.Run("TestJobNotRun", func(t *testing.T) {
			if got := job.GetTimeElapsed(); got != 0 {
				t.Errorf("GetTimeElapsed returned %d when Process wasn't called", got)
			}
		})
	}

}
