package Classification

import (
	"errors"
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
		return ErrRequestTimedOut
	}
	if errors.Is(err, apierrors.ErrAPI) {
		return ErrRequestBadParams
	}
	if errors.Is(err, errBadFileType) {
		return ErrNotPngOrJpg
	}
	if err != nil {
		return ErrGeneral
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
