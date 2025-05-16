package audio

import (
	"bot/util"
	"fmt"
	"sync"
)

type AbstractAudioPlayer interface {
	Add(url string)
	Play(url string)
	PlayAudio()
	Stop()
	Wait()
	Skip() error
	Shuffle() error
	UpdateConnection(service VoiceService)
	UpdateNotifier(service NotificationService)
}

type Player struct {
	IsPlaying bool                // check if we are playing audio
	Done      chan bool           // channel to control the audio playing goroutine
	Q         *util.Queue         // queue for audio to play
	Wg        *sync.WaitGroup     // waitGroup for concurrency control
	yts       StreamService       // Injected Dependency for getting stream urls to play audio
	vc        VoiceService        // Injected Dependency for Discord Voice (audio playing features)
	ns        NotificationService // Injected Dependency for player notification/errors via Discord Messaging
}

func NewAudioPlayer(yts StreamService, vc VoiceService, ns NotificationService, isPlaying bool, queue *util.Queue, wg *sync.WaitGroup) *Player {
	return &Player{
		Done:      make(chan bool, 1),
		yts:       yts,
		vc:        vc,
		ns:        ns,
		IsPlaying: isPlaying,
		Q:         queue,
		Wg:        wg,
	}
}

func (player *Player) Add(url string) {
	player.Q.Enque(url)
	player.Wg.Add(1)
}

func (player *Player) PlayAudio() {

	defer func() {
		player.IsPlaying = false
		player.vc.Disconnect()
		close(player.Done)
	}()

	for player.Q.Size() != 0 {

		player.Wg.Done()

		nextUrl, err := player.Q.Deque()

		if err != nil {
			player.ns.SendError(err.Error())
			continue
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

func (player *Player) Play(url string) {

	player.Add(url)
	if !player.IsPlaying {
		player.IsPlaying = true
		go player.PlayAudio()
	}
}

func (player *Player) Stop() {

	player.Done <- true
	player.IsPlaying = false
	player.Q.Clear()
	player.vc.Disconnect()

}

func (player *Player) Skip() error {

	if !player.IsPlaying {

		return fmt.Errorf("cannot skip, bot is not playing audio")
	}

	if player.Q.Size() == 0 {

		return fmt.Errorf("queue is empty, cannot skip")
	}

	player.Done <- true
	return nil
}

func (player *Player) Shuffle() error {
	if player.Q.Size() == 0 {

		return fmt.Errorf("cannot shuffle empty queue")
	}

	player.Q.Shuffle()
	return nil
}

func (player *Player) Wait() {
	player.Wg.Wait()
}

func (player *Player) UpdateNotifier(service NotificationService) {

	player.ns = service
}

func (player *Player) UpdateConnection(con VoiceService) {
	player.vc = con
}
