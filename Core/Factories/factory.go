package factories

import (
	"bot/Core/Services/Audio"
	"bot/Core/Services/Database"
	"github.com/kkdai/youtube/v2"
)

func CreateNotificationService() *audio.Notifier {

	return &audio.Notifier{}

}
func CreateVoiceService() *audio.Voice {

	return &audio.Voice{}
}

func CreateStreamService() *audio.YTDL {

	return &audio.YTDL{
		Yt_client: youtube.Client{},
	}
}

func CreateDatabaseService() *database.Redis {
	return database.NewRedisClient()
}
