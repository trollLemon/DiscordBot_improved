package Commands

import (
	"bot/Commands/Audio"
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
