package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrImageNotFound     = errors.New("image not found")
	ErrFailedToSaveImage = errors.New("failed to save image")
	ErrFailedToGetImage  = errors.New("failed to get image")
	ErrJobFailed         = errors.New("job failed with error")
)

type RedisStore struct {
	client            *redis.Client
	retentionDuration time.Duration
}

func NewRedisStore(addr, username, password string, db int, retentionDuration time.Duration) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	return &RedisStore{
		client:            client,
		retentionDuration: retentionDuration,
	}
}

// Close releases the underlying Redis client and its background connections.
func (s *RedisStore) Close() error {
	return s.client.Close()
}

// SaveImage saves the processed image bytes with the associated job ID.
// Note: bytes can be compressed, so the caller is responsible for decompressing the data when retrieving it.
func (s *RedisStore) SaveImage(ctx context.Context, jobId string, image []byte) error {
	requestSpanCtx := trace.SpanContextFromContext(ctx)

	tracer := otel.Tracer("goManip-worker")
	ctx, span := tracer.Start(ctx, "saving image", trace.WithLinks(
		trace.Link{SpanContext: requestSpanCtx}))
	span.SetAttributes(attribute.String("jobId", jobId))

	defer span.End()

	pipe := s.client.TxPipeline()
	pipe.HSet(ctx, jobId, "image", image)
	pipe.Expire(ctx, jobId, s.retentionDuration)

	if _, err := pipe.Exec(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("%v: %w", err.Error(), ErrFailedToSaveImage)
	}
	return nil
}

func (s *RedisStore) WriteError(ctx context.Context, jobId string, err error, traceId string) error {
	requestSpanCtx := trace.SpanContextFromContext(ctx)

	tracer := otel.Tracer("goManip-worker")
	ctx, span := tracer.Start(ctx, "writing error", trace.WithLinks(
		trace.Link{SpanContext: requestSpanCtx}))
	span.SetAttributes(attribute.String("jobId", jobId))

	defer span.End()

	pipe := s.client.TxPipeline()
	pipe.HSet(ctx, jobId, "error", err.Error(), "traceId", traceId)
	pipe.Expire(ctx, jobId, s.retentionDuration)

	if _, err := pipe.Exec(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("%v: %w", err.Error(), ErrFailedToSaveImage)
	}
	return nil
}

// GetImage retrieves the processed image bytes for a given job ID.
// Note: the returned bytes may be compressed, so the caller is responsible for decompressing the data.
func (s *RedisStore) GetImage(ctx context.Context, jobId string) ([]byte, string, error) {
	requestSpanCtx := trace.SpanContextFromContext(ctx)

	tracer := otel.Tracer("goManip-worker")
	ctx, span := tracer.Start(ctx, "getting image", trace.WithLinks(
		trace.Link{SpanContext: requestSpanCtx}))
	span.SetAttributes(attribute.String("jobId", jobId))

	defer span.End()

	results, err := s.client.HGetAll(ctx, jobId).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, "", fmt.Errorf("%v: %w", err.Error(), ErrFailedToGetImage)
	}

	if len(results) == 0 {
		return nil, "", ErrImageNotFound
	}

	image := results["image"]
	traceId := results["traceId"]
	errStr := results["error"]

	if errStr != "" {
		return nil, traceId, fmt.Errorf("%v: %w", errStr, ErrJobFailed)
	}

	return []byte(image), "", nil
}
