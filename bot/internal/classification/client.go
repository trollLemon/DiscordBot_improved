package classification

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
)

const (
	SendImageEndpoint         = "/api/v1/images"
	GetClassificationEndpoint = "/api/v1/images/classifications"
)

var (
	ErrBadFileType   = errors.New("unsuported image file type")
	ErrBadRequest    = errors.New("bad request")
	ErrNetwork       = errors.New("network error")
	ErrMakingRequest = errors.New("error creating HTTP request")
	ErrReading       = errors.New("error reading payload")
	ErrWriting       = errors.New("error writing payload")
	ErrUnMarshal     = errors.New("error unmarshaling json")
	ErrServer        = errors.New("server returned 5xx status code")
	ErrRetry         = errors.New("retrying request")
	ErrNotFound      = errors.New("classification job doesn't exist")
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
		return "", ErrBadFileType
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
				log.Err(err).Msg("failed to write multipart data")
				return nil, backoff.Permanent(ErrWriting)
			}

			part.Write(image)
			writer.Close()

			req, err := http.NewRequest("POST", sendImageEndpoint, body)
			if err != nil {
				log.Err(err).Msg("failed to create POST request")
				return nil, ErrMakingRequest
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, err := client.Do(req)
			if err != nil {
				log.Err(err).Msg("failed to perform classification HTTP request")
				return nil, backoff.Permanent(fmt.Errorf("could not send request: %w", ErrNetwork))
			}

			switch {
			case resp.StatusCode == http.StatusCreated:
				{
					return resp, nil
				}
			case resp.StatusCode == http.StatusGatewayTimeout:
				{
					return nil, ErrRetry
				}
			case resp.StatusCode == http.StatusBadRequest:
				{
					var errorResponse ErrorResponse

					respBody, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Err(err).Msg("failed to read classification request error response body")
						return nil, backoff.Permanent(fmt.Errorf("could not read error response body. %w", ErrReading))
					}

					err = json.Unmarshal(respBody, &errorResponse)
					if err != nil {
						log.Err(err).Msg("failed to unmarshal response into go struct")
						return nil, backoff.Permanent(fmt.Errorf("could not unmarshal error response %w", ErrUnMarshal))
					}
					// this endpoint returns 400 if the filetype is invalid. We already check the filetype before making the request, buts its good
					// to check here regardless just in case resp.StatusCode== a change happens upstream.
					log.Error().Msgf("classification service cannot work with given filetype, returned error: %s", errorResponse.Detail)
					return nil, backoff.Permanent(ErrBadFileType)

				}

			case resp.StatusCode >= 500:
				{
					return nil, backoff.Permanent(fmt.Errorf("could not poll job status: %w", ErrServer))
				}
			default:
				return nil, ErrRetry
			}

		},
	)

	if errors.Is(err, context.DeadlineExceeded) {
		log.Err(err).Msg("timed out sending a classification request")
		return "", ErrRetry
	}

	if err != nil {
		log.Err(err).Msg("failed to send classification request")
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Err(err).Msg("failed to read job id response body")
		return "", fmt.Errorf("could not read job id response body. %w", ErrReading)
	}

	jobDetails := JobSend{}

	if err := json.Unmarshal(body, &jobDetails); err != nil {
		log.Err(err).Msg("failed to unmarshal response body")
		return "", fmt.Errorf("error getting jobid. %w", ErrUnMarshal)
	}

	log.Info().Msgf("POST to classification api succeded, recived job-id: %s", jobDetails.JobId)

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
				return nil, backoff.Permanent(ErrMakingRequest)
			}

			switch {
			case resp.StatusCode == http.StatusOK:
				{
					return resp, nil
				}
			case resp.StatusCode == http.StatusNotFound:
				{
					return nil, backoff.Permanent(ErrNotFound)
				}
			case resp.StatusCode == http.StatusGatewayTimeout:
				{
					return nil, ErrRetry
				}
			case resp.StatusCode == http.StatusBadRequest:
				{
					var errorResponse ErrorResponse

					respBody, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Err(err).Msg("failed to read classification poll request error response body")
						return nil, backoff.Permanent(fmt.Errorf("could not read error response body. %w", ErrReading))
					}

					err = json.Unmarshal(respBody, &errorResponse)
					if err != nil {
						log.Err(err).Msg("failed to unmarshal poll error response into go struct")
						return nil, backoff.Permanent(fmt.Errorf("could not unmarshal error response %w", ErrUnMarshal))
					}
					log.Error().Msgf("classification service reported an error processing request: %s", errorResponse.Detail)
					return nil, ErrRetry
				}

			case resp.StatusCode >= 500:
				{
					return nil, backoff.Permanent(fmt.Errorf("could not poll job status: %w", ErrServer))
				}
			default:
				return nil, ErrRetry
			}

		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		log.Err(err).Msg("timed out polling classification service")
		return nil, ErrRetry
	}

	if err != nil {
		log.Err(err).Msg("poll request failed")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Err(err).Msg("failed to read classification response body")
		return nil, fmt.Errorf("could not read classification response body. %w", ErrReading)

	}

	classification := ClassResult{}

	if err := json.Unmarshal(body, &classification); err != nil {
		return nil, fmt.Errorf("could not unmarshal classification response body. %w", ErrUnMarshal)
	}

	log.Info().Msgf("GET to classification api succeded, got classification: %s", classification.Class)

	return &classification, nil
}
