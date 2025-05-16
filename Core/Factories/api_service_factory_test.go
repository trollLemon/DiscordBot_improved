package factories_test

import (
	"bot/Core/Factories"
	"testing"
)

func TestCreateImageManipService(t *testing.T) {
	tests := []struct {
		name       string
		service    factories.Service
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid manip service",
			service:    factories.GoManip,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid manip service",
			service:    factories.DiscordNotification,
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := factories.CreateImageAPIService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAPIService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("CreateAPIService() got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}
