package jobs

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"sync"

	"gocv.io/x/gocv"
)

/*  Helper Functions  */

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/* Colors */

type Invert struct{}

func (_ *Invert) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	white := gocv.NewMatWithSizeFromScalar(gocv.Scalar{255.0, 255.0, 255.0, 255.0}, input.Rows(), input.Cols(), input.Type())

	inverted := gocv.NewMat()

	gocv.Subtract(white, *input, &inverted)

	white.Close()

	return &inverted, nil
}

type Saturate struct {
	value float32
}

func (s *Saturate) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if s.value <= 0.0 {
		return nil, fmt.Errorf("expected saturation value to be greater than 0, got %f", s.value)
	}

	hsvImage := gocv.NewMat()

	gocv.CvtColor(*input, &hsvImage, gocv.ColorBGRToHLSFull)

	chans := gocv.Split(hsvImage)

	hue := chans[0]
	light := chans[1]
	sat := chans[2]

	saturated := gocv.NewMat()

	defer func() {
		hsvImage.Close()
		hue.Close()
		light.Close()
		sat.Close()
		saturated.Close()
	}()

	sameType := -1
	beta := 0

	sat.ConvertToWithParams(&sat, gocv.MatType(sameType), s.value, float32(beta))

	gocv.Merge([]gocv.Mat{hue, light, sat}, &saturated)

	imgSaturated := gocv.NewMat()

	gocv.CvtColor(saturated, &imgSaturated, gocv.ColorHLSToBGR)

	return &imgSaturated, nil

}

/* Edge Detection */

type EdgeDetect struct {
	tLower  float32
	tHigher float32
}

func (e *EdgeDetect) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if e.tLower < 0 || e.tHigher < 0 {
		return nil, fmt.Errorf("expected t_lower and t_higher to be greater than or equal to 0, got %0.2f and %0.2f", e.tLower, e.tHigher)
	}

	edges := gocv.NewMat()

	gocv.Canny(*input, &edges, e.tLower, e.tHigher)

	return &edges, nil

}

/* Morphology */
type Choice string

const (
	Dilate Choice = "Dilate"
	Erode         = "Erode"
)

type Morphology struct {
	kernelSize int
	iterations int
	op         Choice
}

func (m *Morphology) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if m.kernelSize <= 0 || m.iterations <= 0 {
		return nil, fmt.Errorf("expected kernel size and iterations to be greater than 0, got %d and %d", m.kernelSize, m.iterations)
	}

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: m.kernelSize, Y: m.kernelSize})
	morphedImage := gocv.NewMat()

	switch m.op {
	case Dilate:
		gocv.Dilate(*input, &morphedImage, kernel)
	case Erode:
		gocv.Erode(*input, &morphedImage, kernel)
	default:

		defer morphedImage.Close()
		return nil, errors.New("invalid morphology operation")
	}

	return &morphedImage, nil

}

type Reduce struct {
	quality float32
}

func (r *Reduce) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if r.quality <= 0.0 {
		return nil, fmt.Errorf("expected quality to be greater than 0.0, got %0.2f", r.quality)
	}

	resizedImage := gocv.NewMat()
	reducedImage := gocv.NewMat()

	defer resizedImage.Close()

	gocv.Resize(*input, &resizedImage, image.Point{}, float64(r.quality), float64(r.quality), gocv.InterpolationNearestNeighbor)

	gocv.Resize(resizedImage, &reducedImage, image.Point{X: input.Rows(), Y: input.Cols()}, 0.0, 0.0, gocv.InterpolationNearestNeighbor)

	return &reducedImage, nil

}

/* Text */

type AddText struct {
	text      string
	fontScale float64
	x         float64
	y         float64
}

func (a *AddText) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("input image is empty")

	}

	if a.text == "" {
		return nil, errors.New("must be given a non-empty string")
	}

	if a.x < 0.0 || a.y < 0.0 || a.x > 1.0 || a.y > 1.0 {
		return nil, fmt.Errorf("expected x and y percentages to be greater than between 0 and 1, got %0.2f. %0.2f", a.x, a.y)
	}

	if a.fontScale <= 0.0 {
		return nil, fmt.Errorf("expected font scale to be greater than 0, got %0.2f", a.fontScale)
	}

	rows, cols := input.Rows(), input.Cols()

	xPos, yPos := int(float64(rows)*a.x), int(float64(cols)*a.y)

	thickness := 1
	lineType := gocv.LineAA

	gocv.PutTextWithParams(input, a.text, image.Point{X: xPos, Y: yPos}, gocv.FontHersheyPlain, a.fontScale, color.RGBA{255, 255, 255, 255}, thickness, lineType, false)

	return input, nil
}

