// TODO: COVER SEARCH ERROR FOR RANDOMPLAY
package Commands_test

import (
	application "bot/Application"
	"bot/Core/Commands"
	"bot/Core/Interfaces"
	mockinterfaces "bot/Core/Interfaces/Mocks"
	audio "bot/Core/Services/Audio"
	mockaudio "bot/Core/Services/Audio/Mocks"
	mockdatabase "bot/Core/Services/Database/Mocks"
	"bot/util"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"sync"
	"testing"
)

// stubs */
type SearchStub struct {
	val string
	err error
}

func (s *SearchStub) SearchWithQuery(string) (string, error) {
	return s.val, s.err
}

type VoiceConnectionStub struct{}

func (v VoiceConnectionStub) PlayAudioFile(string, chan bool) {
	fmt.Println("PlayAudioFile")

}

func (v VoiceConnectionStub) Disconnect() {
	fmt.Println("Disconnect")
}

type NotificationStub struct{}

func (n NotificationStub) SendNotification(content string) {
	fmt.Println("SendNotification", content)
}

func (n NotificationStub) SendError(error string) {
	fmt.Println("SendError", error)
}

type FakeServiceFactory struct{}

func (f FakeServiceFactory) CreateVoiceService(*discordgo.VoiceConnection) audio.VoiceService {
	return &VoiceConnectionStub{}
}

func (f FakeServiceFactory) CreateNotificationService(Interfaces.DiscordSession, string) audio.NotificationService {
	return &NotificationStub{}
}

// test cases */
type PlayCommandTestCase struct {
	name                string
	userInput           *discordgo.ApplicationCommandInteractionDataOption
	wantSearchError     bool
	wantVoiceStateError bool
	wantJoinVCError     bool
	wantNotInVCError    bool
	guildId             string
	channelId           string
	author              string
	voiceState          *discordgo.VoiceState
	SearchService       *SearchStub
}

type RandomPlayCommandTestCase struct {
	name                string
	numTerms            int
	wantDatabaseError   bool
	wantSearchError     bool
	wantVoiceStateError bool
	wantJoinVCError     bool
	wantNotInVCError    bool
	guildId             string
	channelId           string
	author              string
	SearchService       *SearchStub
	voiceState          *discordgo.VoiceState
	DatabaseResult      []string
}

type StopCommandTestCase struct {
	name                string
	wantVoiceStateError bool
	wantNotInVCError    bool
	author              string
	channelId           string
	guildId             string
	voiceState          *discordgo.VoiceState
}

type SkipCommandTestCase struct {
	queueSize           int
	name                string
	wantVoiceStateError bool
	wantNotInVCError    bool
	wantSkipError       bool
	isPlaying           bool
	author              string
	channelId           string
	guildId             string
	voiceState          *discordgo.VoiceState
}

type ShuffleCommandTestCase struct {
	queueSize           int
	name                string
	wantVoiceStateError bool
	wantNotInVCError    bool
	wantShuffleError    bool
	author              string
	channelId           string
	guildId             string
	voiceState          *discordgo.VoiceState
}

