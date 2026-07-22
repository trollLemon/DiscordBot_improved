package randomwords

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrEmpty     = errors.New("no items in set")
	ErrDuplicate = errors.New("item already in set")
	ErrResult    = errors.New("failed to run operation")
)

var tracer = otel.Tracer("discord-bot/randomwords")

type RandomWords struct {
	ctx     context.Context
	rdb     *redis.Client
	setName string
}

func NewRandomWords(rdb *redis.Client, ctx context.Context, setName string) *RandomWords {
	return &RandomWords{
		rdb:     rdb,
		ctx:     ctx,
		setName: setName,
	}
}

func (r *RandomWords) Insert(item string) error {
	ctx, span := tracer.Start(r.ctx, "randomwords.Insert", trace.WithAttributes(
		attribute.String("randomwords.set", r.setName),
	))
	defer span.End()

	num, err := r.rdb.SAdd(ctx, r.setName, item).Result()
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert into the database", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("%w, %v", ErrResult, err)
	}

	if num == 0 {
		slog.ErrorContext(ctx, "failed to insert, item is already in the set", "item", item)
		return ErrDuplicate
	}

	return nil
}

func (r *RandomWords) Delete(item string) error {
	ctx, span := tracer.Start(r.ctx, "randomwords.Delete", trace.WithAttributes(
		attribute.String("randomwords.set", r.setName),
	))
	defer span.End()

	num, err := r.rdb.SRem(ctx, r.setName, item).Result()
	if err != nil {
		slog.ErrorContext(ctx, "failed to remove from the database", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("%w, %v", ErrResult, err)

	}

	if num == 0 {
		slog.ErrorContext(ctx, "failed to remove, item is not in the set", "item", item)
		return ErrNotFound
	}

	return nil
}

func (r *RandomWords) GetRandom(n int) ([]string, error) {
	ctx, span := tracer.Start(r.ctx, "randomwords.GetRandom", trace.WithAttributes(
		attribute.String("randomwords.set", r.setName),
		attribute.Int("randomwords.count", n),
	))
	defer span.End()

	values, err := r.rdb.SRandMemberN(ctx, r.setName, int64(n)).Result()
	if err != nil {
		slog.ErrorContext(ctx, "failed to get items from the database", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return []string{}, fmt.Errorf("%w, %v", ErrResult, err)
	}

	if len(values) == 0 {
		slog.ErrorContext(ctx, "Cannot fetch random items since the set is empty")
		return []string{}, ErrEmpty
	}

	return values, nil
}

func (r *RandomWords) GetAll() ([]string, error) {
	ctx, span := tracer.Start(r.ctx, "randomwords.GetAll", trace.WithAttributes(
		attribute.String("randomwords.set", r.setName),
	))
	defer span.End()

	values, err := r.rdb.SMembers(ctx, r.setName).Result()
	if err != nil {
		slog.ErrorContext(ctx, "failed to get items from the database", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return []string{}, fmt.Errorf("%w, %v", ErrResult, err)
	}
	if len(values) == 0 {
		slog.ErrorContext(ctx, "Cannot fetch all items since the set is empty")
		return []string{}, ErrEmpty
	}

	return values, nil
}
