package factories

import (
	"testing"

)

func TestCreateStreamService(t *testing.T) {
	tests := []struct {
		name       string
		service    Service
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid YTStream service",
			service:    YTStream,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid stream service",
			service:    DiscordNotification,
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStreamService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStreamService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("CreateStreamService() got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}

func TestCreateVoiceService(t *testing.T) {
	tests := []struct {
		name       string
		service    Service
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid DiscordVoice service",
			service:    DiscordVoice,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid voice service",
			service:    YTStream,
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateVoiceService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVoiceService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("CreateVoiceService() got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}

func TestCreateNotificationService(t *testing.T) {
	tests := []struct {
		name       string
		service    Service
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid DiscordNotification service",
			service:    DiscordNotification,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid notification service",
			service:    Redis,
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateNotificationService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNotificationService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("CreateNotificationService() got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}

func TestCreateDatabaseService(t *testing.T) {
	tests := []struct {
		name       string
		service    Service
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid Redis service",
			service:    Redis,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid database service",
			service:    DiscordVoice,
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDatabaseService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabaseService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("CreateDatabaseService() got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}
