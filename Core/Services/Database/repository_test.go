package database

import (
	"bot/Core/Services/Database/Mocks"
	"fmt"
	"testing"
	"go.uber.org/mock/gomock"
)


func TestAdd(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		isPresent  bool
		wantErr    bool
		expectInsert bool
	}{
		{
			name:       "Add new item",
			input:      "word",
			isPresent:  false,
			wantErr:    false,
			expectInsert: true,
		},
		{
			name:       "Add another new item",
			input:      "aword",
			isPresent:  false,
			wantErr:    false,
			expectInsert: true,
		},
		{
			name:       "Add duplicate item",
			input:      "word",
			isPresent:  true,
			wantErr:    true,
			expectInsert: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mock_database.NewMockDatabaseService(ctrl)
	repo := NewRepository(mockDatabase)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDatabase.EXPECT().IsPresent(tt.input).Return(tt.isPresent).Times(1)
			if tt.expectInsert {
				mockDatabase.EXPECT().Insert(tt.input).Times(1)
			}

			err := repo.Add(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}



func TestRemove(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
	}{
		{
			name:       "Remove existing item",
			input:      "word",
			wantErr:    false,
		},
		{
			name:       "Remove another existing item",
			input:      "aword",
			wantErr:    false,
		},
		{
			name:       "Remove with error",
			input:      "word",
			wantErr:    true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mock_database.NewMockDatabaseService(ctrl)
	repo := NewRepository(mockDatabase)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mockDatabase.EXPECT().Delete(tt.input).Return(fmt.Errorf("Error removing data")).Times(1)
			} else {
				mockDatabase.EXPECT().Delete(tt.input).Times(1)
			}

			err := repo.Remove(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRandN(t *testing.T) {
	tests := []struct {
		name       string
		n          int
		wantData   bool
		wantErr    bool
		mockData   []string
		mockErr    error
	}{
		{
			name:       "Fetch random data with error",
			n:          1,
			wantData:   false,
			wantErr:    true,
			mockData:   nil,
			mockErr:    fmt.Errorf("Could not get data from database"),
		},
		{
			name:       "Fetch random data successfully",
			n:          1,
			wantData:   true,
			wantErr:    false,
			mockData:   []string{"url1"},
			mockErr:    nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mock_database.NewMockDatabaseService(ctrl)
	repo := NewRepository(mockDatabase)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDatabase.EXPECT().FetchRandom(tt.n).Return(tt.mockData, tt.mockErr).Times(1)

			data, err := repo.GetRandN(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (data != nil) != tt.wantData {
				t.Errorf("GetRandN() data = %v, wantData %v", data, tt.wantData)
			}
		})
	}
}
