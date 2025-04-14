package jobs

func NewInvert() Operation {
	return &Invert{}
}

func NewSaturate(value float32) Operation {
	return &Saturate{value: value}
}

func NewEdgeDetection(t_lower, t_higher float32) Operation {

	return &EdgeDetect{t_lower: t_lower, t_higher: t_higher}

}

func NewMorphology(kernel_size, iterations int, op choice) Operation {
	return &Morphology{kernel_size: kernel_size, iterations: iterations, op: op}
}

func NewReduce() Operation {
	return &Reduce{}
}

func NewAddText(text string, fontScale, xPercentage, yPercentage float64) Operation {

	return &AddText{text: text, font_scale: fontScale, x: xPercentage, y: yPercentage}

}

func NewRandomFilter(kernel_size, min, max int, normalize bool) Operation {

	return &RandomFilter{kernel_size: kernel_size, min: min, max: max, normalize: normalize}

}

func NewShuffle(partitions int) Operation {

	return &Shuffle{partitions: partitions}
}
