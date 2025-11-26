package gomanip_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/trollLemon/DiscordBot/internal/gomanip"
)

var (
	readTimeout = time.Second
)

func testEndpoint(t *testing.T, do func(*gomanip.GoManip, []byte, string) ([]byte, error)) {
	tests := []struct {
		name           string
		contentType    string
		wantErr        error
		wantTCPerr     bool
		image          image.Image
		genHandlerFunc func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:        "Success",
			contentType: "image/png",
			image:       image.NewRGBA(image.Rect(0, 0, 100, 100)),
			genHandlerFunc: func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				err := png.Encode(&buf, *img)

				if err != nil {
					t.Fatal(err)
				}

				testImgBytes := buf.Bytes()

				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", contentType)
					w.WriteHeader(http.StatusOK)
					w.Write(testImgBytes)
				}
			},
		},

		{
			name:        "Bad Request",
			contentType: "image/png",
			image:       image.NewRGBA(image.Rect(0, 0, 100, 100)),
			wantErr:     gomanip.ErrBadParams,
			genHandlerFunc: func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", contentType)
					w.WriteHeader(http.StatusBadRequest)

					errResponse := gomanip.ErrorResponse{
						Detail: "bad params",
					}

					bytes, err := json.Marshal(errResponse)
					if err != nil {
						t.Fatal(err)
					}

					w.Write(bytes)
				}
			},
		},
		{
			name:        "Bad Request and error deserializing error json",
			wantErr:     gomanip.ErrGeneral,
			contentType: "image/png",
			image:       image.NewRGBA(image.Rect(0, 0, 100, 100)),
			genHandlerFunc: func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", contentType)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},

		{
			name:        "Server Error",
			contentType: "image/png",
			image:       image.NewRGBA(image.Rect(0, 0, 100, 100)),
			wantErr:     gomanip.ErrGeneral,
			genHandlerFunc: func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", contentType)
					w.WriteHeader(http.StatusInternalServerError)

					errResponse := gomanip.ErrorResponse{
						Detail: "server error",
					}

					bytes, err := json.Marshal(errResponse)
					if err != nil {
						t.Fatal(err)
					}

					w.Write(bytes)
				}
			},
		},

		{
			name:        "Time out",
			contentType: "image/png",
			image:       image.NewRGBA(image.Rect(0, 0, 100, 100)),
			wantErr:     gomanip.ErrTimedOut,
			genHandlerFunc: func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", contentType)
					w.WriteHeader(http.StatusGone)

					errResponse := gomanip.ErrorResponse{
						Detail: "no where to be found...",
					}

					bytes, err := json.Marshal(errResponse)
					if err != nil {
						t.Fatal(err)
					}

					w.Write(bytes)
				}
			},
		},

		{
			name:        "TCP error",
			contentType: "image/png",
			image:       image.NewRGBA(image.Rect(0, 0, 100, 100)),
			wantErr:     gomanip.ErrGeneral,
			wantTCPerr:  true,
			genHandlerFunc: func(img *image.Image, contentType string, t *testing.T) func(w http.ResponseWriter, r *http.Request) {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", contentType)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := png.Encode(&buf, tt.image)

			if err != nil {
				t.Fatal(err)
			}

			testImgBytes := buf.Bytes()

			handlerFunc := tt.genHandlerFunc(&tt.image, tt.contentType, t)

			mockServer := httptest.NewServer(http.HandlerFunc(handlerFunc))

			url := "http://notavalidurl"

			if !tt.wantTCPerr {
				url = mockServer.URL
			}

			manip := gomanip.NewGoManip(url, readTimeout)
			_, err = do(manip, testImgBytes, tt.contentType)
			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr), "returned wrong type of error, want \" %s \", got \" %s \"", tt.wantErr, err)
			} else {
				assert.Nil(t, err, "Expected no error, but got: %v", err)
			}

			mockServer.Close()
		})
	}

}

func TestEndpoints(t *testing.T) {
	tests := []struct {
		about string
		do    func(*gomanip.GoManip, []byte, string) ([]byte, error)
	}{
		{
			about: "Shuffle Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.Shuffle(g, bytes, contentType, 42)
			},
		},

		{
			about: "RandomFilter Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.RandomFilter(g, bytes, contentType, 2, -1, 1, false)
			},
		},
		{
			about: "InvertImage Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.InvertImage(g, bytes, contentType)
			},
		},
		{
			about: "SaturateImage Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.SaturateImage(g, bytes, contentType, 42)
			},
		},

		{
			about: "EdgeDetect Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.EdgeDetect(g, bytes, contentType, 100, 200)
			},
		},

		{
			about: "DilateImage Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.DilateImage(g, bytes, contentType, 42, 1)
			},
		},

		{
			about: "ErodeImage Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.ErodeImage(g, bytes, contentType, 42, 1)
			},
		},

		{
			about: "AddText Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.AddText(g, bytes, contentType, "bleh", 1.0, 0.5, 0.5)
			},
		},

		{
			about: "ReduceImage Endpoint",
			do: func(g *gomanip.GoManip, bytes []byte, contentType string) ([]byte, error) {
				return gomanip.Reduced(g, bytes, contentType, 0.1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.about, func(t *testing.T) {
			testEndpoint(t, tt.do)
		})
	}
}
