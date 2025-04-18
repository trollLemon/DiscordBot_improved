package jobs

import (
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
		op        Invert
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Invert{},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        Invert{},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestSaturate(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        Saturate
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Saturate{value: 0.1},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Saturate{value: 0.5},
		},
		{
			name:      "Handle value = 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        Saturate{value: 0.0},
		},
		{
			name:      "Handle value < 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        Saturate{value: -0.1},
		},

		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        Saturate{value: 0.1},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestEdgeDetect(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        EdgeDetect
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        EdgeDetect{t_lower: 100, t_higher: 200},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        EdgeDetect{t_lower: 10, t_higher: 20},
		},
		{
			name:      "Handle t_lower < 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        EdgeDetect{t_lower: -1, t_higher: 20},
		},
		{
			name:      "Handle t_higher < 0.0 (Should err)",
			wantError: true,
			images:    testImages,
			op:        EdgeDetect{t_lower: 10, t_higher: -20},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        EdgeDetect{t_lower: 10, t_higher: 20},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestMorphology(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        Morphology
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Morphology{kernel_size: 3, iterations: 3, op: Dilate},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Morphology{kernel_size: 3, iterations: 3, op: Erode},
		},
		{
			name:      "Handle invalid kernel_size (less than or equal to 0)",
			wantError: true,
			images:    testImages,
			op:        Morphology{kernel_size: 0, iterations: 3, op: Dilate},
		},
		{
			name:      "Handle invalid iterations (less than or equal to 0)",
			wantError: true,
			images:    testImages,
			op:        Morphology{kernel_size: 3, iterations: 0, op: Dilate},
		},

		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        Morphology{kernel_size: 3, iterations: 3, op: Erode},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestReduce(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        Reduce
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Reduce{quality: 0.8},
		},
		{
			name:      "Handle invalid quality (<=0)",
			wantError: true,
			images:    testImages,
			op:        Reduce{quality: 0.0},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        Reduce{quality: 0.5},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestAddText(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        AddText
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        AddText{text: "text", font_scale: 1.0, x: 0.5, y: 0.5},
		},
		{
			name:      "Handle invalid text (empty)",
			wantError: true,
			images:    testImages,
			op:        AddText{text: "", font_scale: 1.0, x: 0.5, y: 0.5},
		},
		{
			name:      "Handle invalid font_scale (less than or equal to 0.0)",
			wantError: true,
			images:    testImages,
			op:        AddText{text: "text", font_scale: -1.0, x: 0.5, y: 0.5},
		},
		{
			name:      "Handle invalid font_scale (less than or equal to 0.0)",
			wantError: true,
			images:    testImages,
			op:        AddText{text: "text", font_scale: 0.0, x: 0.5, y: 0.5},
		},
		{
			name:      "Handle invalid xy scale",
			wantError: true,
			images:    testImages,
			op:        AddText{text: "text", font_scale: 1.0, x: 0.2, y: 1.5},
		},
		{
			name:      "Handle invalid xy",
			wantError: true,
			images:    testImages,
			op:        AddText{text: "text", font_scale: 1.0, x: -0.5, y: -0.5},
		},
		{
			name:      "Empty Case",
			wantError: true,
			images:    testImages,
			op:        AddText{text: "", font_scale: 1.0, x: 0.5, y: 0.5},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        AddText{text: "text", font_scale: 1.0, x: 0.5, y: 0.5},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestRandomFilter(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        RandomFilter
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        RandomFilter{kernel_size: 3, min: -1, max: 1, normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        RandomFilter{kernel_size: 5, min: -2, max: 2, normalize: false},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        RandomFilter{kernel_size: 7, min: -3, max: 3, normalize: false},
		},

		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        RandomFilter{kernel_size: 3, min: -1, max: 1, normalize: true},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        RandomFilter{kernel_size: 5, min: -2, max: 2, normalize: true},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        RandomFilter{kernel_size: 7, min: -3, max: 3, normalize: true},
		},

		{
			name:      "Handle invalid kernel size",
			wantError: true,
			images:    testImages,
			op:        RandomFilter{kernel_size: 0, min: -1, max: 1, normalize: false},
		},

		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        RandomFilter{kernel_size: 1, min: -1, max: 1, normalize: false},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		images    []*gocv.Mat
		op        Shuffle
	}{
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Shuffle{partitions: 10},
		},
		{
			name:      "test with various image sizes",
			wantError: false,
			images:    testImages,
			op:        Shuffle{partitions: 100},
		},
		{
			name:      "Handle partitions greater than image size, or a ridiculously large parition size",
			wantError: true,
			images:    testImages,
			op:        Shuffle{partitions: 9999999999999},
		},
		{
			name:      "Handle Nil image case",
			wantError: true,
			images:    []*gocv.Mat{nil},
			op:        Shuffle{partitions: 10},
		},
	}

	for _, tt := range tests {

		for _, image := range tt.images {

			_, err := tt.op.Run(image)

			if tt.wantError && err == nil {
				t.Errorf("Test: %s, expected error but got nil", tt.name)
			} else if !tt.wantError && err != nil {
				t.Errorf("Test: %s, error = %v, wantErr %v", tt.name, err.Error(), tt.wantError)
			}
		}

	}

}
