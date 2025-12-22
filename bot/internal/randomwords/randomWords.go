package randomwords

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)


var (
	ErrNotFound  = errors.New("not found")
	ErrEmpty     = errors.New("no items in set")
	ErrDuplicate = errors.New("item already in set")
	ErrResult    = errors.New("failed to run operation")
)

type RandomWords struct {
	ctx     context.Context
	rdb     *redis.Client
	setName string
}


func NewRandomWords( rdb *redis.Client, ctx context.Context, setName string ) *RandomWords {
	return &RandomWords{
		rdb: rdb,
		ctx: ctx,
		setName: setName,
	}
}

func (r *RandomWords) Insert(item string) error {
	num, err := r.rdb.SAdd(r.ctx, r.setName, item).Result()
	if err != nil {
		log.Err(err).Msg("failed to insert into the database")
		return fmt.Errorf("%w, %v", ErrResult, err) 
	}

	if num == 0 {
		log.Error().Msgf("failed to insert, item %s is already in the set", item)
		return ErrDuplicate
	}

	return nil
}

func (r *RandomWords) Delete(item string) error {
	num, err := r.rdb.SRem(r.ctx, r.setName, item).Result()
	if err != nil {
		log.Err(err).Msg("failed to remove from the database")
		return fmt.Errorf("%w, %v", ErrResult, err) 

	}

	if num == 0 {
		log.Error().Msgf("failed to remove, item %s is not in the set", item)
		return ErrNotFound
	}

	return nil
}

func (r *RandomWords) GetRandom(n int) ([]string, error) {
	values, err := r.rdb.SRandMemberN(r.ctx, r.setName, int64(n)).Result()
	if err != nil {
		log.Err(err).Msg("failed to get items from the database")
		return []string{}, fmt.Errorf("%w, %v", ErrResult, err) 
	}

	if len(values) == 0 {
		log.Error().Msg("Cannot fetch random items since the set is empty")
		return []string{}, ErrEmpty
	}

	return values, nil
}

func (r *RandomWords) GetAll() ([]string, error) {
	values, err := r.rdb.SMembers(r.ctx, r.setName).Result()
	if err != nil {
		log.Err(err).Msg("failed to get items from the database")
		return []string{}, fmt.Errorf("%w, %v", ErrResult, err) 
	}
	if len(values) == 0 {
		log.Error().Msg("Cannot fetch all items since the set is empty")
		return []string{}, ErrEmpty
	}

	return values, nil
}
