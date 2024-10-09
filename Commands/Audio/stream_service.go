package audio

import (
	"github.com/kkdai/youtube/v2"
)

type StreamService interface {
	GetAudioStream(url string) (string, error)
}

type YTDL struct {
	Yt_client youtube.Client
}

func (ytdl *YTDL) GetAudioStream(url string) (string, error) {

	video, err := ytdl.Yt_client.GetVideo(url) // get the url to the video source, not the one on the webpage

	if err != nil {
		return "", err
	}

	format := video.Formats.WithAudioChannels()[0]
	streamUrl, err := ytdl.Yt_client.GetStreamURL(video, &format)

	if err != nil {
		return "", err
	}

	return streamUrl, nil

}
