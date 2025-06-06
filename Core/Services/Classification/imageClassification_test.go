package Classification_test

import (
	"bot/Core/Services/Classification"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestImageClassification_ClassifyImage(t *testing.T) {

	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}

	testImgBytes := buf.Bytes()
	testContentType := "image/png"

	testHttpClient := &http.Client{}
	mockJobId := "mockJobId"
	sendImageEndpoint := "/api/v1/image"
	getClassificationEndpoint := "/api/v1/image/classification"
	getClassificationEndpointWithId := getClassificationEndpoint + "/" + mockJobId
	tests := []struct {
		name            string
		wantErr         bool
		postHandlerFunc func(w http.ResponseWriter, r *http.Request)
		getHandlerFunc  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "success",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				resp := map[string]string{"Class": "image"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			wantErr: false,
		},

		{
			name: "Server Side Error (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: true,
		},
		{
			name: "Bad Request (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: true,
		},

		{
			name: "Service Unavailable (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: true,
		},
		{
			name: "Gateway timeout (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: true,
		},

		{
			name: "Server Side Error (Polling image classification)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "Bad Request (Polling image classification)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
			},
			wantErr: true,
		},

		{
			name: "Service Unavailable (Polling image classification)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			wantErr: true,
		},
		{
			name: "Gateway timeout (Polling image classification)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)
			},
			wantErr: true,
		},

		{
			name: "Poll Timeout",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mux := http.NewServeMux()
			mux.HandleFunc(sendImageEndpoint, tt.postHandlerFunc)
			mux.HandleFunc(getClassificationEndpointWithId, tt.getHandlerFunc)
			mockServer := httptest.NewServer(mux)
			waitTime := 3 * time.Second
			clsApi := Classification.NewImageClassification(testHttpClient, waitTime, mockServer.URL, sendImageEndpoint, getClassificationEndpoint)

			_, err := clsApi.ClassifyImage(testImgBytes, testContentType)

			assert.Equal(t, tt.wantErr, err != nil)

			mockServer.Close()

		})
	}

}