// tests */
func TestPlay(t *testing.T) {

	tests := []PlayCommandTestCase{
		{
			name: "Test play successful",
			userInput: &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "audio",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "https://video.link.com",
			},
			wantSearchError:     false,
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantJoinVCError:     false,
			guildId:             "guildId",
			channelId:           "channelId",
			author:              "author",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: nil,
		},
		{
			name: "Test play successful non url input",
			userInput: &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "audio",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "funny video",
			},
			wantSearchError:     false,
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantJoinVCError:     false,
			guildId:             "guildId",
			channelId:           "channelId",
			author:              "author",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: &SearchStub{
				val: "https://result.video.link.com",
				err: nil,
			},
		},
		{
			name: "Test play successful non url input search error",
			userInput: &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "audio",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "funny video",
			},
			wantSearchError:     true,
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantJoinVCError:     false,
			guildId:             "guildId",
			channelId:           "channelId",
			author:              "author",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: &SearchStub{
				val: "",
				err: errors.New("search error"),
			},
		},
		{
			name: "Test user not in vc",
			userInput: &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "audio",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "funny video",
			},
			wantSearchError:     false,
			wantVoiceStateError: false,
			wantNotInVCError:    true,
			wantJoinVCError:     false,
			guildId:             "guildId",
			channelId:           "channelId",
			author:              "author",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "", // if this is empty the user is not in a vc
			},
			SearchService: &SearchStub{
				val: "https://result.video.link.com",
				err: nil,
			},
		},
		{
			name: "Test Voice State Error",
			userInput: &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "audio",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "funny video",
			},
			wantSearchError:     false,
			wantVoiceStateError: true,
			wantNotInVCError:    false,
			wantJoinVCError:     false,
			guildId:             "guildId",
			channelId:           "channelId",
			author:              "author",
			voiceState:          nil,
			SearchService:       nil,
		},
		{
			name: "Test Join VC Error",
			userInput: &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "audio",
				Type:  discordgo.ApplicationCommandOptionString,
				Value: "funny video",
			},
			wantSearchError:     false,
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantJoinVCError:     true,
			guildId:             "guildId",
			channelId:           "channelId",
			author:              "author",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: &SearchStub{
				val: "https://result.video.link.com",
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			applicationData := &discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{tt.userInput},
			}

			mockStream := mockaudio.NewMockStreamService(ctrl)
			mockVoice := mockaudio.NewMockVoiceService(ctrl)
			mockNotif := mockaudio.NewMockNotificationService(ctrl)

			serviceFactory := &FakeServiceFactory{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(applicationData).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockInteractionCreate.EXPECT().GetInteractionAuthor().Times(1).Return(tt.author)
			mockInteractionCreate.EXPECT().GetChannel().Times(1).Return(tt.channelId)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			// flow specific
			if tt.wantVoiceStateError {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(nil, errors.New("voice State error")).Times(1)
			} else {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(tt.voiceState, nil).Times(1)
			}

			if tt.wantJoinVCError {
				mockSession.EXPECT().ChannelVoiceJoin(tt.guildId, tt.voiceState.ChannelID, false, false).Return(nil, errors.New("join vc error")).Times(1)
			}

			if !(tt.wantSearchError || tt.wantVoiceStateError || tt.wantNotInVCError || tt.wantJoinVCError) {
				mockSession.EXPECT().ChannelVoiceJoin(tt.guildId, tt.voiceState.ChannelID, false, false).Times(1).Return(&discordgo.VoiceConnection{}, nil)
			}

			if !(tt.wantJoinVCError || tt.wantVoiceStateError || tt.wantSearchError || tt.wantNotInVCError) {
				mockStream.EXPECT().GetAudioStream(gomock.Any()).Times(1).Return("stream url", nil)
			}

			audioPlayer := audio.NewAudioPlayer(mockStream, mockVoice, mockNotif, false, util.NewQueue(), &sync.WaitGroup{})

			app := &application.Application{
				ImageApi:       nil,
				WordDatabase:   nil,
				GuildID:        tt.guildId,
				Search:         tt.SearchService,
				ServiceFactory: serviceFactory,
				AudioPlayer:    audioPlayer,
			}

			err := Commands.Play(mockSession, mockInteractionCreate, app)
			audioPlayer.Wait()
			assert.Equal(t, tt.wantSearchError || tt.wantJoinVCError || tt.wantVoiceStateError || tt.wantNotInVCError, err != nil)
			ctrl.Finish()
		})
	}

}

func TestStop(t *testing.T) {

	tests := []StopCommandTestCase{
		{
			name:                "Test Stop playing audio",
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			author:              "author",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
		},
		{
			name:                "Test Stop playing audio User not in vc",
			wantVoiceStateError: false,
			wantNotInVCError:    true,
			author:              "author",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "",
			},
		},
		{
			name:                "Test Stop playing audio Voice State Error",
			wantVoiceStateError: true,
			wantNotInVCError:    false,
			author:              "author",
			channelId:           "channelId",
			voiceState:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			mockStream := mockaudio.NewMockStreamService(ctrl)
			mockVoice := mockaudio.NewMockVoiceService(ctrl)
			mockNotif := mockaudio.NewMockNotificationService(ctrl)

			serviceFactory := &FakeServiceFactory{}
			audioPlayer := audio.NewAudioPlayer(mockStream, mockVoice, mockNotif, false, util.NewQueue(), &sync.WaitGroup{})

			app := &application.Application{
				ImageApi:       nil,
				WordDatabase:   nil,
				GuildID:        tt.guildId,
				Search:         nil,
				ServiceFactory: serviceFactory,
				AudioPlayer:    audioPlayer,
			}
			//common expectations
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockInteractionCreate.EXPECT().GetInteractionAuthor().Times(1).Return(tt.author)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			//flow specific
			if tt.wantVoiceStateError {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(nil, errors.New("voice State error")).Times(1)
			} else {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(tt.voiceState, nil).Times(1)
			}

			if !(tt.wantNotInVCError || tt.wantVoiceStateError) {
				mockVoice.EXPECT().Disconnect().Times(1)
			}

			err := Commands.Stop(mockSession, mockInteractionCreate, app)

			assert.Equal(t, tt.wantVoiceStateError || tt.wantNotInVCError, err != nil)
			ctrl.Finish()
		})
	}

}

