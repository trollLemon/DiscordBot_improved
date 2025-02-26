package audio

import (
	"bot/util"
	"fmt"
)

type AudioPlayer struct {
	Done      chan bool           // channel to control the audio playing goroutine
	q         *util.Queue         // queue for audio to play
	yts       StreamService       // Injected Dependency for getting stream urls to play audio
	vc        VoiceService        // Injected Dependency for Discord Voice (audio playing features)
	ns        NotificationService // Injected Dependency for player notification/errors via Discord Messaging
	isPlaying bool                // check if we are playing audio
}

func NewAudioPlayer(yts StreamService, vc VoiceService, ns NotificationService) *AudioPlayer {
	return &AudioPlayer{
		Done:      nil,
		yts:       yts,
		q:         util.NewQueue(),
		vc:        vc,
		ns:        ns,
		isPlaying: false,
	}
}

func (player *AudioPlayer) add(url string) {
	player.q.Enque(url)
}

func (player *AudioPlayer) playAudio() {

	defer func() {
		player.isPlaying = false
		player.vc.Disconnect()
		close(player.Done)
	}()

	for player.q.Size() != 0 {

		nextUrl, err := player.q.Deque()

		if err != nil {
			println(err.Error())
		}

		player.ns.SendNotification(nextUrl)

		streamUrl, err := player.yts.GetAudioStream(nextUrl)
		if err != nil {
			player.ns.SendError(err.Error())
			continue
		}
		player.vc.PlayAudioFile(streamUrl, player.Done)

	}

}

func (player *AudioPlayer) Play(url string) {

	player.add(url)

	if !player.isPlaying {
		player.isPlaying = true

		player.Done = make(chan bool)

		go player.playAudio()
	}
}

func (player *AudioPlayer) Stop() {

	player.Done <- true
	player.isPlaying = false
	player.q.Clear()
	player.vc.Disconnect()

}

func (player *AudioPlayer) Skip() error {

	if !player.isPlaying {

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

func (player *AudioPlayer) UpdateNotifier(service NotificationService) {

	player.ns = service
}

func (player *AudioPlayer) SetConnection(con VoiceService) {
	player.vc = con
}
