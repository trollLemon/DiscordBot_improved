package Commands_test

import (
	application "bot/Application"
	"bot/Core/Commands"
	mockinterfaces "bot/Core/Interfaces/Mocks"
	"bot/Core/Services/Classification"
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
	"time"
)

type ClassificationTestCase struct {
	name               string
	wantApiErr         bool
	wantCDNErr         bool
	userInput          discordgo.ApplicationCommandInteractionData
	postHandlerFunc    func(w http.ResponseWriter, r *http.Request)
	getHandlerFunc     func(w http.ResponseWriter, r *http.Request)
	mockCDNHandlerFunc func(w http.ResponseWriter, r *http.Request)
}

func TestClassify(t *testing.T) {
	mockJobId := "mockJobId"
	sendImageEndpoint := "/api/v1/image"
	getClassificationEndpoint := "/api/v1/image/classification"
	getClassificationEndpointWithId := getClassificationEndpoint + "/" + mockJobId
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	var buf bytes.Buffer
	err := png.Encode(&buf, testImg)

	if err != nil {
		t.Fatal(err)
	}
	testImgBytes := buf.Bytes()
	testContentType := "image/png"

	tests := []ClassificationTestCase{
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
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				resp := map[string]string{"Class": "image"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
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
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				resp := map[string]string{"Class": "image"}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:       "Test API Error (polling timeout)",
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
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)

				resp := map[string]string{"jobId": mockJobId}

				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
			},
			mockCDNHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", testContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(testImgBytes)
			},
		},
		{
			name:       "Test API Error (sending image timeout)",
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
			postHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
			},
			getHandlerFunc: func(w http.ResponseWriter, r *http.Request) {

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
			mux := http.NewServeMux()
			mux.HandleFunc(sendImageEndpoint, tt.postHandlerFunc)
			mux.HandleFunc(getClassificationEndpointWithId, tt.getHandlerFunc)
			testServer := httptest.NewServer(mux)

			//api := imagemanip.NewGoManip(testHttpClient, testClassificationServer.URL)

			api := Classification.NewImageClassification(testHttpClient, 3*time.Second, testServer.URL, sendImageEndpoint, getClassificationEndpoint)
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
				ClassificationApi: api,
				WordDatabase:      nil,
				GuildID:           "guildID",
			}

			err := Commands.Classify(mockSession, mockInteractionCreate, mockApplication)

			assert.Equal(t, tt.wantApiErr || tt.wantCDNErr, err != nil)
			ctrl.Finish()
		})
	}
}
