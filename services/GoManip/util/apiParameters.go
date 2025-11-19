package util

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

func ParseSaturation(c echo.Context) (float32, error) {
	saturationStr := c.QueryParam("saturation")
	saturation, err := strconv.ParseFloat(saturationStr, 32)

	if err != nil {
		return 0.0, err
	}
	return float32(saturation), nil
}

func ParseEdgeDetection(c echo.Context) (float32, float32, error) {
	tLowerStr := c.QueryParam("lower")
	tHighStr := c.QueryParam("higher")

	tLower, err := strconv.ParseFloat(tLowerStr, 32)
	if err != nil {
		return 0.0, 0.0, err
	}
	tHigh, err := strconv.ParseFloat(tHighStr, 32)
	if err != nil {
		return 0.0, 0.0, err
	}
	return float32(tLower), float32(tHigh), nil
}

func ParseMorphology(c echo.Context) (string, int, int, error) {
	morphType := c.QueryParam("type")
	kernelSizeStr := c.QueryParam("kernelSize")
	iterationsStr := c.QueryParam("iterations")

	if morphType == "" {
		return "", 0, 0, errors.New("morphology type is required")
	}

	kernelSize, err := strconv.Atoi(kernelSizeStr)
	if err != nil {
		return "", 0, 0, err
	}

	iterations, err := strconv.Atoi(iterationsStr)
	if err != nil {
		return "", 0, 0, err
	}

	return morphType, kernelSize, iterations, nil
}

func ParseReduce(c echo.Context) (float32, error) {
	qualityStr := c.QueryParam("quality")
	quality, err := strconv.ParseFloat(qualityStr, 32)

	if err != nil {
		return 0.0, err
	}
	return float32(quality), nil

}

func ParseAddText(c echo.Context) (string, float64, float64, float64, error) {
	text := c.QueryParam("text")
	fontScaleStr := c.QueryParam("fontScale")
	xPercStr := c.QueryParam("xPerc")
	yPercStr := c.QueryParam("yPerc")

	fontScale, err := strconv.ParseFloat(fontScaleStr, 32)
	if err != nil {
		return "", 0, 0, 0, err
	}
	xPerc, err := strconv.ParseFloat(xPercStr, 32)
	if err != nil {
		return "", 0, 0, 0, err
	}
	yPerc, err := strconv.ParseFloat(yPercStr, 32)
	if err != nil {
		return "", 0, 0, 0, err
	}

	return text, fontScale, xPerc, yPerc, nil
}

func ParseRandomFilter(c echo.Context) (int, int, int, bool, error) {
	minValStr := c.QueryParam("minVal")
	maxValStr := c.QueryParam("maxVal")
	kernelSizeStr := c.QueryParam("kernelSize")
	normalize := c.QueryParam("normalize") == "true"

	minVal, err := strconv.Atoi(minValStr)

	if err != nil {
		return 0, 0, 0, false, err
	}

	maxVal, err := strconv.Atoi(maxValStr)
	if err != nil {
		return 0, 0, 0, false, err
	}
	kernelSize, err := strconv.Atoi(kernelSizeStr)
	if err != nil {
		return 0, 0, 0, false, err
	}

	return minVal, maxVal, kernelSize, normalize, nil
}

func ParseShuffle(c echo.Context) (int, error) {
	partitionsStr := c.QueryParam("partitions")
	partitions, err := strconv.Atoi(partitionsStr)
	if err != nil {
		return 0, err
	}
	return partitions, nil
}
