/* Wrapper for Image Manipulation microservice. Assumes all text (URLs and text parameters) are encoded prior to the function calls

   Server provides a RESTful api with the following endpoints:
    - /api/randomFilteredImage/{image_link:path}/{kernel_size}/{low}/{high}/{kernel_type}
    - /api/invertedImage/{image_link:path}
    - /api/saturatedImage/{image_link:path}/{saturation}
    - /api/edgeImage/{image_link:path}/{lower}/{higher}
    - /api/dilatedImage/{image_link:path}/{box_size}/{iterations}
    - /api/erodedImage/{image_link:path}/{box_size}/{iterations}
    - /api/textImage/{image_link:path}/{text}/{font_scale}/{x}/{y}
    - /api/reducedImage/{image_link:path}/{quality}
*/

package imagemanip

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cenkalti/backoff/v5"
)

type ImageAPI interface {
	get_image(string) ([]byte, error)
}

type ImageAPIWrapper struct {
	client       *http.Client
	api_endpoint string
	backoff      *backoff.ExponentialBackOff
}

func NewImageAPIWrapper(api_endpoint string, max_backoff_time int8) ImageAPIWrapper {

	backoff := backoff.NewExponentialBackOff()

	return ImageAPIWrapper{client: &http.Client{}, api_endpoint: api_endpoint, backoff: backoff}
}

func (i *ImageAPIWrapper) get_image(url string) ([]byte, error) {

	operation := func() (*http.Response, error) {

		resp, err := i.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Error fetching image: %s", err)

		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil

		}
		log.Printf("Error in image API, status code = %d", resp.StatusCode)
		return nil, fmt.Errorf("Error in image API, status code = %d ", resp.StatusCode)
	}

	resp, err := backoff.Retry(context.TODO(), operation, backoff.WithBackOff(i.backoff))

	defer resp.Body.Close()
	imageBytes, err := io.ReadAll(resp.Body)
	return imageBytes, err

}

func (i *ImageAPIWrapper) RandomFilter(url string, kernel_size int, lower_bound int, upper_bound int) ([]byte, error) {

	imageType := "randomFilteredImage"
	api_url := fmt.Sprintf("%s/%s/%s/%d/%d/%d/raw", i.api_endpoint, imageType, url, kernel_size, lower_bound, upper_bound)

	return i.get_image(api_url)

}

func (i *ImageAPIWrapper) InvertImage(url string) ([]byte, error) {

	imageType := "invertedImage"
	api_url := fmt.Sprintf("%s/%s/%s", i.api_endpoint, imageType, url)

	return i.get_image(api_url)

}

func (i *ImageAPIWrapper) SaturateImage(url string, magnitude int) ([]byte, error) {

	imageType := "saturatedImage"
	magnitude_norm := float32(magnitude) / 100.0

	api_url := fmt.Sprintf("%s/%s/%s/%0.2f", i.api_endpoint, imageType, url, magnitude_norm)

	return i.get_image(api_url)
}
func (i *ImageAPIWrapper) EdgeDetect(url string, lower_bound int, upper_bound int) ([]byte, error) {

	imageType := "edgeImage"

	api_url := fmt.Sprintf("%s/%s/%s/%d/%d", i.api_endpoint, imageType, url, lower_bound, upper_bound)

	return i.get_image(api_url)
}

func (i *ImageAPIWrapper) DilateImage(url string, box_size int, iterations int) ([]byte, error) {

	imageType := "dilatedImage"

	api_url := fmt.Sprintf("%s/%s/%s/%d/%d", i.api_endpoint, imageType, url, box_size, iterations)

	return i.get_image(api_url)
}

func (i *ImageAPIWrapper) ErodeImage(url string, box_size int, iterations int) ([]byte, error) {

	imageType := "erodedImage"

	api_url := fmt.Sprintf("%s/%s/%s/%d/%d", i.api_endpoint, imageType, url, box_size, iterations)

	return i.get_image(api_url)
}

func (i *ImageAPIWrapper) AddText(url string, text string, font_scale float32, x_percentage float32, y_percentage float32) ([]byte, error) {

	imageType := "textImage"

	api_url := fmt.Sprintf("%s/%s/%s/%s/%0.2f/%0.2f/%0.2f", i.api_endpoint, imageType, url, text, font_scale, x_percentage, y_percentage)

	return i.get_image(api_url)
}

func (i *ImageAPIWrapper) Reduced(url string, quality float32) ([]byte, error) {

	imageType := "reducedImage"

	api_url := fmt.Sprintf("%s/%s/%s/%0.2f", i.api_endpoint, imageType, url, quality)

	return i.get_image(api_url)
}
