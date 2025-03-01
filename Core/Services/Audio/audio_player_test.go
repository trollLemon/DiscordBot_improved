package audio

import (
	"bot/Core/Services/Audio/Mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)



func TestPlay(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		queueSize  int
		wantQueue  int
		wantNotif  int
		wantStream int
		wantPlay   int
	}{
		{
			name:       "Play with empty queue",
			url:        "http://example.com/testvideo",
			queueSize:  0,
			wantQueue:  0,
			wantNotif:  1,
			wantStream: 1,
			wantPlay:   1,
		},
		{
			name:       "Play with populated queue",
			url:        "http://example.com/testvideo",
			queueSize:  5,
			wantQueue:  0,
			wantNotif:  6,
			wantStream: 6,
			wantPlay:   6,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotif := mock_audio.NewMockNotificationService(ctrl)
	mockVoice := mock_audio.NewMockVoiceService(ctrl)
	mockStream := mock_audio.NewMockStreamService(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNotif.EXPECT().SendNotification(gomock.Any()).Times(tt.wantNotif)
			mockVoice.EXPECT().Disconnect().Times(1)
			mockStream.EXPECT().GetAudioStream(tt.url).Return("stream_url", nil).Times(tt.wantStream)
			mockVoice.EXPECT().PlayAudioFile(gomock.Any(), gomock.Any()).Times(tt.wantPlay)

			player := NewAudioPlayer(mockStream, mockVoice, mockNotif)
			for i := 0; i < tt.queueSize; i++ {
				player.add(tt.url)
			}

			player.Play(tt.url)
			time.Sleep(10*time.Millisecond)
			assert.Equal(t, tt.wantQueue, player.q.Size())

		})
	}
}


func TestStop(t *testing.T) {
	tests := []struct {
		name       string
		queueSize  int
		wantQueue  int
		wantDisconnect bool
	}{
		{
			name:       "Stop with empty queue",
			queueSize:  0,
			wantQueue:  0,
			wantDisconnect: true,
		},
		{
			name:       "Stop with populated queue",
			queueSize:  7,
			wantQueue:  0,
			wantDisconnect: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)
			for i := 0; i < tt.queueSize; i++ {
				player.q.Enque("http://example.com/audio.mp3")
			}
			player.Done = make(chan bool)
			go func() {
				<-player.Done
			}()

			if tt.wantDisconnect {
				mockVoiceService.EXPECT().Disconnect()
			}
			player.Stop()

			assert.Equal(t, tt.wantQueue, player.q.Size())
		})
	}
}



func TestSkip(t *testing.T) {
	tests := []struct {
		name       string
		queueSize  int
		wantErr    bool
	}{
		{
			name:       "Skip with items in queue",
			queueSize:  1,
			wantErr:    false,
		},
		{
			name:       "Skip with empty queue",
			queueSize:  0,
			wantErr:    true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)
			for i := 0; i < tt.queueSize; i++ {
				player.q.Enque("http://example.com/audio.mp3")
			}
			player.Done = make(chan bool)
			go func() {
				<-player.Done
			}()
			player.isPlaying = true

			err := player.Skip()
			if (err != nil) != tt.wantErr {
				t.Errorf("Skip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestShuffle(t *testing.T) {
	tests := []struct {
		name       string
		queueSize  int
		wantErr    bool
	}{
		{
			name:       "Shuffle with empty queue",
			queueSize:  0,
			wantErr:    true,
		},
		{
			name:       "Shuffle with populated queue",
			queueSize:  4,
			wantErr:    false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVoiceService := mock_audio.NewMockVoiceService(ctrl)
	mockNotifService := mock_audio.NewMockNotificationService(ctrl)
	mockStreamService := mock_audio.NewMockStreamService(ctrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewAudioPlayer(mockStreamService, mockVoiceService, mockNotifService)
			for i := 0; i < tt.queueSize; i++ {
				player.q.Enque("https://www.youtube.com/video")
			}

			err := player.Shuffle()
			if (err != nil) != tt.wantErr {
				t.Errorf("Shuffle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
