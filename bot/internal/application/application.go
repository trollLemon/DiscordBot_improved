package application

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"github.com/trollLemon/DiscordBot/internal/classification"
	"github.com/trollLemon/DiscordBot/internal/common"
	"github.com/trollLemon/DiscordBot/internal/gomanip"
	"github.com/trollLemon/DiscordBot/internal/randomwords"
)

type Application struct {
	Gomanip        *gomanip.GoManip
	Classification *classification.ImageClassification
	RandomWords    *randomwords.RandomWords
	GuildID        string
}

func InitializeApplication(conf *common.BotConfig, ctx context.Context) *Application {

	gomanip := gomanip.NewGoManip(conf.GomanipURL, conf.GomanipTimeout)

	classifier := classification.NewImageClassification(conf.ClassificationTimeout, conf.ClassificationURL, classification.SendImageEndpoint, classification.GetClassificationEndpoint)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.DbDSN,
		Password: conf.RedisPass,
		DB:       conf.RedisSetNumber,
	})

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		slog.Error("failed to instrument redis client with tracing", "error", err)
	}

	randomWordsStore := randomwords.NewRandomWords(redisClient, ctx, conf.RandomWordsSetName)

	return &Application{
		Gomanip:        gomanip,
		Classification: classifier,
		RandomWords:    randomWordsStore,
	}
}