func TestSkip(t *testing.T) {

	tests := []SkipCommandTestCase{
		{
			queueSize:           3,
			name:                "Test Skip playing audio",
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantSkipError:       false,
			isPlaying:           true,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
		},
		{
			queueSize:           3,
			name:                "Test Skip not playing audio",
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantSkipError:       true,
			isPlaying:           false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
		},
		{
			queueSize:           0,
			name:                "Test Skip queue empty",
			wantVoiceStateError: false,
			wantNotInVCError:    false,
			wantSkipError:       true,
			isPlaying:           true,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
		},
		{
			queueSize:           4,
			name:                "Test Skip user not in vc",
			wantVoiceStateError: false,
			wantNotInVCError:    true,
			wantSkipError:       false,
			isPlaying:           true,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "",
			},
		},
		{
			queueSize:           3,
			name:                "Test Skip Voice State Error",
			wantVoiceStateError: true,
			wantNotInVCError:    false,
			wantSkipError:       false,
			isPlaying:           true,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			mockStream := mockaudio.NewMockStreamService(ctrl)
			mockVoice := mockaudio.NewMockVoiceService(ctrl)
			mockNotif := mockaudio.NewMockNotificationService(ctrl)

			serviceFactory := &FakeServiceFactory{}

			queue := util.NewQueue()

			for _ = range tt.queueSize {
				queue.Enque("http://video/url.com")
			}

			audioPlayer := audio.NewAudioPlayer(mockStream, mockVoice, mockNotif, tt.isPlaying, queue, &sync.WaitGroup{})

			app := &application.Application{
				ImageApi:       nil,
				WordDatabase:   nil,
				GuildID:        tt.guildId,
				Search:         nil,
				ServiceFactory: serviceFactory,
				AudioPlayer:    audioPlayer,
			}
			//common expectations
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockInteractionCreate.EXPECT().GetInteractionAuthor().Times(1).Return(tt.author)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			//flow specific
			if tt.wantVoiceStateError {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(nil, errors.New("voice State error")).Times(1)
			} else {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(tt.voiceState, nil).Times(1)
			}

			err := Commands.Skip(mockSession, mockInteractionCreate, app)
			assert.Equal(t, tt.wantVoiceStateError || tt.wantNotInVCError || tt.wantSkipError, err != nil)
			ctrl.Finish()
		})
	}

}

func TestShuffle(t *testing.T) {

	tests := []ShuffleCommandTestCase{
		{
			name:                "Test Shuffle",
			queueSize:           3,
			wantShuffleError:    false,
			wantNotInVCError:    false,
			wantVoiceStateError: false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
		},
		{
			name:                "Test Shuffle empty queue",
			queueSize:           0,
			wantShuffleError:    true,
			wantNotInVCError:    false,
			wantVoiceStateError: false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
		},
		{
			name:                "Test Shuffle not in vc",
			queueSize:           3,
			wantShuffleError:    false,
			wantNotInVCError:    true,
			wantVoiceStateError: false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "",
			},
		},
		{
			name:                "Test Shuffle voice state error",
			queueSize:           3,
			wantShuffleError:    false,
			wantNotInVCError:    false,
			wantVoiceStateError: true,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			mockStream := mockaudio.NewMockStreamService(ctrl)
			mockVoice := mockaudio.NewMockVoiceService(ctrl)
			mockNotif := mockaudio.NewMockNotificationService(ctrl)

			serviceFactory := &FakeServiceFactory{}

			queue := util.NewQueue()

			for _ = range tt.queueSize {
				queue.Enque("http://video/url.com")
			}

			audioPlayer := audio.NewAudioPlayer(mockStream, mockVoice, mockNotif, true, queue, &sync.WaitGroup{})

			app := &application.Application{
				ImageApi:       nil,
				WordDatabase:   nil,
				GuildID:        tt.guildId,
				Search:         nil,
				ServiceFactory: serviceFactory,
				AudioPlayer:    audioPlayer,
			}
			//common expectations
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockInteractionCreate.EXPECT().GetInteractionAuthor().Times(1).Return(tt.author)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			//flow specific
			if tt.wantVoiceStateError {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(nil, errors.New("voice State error")).Times(1)
			} else {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(tt.voiceState, nil).Times(1)
			}

			err := Commands.Shuffle(mockSession, mockInteractionCreate, app)

			assert.Equal(t, tt.wantShuffleError || tt.wantVoiceStateError || tt.wantNotInVCError, err != nil)
		})
	}

}

