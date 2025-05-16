package Commands_test

import (
	application "bot/Application"
	"bot/Core/Commands"
	"bot/Core/Interfaces/Mocks"
	mockdatabase "bot/Core/Services/Database/Mocks"
	"errors"
	"go.uber.org/mock/gomock"
	"testing"
)

type InsertTestCase struct {
	name    string
	data    string
	wantErr bool
}

/*
Expect calling Insert to insert the passed item
*/
func ExpectInsert(mockDB *mockdatabase.MockDatabaseService, item string) {
	mockDB.EXPECT().Insert(item).Return(nil)
}

/*
Expect the following:
  - Calling Insert resulting in an error (i.e. item already in the database, or other database error)
  - Calling GetInteraction and InteractionRespond to inform the user of the error
*/
func ExpectInsertError(mockDB *mockdatabase.MockDatabaseService, mockInteractionCreate *mock_Interfaces.MockDiscordInteraction, mockSession *mock_Interfaces.MockDiscordSession, item string) {
	mockDB.EXPECT().Insert(item).Return(errors.New("error"))
	mockInteractionCreate.EXPECT().GetInteraction()
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
}

/*
Expect calling Remove to remove 'item' from the database
*/
func ExpectRemove(mockDB *mockdatabase.MockDatabaseService, item string) {
	mockDB.EXPECT().Delete(item).Return(nil)
}

/*
Expect the following:
  - Calling Insert resulting in an error (i.e. item already in the database, or other database error)
  - Calling GetInteraction and InteractionRespond to inform the user of the error
*/
func ExpectRemoveError(mockDB *mockdatabase.MockDatabaseService, mockInteractionCreate *mock_Interfaces.MockDiscordInteraction, mockSession *mock_Interfaces.MockDiscordSession, item string) {
	mockDB.EXPECT().Insert(item).Return(errors.New("error"))
	mockInteractionCreate.EXPECT().GetInteraction()
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
}

/*
Expect calling GetAll to return all items in a database
*/
func ExpectGetAll(mockDB *mockdatabase.MockDatabaseService) {
	mockDB.EXPECT().GetAll().Return([]string{}, nil)
}

/*
Expect the following:
  - Calling GetAll resulting in an error
  - Calling GetInteraction and InteractionRespond to inform the user of the error
*/
func ExpectGetAllError(mockDB *mockdatabase.MockDatabaseService, mockInteractionCreate *mock_Interfaces.MockDiscordInteraction, mockSession *mock_Interfaces.MockDiscordSession, item string) {
	mockDB.EXPECT().GetAll().Return(nil, errors.New("error"))
	mockInteractionCreate.EXPECT().GetInteraction()
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
}

func setUpAddExpectations(tt *InsertTestCase, mockInteractionCreate *mock_Interfaces.MockDiscordInteraction, mockSession *mock_Interfaces.MockDiscordSession, mockDB *mockdatabase.MockDatabaseService) {

	if tt.wantErr {
		ExpectInsertError(mockDB, mockInteractionCreate, mockSession, tt.data)
		return
	}
	ExpectInsert(mockDB, tt.data)

}

func TestAdd(t *testing.T) {
	tests := []InsertTestCase{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			mockSession := mock_Interfaces.NewMockDiscordSession(ctrl)
			mockInteractionCreate := mock_Interfaces.NewMockDiscordInteraction(ctrl)

			mockDBService := mockdatabase.NewMockDatabaseService(ctrl)

			mockApplication := &application.Application{
				ImageApi:     nil,
				AudioPlayer:  nil,
				WordDatabase: mockDBService,
				Search:       nil,
				GuildID:      "guildID",
			}

			setUpAddExpectations(&tt, mockInteractionCreate, mockSession, mockDBService)
			Commands.Add(mockSession, mockInteractionCreate, mockApplication)
		})
	}
}
