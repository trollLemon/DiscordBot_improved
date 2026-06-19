package jobs

import "fmt"

// NewInvert creates a new Invert operation.
func NewInvert() Operation {
	return &Invert{}
}

// NewSaturate creates a new Saturate operation with the specified value.
func NewSaturate(value float32) (Operation, error) {

	if value <= 0.0 {
		return nil, fmt.Errorf("expected saturation value to be greater than 0, got %0.1f", value)
	}
	return &Saturate{Value: value}, nil
}

// NewEdgeDetection creates a new EdgeDetect operation with the specified lower and higher thresholds for Canny edge detection.
func NewEdgeDetection(tLower, tHigher float32) (Operation, error) {
	if tLower < 0 || tHigher < 0 {
		return nil, fmt.Errorf("expected t_lower and t_higher to be greater than or equal to 0, got %0.1f and %0.1f", tLower, tHigher)
	}
	return &EdgeDetect{TLower: tLower, THigher: tHigher}, nil

}

// NewMorphology creates a new Morphology operation with the specified kernel size, iterations, and choice of dilation or erosion.
func NewMorphology(kernelSize, iterations int, op Choice) (Operation, error) {
	if kernelSize <= 0 || iterations <= 0 {
		return nil, fmt.Errorf("expected kernel size and iterations to be greater than 0, got %d and %d", kernelSize, iterations)
	}

	if op != Dilate && op != Erode {
		return nil, fmt.Errorf("invalid morphology operation: %s, expected 'dilate' or 'erode'", op)
	}

	return &Morphology{KernelSize: kernelSize, Iterations: iterations, Op: op}, nil
}

// NewReduce creates a new Reduce operation with the specified quality factor.
func NewReduce(quality float32) (Operation, error) {
	if quality <= 0.0 {
		return nil, fmt.Errorf("expected quality to be greater than 0, got %0.1f", quality)
	}
	return &Reduce{Quality: quality}, nil
}

// NewAddText creates a new AddText operation with the specified text, font scale, and position percentages.
func NewAddText(text string, fontScale, xPercentage, yPercentage float64) (Operation, error) {
	if text == "" {
		return nil, fmt.Errorf("must be given a non-empty string")

	}

	if xPercentage < 0.0 || yPercentage < 0.0 || xPercentage > 1.0 || yPercentage > 1.0 {
		return nil, fmt.Errorf("expected x and y percentages to be between 0 and 1, got %0.1f and %0.1f", xPercentage, yPercentage)
	}

	if fontScale <= 0.0 {
		return nil, fmt.Errorf("expected font scale to be greater than 0, got %0.1f", fontScale)
	}
	return &AddText{Text: text, FontScale: fontScale, X: xPercentage, Y: yPercentage}, nil

}

// NewRandomFilter creates a new RandomFilter operation with the specified kernel size, min and max values.
func NewRandomFilter(kernelSize, min, max int, normalize bool) (Operation, error) {
	if kernelSize <= 0 {
		return nil, fmt.Errorf("expected kernel size to be greater than 0, got %d", kernelSize)
	}
	return &RandomFilter{KernelSize: kernelSize, Min: min, Max: max, Normalize: normalize}, nil

}

// NewShuffle creates a new Shuffle operation with the specified number of partitions.
// imgRows and imgCols are used to validate that the number of partitions can fit within the image dimensions.
func NewShuffle(partitions, imgRows, imgCols int) (Operation, error) {
	if partitions <= 1 {
		return nil, fmt.Errorf("expected partitions to be greater than 1, got %d", partitions)
	}

	if partitions >= imgRows*imgCols {
		return nil, fmt.Errorf("cannot fit %d partitions in a %d by %d image", partitions, imgRows, imgCols)
	}
	return &Shuffle{Partitions: partitions}, nil
}
