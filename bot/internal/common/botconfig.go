package common

import (
	"os"
	"time"
	"strconv"

	"github.com/rs/zerolog/log"
)

type BotConfig struct {
	BotToken  string
	RedisPass string

	GomanipURL        string
	ClassificationURL string
	DbDSN             string

	GomanipTimeout        time.Duration
	ClassificationTimeout time.Duration

	RandomWordsSetName string
	RedisSetNumber     int
}






func GetBotConfig() *BotConfig {
	conf := BotConfig{
		BotToken:           os.Getenv("DISCORD_TOKEN"),
		RedisPass:          os.Getenv("REDIS_PASS"),
		GomanipURL:         os.Getenv("GOMANIP_URL"),
		ClassificationURL:  os.Getenv("CLASSIFICATION_URL"),
		DbDSN:              os.Getenv("DATABASE_DSN"),
		RandomWordsSetName: os.Getenv("RANDOM_WORD_SET_NAME"),
	}

	gomanipTimeout, err := time.ParseDuration(os.Getenv("GOMANIP_TIMEOUT"))
	if err != nil {
		log.Error().Msgf("failed to parse duration in provided variable GOMANIP_TIMEOUT=%s. Using default value of 30 seconds.", os.Getenv("GOMANIP_TIMEOUT"))
		gomanipTimeout = time.Second * 30
	}

	classificationTimeout, err := time.ParseDuration(os.Getenv("CLASSIFICATION_TIMEOUT"))
	if err != nil {
		log.Error().Msgf("failed to parse duration in provided variable CLASSIFICATION_TIMEOUT=%s. Using default value of 5 minutes.", os.Getenv("CLASSIFICATION_TIMEOUT"))
		gomanipTimeout = time.Minute * 5
	}

	conf.GomanipTimeout = gomanipTimeout
	conf.ClassificationTimeout = classificationTimeout

	setNum, err := strconv.Atoi(os.Getenv("REDIS_SET_NUM"))
	if err != nil {
		log.Error().Msgf("failed to parse string in provided variable REDIS_SET_NUM=%s. Using default value of 0.", os.Getenv("REDIS_SET_NUM"))
	}

	conf.RedisSetNumber = setNum

	return &conf
}
