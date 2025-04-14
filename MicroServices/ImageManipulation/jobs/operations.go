package jobs

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand/v2"

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
		return nil, fmt.Errorf("Expected saturation value to be greater than 0, got %f", s.value)
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
	t_lower  float32
	t_higher float32
}

func (e *EdgeDetect) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if e.t_lower < 0 || e.t_higher < 0 {
		return nil, fmt.Errorf("Expected t_lower and t_higher to be greater than or equal to 0, got %0.2f and %0.2f", e.t_lower, e.t_higher)
	}

	edges := gocv.NewMat()

	gocv.Canny(*input, &edges, e.t_lower, e.t_higher)

	return &edges, nil

}

/* Morphology */
type choice int

const (
	Dilate choice = iota
	Erode
)

type Morphology struct {
	kernel_size int
	iterations  int
	op          choice
}

func (m *Morphology) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {
		return nil, errors.New("input image is empty")
	}

	if m.kernel_size <= 0 || m.iterations <= 0 {
		return nil, fmt.Errorf("Expected kernel size and iterations to be greater than 0, got %d and %d", m.kernel_size, m.iterations)
	}

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: m.kernel_size, Y: m.kernel_size})
	morphedImage := gocv.NewMat()

	switch m.op {
	case Dilate:
		gocv.Dilate(*input, &morphedImage, kernel)
	case Erode:
		gocv.Erode(*input, &morphedImage, kernel)
	default:

		defer morphedImage.Close()
		return nil, errors.New("Invalid morphology operation")
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
		return nil, fmt.Errorf("Expected quality to be greater than 0.0, got %0.2f", r.quality)
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
	text       string
	font_scale float64
	x          float64
	y          float64
}

func (a *AddText) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("Input image is empty")

	}

	if a.text == "" {
		return nil, errors.New("Must be given a non-empty string")
	}

	if a.x < 0.0 || a.y < 0.0 || a.x > 1.0 || a.y > 1.0 {
		return nil, fmt.Errorf("Expected x and y percentages to be greater than between 0 and 1, got %0.2f. %0.2f", a.x, a.y)
	}

	if a.font_scale <= 0.0 {
		return nil, fmt.Errorf("Expected font scale to be greater than 0, got %0.2f", a.font_scale)
	}

	rows, cols := input.Rows(), input.Cols()

	xPos, yPos := int(float64(rows)*a.x), int(float64(cols)*a.y)

	thickness := 1
	lineType := gocv.LineAA

	gocv.PutTextWithParams(input, a.text, image.Point{X: xPos, Y: yPos}, gocv.FontHersheyPlain, a.font_scale, color.RGBA{255, 255, 255, 255}, thickness, lineType, false)

	return input, nil
}

/* Misc */

type RandomFilter struct {
	kernel_size int
	min         int
	max         int
	normalize   bool
}

func (r *RandomFilter) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("Input image is empty")

	}

	if r.kernel_size <= 0 {

		return nil, fmt.Errorf("Expected kernel size to be greater than 0, got %d", r.kernel_size)

	}

	kernels := []gocv.Mat{
		gocv.NewMatWithSize(r.kernel_size, r.kernel_size, gocv.MatTypeCV64F),
		gocv.NewMatWithSize(r.kernel_size, r.kernel_size, gocv.MatTypeCV64F),
		gocv.NewMatWithSize(r.kernel_size, r.kernel_size, gocv.MatTypeCV64F),
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
	partitions int
}

func (s *Shuffle) Run(input *gocv.Mat) (*gocv.Mat, error) {

	if input == nil {

		return nil, errors.New("Input image is empty")

	}
	if s.partitions <= 1 {

		return nil, fmt.Errorf("Expected partitions to be greater than 1, got %d", s.partitions)

	}

	if s.partitions >= input.Rows()*input.Cols() {
		return nil, fmt.Errorf("Cannot fit %d partitions in a %d by %d image", s.partitions, input.Rows(), input.Cols())
	}

	rows := input.Rows()
	cols := input.Cols()
	dataType := input.Type()

	part_rows_flr := math.Floor(math.Sqrt(float64(s.partitions)))
	part_cols_flr := math.Floor(float64(s.partitions) / part_rows_flr)

	part_rows := int(part_rows_flr)
	part_cols := int(part_cols_flr)

	slice_width := cols / int(part_cols)
	slice_height := cols / int(part_rows)

	slices := []gocv.Mat{}

	for r := range part_rows {
		for c := range part_cols {

			row_range := r * slice_height
			col_range := c * slice_width

			//select a partition, but keep it within the bounds of the image
			roiRect := image.Rect(col_range, row_range,
				min(col_range+slice_width, cols),
				min(row_range+slice_height, rows))

			img_slice := input.Region(roiRect)
			slices = append(slices, img_slice)
		}

	}

	rand.Shuffle(len(slices), func(i, j int) {
		slices[i], slices[j] = slices[j], slices[i]
	})

	new_height := min(part_rows*slice_height, rows)
	new_width := min(part_cols*slice_width, cols)

	shuffled_image := gocv.NewMatWithSize(new_height, new_width, dataType)

	for idx, slice := range slices {

		row_idx := idx / part_cols
		col_idx := idx % part_cols

		rowStart := row_idx * slice_height
		colStart := col_idx * slice_width

		sliceRows := slice.Rows()
		sliceCols := slice.Cols()

		roiRect := image.Rect(colStart, rowStart, colStart+sliceCols, rowStart+sliceRows)
		roi := shuffled_image.Region(roiRect)
		slice.CopyTo(&roi)
		roi.Close()
	}

	return &shuffled_image, nil
}
