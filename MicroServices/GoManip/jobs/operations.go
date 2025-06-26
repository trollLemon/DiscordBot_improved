package jobs

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"time"
)

/*  Helper Functions  */

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	Value float32
}

func (s *Saturate) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if s.Value <= 0.0 {
		return nil, fmt.Errorf("expected saturation value to be greater than 0, got %f", s.Value)
	}

	hsvImage := gocv.NewMat()

	expectedChannels := 3

	if input.Channels() != expectedChannels {
		converted := gocv.NewMat()
		defer converted.Close()

		err := gocv.CvtColor(*input, &converted, gocv.ColorGrayToBGR)
		if err != nil {
			return nil, err
		}

		*input = converted.Clone()
	}

	err := gocv.CvtColor(*input, &hsvImage, gocv.ColorBGRToHLSFull)

	if err != nil {
		return nil, fmt.Errorf("failed to convert to HSV: %v", err)
	}

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

	sat.ConvertToWithParams(&sat, gocv.MatType(sameType), s.Value, float32(beta))

	gocv.Merge([]gocv.Mat{hue, light, sat}, &saturated)

	imgSaturated := gocv.NewMat()

	gocv.CvtColor(saturated, &imgSaturated, gocv.ColorHLSToBGR)

	return &imgSaturated, nil

}

type EdgeDetect struct {
	TLower  float32
	THigher float32
}

func (e *EdgeDetect) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if e.TLower < 0 || e.THigher < 0 {
		return nil, fmt.Errorf("expected t_lower and t_higher to be greater than or equal to 0, got %0.2f and %0.2f", e.TLower, e.THigher)
	}

	edges := gocv.NewMat()

	gocv.Canny(*input, &edges, e.TLower, e.THigher)

	return &edges, nil

}

type Choice string

const (
	Dilate Choice = "Dilate"
	Erode         = "Erode"
)

type Morphology struct {
	KernelSize int
	Iterations int
	Op         Choice
}

func (m *Morphology) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if m.KernelSize <= 0 || m.Iterations <= 0 {
		return nil, fmt.Errorf("expected kernel size and iterations to be greater than 0, got %d and %d", m.KernelSize, m.Iterations)
	}

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: m.KernelSize, Y: m.KernelSize})
	morphedImage := gocv.NewMat()

	switch m.Op {
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
	Quality float32
}

func (r *Reduce) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if r.Quality <= 0.0 {
		return nil, fmt.Errorf("expected quality to be greater than 0.0, got %0.2f", r.Quality)
	}

	resizedImage := gocv.NewMat()
	reducedImage := gocv.NewMat()

	defer resizedImage.Close()

	gocv.Resize(*input, &resizedImage, image.Point{}, float64(r.Quality), float64(r.Quality), gocv.InterpolationNearestNeighbor)

	gocv.Resize(resizedImage, &reducedImage, image.Point{X: input.Rows(), Y: input.Cols()}, 0.0, 0.0, gocv.InterpolationNearestNeighbor)

	return &reducedImage, nil

}

type AddText struct {
	Text      string
	FontScale float64
	X         float64
	Y         float64
}

func (a *AddText) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("input image is empty")

	}

	if a.Text == "" {
		return nil, errors.New("must be given a non-empty string")
	}

	if a.X < 0.0 || a.Y < 0.0 || a.X > 1.0 || a.Y > 1.0 {
		return nil, fmt.Errorf("expected x and y percentages to be greater than between 0 and 1, got %0.2f. %0.2f", a.X, a.Y)
	}

	if a.FontScale <= 0.0 {
		return nil, fmt.Errorf("expected font scale to be greater than 0, got %0.2f", a.FontScale)
	}

	rows, cols := input.Rows(), input.Cols()

	xPos, yPos := int(float64(rows)*a.X), int(float64(cols)*a.Y)

	thickness := 1
	lineType := gocv.LineAA

	gocv.PutTextWithParams(input, a.Text, image.Point{X: xPos, Y: yPos}, gocv.FontHersheyPlain, a.FontScale, color.RGBA{255, 255, 255, 255}, thickness, lineType, false)

	return input, nil
}

type RandomFilter struct {
	KernelSize int
	Min        int
	Max        int
	Normalize  bool
}

func (r *RandomFilter) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("input image is empty")

	}

	if r.KernelSize <= 0 {

		return nil, fmt.Errorf("expected kernel size to be greater than 0, got %d", r.KernelSize)

	}

	kernels := make([]gocv.Mat, input.Channels())

	for i := 0; i < input.Channels(); i++ {
		kernels[i] = gocv.NewMatWithSize(r.KernelSize, r.KernelSize, gocv.MatTypeCV32F)
	}

	gocv.SetRNGSeed(int(time.Now().UnixNano()))
	rng := gocv.TheRNG()

	for i := range kernels {
		rng.Fill(&kernels[i], gocv.RNGDistUniform, float64(r.Min), float64(r.Max), false)

		if r.Normalize {
			gocv.Normalize(kernels[i], &kernels[i], 1, 0, gocv.NormL2)
		}
	}

	ddepth := -1

	channels := gocv.Split(*input)

	filteredChannels := make([]gocv.Mat, input.Channels())

	for i := 0; i < input.Channels(); i++ {
		filteredChannels[i] = gocv.NewMat()
	}

	// convolve the 3D filter over the RBG image
	for idx, kernel := range kernels {

		gocv.Filter2D(channels[idx], &filteredChannels[idx], gocv.MatType(ddepth), kernel, image.Point{-1, -1}, 0, gocv.BorderDefault)
		kernel.Close()
		channels[idx].Close()
	}

	filteredImage := gocv.NewMat()

	gocv.Merge(filteredChannels, &filteredImage)

	for _, imgChan := range filteredChannels {

		imgChan.Close()
	}

	return &filteredImage, nil
}

type Shuffle struct {
	Partitions int
}

func (s *Shuffle) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("input image is empty")

	}
	if s.Partitions <= 1 {

		return nil, fmt.Errorf("expected partitions to be greater than 1, got %d", s.Partitions)

	}

	if s.Partitions >= input.Rows()*input.Cols() {
		return nil, fmt.Errorf("cannot fit %d partitions in a %d by %d image", s.Partitions, input.Rows(), input.Cols())
	}

	rows := input.Rows()
	cols := input.Cols()
	dataType := input.Type()

	partRowsFlr := math.Floor(math.Sqrt(float64(s.Partitions)))
	partColsFlr := math.Floor(float64(s.Partitions) / partRowsFlr)

	partRows := int(partRowsFlr)
	partCols := int(partColsFlr)

	sliceWidth := cols / partCols
	sliceHeight := rows / partRows

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
