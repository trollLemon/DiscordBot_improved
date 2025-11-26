package util_test

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trollLemon/DiscordBot/internal/util"
)

const maxDiff = 256 // allow for error for jpeg compression

func newTestServer(buf bytes.Buffer, contentType string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(buf.Bytes())
		if err != nil {
			panic(err) // if we fail to write during the test we shouldn't continue running
		}
	}))
}

func bytesDifference(bytes1, bytes2 []byte) (int64, error) {
	if len(bytes1) != len(bytes2) {
		return -1, fmt.Errorf("bytes length not equal")
	}

	var accumErr int32 = 0

	for i := 0; i < len(bytes1); i++ {
		accumErr += squareDiff(bytes1[i], bytes2[i])
	}

	return int64(math.Sqrt(float64(accumErr))), nil
}

func squareDiff(x, y uint8) int32 {
	d := x - y
	return int32(d * d)
}

func TestImageFromURL(t *testing.T) {

	testImage := image.NewNRGBA(image.Rect(0, 0, 256, 256))

	var jpegBuf bytes.Buffer
	var pngBuf bytes.Buffer

	if err := png.Encode(&pngBuf, testImage); err != nil {
		t.Fatal(err)
	}

	if err := jpeg.Encode(&jpegBuf, testImage, nil); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		contentType string
		image       image.Image
		buf         bytes.Buffer
		wantErr     bool
	}{
		{
			name:        "ImageFromURL .PNG",
			image:       testImage,
			buf:         pngBuf,
			contentType: "image/png",
			wantErr:     false,
		},
		{
			name:        "ImageFromURL .JPG",
			image:       testImage,
			buf:         jpegBuf,
			contentType: "image/jpeg",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testServer := newTestServer(tt.buf, tt.contentType)
			defer testServer.Close()

			retrievedImg, format, err := util.GetImageFromURL(testServer.URL)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.contentType, format)

			diff, err := bytesDifference(retrievedImg, tt.buf.Bytes())

			if err != nil {
				t.Fatal(err)
			}

			if diff > maxDiff {
				t.Errorf("expected image difference to be less than %d, got %d", maxDiff, diff)
			}

		})
	}

}
