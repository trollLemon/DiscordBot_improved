package classification

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("discord-bot/classification")

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
	ctx, span := tracer.Start(context.Background(), "classification.ClassifyImage", trace.WithAttributes(
		attribute.String("classification.content_type", contentType),
		attribute.Int("classification.image_size_bytes", len(image)),
	))
	defer span.End()

	jobId, err := i.do(ctx, image, contentType)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}

	span.SetAttributes(attribute.String("classification.job_id", jobId))

	classification, err := i.poll(ctx, jobId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}

	span.SetAttributes(attribute.String("classification.label", classification.Class))

	return classification.Class, nil
}
