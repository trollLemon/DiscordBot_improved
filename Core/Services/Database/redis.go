package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	ctx     context.Context
	rdb     *redis.Client
	setName string
}

func NewRedisClient(ctx context.Context, redisClient *redis.Client, defaultSetName string) *Redis {
	return &Redis{
		ctx:     ctx,
		rdb:     redisClient,
		setName: defaultSetName,
	}
}

func (r *Redis) Insert(item string) error {
	num, err := r.rdb.SAdd(r.ctx, r.setName, item).Result()
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New(fmt.Sprintf("item %s already in set", item))
	}

	return nil
}

func (r *Redis) Delete(item string) error {
	num, err := r.rdb.SRem(r.ctx, r.setName, item).Result()

	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New(fmt.Sprintf("item %s not in set", item))
	}
	return nil
}

func (r *Redis) FetchRandom(n int) ([]string, error) {

	values, err := r.rdb.SRandMemberN(r.ctx, r.setName, int64(n)).Result()
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return []string{}, errors.New("error fetching random data, got 0 items")
	}

	return values, nil

}

func (r *Redis) GetAll() ([]string, error) {

	values, err := r.rdb.SMembers(r.ctx, r.setName).Result()

	if len(values) == 0 {
		return []string{}, errors.New("error fetching all data, got 0 items")
	}

	if err != nil {
		return []string{}, err
	}

	return values, nil

}
