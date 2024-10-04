package audio

import (
	"bot/Commands/Audio/Mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)


func TestPlay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	

	player := NewAudioPlayer(mockVoiceService)
	

	//TODO: use acutal youtube links for testing

	url := "https://www.youtube.com/video1"
	url2 := "https://www.youtube.com/video2"
	url3 := "https://www.youtube.com/video3"
	url4 := "https://www.youtube.com/video4"

	go func() {
		<-player.Done
	}()

	player.Play(url)

	player.Done <- true
	//player should pop the stack when done playing
	assert.Equal(t, 0, player.q.Size())
	player.Done = nil
	

	player.q.Enque(url3)
	player.q.Enque(url4)



	player.Play(url2)
	//at this point the player has played all urls in the queue, so it should be empty
	assert.Equal(t,0,player.q.Size())

}

func TestStop(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockVoiceService.EXPECT().Disconnect()
	player := NewAudioPlayer(mockVoiceService)

	url := "http://example.com/audio.mp3"

	//simulate audio playing, and a song in the queue

	player.q.Enque(url)
	player.q.Enque(url)
	player.q.Enque(url)
	player.q.Enque(url)
	player.q.Enque(url)
	player.q.Enque(url)
	player.q.Enque(url)
	player.Done = make(chan bool)

	go func() {
		//when audio is playing, there is a channel with a reciever,
		//in this case for testing its just a simple function that exits when signaled
		<-player.Done
	}()
	player.Stop()

	assert.Equal(t, 0, player.q.Size())
	assert.Nil(t, player.Done)

}

func TestSkip(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	player := NewAudioPlayer(mockVoiceService)

	url := "http://example.com/audio.mp3"

	player.q.Enque(url)
	player.Done = make(chan bool)
	go func() {
		<-player.Done
	}()
	err := player.Skip()
	player.q.Deque()
	assert.Nil(t, err)
	go func() {
		<-player.Done
	}()
	err = player.Skip()
	assert.NotNil(t, err)

}



func TestShuffle(t *testing.T) {
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	player := NewAudioPlayer(mockVoiceService)
	
	//cannot shuffle if queue is empty
	err := player.Shuffle()
	assert.NotNil(t,err)
	url := "https://www.youtube.com/video1"
	url2 := "https://www.youtube.com/video2"
	url3 := "https://www.youtube.com/video3"
	url4 := "https://www.youtube.com/video4"	
	
	player.q.Enque(url)
	player.q.Enque(url2)
	player.q.Enque(url3)
	player.q.Enque(url4)

	err = player.Shuffle()
	assert.Nil(t,err)

}
