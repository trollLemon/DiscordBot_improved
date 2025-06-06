package Classification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v5"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

type AbstractClassificationAPI interface {
	/*
		_sendImage
		sends image bytes as multipart form data to the classification API
		returns a jobId to pollEndpoint, or an error. An error shall be returned if:
		  - There is an error in job queue pipeline (i.e. unsupported image filetype)
		  - Context deadline is reached before a successful result.
		  - Any error from sending http requests or serializing responses
	*/
	sendImage(image []byte, contentType string) (string, error)
	/*
		_pollResponse
		polls the classification api for the job status given a jobId. Returns the predicted classification, or an error.
		An error shall be returned if:
		  - There is an error in the inference pipeline from api.
		  - Context deadline is reached before a successful result.
		  - Any error from sending http requests or serializing responses
	*/
	pollResponse(jobId string) (string, error)
	// ClassifyImage returns the predicted classification of an image, or an error, using the classification api
	ClassifyImage(image []byte, contentType string) (string, error)
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type JobSend struct {
	JobId string `json:"jobId"`
}

type JobPoll struct {
	Class string `json:"Class"`
}

type ImageClassification struct {
	client            *http.Client
	pollEndpoint      string
	sendImageEndpoint string
	apiURL            string
	maxWaitTime       time.Duration
}

func NewImageClassification(client *http.Client, maxWaitTIme time.Duration, endpoint, sendImage, pollImage string) *ImageClassification {
	return &ImageClassification{
		client:            client,
		apiURL:            endpoint,
		pollEndpoint:      pollImage,
		sendImageEndpoint: sendImage,
		maxWaitTime:       maxWaitTIme,
	}
}

func serializeError(resp *http.Response) error {
	var errorResponse ErrorResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return fmt.Errorf("error unmarshalling response body: %s", err)
	}
	return errors.New(errorResponse.Detail)
}

/*
Returns error (if any) based on the StatusCode from a POST request to the API
- Serverside errors (>=500) are returned as standard errors
- 404 returns a permanent error to stop exponential backoff
- 400 returns the serialized result from an error json returned by the server
*/
func handleSendResponse(resp *http.Response) error {

	if resp.StatusCode == http.StatusBadRequest {

		errorDetail := serializeError(resp)

		return backoff.Permanent(errorDetail)
	}

	if resp.StatusCode == http.StatusNotFound {
		return backoff.Permanent(fmt.Errorf("endpoint not found"))
	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusInternalServerError {
		return fmt.Errorf("server Error: %s", resp.Status)
	}

	return nil

}

/*
handlePoll
Returns error (if any) based on the StatusCode from a GET request to the API:
  - Serverside errors (>=500) are returned as standard errors
  - 404 returns a permanent error to stop exponential backoff
  - 400 returns the serialized result from an error json returned by the server
  - 202 is treated as an error, since the task server-side may be in a queue, or not exists.
    If the latter is true, then we will eventually reach the context deadline
*/
func handlePoll(resp *http.Response) error {

	if resp.StatusCode == http.StatusBadRequest {
		errorDetail := serializeError(resp)

		return backoff.Permanent(errorDetail)
	}

	if resp.StatusCode == http.StatusNotFound {
		return backoff.Permanent(fmt.Errorf("endpoint not found"))
	}

	if resp.StatusCode == http.StatusAccepted {
		// api's default state for tasks is pending.
		// However, a job could exist but also be pending. Regardless, we will exit
		// and retry until the context deadline.
		return errors.New("classification pending, retrying")
	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusInternalServerError {
		return fmt.Errorf("server Error: %s", resp.Status)
	}

	return nil
}

/*
_sendImage
sends image bytes as multipart form data to the classification API
returns a jobId to pollEndpoint, or an error. An error shall be returned if:
  - There is an error in job queue pipeline (i.e. unsupported image filetype)
  - Context deadline is reached before a successful result.
  - Any error from sending http requests or serializing responses
*/
func (i *ImageClassification) sendImage(image []byte, contentType string) (string, error) {
	var filename string
	if strings.Contains(contentType, "image") {
		fileType := contentType[6:] //extracts png from image/png, or jpg from image/jpg
		filename = strconv.FormatInt(time.Now().UnixMilli(), 10) + "." + fileType
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), i.maxWaitTime)

	defer cancel()
	sendImageEndpoint := i.apiURL + i.sendImageEndpoint

	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			header := textproto.MIMEHeader{}
			header.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
			header.Set("Content-Type", contentType)

			part, err := writer.CreatePart(header)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("failed to create form file: %w", err))
			}

			_, err = part.Write(image)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("failed to write image bytes: %w", err))
			}

			writer.Close()

			req, err := http.NewRequest("POST", sendImageEndpoint, body)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("failed to create request: %w", err))
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			resp, err := i.client.Do(req)
			if err != nil {
				return nil, backoff.Permanent(fmt.Errorf("http error posting image: %s", err))
			}

			err = handleSendResponse(resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	jobDetails := JobSend{}

	if err := json.Unmarshal(body, &jobDetails); err != nil {
		return "", fmt.Errorf("error getting jobid: %s", err.Error())
	}

	return jobDetails.JobId, nil

}

/*
_pollResponse
polls the classification api for the job status given a jobId. Returns the predicted classification, or an error.
An error shall be returned if:
  - There is an error in the inference pipeline from api.
  - Context deadline is reached before a successful result.
  - Any error from sending http requests or serializing responses
*/
func (i *ImageClassification) pollResponse(jobId string) (string, error) {

	timeoutCtx, cancel := context.WithTimeout(context.Background(), i.maxWaitTime)
	defer cancel()

	getClassEndpoint := i.apiURL + i.pollEndpoint + "/" + jobId

	resp, err := backoff.Retry(
		timeoutCtx,
		func() (*http.Response, error) {

			resp, err := i.client.Get(getClassEndpoint)
			if err != nil {
				return nil, backoff.Permanent(errors.New("an unknown error occurred"))

			}

			err = handlePoll(resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	classification := JobPoll{}

	if err := json.Unmarshal(body, &classification); err != nil {
		return "", errors.New("error getting image classification from api response")
	}

	return classification.Class, nil

}

// ClassifyImage returns the predicted classification of an image, or an error, using the classification api
func (i *ImageClassification) ClassifyImage(image []byte, contentType string) (string, error) {

	jobId, err := i.sendImage(image, contentType)
	if err != nil {
		return "", err
	}

	classification, err := i.pollResponse(jobId)
	if err != nil {

		return "", err
	}
	return classification, nil

}
