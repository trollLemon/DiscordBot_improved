package factories

import (
	database "bot/Core/Services/Database"
	"bot/Core/Services/ImageManip"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type ImageApiService int
type DatabaseService int

const (
	GoManip            ImageApiService = 0
	Redis              DatabaseService = 1
	ImageManipEndpoint                 = "http://image:8080/api/image"
)

func CreateImageAPIService(service ImageApiService) (imagemanip.AbstractImageAPI, error) {

	switch service {

	case GoManip:
		return imagemanip.NewGoManip(&http.Client{}, ImageManipEndpoint), nil

	default:
		return nil, fmt.Errorf("invalid image manipulation service")

	}

}

func CreateDatabaseService(service DatabaseService) (database.AbstractDatabaseService, error) {
	switch service {
	case Redis:
		ctx := context.WithoutCancel(context.Background())
		setName := "items"

		rdb := redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		return database.NewRedisClient(ctx, rdb, setName), nil
	default:
		return nil, fmt.Errorf("invalid database service")
	}

}