func TestRandomPlay(t *testing.T) {
	tests := []RandomPlayCommandTestCase{
		{
			name:                "Test Random Play",
			wantNotInVCError:    false,
			wantVoiceStateError: false,
			wantSearchError:     false,
			wantJoinVCError:     false,
			wantDatabaseError:   false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: &SearchStub{
				val: "http://videolink",
				err: nil,
			},
			DatabaseResult: []string{"A", "search", "Query"},
		},
		{
			name:                "Test Random Play user not in vc",
			wantNotInVCError:    true,
			wantVoiceStateError: false,
			wantSearchError:     false,
			wantJoinVCError:     false,
			wantDatabaseError:   false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "",
			},
			SearchService:  nil,
			DatabaseResult: []string{},
		},
		{
			name:                "Test Random Play Voice State error",
			wantNotInVCError:    false,
			wantVoiceStateError: true,
			wantSearchError:     false,
			wantJoinVCError:     false,
			wantDatabaseError:   false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState:          nil,
			SearchService:       nil,
			DatabaseResult:      []string{},
		},
		{
			name:                "Test Random Play Database Error",
			wantNotInVCError:    false,
			wantVoiceStateError: false,
			wantSearchError:     false,
			wantJoinVCError:     false,
			wantDatabaseError:   true,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService:  nil,
			DatabaseResult: []string{},
		},
		{
			name:                "Test Random Play join VC error",
			wantNotInVCError:    false,
			wantVoiceStateError: false,
			wantSearchError:     false,
			wantJoinVCError:     true,
			wantDatabaseError:   false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: &SearchStub{
				val: "http://videolink",
				err: nil,
			},
			DatabaseResult: []string{"A", "search", "Query"},
		},
		{
			name:                "Test Random Play",
			wantNotInVCError:    false,
			wantVoiceStateError: false,
			wantSearchError:     true,
			wantJoinVCError:     false,
			wantDatabaseError:   false,
			author:              "author",
			guildId:             "guildId",
			channelId:           "channelId",
			voiceState: &discordgo.VoiceState{
				GuildID:   "guildId",
				ChannelID: "channelId",
			},
			SearchService: &SearchStub{
				val: "",
				err: errors.New("search error"),
			},
			DatabaseResult: []string{"A", "search", "Query"},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			userInput := &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "terms",
				Type:  discordgo.ApplicationCommandOptionInteger,
				Value: float64(tt.numTerms),
			}

			applicationData := &discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{userInput},
			}

			mockStream := mockaudio.NewMockStreamService(ctrl)
			mockVoice := mockaudio.NewMockVoiceService(ctrl)
			mockNotif := mockaudio.NewMockNotificationService(ctrl)

			serviceFactory := &FakeServiceFactory{}

			mockDatabase := mockdatabase.NewMockDatabaseService(ctrl)

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(applicationData).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Times(1).Return(interaction)
			mockInteractionCreate.EXPECT().GetInteractionAuthor().Times(1).Return(tt.author)
			mockInteractionCreate.EXPECT().GetChannel().Times(1).Return(tt.channelId)
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			// flow specific
			if tt.wantVoiceStateError {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(nil, errors.New("voice State error")).Times(1)
			} else {
				mockSession.EXPECT().VoiceState(tt.guildId, tt.author).Return(tt.voiceState, nil).Times(1)
			}

			if !tt.wantNotInVCError && !tt.wantVoiceStateError {
				if tt.wantDatabaseError {
					mockDatabase.EXPECT().FetchRandom(tt.numTerms).Return(nil, errors.New("database error")).Times(1)
				} else {
					mockDatabase.EXPECT().FetchRandom(tt.numTerms).Return(tt.DatabaseResult, nil).Times(1)
				}
			}

			if tt.wantJoinVCError {
				mockSession.EXPECT().ChannelVoiceJoin(tt.guildId, tt.voiceState.ChannelID, false, false).Return(nil, errors.New("join vc error")).Times(1)
			}

			if !(tt.wantDatabaseError || tt.wantVoiceStateError || tt.wantNotInVCError || tt.wantJoinVCError) {
				mockSession.EXPECT().ChannelVoiceJoin(tt.guildId, tt.voiceState.ChannelID, false, false).Times(1).Return(&discordgo.VoiceConnection{}, nil)
			}

			if !(tt.wantJoinVCError || tt.wantVoiceStateError || tt.wantDatabaseError || tt.wantNotInVCError) {
				mockStream.EXPECT().GetAudioStream(gomock.Any()).Times(1).Return("stream url", nil)
			}

			audioPlayer := audio.NewAudioPlayer(mockStream, mockVoice, mockNotif, false, util.NewQueue(), &sync.WaitGroup{})

			app := &application.Application{
				ImageApi:       nil,
				WordDatabase:   mockDatabase,
				GuildID:        tt.guildId,
				Search:         tt.SearchService,
				ServiceFactory: serviceFactory,
				AudioPlayer:    audioPlayer,
			}

			err := Commands.RandomPlay(mockSession, mockInteractionCreate, app)
			audioPlayer.Wait()
			assert.Equal(t, tt.wantDatabaseError || tt.wantJoinVCError || tt.wantVoiceStateError || tt.wantNotInVCError || tt.wantSearchError, err != nil)
			ctrl.Finish()
		})
	}
}
