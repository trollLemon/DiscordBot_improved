package Classification_test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trollLemon/DiscordBot/internal/classification"
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

	mockJobId := "mockJobId"

	getClassificationEndpointWithId := Classification.GetClassificationEndpoint + "/" + mockJobId
	tests := []struct {
		name            string
		wantErr         error
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
		},

		{
			name: "Server Side Error (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				resp := map[string]string{"detail": "error"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: Classification.ErrGeneral,
		},
		{
			name: "Bad Request (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				resp := map[string]string{"detail": "error"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: Classification.ErrRequestBadParams,
		},

		{
			name: "Service Unavailable (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: Classification.ErrGeneral,
		},
		{
			name: "Gateway timeout (Sending Image)",
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: Classification.ErrGeneral,
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
				resp := map[string]string{"detail": "error"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			},
			wantErr: Classification.ErrGeneral,
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
				resp := map[string]string{"detail": "error"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			},
			wantErr: Classification.ErrRequestBadParams,
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
			wantErr: Classification.ErrGeneral,
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
			wantErr: Classification.ErrGeneral,
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
			wantErr: Classification.ErrRequestTimedOut,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mux := http.NewServeMux()
			mux.HandleFunc(Classification.SendImageEndpoint, tt.postHandlerFunc)
			mux.HandleFunc(getClassificationEndpointWithId, tt.getHandlerFunc)
			mockServer := httptest.NewServer(mux)
			waitTime := 3 * time.Second
			clsApi := Classification.NewImageClassification(waitTime, mockServer.URL, Classification.SendImageEndpoint, Classification.GetClassificationEndpoint)

			_, err := clsApi.ClassifyImage(testImgBytes, testContentType)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.Nil(t, err)
			}

			mockServer.Close()

		})
	}

}
