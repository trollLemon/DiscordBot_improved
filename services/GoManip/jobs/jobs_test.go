package jobs

import (
	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
	"testing"
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
		wantId  uint32
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

			assert.Equal(t, tt.wantErr, err != nil)

			if !tt.job.endTime.After(tt.job.startTime) {
				t.Error("endTime must be after startTime")
			}

			elapsed := tt.job.GetTimeElapsed()
			assert.True(t, tt.job.GetStartTime().UnixMilli() != 0)
			assert.True(t, tt.job.GetEndTime().UnixMilli() != 0)
			assert.True(t, elapsed > 0)
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
