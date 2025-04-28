package util_test

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"goManip/util"
	"gocv.io/x/gocv"
	"net/http"
	"net/http/httptest"
	"testing"
)

const THRESHOLD = 0

func newTestContextWithImage(image *gocv.Mat, fileExt, contentType string) echo.Context {

	e := echo.New()

	encodedImage, _ := gocv.IMEncode(gocv.FileExt(fileExt), *image)
	defer encodedImage.Close()
	imgBytes := encodedImage.GetBytes()

	imageReader := bytes.NewReader(imgBytes)

	req := httptest.NewRequest(http.MethodPost, "/", imageReader)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c
}

func newTestContextInvalidImage(fileExt, contentType string) echo.Context {

	e := echo.New()

	emptyBytes := []byte("")

	imageReader := bytes.NewReader(emptyBytes)

	req := httptest.NewRequest(http.MethodPost, "/", imageReader)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c

}

func isEqual(image, other *gocv.Mat) bool {

	if image.Cols() != other.Cols() || image.Rows() != other.Rows() || image.Type() != other.Type() {
		return false
	}

	diff := gocv.NewMat()
	defer diff.Close()

	grayImage := gocv.NewMat()
	defer grayImage.Close()
	gocv.CvtColor(*image, &grayImage, gocv.ColorBGRToGray)

	grayOther := gocv.NewMat()
	defer grayOther.Close()
	gocv.CvtColor(*other, &grayOther, gocv.ColorBGRToGray)

	gocv.AbsDiff(grayImage, grayOther, &diff)

	nonZero := gocv.CountNonZero(diff)
	return nonZero <= THRESHOLD
}

func TestGetImageFromBody(t *testing.T) {
	fileExtToContentType := map[string]string{
		".png": "image/png",
		".jpg": "image/jpeg",
	}
	smallImage := gocv.NewMatWithSize(700, 400, gocv.MatTypeCV8UC3)
	largeImage := gocv.NewMatWithSize(1920, 1080, gocv.MatTypeCV8UC3)
	//emptyImage := gocv.NewMat()

	tests := []struct {
		name      string
		image     *gocv.Mat
		fileExt   string
		wantError bool
	}{

		{
			name:      "Test with small image (png)",
			image:     &smallImage,
			fileExt:   ".png",
			wantError: false,
		},
		{
			name:      "Test with small image (jpg)",
			image:     &smallImage,
			fileExt:   ".jpg",
			wantError: false,
		},
		{
			name:      "Test with large image (png)",
			image:     &largeImage,
			fileExt:   ".png",
			wantError: false,
		},
		{
			name:      "Test with large image (jpg)",
			image:     &largeImage,
			fileExt:   ".jpg",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := newTestContextWithImage(tt.image, tt.fileExt, fileExtToContentType[tt.fileExt])

			image, err := util.GetImageFromBody(ctx)

			assert.Equal(t, tt.wantError, err != nil)
			assert.True(t, isEqual(image, tt.image))

		})
	}

}

func TestGetImageFromBodyEmptyMat(t *testing.T) {
	testImage := gocv.NewMat()
	ctx := newTestContextWithImage(&testImage, ".png", "image/png")

	_, err := util.GetImageFromBody(ctx)
	assert.NotNil(t, err)
}

func TestGetImageFromBodyInvalidImage(t *testing.T) {

	ctx := newTestContextInvalidImage(".png", "image/png")

	_, err := util.GetImageFromBody(ctx)
	assert.NotNil(t, err)
}
