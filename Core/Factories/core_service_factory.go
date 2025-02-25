package factories

import (
	"bot/Core/Services"
	"bot/Core/Services/Audio"
	"bot/Core/Services/Database"
	"fmt"
	"github.com/kkdai/youtube/v2"
)

type Service int8

const (
	DiscordNotification Service = 0
	DiscordVoice        Service = 1
	YTStream            Service = 2
	Redis               Service = 3
)

func CreateService(service Service) (services.IService, error) {

	switch service {

	case DiscordNotification:
		return &audio.Notifier{}, nil
	case DiscordVoice:
		return &audio.Voice{}, nil
	case YTStream:
		return &audio.YTDL{
			Yt_client: youtube.Client{},
		}, nil

	case Redis:
		return database.NewRedisClient(), nil
	default:
		return nil, fmt.Errorf("Cannot create non-existant service")

	}

}
