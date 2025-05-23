package Commands_test

import (
	application "bot/Application"
	"bot/Core/Commands"
	mockinterfaces "bot/Core/Interfaces/Mocks"
	database "bot/Core/Services/Database"
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

type AddCommandTestCase struct {
	name               string
	wantDBError        bool
	wantDuplicateError bool
	userInput          discordgo.ApplicationCommandInteractionDataOption
}

type ShowTestCase struct {
	name             string
	wantDBError      bool
	wantNoItemsError bool
}
type RemoveTestCase struct {
	name                string
	wantDBError         bool
	wantNotPresentError bool
	userInput           discordgo.ApplicationCommandInteractionDataOption
}

func TestAdd(t *testing.T) {
	tests := []AddCommandTestCase{
		{
			name:               "Test add to db",
			wantDBError:        false,
			wantDuplicateError: false,
			userInput: discordgo.ApplicationCommandInteractionDataOption{
				Name:  "term",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "word",
			},
		},
		{
			name:               "Test duplicate error",
			wantDBError:        false,
			wantDuplicateError: true,
			userInput: discordgo.ApplicationCommandInteractionDataOption{
				Name:  "term",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "word",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			applicationData := discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{&tt.userInput},
			}

			s := miniredis.RunT(t)
			defer s.Close()
			rdb := redis.NewClient(&redis.Options{
				Addr: s.Addr(),
			})

			redisClient := database.NewRedisClient(context.Background(), rdb, "set")
			if tt.wantDuplicateError {
				redisClient.Insert(tt.userInput.Value.(string))
			}
			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(applicationData).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     nil,
				WordDatabase: redisClient,
				GuildID:      "guildID",
			}

			err := Commands.Add(mockSession, mockInteractionCreate, mockApplication)

			if !tt.wantDBError || !tt.wantDuplicateError {
				assert.Equal(t, tt.wantDBError || tt.wantDuplicateError, err != nil)
			}

			if tt.wantDuplicateError {
				assert.EqualError(t, err, fmt.Sprintf("item %s already in set", tt.userInput.Value.(string)))
			}

			if tt.wantDBError {
				assert.EqualError(t, err, "database error")
			}
			ctrl.Finish()

		})
	}
}

func TestShow(t *testing.T) {
	tests := []ShowTestCase{
		{
			name:             "Test show words",
			wantDBError:      false,
			wantNoItemsError: false,
		},
		{
			name:             "Test database no items",
			wantDBError:      false,
			wantNoItemsError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			s := miniredis.RunT(t)
			defer s.Close()
			rdb := redis.NewClient(&redis.Options{
				Addr: s.Addr(),
			})

			redisClient := database.NewRedisClient(context.Background(), rdb, "set")

			if !tt.wantNoItemsError {
				redisClient.Insert("hello")
				redisClient.Insert("world")
			}

			mockApplication := &application.Application{
				ImageApi:     nil,
				WordDatabase: redisClient,
				GuildID:      "guildID",
			}

			err := Commands.Show(mockSession, mockInteractionCreate, mockApplication)

			if !tt.wantDBError || !tt.wantNoItemsError {
				assert.Equal(t, tt.wantDBError || tt.wantNoItemsError, err != nil)
			}

			if tt.wantNoItemsError {
				assert.EqualError(t, err, "error fetching all data, got 0 items")
			}

			if tt.wantDBError {
				assert.EqualError(t, err, "database error")
			}
			ctrl.Finish()

		})
	}
}

func TestRemove(t *testing.T) {
	tests := []RemoveTestCase{
		{
			name:                "Test remove from db",
			wantDBError:         false,
			wantNotPresentError: false,
			userInput: discordgo.ApplicationCommandInteractionDataOption{
				Name:  "term",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "word",
			},
		},
		{
			name:                "Test not in database error",
			wantDBError:         false,
			wantNotPresentError: true,
			userInput: discordgo.ApplicationCommandInteractionDataOption{
				Name:  "term",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "word",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			applicationData := discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{&tt.userInput},
			}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(applicationData).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			s := miniredis.RunT(t)
			defer s.Close()
			rdb := redis.NewClient(&redis.Options{
				Addr: s.Addr(),
			})
			redisClient := database.NewRedisClient(context.Background(), rdb, "set")

			if !tt.wantNotPresentError {
				redisClient.Insert(tt.userInput.Value.(string))
			}

			mockApplication := &application.Application{
				ImageApi:     nil,
				WordDatabase: redisClient,
				GuildID:      "guildID",
			}

			err := Commands.Remove(mockSession, mockInteractionCreate, mockApplication)

			if !tt.wantDBError || !tt.wantNotPresentError {
				assert.Equal(t, tt.wantDBError || tt.wantNotPresentError, err != nil)
			}

			if tt.wantNotPresentError {
				assert.EqualError(t, err, fmt.Sprintf("item %s not in set", tt.userInput.Value.(string)))
			}

			if tt.wantDBError {
				assert.EqualError(t, err, "database error")
			}
			ctrl.Finish()

		})
	}
}
