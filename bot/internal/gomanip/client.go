package gomanip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/rs/zerolog/log"

	"github.com/trollLemon/DiscordBot/internal/apiErrors"
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type GoManip struct {
	apiEndpoint string
	readTimeout time.Duration
}

func NewGoManip(apiEndpoint string, readTimeout time.Duration) *GoManip {
	return &GoManip{
		apiEndpoint: apiEndpoint,
		readTimeout: readTimeout,
	}
}

func (g *GoManip) try(apiURI, contentType string, imageBytesBuffer *bytes.Buffer) (*http.Response, error) {
	client := http.Client{
		Timeout: g.readTimeout,
	}
	resp, err := client.Post(apiURI, contentType, imageBytesBuffer)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("error sending request: %v; %w", err, apierrors.ErrNetwork))

	}

	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	var errorResponse ErrorResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("error reading response body: %v; %w", err, apierrors.ErrResp))
	}

	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("failed to unmarshal response body: %v; %w", err, apierrors.ErrResp))
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, backoff.Permanent(fmt.Errorf("%s; %w", errorResponse.Detail, apierrors.ErrAPI))
	}

	if resp.StatusCode >= 500 {
		return nil, backoff.Permanent(fmt.Errorf("server reported error in response: %s; %w", errorResponse.Detail, apierrors.ErrServer))
	}

	log.Warn().Msg("status not OK after calling gomanip endpoint, attempting to retry http call")

	return nil, fmt.Errorf("status not OK, got: %s. %w", errorResponse.Detail, apierrors.ErrRetry)
}

func (g *GoManip) Do(image []byte, contentType, endpoint, queries string) ([]byte, error) {
	apiURI := fmt.Sprintf("%s/%s/%s", g.apiEndpoint, endpoint, queries)

	var imageBytesBuffer *bytes.Buffer

	timeoutCtx, cancel := context.WithTimeout(context.Background(), g.readTimeout)

	defer cancel()

	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {
			imageBytesBuffer = bytes.NewBuffer(image)
			return g.try(apiURI, contentType, imageBytesBuffer)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		log.Err(err).Msg("timed out calling gomanip service")
		return nil, apierrors.ErrRetry
	}

	if err != nil {
		log.Err(err).Msg("error calling gomanip service")
		return nil, err
	}

	defer resp.Body.Close()

	resultBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("failed reading response body")
		return nil, fmt.Errorf("failed reading response body; %w", apierrors.ErrResp)
	}

	return resultBytes, nil
}
