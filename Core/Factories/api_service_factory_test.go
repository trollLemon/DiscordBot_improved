package factories_test

import (
	"bot/Core/Factories"
	"testing"
)

func TestCreateImageManipService(t *testing.T) {
	tests := []struct {
		name       string
		service    factories.ImageApiService
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
			service:    -1, //simulate invalid service
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := factories.CreateImageAPIService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create Service error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("Create Service got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}

func TestCreateDatabaseService(t *testing.T) {
	tests := []struct {
		name       string
		service    factories.DatabaseService
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid database service",
			service:    factories.Redis,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid database service",
			service:    -1, //simulate invalid service
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := factories.CreateDatabaseService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create Service error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("Create Service got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}

func TestCreateClassificationService(t *testing.T) {
	tests := []struct {
		name       string
		service    factories.ClassificationService
		wantErr    bool
		wantNotNil bool
	}{
		{
			name:       "Valid classification service",
			service:    factories.VitClassification,
			wantErr:    false,
			wantNotNil: true,
		},
		{
			name:       "Invalid classification service",
			service:    -1, //simulate invalid service
			wantErr:    true,
			wantNotNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := factories.CreateClassificationAPIService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create Service error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantNotNil {
				t.Errorf("Create Service got = %v, wantNotNil %v", got, tt.wantNotNil)
			}
		})
	}
}
