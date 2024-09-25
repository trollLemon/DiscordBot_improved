package audio

import (
	"fmt"

	"bot/util"

	"github.com/kkdai/youtube/v2"
)



type AudioPlayer struct {
	Done      chan bool
	vc        VoiceService
	yt_client youtube.Client
	q         *util.Queue
}


func NewAudioPlayer(vc VoiceService) *AudioPlayer {
    return &AudioPlayer{
        Done:      nil,
        vc:        vc,
        yt_client: youtube.Client{},
        q:         util.NewQueue(),
    }
}


func (player *AudioPlayer) add(url string) {
	player.q.Enque(url)
}

func (player *AudioPlayer) playAudio() error {
	for {
		if player.q.Size() == 0 {
			break
		}

		url, _ := player.q.Deque()

		video, err := player.yt_client.GetVideo(url)
		if err != nil {
			continue
		}

		format := video.Formats.WithAudioChannels()[0]
		streamUrl, _ := player.yt_client.GetStreamURL(video, &format)
		player.vc.PlayAudioFile(streamUrl,player.Done)

	}

	return nil

}

func (player *AudioPlayer) SetConnection(con VoiceService) {
	player.vc = con
}

func (player *AudioPlayer) Play(url string) error {
	player.add(url)
	if player.Done == nil {
		newChan := make(chan bool)
		player.Done = newChan
		player.playAudio()
	}
	return nil
}

func (player *AudioPlayer) Stop() error {

	player.q.Clear()
	player.vc.Disconnect()

	//signal to stop playing audio
	if player.Done != nil {
		player.Done <- true
		player.Done = nil

	}

	return nil
}

func (player *AudioPlayer) Skip() error {

	if player.Done == nil {
		return fmt.Errorf("No song is playing, cannot skip")
	}

	if player.q.Size() == 0 {
		return fmt.Errorf("Queue is empty, cannot skip")
	}

	player.Done <- true
	return nil

}
