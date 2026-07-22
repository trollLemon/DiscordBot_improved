package gomanip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	ErrNetwork   = errors.New("network error")
	ErrReading   = errors.New("error reading payload")
	ErrUnMarshal = errors.New("error unmarshaling json")
	ErrServer    = errors.New("server returned 5xx status code")
	ErrBadInput  = errors.New("user input was invalid")
	ErrRetry     = errors.New("retrying request")
)

// UserError implements the error interface and provides a way to store a user-friendly message from
// the gomanip API while keeping the error chain intact.
type UserError struct {
	// Message is a custom message to be passed through the error chain.
	Msg string
	// Err is the underlying error
	Err error
}

func (e *UserError) Error() string {
	return e.Err.Error()
}

func (e *UserError) Message() string {
	return e.Msg
}

func (e *UserError) Unwrap() error {
	return e.Err
}

type GomanipError struct {
	Status string `json:"status"`
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

func (g *GoManip) try(ctx context.Context, apiURI, contentType string, imageBytesBuffer *bytes.Buffer) (*http.Response, error) {
	client := http.Client{
		Timeout:   g.readTimeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURI, imageBytesBuffer)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build POST request to gomanip service", "error", err)
		return nil, backoff.Permanent(fmt.Errorf("could not build POST request: %w", ErrNetwork))
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send POST request to gomanip service", "error", err)
		return nil, backoff.Permanent(fmt.Errorf("could not send POST request: %w", ErrNetwork))

	}

	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	var errorResponse GomanipError
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read gomanip error response body", "error", err)
		return nil, backoff.Permanent(fmt.Errorf("%w, %v", ErrReading, err))
	}

	if resp.StatusCode >= 500 {
		return nil, backoff.Permanent(ErrServer)
	}

	if resp.StatusCode == http.StatusBadRequest {
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal gomanip error response", "error", err)
			return nil, backoff.Permanent(fmt.Errorf("%w, %v", ErrUnMarshal, err))
		}

		return nil, backoff.Permanent(&UserError{
			Msg: errorResponse.Detail,
			Err: fmt.Errorf("%w, %s, %s", ErrBadInput, errorResponse.Detail, errorResponse.Status),
		})
	}

	return nil, ErrRetry
}

func (g *GoManip) Do(ctx context.Context, image []byte, contentType, endpoint, queries string) ([]byte, error) {
	apiURI := fmt.Sprintf("%s/%s/%s", g.apiEndpoint, endpoint, queries)

	var imageBytesBuffer *bytes.Buffer

	timeoutCtx, cancel := context.WithTimeout(ctx, g.readTimeout)

	defer cancel()

	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {
			imageBytesBuffer = bytes.NewBuffer(image)
			return g.try(timeoutCtx, apiURI, contentType, imageBytesBuffer)
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		slog.ErrorContext(ctx, "timed out calling gomanip service", "error", err)
		return nil, ErrRetry
	}

	if errors.Is(err, ErrBadInput) {
		slog.ErrorContext(ctx, "user provided invalid parameters", "error", err)
		return nil, err
	}

	if err != nil {
		slog.ErrorContext(ctx, "error calling gomanip service", "error", err)
		return nil, err
	}

	defer resp.Body.Close()

	resultBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed reading image from response body", "error", err)
		return nil, fmt.Errorf("%w, could not read image", ErrReading)
	}

	return resultBytes, nil
}
