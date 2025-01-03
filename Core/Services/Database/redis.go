package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	ctx      context.Context
	rdb      *redis.Client
	set_name string
}

func NewRedisClient() *Redis {

	return &Redis{
		ctx: context.TODO(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		set_name: "items",
	}
}

func (r *Redis) Insert(item string) error {
	err := r.rdb.SAdd(r.ctx, r.set_name, item).Err()
	return err
}

func (r *Redis) Delete(item string) error {
	res, err := r.rdb.SRem(r.ctx, r.set_name, item).Result()
	if err != nil {
		return fmt.Errorf("Error deleting from Redis: %s", err.Error())
	}

	if res == 0 {
		return fmt.Errorf("Item is not present in Redis")
	}

	return nil
}

func (r *Redis) FetchRandom(n int) ([]string, error) {


	values, err := r.rdb.SRandMemberN(r.ctx, r.set_name, int64(n)).Result()
	if err != nil {
		return nil, err
	}

	return values, nil

}

func (r *Redis) IsPresent(item string) bool {

	exists, err := r.rdb.SIsMember(r.ctx, r.set_name, item).Result()

	if err != nil {
		return false //todo: better error checking
	}

	return exists
}

func (r *Redis) GetAll() ([]string, error) {

	values, err := r.rdb.SMembers(r.ctx, r.set_name).Result()

	return values, err

}