/* Misc */

type RandomFilter struct {
	kernelSize int
	min        int
	max        int
	normalize  bool
}

func (r *RandomFilter) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("input image is empty")

	}

	if r.kernelSize <= 0 {

		return nil, fmt.Errorf("expected kernel size to be greater than 0, got %d", r.kernelSize)

	}

	kernels := []gocv.Mat{
		gocv.NewMatWithSize(r.kernelSize, r.kernelSize, gocv.MatTypeCV64F),
		gocv.NewMatWithSize(r.kernelSize, r.kernelSize, gocv.MatTypeCV64F),
		gocv.NewMatWithSize(r.kernelSize, r.kernelSize, gocv.MatTypeCV64F),
	}
	rng := gocv.TheRNG()

	for _, mat := range kernels {
		rng.Fill(&mat, gocv.RNGDistUniform, float64(r.min), float64(r.max), false)
	}

	ddepth := -1

	channels := gocv.Split(*input)

	filteredChannels := []gocv.Mat{
		gocv.NewMat(),
		gocv.NewMat(),
		gocv.NewMat(),
	}

	wg := sync.WaitGroup{}

	// convolve the 3D filter over the RBG image concurrently
	for idx, kernel := range kernels {

		wg.Add(1)
		go func() {
			defer wg.Done()
			gocv.Filter2D(channels[idx], &filteredChannels[idx], gocv.MatType(ddepth), kernel, image.Point{-1, -1}, 0, gocv.BorderDefault)
			kernel.Close()
			channels[idx].Close()
		}()
	}

	wg.Wait()

	filteredImage := gocv.NewMat()

	gocv.Merge(filteredChannels, &filteredImage)

	for _, imgChan := range filteredChannels {

		imgChan.Close()
	}

	return &filteredImage, nil
}

type Shuffle struct {
	partitions int
}

func (s *Shuffle) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("input image is empty")

	}
	if s.partitions <= 1 {

		return nil, fmt.Errorf("expected partitions to be greater than 1, got %d", s.partitions)

	}

	if s.partitions >= input.Rows()*input.Cols() {
		return nil, fmt.Errorf("cannot fit %d partitions in a %d by %d image", s.partitions, input.Rows(), input.Cols())
	}

	rows := input.Rows()
	cols := input.Cols()
	dataType := input.Type()

	partRowsFlr := math.Floor(math.Sqrt(float64(s.partitions)))
	partColsFlr := math.Floor(float64(s.partitions) / partRowsFlr)

	partRows := int(partRowsFlr)
	partCols := int(partColsFlr)

	sliceWidth := cols / int(partCols)
	sliceHeight := rows / int(partRows)

	var slices []gocv.Mat

	for r := range partRows {
		for c := range partCols {

			rowRange := r * sliceHeight
			colRange := c * sliceWidth

			rowStart := rowRange
			rowEnd := min(rowRange+sliceHeight, rows)

			colStart := colRange
			colEnd := min(colRange+sliceWidth, cols)

			roiRect := image.Rect(colStart, rowStart, colEnd, rowEnd)

			imgSlice := input.Region(roiRect)
			slices = append(slices, imgSlice)
		}

	}

	rand.Shuffle(len(slices), func(i, j int) {
		slices[i], slices[j] = slices[j], slices[i]
	})

	newHeight := min(partRows*sliceHeight, rows)
	newWidth := min(partCols*sliceWidth, cols)

	shuffledImage := gocv.NewMatWithSize(newHeight, newWidth, dataType)

	for idx, slice := range slices {

		rowIdx := idx / partCols
		colIdx := idx % partCols

		rowStart := rowIdx * sliceHeight
		colStart := colIdx * sliceWidth

		sliceRows := slice.Rows()
		sliceCols := slice.Cols()
		roiRect := image.Rect(colStart, rowStart, colStart+sliceCols, rowStart+sliceRows)
		roi := shuffledImage.Region(roiRect)
		slice.CopyTo(&roi)
		roi.Close()
	}

	return &shuffledImage, nil
}
