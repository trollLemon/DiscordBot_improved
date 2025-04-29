package imagemanip

import (
	"bot/util"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type ImageAPIWrapper struct {
	client      *http.Client
	apiEndpoint string
	backoff     *backoff.ExponentialBackOff
}

func NewImageAPIWrapper(apiEndpoint string) *ImageAPIWrapper {

	return &ImageAPIWrapper{client: &http.Client{}, apiEndpoint: apiEndpoint, backoff: backoff.NewExponentialBackOff()}
}

func (i *ImageAPIWrapper) getImage(imageBytes []byte, contentType string, apiURI string) ([]byte, error) {

	imageBytesBuffer := bytes.NewBuffer(imageBytes)

	operation := func() (*http.Response, error) {

		resp, err := i.client.Post(apiURI, contentType, imageBytesBuffer)
		if err != nil {
			return nil, backoff.Permanent(fmt.Errorf("error fetching image: %s", err))

		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil

		}

		if resp.StatusCode == http.StatusBadRequest {

			var errorResponse ErrorResponse

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v", err)
				return nil, fmt.Errorf("error: internal error calling image api")
			}

			if err := json.Unmarshal(body, &errorResponse); err != nil {
				log.Printf("Error decoding JSON: %v", err)
				return nil, fmt.Errorf("error: internal error calling image api")
			}

			errMsg := "Invalid parameters: " + errorResponse.Detail
			log.Println(errMsg)
			return nil, backoff.Permanent(errors.New(errMsg))

		}

		log.Printf("Error in image API, status code = %d", resp.StatusCode)
		return nil, fmt.Errorf("Error in image API, status code = %d ", resp.StatusCode)
	}

	maxDuration := 30 * time.Second

	timeoutCtx, cancel := context.WithTimeout(context.Background(), maxDuration)

	defer cancel()

	resp, err := backoff.Retry(timeoutCtx, operation, backoff.WithBackOff(i.backoff))

	if err != nil {

		return nil, err
	}

	defer resp.Body.Close()

	resultBytes, err := io.ReadAll(resp.Body)
	return resultBytes, err

}

func (i *ImageAPIWrapper) RandomFilter(image []byte, contentType string, kernelSize, lower, higher int64, normalize bool) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/randomFilter/%s", i.apiEndpoint, util.RandomFilterQuery(kernelSize, lower, higher, normalize))

	return i.getImage(image, contentType, apiURI)

}

func (i *ImageAPIWrapper) InvertImage(image []byte, contentType string) ([]byte, error) {
	apiURI := fmt.Sprintf("%s/invert/", i.apiEndpoint)
	return i.getImage(image, contentType, apiURI)

}

func (i *ImageAPIWrapper) SaturateImage(image []byte, contentType string, saturation int64) ([]byte, error) {

	saturationNorm := float32(saturation) / 100.0

	apiURI := fmt.Sprintf("%s/saturate/%s", i.apiEndpoint, util.SaturateQuery(saturationNorm))

	return i.getImage(image, contentType, apiURI)

}
func (i *ImageAPIWrapper) EdgeDetect(image []byte, contentType string, lower, higher int64) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/edgeDetection/%s", i.apiEndpoint, util.EdgeDetectQuery(lower, higher))

	return i.getImage(image, contentType, apiURI)

}

func (i *ImageAPIWrapper) DilateImage(image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/morphology/%s", i.apiEndpoint, util.DilateQuery(kernelSize, iterations))
	return i.getImage(image, contentType, apiURI)
}

func (i *ImageAPIWrapper) ErodeImage(image []byte, contentType string, kernelSize, iterations int64) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/morphology/%s", i.apiEndpoint, util.ErodeQuery(kernelSize, iterations))
	return i.getImage(image, contentType, apiURI)
}

func (i *ImageAPIWrapper) AddText(image []byte, contentType string, text string, fontScale float32, xPercentage float32, yPercentage float32) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/text/%s", i.apiEndpoint, util.AddTextQuery(text, fontScale, xPercentage, yPercentage))
	return i.getImage(image, contentType, apiURI)
}

func (i *ImageAPIWrapper) Reduced(image []byte, contentType string, quality float32) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/reduced/%s", i.apiEndpoint, util.ReduceQuery(quality))

	return i.getImage(image, contentType, apiURI)
}

func (i *ImageAPIWrapper) Shuffle(image []byte, contentType string, partitions int64) ([]byte, error) {

	apiURI := fmt.Sprintf("%s/shuffle/%s", i.apiEndpoint, util.ShuffleQuery(partitions))

	return i.getImage(image, contentType, apiURI)
}
