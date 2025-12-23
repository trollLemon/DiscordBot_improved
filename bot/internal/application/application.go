package application

import (

	"context"

	"github.com/redis/go-redis/v9"

	"github.com/trollLemon/DiscordBot/internal/classification"
	"github.com/trollLemon/DiscordBot/internal/gomanip"
	"github.com/trollLemon/DiscordBot/internal/randomwords"
	"github.com/trollLemon/DiscordBot/internal/common"
)

type Application struct {
	Gomanip        *gomanip.GoManip
	Classification *Classification.ImageClassification
	RandomWords    *randomwords.RandomWords
	GuildID        string
}



func InitializeApplication(conf *common.BotConfig, ctx context.Context) *Application {

	gomanip := gomanip.NewGoManip(conf.GomanipURL, conf.GomanipTimeout)

	classifier := Classification.NewImageClassification(conf.ClassificationTimeout, conf.ClassificationURL, Classification.SendImageEndpoint, Classification.GetClassificationEndpoint)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.DbDSN,
		Password: conf.RedisPass,
		DB:       conf.RedisSetNumber,
	})

	randomWordsStore := randomwords.NewRandomWords(redisClient, ctx, conf.RandomWordsSetName)


	return &Application{
		Gomanip:        gomanip,
		Classification: classifier,
		RandomWords:    randomWordsStore,
	}
}
