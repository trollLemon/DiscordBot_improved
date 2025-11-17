package Classification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/rs/zerolog/log"

	"github.com/trollLemon/DiscordBot/internal/apiErrors"
)

const (
	SendImageEndpoint         = "/api/v1/images"
	GetClassificationEndpoint = "/api/v1/images/classifications"
)

var (
	errBadFileType = errors.New("unsuported image file type")
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type JobSend struct {
	JobId string `json:"jobId"`
}

type ClassResult struct {
	Class string `json:"Class"`
}

type ImageClassification struct {
	pollEndpoint      string
	sendImageEndpoint string
	apiURL            string
	maxWaitTime       time.Duration
}

func (i *ImageClassification) do(image []byte, contentType string) (string, error) {
	var filename string

	filename += time.Millisecond.String()

	if strings.Contains(contentType, "png") {
		filename += ".png"
	} else if strings.Contains(contentType, "jpeg") {
		filename += ".jpeg"
	} else {
		return "", errBadFileType
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), i.maxWaitTime)

	defer cancel()
	sendImageEndpoint := i.apiURL + i.sendImageEndpoint
	client := http.Client{}
	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {
			body := &bytes.Buffer{}

			mimeType := http.DetectContentType(image)
			header := textproto.MIMEHeader{}
			header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
			header.Set("Content-Type", mimeType)
			writer := multipart.NewWriter(body)
			part, err := writer.CreatePart(header)
			if err != nil {
				return nil, backoff.Permanent(apierrors.ErrWriting)
			}

			part.Write(image)
			writer.Close()

			req, err := http.NewRequest("POST", sendImageEndpoint, body)
			if err != nil {
				return nil, apierrors.ErrMakingRequest
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, err := client.Do(req)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("error sending request: %v; %w", err, apierrors.ErrNetwork))
			}
			if resp.StatusCode == http.StatusAccepted {
				return resp, nil
			}

			var errorResponse ErrorResponse

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("failed to read response body. %w", apierrors.ErrReading))
			}
			println(string(contentType))
			println(string(respBody))
			err = json.Unmarshal(respBody, &errorResponse)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("failed to unmarshal response body: %v; %w", err, apierrors.ErrResp))
			}

			if resp.StatusCode == http.StatusBadRequest {
				return nil, backoff.Permanent(fmt.Errorf("%s; %w", errorResponse.Detail, apierrors.ErrAPI))
			}

			if resp.StatusCode >= 500 {
				return nil, backoff.Permanent(fmt.Errorf("server reported error in response: %s; %w", errorResponse.Detail, apierrors.ErrServer))
			}

			log.Warn().Msg("status not OK after calling classification endpoint, attempting to retry http call")
			return nil, fmt.Errorf("status not OK, got: %s. %w", errorResponse.Detail, apierrors.ErrRetry)

		},
	)

	if err != nil {
		log.Err(err).Msg("send request failed")
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	jobDetails := JobSend{}

	if err := json.Unmarshal(body, &jobDetails); err != nil {
		return "", fmt.Errorf("error getting jobid. %w", apierrors.ErrReading)
	}

	return jobDetails.JobId, nil

}

func (i *ImageClassification) poll(jobId string) (*ClassResult, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), i.maxWaitTime)
	defer cancel()

	getClassEndpoint := i.apiURL + i.pollEndpoint + "/" + jobId
	client := http.Client{}
	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {

			resp, err := client.Get(getClassEndpoint)
			if err != nil {
				return nil, backoff.Permanent(apierrors.ErrMakingRequest)
			}

			if resp.StatusCode == http.StatusOK {
				return resp, nil
			}

			if resp.StatusCode == http.StatusAccepted {
				return nil, apierrors.ErrRetry
			}

			var errorResponse ErrorResponse
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("error reading response body: %v; %w", err, apierrors.ErrReading))
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

			log.Warn().Msg("status not OK after calling classification polling endpoint, attempting to retry http call")

			return nil, fmt.Errorf("status not OK, got: %s. %w", errorResponse.Detail, apierrors.ErrRetry)

		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		log.Err(err).Msg("timed out calling classification service")
		return nil, apierrors.ErrRetry
	}

	if err != nil {
		log.Err(err).Msg("poll request failed")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	classification := ClassResult{}

	if err := json.Unmarshal(body, &classification); err != nil {
		return nil, fmt.Errorf("failed to read json body. %w", apierrors.ErrReading)
	}

	return &classification, nil

}
