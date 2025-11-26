package jobs

func NewInvert() Operation {
	return &Invert{}
}

func NewSaturate(value float32) Operation {
	return &Saturate{Value: value}
}

func NewEdgeDetection(tLower, tHigher float32) Operation {

	return &EdgeDetect{TLower: tLower, THigher: tHigher}

}

func NewMorphology(kernelSize, iterations int, op Choice) Operation {
	return &Morphology{KernelSize: kernelSize, Iterations: iterations, Op: op}
}

func NewReduce(quality float32) Operation {
	return &Reduce{Quality: quality}
}

func NewAddText(text string, fontScale, xPercentage, yPercentage float64) Operation {

	return &AddText{Text: text, FontScale: fontScale, X: xPercentage, Y: yPercentage}

}

func NewRandomFilter(kernelSize, min, max int, normalize bool) Operation {

	return &RandomFilter{KernelSize: kernelSize, Min: min, Max: max, Normalize: normalize}

}

func NewShuffle(partitions int) Operation {

	return &Shuffle{Partitions: partitions}
}
