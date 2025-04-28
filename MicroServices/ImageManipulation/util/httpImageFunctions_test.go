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

const THRESHOLD = 10

func newTestContextWithImage(image *gocv.Mat, fileExt, contentType string) echo.Context {

	e := echo.New()

	encodedImage, err := gocv.IMEncode(gocv.FileExt(fileExt), *image)

	if err != nil {
		panic(err)
	}

	imgBytes := encodedImage.GetBytes()

	imageReader := bytes.NewReader(imgBytes)

	req := httptest.NewRequest(http.MethodPost, "/", imageReader)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	encodedImage.Close()
	return c
}

func newTestContextInvalidImage(contentType string) echo.Context {

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
			imgCopy := tt.image.Clone()
			ctx := newTestContextWithImage(&imgCopy, tt.fileExt, fileExtToContentType[tt.fileExt])

			image, err := util.GetImageFromBody(ctx)

			if imgCopy.Rows() != image.Rows() || imgCopy.Cols() != image.Cols() || imgCopy.Channels() != image.Channels() {
				t.Errorf("Expected input image and output image to have the same dims. Input: %d, %d, %d. Output: %d, %d, %d", imgCopy.Rows(), imgCopy.Cols(), imgCopy.Channels(), image.Rows(), image.Cols(), image.Channels())
			}

			assert.Equal(t, tt.wantError, err != nil)
			assert.True(t, isEqual(image, &imgCopy))
			image.Close()
			imgCopy.Close()
		})
	}

}

func TestGetImageFromBodyInvalidImage(t *testing.T) {

	ctx := newTestContextInvalidImage("image/png")

	_, err := util.GetImageFromBody(ctx)
	assert.NotNil(t, err)
}
