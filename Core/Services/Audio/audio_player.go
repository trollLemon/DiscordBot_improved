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

	}()

	for {
		if player.q.Size() == 0 {

			break
		}

		nextUrl, _ := player.q.Deque()
		player.ns.SendNotification("Now playing " + nextUrl)

		streamUrl, err := player.yts.GetAudioStream(nextUrl)

		if err != nil {
			player.ns.SendError(err.Error())
			continue
		}
		player.Done = make(chan bool)
		player.vc.PlayAudioFile(streamUrl, player.Done)
	}

}

func (player *AudioPlayer) Play(url string) {

	player.add(url)

	if !player.isPlaying {
		player.isPlaying = true
		go func() {
			player.playAudio()
			// Only send if the channel is open
			select {
			case player.Done <- true:
			default:
			}
		}()
	}
}

func (player *AudioPlayer) Stop() {

	player.Done <- true
	player.Done = nil
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
