package Classification

import (
	"errors"
	"fmt"
	"time"

	"github.com/trollLemon/DiscordBot/internal/apiErrors"
)

var (
	ErrRequestTimedOut  = errors.New("the service timed out")
	ErrRequestBadParams = errors.New("given parameters are invalid")
	ErrNotPngOrJpg      = errors.New("image format is invalid")
	ErrGeneral          = errors.New("error while calling api")
)

func errorChecker(err error) error {
	if errors.Is(err, apierrors.ErrRetry) {
		return fmt.Errorf("request timed out. %w", ErrRequestTimedOut)
	}
	if errors.Is(err, apierrors.ErrAPI) {
		return fmt.Errorf("given parameters were invalid. %w", ErrRequestBadParams)
	}
	if errors.Is(err, errBadFileType) {
		return fmt.Errorf("given image is not a png or jpeg image. %w", ErrNotPngOrJpg)
	}
	if err != nil {
		return fmt.Errorf("failed to get image from api. %w", ErrGeneral)
	}
	return nil
}

func NewImageClassification(maxWaitTIme time.Duration, url, sendEndpoint, pollEndpoint string) *ImageClassification {
	return &ImageClassification{
		apiURL:            url,
		pollEndpoint:      pollEndpoint,
		sendImageEndpoint: sendEndpoint,
		maxWaitTime:       maxWaitTIme,
	}
}

// ClassifyImage returns the predicted classification of an image, or an error, using the classification api
func (i *ImageClassification) ClassifyImage(image []byte, contentType string) (string, error) {
	jobId, err := i.do(image, contentType)
	if err != nil {
		return "", errorChecker(err)
	}

	classification, err := i.poll(jobId)
	if err != nil {

		return "", errorChecker(err)
	}

	return classification.Class, nil
}
