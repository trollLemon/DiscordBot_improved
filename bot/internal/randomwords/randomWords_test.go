package randomwords_test

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
		wantErr error
	}{
		{
			name:    "Insert",
			item:    "hello world",
			setName: "set",
			wantErr: nil,
		},
		{
			name:    "Insert (duplicate)",
			item:    "hello world",
			setName: "set",
			wantErr: randomwords.ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()
			redis := randomwords.NewRandomWords(db, testContext, tt.setName)
			
			if tt.wantErr != nil {
				mock.ExpectSAdd(tt.setName, tt.item).SetVal(0)

			} else {
				mock.ExpectSAdd(tt.setName, tt.item).SetVal(1)
			}
			err := redis.Insert(tt.item)

			assert.ErrorIs(t, tt.wantErr, err)

		})
	}
}
func TestRedis_Delete(t *testing.T) {
	tests := []struct {
		name    string
		item    string
		setName string
		wantErr error
	}{
		{
			name:    "Delete",
			item:    "hello world",
			setName: "set",
			wantErr: nil,
		},
		{
			name:    "Delete (not present)",
			item:    "hello world",
			setName: "set",
			wantErr: randomwords.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()

			redis := randomwords.NewRandomWords(db, testContext, tt.setName)

			if tt.wantErr != nil {
				mock.ExpectSRem(tt.setName, tt.item).SetVal(0)

			} else {
				mock.ExpectSRem(tt.setName, tt.item).SetVal(1)
			}
			err := redis.Delete(tt.item)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func TestRedis_FetchRandom(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		numItems int
		setName  string
		wantErr  error
	}{
		{
			name:     "Fetch Random",
			items:    []string{"hello", "world"},
			numItems: 2,
			setName:  "set",
			wantErr:  nil,
		},
		{
			name:     "Fetch Random (no items)",
			items:    []string{},
			numItems: 0,
			setName:  "set",
			wantErr:  randomwords.ErrEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()

			redis := randomwords.NewRandomWords(db, testContext, tt.setName)
			
			mock.ExpectSRandMemberN(tt.setName, int64(tt.numItems)).SetVal(tt.items)

			items, err := redis.GetRandom(tt.numItems)

			assert.ErrorIs(t, tt.wantErr, err)
			assert.Equal(t, tt.items, items)

		})
	}
}

func TestRedis_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		items   []string
		setName string
		wantErr error
	}{
		{
			name:    "Get All",
			items:   []string{"hello", "world"},
			setName: "set",
			wantErr: nil,
		},
		{
			name:    "Get All (no items)",
			items:   []string{},
			setName: "set",
			wantErr: randomwords.ErrEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			testContext := context.Background()

			redis := randomwords.NewRandomWords(db, testContext, tt.setName)
			mock.ExpectSMembers(tt.setName).SetVal(tt.items)

			items, err := redis.GetAll()

			assert.ErrorIs(t, tt.wantErr, err)
			assert.Equal(t, tt.items, items)

		})
	}
}
