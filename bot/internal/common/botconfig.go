package common

import (
	"fmt"
	"os"
	"strconv"
	"time"
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

func GetBotConfig() (*BotConfig, error) {
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
		return nil, fmt.Errorf("invalid GOMANIP_TIMEOUT=%q: %w", os.Getenv("GOMANIP_TIMEOUT"), err)
	}

	classificationTimeout, err := time.ParseDuration(os.Getenv("CLASSIFICATION_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid CLASSIFICATION_TIMEOUT=%q: %w", os.Getenv("CLASSIFICATION_TIMEOUT"), err)
	}

	conf.GomanipTimeout = gomanipTimeout
	conf.ClassificationTimeout = classificationTimeout

	setNum, err := strconv.Atoi(os.Getenv("REDIS_SET_NUM"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_SET_NUM=%q: %w", os.Getenv("REDIS_SET_NUM"), err)
	}

	conf.RedisSetNumber = setNum

	return &conf, nil
}
