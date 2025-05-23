package Commands_test

import (
	application "bot/Application"
	"bot/Core/Commands"
	mockinterfaces "bot/Core/Interfaces/Mocks"
	imagemanip "bot/Core/Services/ImageManip"
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GeneralImageTestCase struct {
	name               string
	wantApiErr         bool
	wantCDNErr         bool
	userInput          discordgo.ApplicationCommandInteractionData
	apiHandlerFunc     func(w http.ResponseWriter, r *http.Request)
	mockCDNHandlerFunc func(w http.ResponseWriter, r *http.Request)
}

func TestRandomImageFilter(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Test Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},

					{
						Name:  "kernel",
						Value: 3.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "lowerbound",
						Value: -1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "upperbound",
						Value: 1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "normalize",
						Value: true,
						Type:  discordgo.ApplicationCommandOptionBoolean,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "Test CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},

					{
						Name:  "kernel",
						Value: 3.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "lowerbound",
						Value: -1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "upperbound",
						Value: 1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "normalize",
						Value: true,
						Type:  discordgo.ApplicationCommandOptionBoolean,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Test Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},

					{
						Name:  "kernel",
						Value: 3.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "lowerbound",
						Value: -1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "upperbound",
						Value: 1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "normalize",
						Value: true,
						Type:  discordgo.ApplicationCommandOptionBoolean,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.RandomImageFilter(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestInvertImage(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Test Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "Test CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Test Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.InvertImage(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestSaturateImage(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "magnitude",
						Value: 10.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "magnitude",
						Value: 10.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "magnitude",
						Value: 10.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.SaturateImage(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestEdgeDetection(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "lowerbound",
						Value: 100.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "upperbound",
						Value: 200.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "Test Dilate CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "lowerbound",
						Value: 100.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "upperbound",
						Value: 200.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Test Dilate Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "lowerbound",
						Value: 100.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "upperbound",
						Value: 200.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.EdgeDetection(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestDilate(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "boxsize",
						Value: 5.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "iterations",
						Value: 2.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "boxsize",
						Value: 5.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "iterations",
						Value: 2.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "boxsize",
						Value: 5.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "iterations",
						Value: 2.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.Dilate(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestErode(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "boxsize",
						Value: 5.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "iterations",
						Value: 2.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "boxsize",
						Value: 5.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "iterations",
						Value: 2.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "boxsize",
						Value: 5.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "iterations",
						Value: 2.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.Erode(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestAddText(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "text",
						Value: "sample text",
						Type:  discordgo.ApplicationCommandOptionString,
					},
					{
						Name:  "fontScale",
						Value: 1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "x",
						Value: 50.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "y",
						Value: 50.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "text",
						Value: "sample text",
						Type:  discordgo.ApplicationCommandOptionString,
					},
					{
						Name:  "fontScale",
						Value: 1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "x",
						Value: 50.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "y",
						Value: 50.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "text",
						Value: "sample text",
						Type:  discordgo.ApplicationCommandOptionString,
					},
					{
						Name:  "fontScale",
						Value: 1.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "x",
						Value: 50.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
					{
						Name:  "y",
						Value: 50.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.AddText(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestReduceImage(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "quality",
						Value: 10.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "quality",
						Value: 10.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "quality",
						Value: 10.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.ReduceImage(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}

func TestShuffleImage(t *testing.T) {
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()

	testContentType := "image/png"
	tests := []GeneralImageTestCase{
		{
			name:       "Success",
			wantApiErr: false,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "partitions",
						Value: 64.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "CDN error",
			wantApiErr: false,
			wantCDNErr: true,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "partitions",
						Value: 64.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Gomanip error",
			wantApiErr: true,
			wantCDNErr: false,
			userInput: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "image",
						Value: "imageId",
						Type:  discordgo.ApplicationCommandOptionAttachment,
					},
					{
						Name:  "partitions",
						Value: 64.0,
						Type:  discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			apiHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // or another error status code

				errResp := imagemanip.ErrorResponse{
					Detail: "api error",
				}

				json.NewEncoder(w).Encode(errResp)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			testHttpClient := &http.Client{}

			testCDN := httptest.NewServer(http.HandlerFunc(tt.mockCDNHandlerFunc))
			testGoManipServer := httptest.NewServer(http.HandlerFunc(tt.apiHandlerFunc))

			api := imagemanip.NewGoManip(testHttpClient, testGoManipServer.URL)

			mockSession := mockinterfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mockinterfaces.NewMockDiscordInteraction(ctrl)

			interaction := &discordgo.Interaction{}

			//common expectations
			mockInteractionCreate.EXPECT().ApplicationCommandData().Return(tt.userInput).Times(1)
			mockInteractionCreate.EXPECT().GetImageURLFromAttachmentID("imageId").Return(testCDN.URL).Times(1)
			mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)

			if !tt.wantCDNErr {
				mockInteractionCreate.EXPECT().GetInteraction().Return(interaction)
				mockSession.EXPECT().InteractionResponseEdit(interaction, gomock.Any()).Times(1)

			}
			mockSession.EXPECT().InteractionRespond(interaction, gomock.Any()).Times(1)

			mockApplication := &application.Application{
				ImageApi:     api,
				WordDatabase: nil,
				GuildID:      "guildID",
			}

			err := Commands.ShuffleImage(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}
