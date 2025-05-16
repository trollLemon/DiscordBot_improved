package imagemanip_test

import (
	"bot/Core/Services/ImageManip"
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoManip_ProcessImage(t *testing.T) {

	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}

	testImgBytes := buf.Bytes()
	testContentType := "image/png"

	testHttpClient := &http.Client{}

	tests := []struct {
		name        string
		wantErr     bool
		handlerFunc func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:    "Success",
			wantErr: false,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:    "Bad Request",
			wantErr: true,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"detail":"error"}`))
			},
		},
		{
			name:    "Bad Request and error deserializing error json",
			wantErr: true,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				var empty []byte
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(empty)
			},
		},

		{
			name:    "Server Side Error",
			wantErr: true,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"detail":"error"}`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockServer := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))

			testApi := imagemanip.NewGoManip(testHttpClient, mockServer.URL)

			_, err := testApi.ProcessImage(testImgBytes, testContentType, "", "")
			if tt.wantErr {
				// If we expect an error, make sure err is not nil
				assert.NotNil(t, err, "Expected an error, but got none")
			} else {
				// If we do not expect an error, ensure err is nil
				assert.Nil(t, err, "Expected no error, but got: %v", err)
			}

			mockServer.Close()

		})
	}

}
