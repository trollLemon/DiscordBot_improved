package store_test

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/trollLemon/DiscordBot/internal/randomwords"
)

func TestRedis_Insert(t *testing.T) {

	tests := []struct {
		name    string
		item    string
		setName string
		wantErr bool
	}{
		{
			name:    "Insert",
			item:    "hello world",
			setName: "set",
			wantErr: false,
		},
		{
			name:    "Insert (duplicate)",
			item:    "hello world",
			setName: "set",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()
			redis := store.NewRedisClient(testContext, db, tt.setName)

			if tt.wantErr {
				mock.ExpectSAdd(tt.setName, tt.item).SetVal(0)

			} else {
				mock.ExpectSAdd(tt.setName, tt.item).SetVal(1)
			}
			err := redis.Insert(tt.item)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}
func TestRedis_Delete(t *testing.T) {

	tests := []struct {
		name    string
		item    string
		setName string
		wantErr bool
	}{
		{
			name:    "Delete",
			item:    "hello world",
			setName: "set",
			wantErr: false,
		},
		{
			name:    "Delete (not present)",
			item:    "hello world",
			setName: "set",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()
			redis := store.NewRedisClient(testContext, db, tt.setName)

			if tt.wantErr {
				mock.ExpectSRem(tt.setName, tt.item).SetVal(0)

			} else {
				mock.ExpectSRem(tt.setName, tt.item).SetVal(1)
			}
			err := redis.Delete(tt.item)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestRedis_FetchRandom(t *testing.T) {

	tests := []struct {
		name     string
		items    []string
		numItems int
		setName  string
		wantErr  bool
	}{
		{
			name:     "Fetch Random",
			items:    []string{"hello", "world"},
			numItems: 2,
			setName:  "set",
			wantErr:  false,
		},
		{
			name:     "Fetch Random (no items)",
			items:    []string{},
			numItems: 0,
			setName:  "set",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()

			redis := store.NewRedisClient(testContext, db, tt.setName)
			mock.ExpectSRandMemberN(tt.setName, int64(tt.numItems)).SetVal(tt.items)

			items, err := redis.FetchRandom(tt.numItems)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.items, items)

		})
	}
}

func TestRedis_GetAll(t *testing.T) {

	tests := []struct {
		name    string
		items   []string
		setName string
		wantErr bool
	}{
		{
			name:    "Get All",
			items:   []string{"hello", "world"},
			setName: "set",
			wantErr: false,
		},
		{
			name:    "Get All (no items)",
			items:   []string{},
			setName: "set",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()

			redis := store.NewRedisClient(testContext, db, tt.setName)
			mock.ExpectSMembers(tt.setName).SetVal(tt.items)

			items, err := redis.GetAll()

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.items, items)

		})
	}
}
