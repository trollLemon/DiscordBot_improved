package jobs

func NewInvert() Operation {
	return &Invert{}
}

func NewSaturate(value float32) Operation {
	return &Saturate{value: value}
}

func NewEdgeDetection(tLower, tHigher float32) Operation {

	return &EdgeDetect{tLower: tLower, tHigher: tHigher}

}

func NewMorphology(kernelSize, iterations int, op Choice) Operation {
	return &Morphology{kernelSize: kernelSize, iterations: iterations, op: op}
}

func NewReduce(quality float32) Operation {
	return &Reduce{quality: quality}
}

func NewAddText(text string, fontScale, xPercentage, yPercentage float64) Operation {

	return &AddText{text: text, fontScale: fontScale, x: xPercentage, y: yPercentage}

}

func NewRandomFilter(kernelSize, min, max int, normalize bool) Operation {

	return &RandomFilter{kernelSize: kernelSize, min: min, max: max, normalize: normalize}

}

func NewShuffle(partitions int) Operation {

	return &Shuffle{partitions: partitions}
}
