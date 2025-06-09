package imagemanip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v5"
	"io"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type AbstractImageAPI interface {
	ProcessImage(image []byte, contentType string, endpoint string, query string) ([]byte, error)
	_TryApi(apiURI, contentType string, imageBytesBuffer *bytes.Buffer) (*http.Response, error)
}

type GoManip struct {
	client      *http.Client
	apiEndpoint string
}

func NewGoManip(client *http.Client, apiEndpoint string) *GoManip {
	return &GoManip{
		client:      client,
		apiEndpoint: apiEndpoint,
	}
}

func (g *GoManip) _TryApi(apiURI, contentType string, imageBytesBuffer *bytes.Buffer) (*http.Response, error) {

	resp, err := g.client.Post(apiURI, contentType, imageBytesBuffer)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("error fetching image: %s", err))

	}

	if resp.StatusCode == http.StatusOK {
		return resp, nil

	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusInternalServerError {
		return nil, backoff.Permanent(fmt.Errorf("server Error: %s", resp.Status))
	}

	if resp.StatusCode == http.StatusBadRequest {

		var errorResponse ErrorResponse

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %s", err)
		}

		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, fmt.Errorf("error reading response body: %s", err)
		}

		return nil, backoff.Permanent(errors.New("Invalid parameters: " + errorResponse.Detail))

	}

	return nil, fmt.Errorf("Error in image API, status code = %d ", resp.StatusCode)
}

func (g *GoManip) apiRequest(imageBytes []byte, contentType string, apiURI string) ([]byte, error) {

	var imageBytesBuffer *bytes.Buffer

	maxDuration := 30 * time.Second
	timeoutCtx, cancel := context.WithTimeout(context.Background(), maxDuration)

	defer cancel()

	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {
			imageBytesBuffer = bytes.NewBuffer(imageBytes)
			return g._TryApi(apiURI, contentType, imageBytesBuffer)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if err != nil {

		return nil, err
	}

	defer resp.Body.Close()

	resultBytes, err := io.ReadAll(resp.Body)
	return resultBytes, err

}

func (g *GoManip) ProcessImage(image []byte, contentType, endpoint, queries string) ([]byte, error) {
	apiURI := fmt.Sprintf("%s/%s/%s", g.apiEndpoint, endpoint, queries)
	return g.apiRequest(image, contentType, apiURI)
}
