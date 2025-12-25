package classification

import (
	"time"
)

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
		return "", err
	}

	classification, err := i.poll(jobId)
	if err != nil {

		return "", err
	}

	return classification.Class, nil
}
