package jobs

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"gocv.io/x/gocv"
)

var (
	ErrOpenCV = errors.New("opencv error")
)

// Choice represents which morphology operation to perform.
type Choice string

const (
	Dilate Choice = "Dilate"
	Erode         = "Erode"
)

type Invert struct {
	name string
}

func (_ *Invert) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "invert"),
	}
}

func (_ *Invert) Run(input *gocv.Mat) (*gocv.Mat, error) {
	white := gocv.NewMatWithSizeFromScalar(gocv.Scalar{255.0, 255.0, 255.0, 255.0}, input.Rows(), input.Cols(), input.Type())
	defer white.Close()

	inverted := gocv.NewMat()

	gocv.Subtract(white, *input, &inverted)

	return &inverted, nil
}

type Saturate struct {
	Value float32
}

func (s *Saturate) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "saturate"),
		attribute.Float64("value", float64(s.Value)),
	}
}

func (s *Saturate) Run(input *gocv.Mat) (*gocv.Mat, error) {
	hsvImage := gocv.NewMat()
	defer hsvImage.Close()

	var bgr gocv.Mat
	switch input.Channels() {
	case 3:
		bgr = *input
	case 1, 4:
		bgr = gocv.NewMat()
		defer bgr.Close()

		conversion := gocv.ColorGrayToBGR
		if input.Channels() == 4 {
			conversion = gocv.ColorBGRAToBGR
		}

		if err := gocv.CvtColor(*input, &bgr, conversion); err != nil {
			return nil, fmt.Errorf("%w, gocv cvtColor failed: %v", ErrOpenCV, err)
		}
	default:
		return nil, fmt.Errorf("image has unexpected number of channels: %d, %w", input.Channels(), ErrOpenCV)
	}

	err := gocv.CvtColor(bgr, &hsvImage, gocv.ColorBGRToHLSFull)

	if err != nil {
		return nil, fmt.Errorf("%w, gocv cvtColor failed: %v", ErrOpenCV, err)
	}

	chans := gocv.Split(hsvImage)

	hue := chans[0]
	light := chans[1]
	sat := chans[2]
	defer func() {
		hue.Close()
		light.Close()
		sat.Close()
	}()

	saturated := gocv.NewMat()
	defer saturated.Close()

	sameType := -1
	beta := 0

	sat.ConvertToWithParams(&sat, gocv.MatType(sameType), s.Value, float32(beta))

	gocv.Merge([]gocv.Mat{hue, light, sat}, &saturated)

	imgSaturated := gocv.NewMat()

	err = gocv.CvtColor(saturated, &imgSaturated, gocv.ColorHLSToBGR)
	if err != nil {
		imgSaturated.Close()
		return nil, fmt.Errorf("%w, gocv cvtColor failed: %v", ErrOpenCV, err)
	}
	return &imgSaturated, nil
}

type EdgeDetect struct {
	TLower  float32
	THigher float32
}

func (e *EdgeDetect) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "edge_detect"),
		attribute.Float64("t_lower", float64(e.TLower)),
		attribute.Float64("t_higher", float64(e.THigher)),
	}
}

func (e *EdgeDetect) Run(input *gocv.Mat) (*gocv.Mat, error) {
	edges := gocv.NewMat()

	gocv.Canny(*input, &edges, e.TLower, e.THigher)

	return &edges, nil
}

type Morphology struct {
	KernelSize int
	Iterations int
	Op         Choice
}

func (m *Morphology) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "morphology"),
		attribute.Int("kernel_size", m.KernelSize),
		attribute.Int("iterations", m.Iterations),
		attribute.String("morphology_operation", string(m.Op)),
	}
}

func (m *Morphology) Run(input *gocv.Mat) (*gocv.Mat, error) {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: m.KernelSize, Y: m.KernelSize})
	defer kernel.Close()
	morphedImage := gocv.NewMat()

	if m.Op == Dilate {
		gocv.Dilate(*input, &morphedImage, kernel)
	} else {
		gocv.Erode(*input, &morphedImage, kernel)
	}

	return &morphedImage, nil
}

type Reduce struct {
	Quality float32
}

func (r *Reduce) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "reduce"),
		attribute.Float64("quality", float64(r.Quality)),
	}
}

func (r *Reduce) Run(input *gocv.Mat) (*gocv.Mat, error) {
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

func (a *AddText) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "add_text"),
		attribute.String("text", a.Text),
		attribute.Float64("font_scale", a.FontScale),
		attribute.Float64("x", a.X),
		attribute.Float64("y", a.Y),
	}
}

func (a *AddText) Run(input *gocv.Mat) (*gocv.Mat, error) {
	if input == nil {
		return nil, ErrImgEmpty

	}

	result := input.Clone()

	rows, cols := input.Rows(), input.Cols()

	xPos, yPos := int(float64(rows)*a.X), int(float64(cols)*a.Y)

	thickness := 1
	lineType := gocv.LineAA

	gocv.PutTextWithParams(&result, a.Text, image.Point{X: xPos, Y: yPos}, gocv.FontHersheyPlain, a.FontScale, color.RGBA{255, 255, 255, 255}, thickness, lineType, false)

	return &result, nil
}

type RandomFilter struct {
	KernelSize int
	Min        int
	Max        int
	Normalize  bool
}

func (r *RandomFilter) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "random_filter"),
		attribute.Int("kernel_size", r.KernelSize),
		attribute.Int("min", r.Min),
		attribute.Int("max", r.Max),
		attribute.Bool("normalize", r.Normalize),
	}
}

func (r *RandomFilter) Run(input *gocv.Mat) (*gocv.Mat, error) {
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

func (s *Shuffle) ToAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("operation", "shuffle"),
		attribute.Int("partitions", s.Partitions),
	}
}

func (s *Shuffle) Run(input *gocv.Mat) (*gocv.Mat, error) {
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
		slice.Close()
		roi.Close()
	}

	return &shuffledImage, nil
}
