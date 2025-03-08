package factories

import (
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
	YTDLP		    Service = 5
)

func CreateStreamService(service Service) (audio.StreamService, error) {

	switch service {
	case YTStream:
		return &audio.YTDL{
			Yt_client: youtube.Client{},
		}, nil
	case YTDLP:
		return &audio.YtDLP{}, nil

	default:
		return nil, fmt.Errorf("Invalid stream service type")
	}
}

func CreateVoiceService(service Service) (audio.VoiceService, error) {

	switch service {
	case DiscordVoice:
		return &audio.Voice{}, nil

	default:
		return nil, fmt.Errorf("Invalid voice service type")
	}

}

func CreateNotificationService(service Service) (audio.NotificationService, error) {

	switch service {
	case DiscordNotification:
		return &audio.Notifier{}, nil

	default:
		return nil, fmt.Errorf("Invalid notification service type")
	}

}

func CreateDatabaseService(service Service) (database.DatabaseService, error) {

	switch service {
	case Redis:
		return database.NewRedisClient(), nil

	default:
		return nil, fmt.Errorf("Invalid database service type")
	}

}
