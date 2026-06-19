package jobs_test

import (
	"goManip/jobs"
	"testing"

	"gocv.io/x/gocv"
)

func generateTestImages() []*gocv.Mat {

	rng := gocv.TheRNG()

	var images []*gocv.Mat
	sizes := []struct {
		width, height int
	}{
		{width: 64, height: 64},
		{width: 64, height: 32},
		{width: 256, height: 256},
	}

	for _, size := range sizes {
		image := gocv.NewMatWithSize(size.width, size.height, gocv.MatTypeCV8UC3)
		rng.Fill(&image, gocv.RNGDistNormal, 1.0, 0.0, false)
		images = append(images, &image)
	}

	return images
}

func closeTestImages(images []*gocv.Mat) {
	for _, img := range images {
		img.Close()
	}
}

func TestInvert(t *testing.T) {

	tests := []struct {
		name      string
		wantError bool
		op        jobs.Invert
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Invert{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}

			}

		})

	}

}

func TestSaturate(t *testing.T) {

	tests := []struct {
		name      string
		wantError bool
		op        jobs.Saturate
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Saturate{Value: 0.1},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Saturate{Value: 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}

func TestEdgeDetect(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		op        jobs.EdgeDetect
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.EdgeDetect{TLower: 100, THigher: 200},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.EdgeDetect{TLower: 10, THigher: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}

func TestMorphology(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		op        jobs.Morphology
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Morphology{KernelSize: 3, Iterations: 3, Op: jobs.Dilate},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Morphology{KernelSize: 3, Iterations: 3, Op: jobs.Erode},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}

func TestReduce(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		op        jobs.Reduce
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Reduce{Quality: 0.8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}

func TestAddText(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		op        jobs.AddText
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.AddText{Text: "text", FontScale: 1.0, X: 0.5, Y: 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}
		})

	}

}

func TestRandomFilter(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		op        jobs.RandomFilter
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.RandomFilter{KernelSize: 3, Min: -1, Max: 1, Normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.RandomFilter{KernelSize: 5, Min: -2, Max: 2, Normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.RandomFilter{KernelSize: 7, Min: -3, Max: 3, Normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.RandomFilter{KernelSize: 3, Min: -1, Max: 1, Normalize: true},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.RandomFilter{KernelSize: 5, Min: -2, Max: 2, Normalize: true},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.RandomFilter{KernelSize: 7, Min: -3, Max: 3, Normalize: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		op        jobs.Shuffle
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			op:        jobs.Shuffle{Partitions: 10},
		},
		{
			name:      "test with various image sizes large partition",
			wantError: false,
			op:        jobs.Shuffle{Partitions: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			images := generateTestImages()
			defer closeTestImages(images)
			for _, image := range images {

				result, err := tt.op.Run(image)
				if result != nil {
					result.Close()
				}

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}
