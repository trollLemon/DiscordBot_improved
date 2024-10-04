package audio

import (
	"bot/util"
	"fmt"

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

func (player *AudioPlayer) playAudio() {

	for {
		if player.q.Size() == 0 {
			break;
		}

		nextUrl, _ := player.q.Deque()
		video, err := player.yt_client.GetVideo(nextUrl) // get the url to the video source, not the one on the webpage
		
		if err != nil {
			//logger
			continue
		}

		format := video.Formats.WithAudioChannels()[0]
		streamUrl, err := player.yt_client.GetStreamURL(video, &format)
		
		if err != nil {
			continue
		}

		player.vc.PlayAudioFile(streamUrl,player.Done)


	}
}

func (player *AudioPlayer) SetConnection(con VoiceService) {
	player.vc = con
}

func (player *AudioPlayer) isPlaying() bool {
	return player.Done != nil
}

func (player *AudioPlayer) Play(url string) {
	
	player.add(url)
	if !player.isPlaying() {
		player.Done = make(chan bool)
		player.playAudio()
	}
}





func (player *AudioPlayer) Stop() {
	player.Done <-true
	player.Done = nil
	player.q.Clear()
	player.vc.Disconnect()
	
	
}

func (player *AudioPlayer) Skip() error {
	
	if player.Done == nil {
		return fmt.Errorf("Cannot skip, bot is not playing audio")
	}

	if player.q.Size() == 0 {
		return fmt.Errorf("Queue is empty, cannot skip")
	}
	
	player.Done <- true
	return nil
}

func (player *AudioPlayer) Shuffle() error {
	if player.q.Size() == 0 {
		
		return fmt.Errorf("Cannot shuffle empty queue")
	}
	
	player.q.Shuffle()
	return nil
}

