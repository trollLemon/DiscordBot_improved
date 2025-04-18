package jobs

import (
	"fmt"
	"testing"

	"gocv.io/x/gocv"
)

type TestBasicOP struct{}

func (t *TestBasicOP) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return input, nil
}

type TestFailureOP struct{}

func (t *TestFailureOP) Run(input *gocv.Mat) (*gocv.Mat, error) {
	return nil, fmt.Errorf("Operation Failed")
}

func TestJobs(t *testing.T) {
	inputImage := gocv.NewMatWithSize(32, 32, gocv.MatTypeCV8U)
	defer inputImage.Close()

	outputImage := gocv.NewMatWithSize(32, 32, gocv.MatTypeCV8U)
	defer outputImage.Close()

	tests := []struct {
		name          string
		inputImage    gocv.Mat
		expectedImage gocv.Mat
		operation     Operation
		expectError   bool
	}{
		{
			name:          "Basic Operation Test",
			inputImage:    inputImage,
			expectedImage: outputImage,
			operation:     &TestBasicOP{},
			expectError:   false,
		},
		{
			name:          "Failure Operation Test",
			inputImage:    inputImage,
			expectedImage: outputImage,
			operation:     &TestFailureOP{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			job := Job{
				inputImage: &inputImage,
				jobId:      1,
				Operation:  tt.operation,
			}

			img, err := job.Process()

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			} else if !tt.expectError {

				diff := gocv.NewMat()
				gocv.AbsDiff(*img, tt.expectedImage, &diff)

				if gocv.CountNonZero(diff) != 0 {
					t.Errorf("Expected Image differs from actual image")
				}

			}

		})
	}
}
