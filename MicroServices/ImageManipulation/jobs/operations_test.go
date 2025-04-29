package jobs_test

import (
	"goManip/jobs"
	"gocv.io/x/gocv"
	"testing"
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
		{width: 128, height: 128},
		{width: 359, height: 474},
		{width: 512, height: 512},
		{width: 1280, height: 720},
		{width: 1920, height: 1080},
		{width: 3440, height: 1440},
	}

	for _, size := range sizes {
		image := gocv.NewMatWithSize(size.width, size.height, gocv.MatTypeCV8UC3)
		rng.Fill(&image, gocv.RNGDistNormal, 1.0, 0.0, false)
		images = append(images, &image)
	}

	return images
}

var (
	testImages = generateTestImages()
)

func TestInvert(t *testing.T) {

	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        jobs.Invert
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Invert{},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.Invert{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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

	nonRgbTestImage := gocv.NewMatWithSize(64, 64, gocv.MatTypeCV8UC3)
	defer nonRgbTestImage.Close()
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        jobs.Saturate
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Saturate{Value: 0.1},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Saturate{Value: 0.5},
		},
		{
			name:      "Handle value = 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        jobs.Saturate{Value: 0.0},
		},
		{
			name:      "Handle value < 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        jobs.Saturate{Value: -0.1},
		},

		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.Saturate{Value: 0.1},
		},
		{
			name:      "Handle non RGB image case",
			wantError: false,
			images:    []*gocv.Mat{&nonRgbTestImage},
			op:        jobs.Saturate{Value: 0.1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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
		images    []*gocv.Mat
		op        jobs.EdgeDetect
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.EdgeDetect{TLower: 100, THigher: 200},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.EdgeDetect{TLower: 10, THigher: 20},
		},
		{
			name:      "Handle tLower < 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        jobs.EdgeDetect{TLower: -1, THigher: 20},
		},
		{
			name:      "Handle tHigher < 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        jobs.EdgeDetect{TLower: 10, THigher: -20},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.EdgeDetect{TLower: 10, THigher: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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
		images    []*gocv.Mat
		op        jobs.Morphology
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Morphology{KernelSize: 3, Iterations: 3, Op: jobs.Dilate},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Morphology{KernelSize: 3, Iterations: 3, Op: jobs.Erode},
		},
		{
			name:      "Handle invalid kernelSize (less than or equal to 0)",
			wantError: true,
			images:    testImages,
			op:        jobs.Morphology{KernelSize: 0, Iterations: 3, Op: jobs.Dilate},
		},
		{
			name:      "Handle invalid iterations (less than or equal to 0)",
			wantError: true,
			images:    testImages,
			op:        jobs.Morphology{KernelSize: 3, Iterations: 0, Op: jobs.Dilate},
		},

		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.Morphology{KernelSize: 3, Iterations: 3, Op: jobs.Erode},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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
		images    []*gocv.Mat
		op        jobs.Reduce
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Reduce{Quality: 0.8},
		},
		{
			name:      "Handle invalid quality (<=0)",
			wantError: true,
			images:    testImages,
			op:        jobs.Reduce{Quality: 0.0},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.Reduce{Quality: 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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
		images    []*gocv.Mat
		op        jobs.AddText
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.AddText{Text: "text", FontScale: 1.0, X: 0.5, Y: 0.5},
		},
		{
			name:      "Handle invalid text (empty)",
			wantError: true,
			images:    testImages,
			op:        jobs.AddText{Text: "", FontScale: 1.0, X: 0.5, Y: 0.5},
		},
		{
			name:      "Handle invalid fontScale (less than or equal to 0.0)",
			wantError: true,
			images:    testImages,
			op:        jobs.AddText{Text: "text", FontScale: -1.0, X: 0.5, Y: 0.5},
		},
		{
			name:      "Handle invalid fontScale (less than or equal to 0.0)",
			wantError: true,
			images:    testImages,
			op:        jobs.AddText{Text: "text", FontScale: 0.0, X: 0.5, Y: 0.5},
		},
		{
			name:      "Handle invalid xy scale",
			wantError: true,
			images:    testImages,
			op:        jobs.AddText{Text: "text", FontScale: 1.0, X: 0.2, Y: 1.5},
		},
		{
			name:      "Handle invalid xy",
			wantError: true,
			images:    testImages,
			op:        jobs.AddText{Text: "text", FontScale: 1.0, X: -0.5, Y: -0.5},
		},
		{
			name:      "Empty Case",
			wantError: true,
			images:    testImages,
			op:        jobs.AddText{Text: "", FontScale: 1.0, X: 0.5, Y: 0.5},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.AddText{Text: "text", FontScale: 1.0, X: 0.5, Y: 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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

	nonRgbTestImage := gocv.NewMatWithSize(64, 64, gocv.MatTypeCV8UC3)
	defer nonRgbTestImage.Close()

	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        jobs.RandomFilter
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 3, Min: -1, Max: 1, Normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 5, Min: -2, Max: 2, Normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 7, Min: -3, Max: 3, Normalize: false},
		},

		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 3, Min: -1, Max: 1, Normalize: true},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 5, Min: -2, Max: 2, Normalize: true},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 7, Min: -3, Max: 3, Normalize: true},
		},

		{
			name:      "Handle invalid kernel size",
			wantError: true,
			images:    testImages,
			op:        jobs.RandomFilter{KernelSize: 0, Min: -1, Max: 1, Normalize: false},
		},

		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.RandomFilter{KernelSize: 1, Min: -1, Max: 1, Normalize: false},
		},
		{
			name:      "Handle non RGB image case",
			wantError: false,
			images:    []*gocv.Mat{&nonRgbTestImage},
			op:        jobs.RandomFilter{KernelSize: 1, Min: -1, Max: 1, Normalize: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

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
		images    []*gocv.Mat
		op        jobs.Shuffle
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Shuffle{Partitions: 10},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        jobs.Shuffle{Partitions: 100},
		},
		{
			name:      "Handle partitions greater than image size, or a ridiculously large parition size",
			wantError: true,
			images:    testImages,
			op:        jobs.Shuffle{Partitions: 9999999999999},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        jobs.Shuffle{Partitions: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, image := range tt.images {

				_, err := tt.op.Run(image)

				if tt.wantError && err == nil {
					t.Errorf("Test: %s, expected error but got nil", tt.name)
				} else if !tt.wantError && err != nil {
					t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
				}
			}

		})
	}

}
