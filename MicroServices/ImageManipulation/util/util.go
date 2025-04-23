package util

import (
	"github.com/labstack/echo/v4"
	"gocv.io/x/gocv"
	"io"
)

func bytesToMat(data []byte) (*gocv.Mat, error) {

	mat, err := gocv.IMDecode(data, gocv.IMReadUnchanged)

	return &mat, err
}

func GetImageFromBody(c echo.Context) (*gocv.Mat, error) {
	imageBytes, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()

	if err != nil {
		return nil, err
	}

	mat, err := bytesToMat(imageBytes)

	return mat, err
}
