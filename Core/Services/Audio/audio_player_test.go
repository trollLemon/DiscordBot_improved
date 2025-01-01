package audio

import (
	"bot/Services/Audio/Mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestPlay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)
	player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)

	mockNotifService.EXPECT().SendNotification("Now playing http://example.com/testvideo").Times(6)
	mockVoiceService.EXPECT().Disconnect().Times(2)
	testUrl := "http://example.com/testvideo"

	mockStreamService.EXPECT().GetAudioStream(testUrl).Times(6)
	mockVoiceService.EXPECT().PlayAudioFile(gomock.Any(), gomock.Any()).Times(6)
	// Start playing the audio in a separate goroutine
	player.Play(testUrl)
	time.Sleep(100 * time.Millisecond)
	
	assert.Equal(t,player.q.Size(), 0)

	player.add(testUrl)
	player.add(testUrl)
	player.add(testUrl)
	player.add(testUrl)
	player.add(testUrl)

	player.playAudio()

	close(player.Done)

}

func TestStop(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)
	player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)

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

	mockVoiceService.EXPECT().Disconnect()
	player.Stop()

	assert.Equal(t, 0, player.q.Size())
	assert.Nil(t, player.Done)

}

func TestSkip(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)
	player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)

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
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)
	player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)

	//cannot shuffle if queue is empty
	err := player.Shuffle()
	assert.NotNil(t, err)
	url := "https://www.youtube.com/video1"
	url2 := "https://www.youtube.com/video2"
	url3 := "https://www.youtube.com/video3"
	url4 := "https://www.youtube.com/video4"

	player.q.Enque(url)
	player.q.Enque(url2)
	player.q.Enque(url3)
	player.q.Enque(url4)

	err = player.Shuffle()
	assert.Nil(t, err)

}
