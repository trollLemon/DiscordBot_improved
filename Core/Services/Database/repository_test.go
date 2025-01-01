package database

import (
	"bot/Core/Services/Database/Mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)



func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mock_database.NewMockDatabaseService(ctrl)
	
	repo := NewRepository(mockDatabase)

	input1 := "word"
	input2 := "aword"	
	input3 := "word"
	

	mockDatabase.EXPECT().IsPresent(input1).Return(false).Times(1) 
	mockDatabase.EXPECT().IsPresent(input2).Return(false).Times(1)
	
	mockDatabase.EXPECT().Insert(input1).Times(1)
	mockDatabase.EXPECT().Insert(input2).Times(1)
	
	err := repo.Add(input1)
	

	assert.Nil(t,err)
	
	err = repo.Add(input2)
	
	assert.Nil(t,err)

	// duplicate item case
	mockDatabase.EXPECT().IsPresent(input3).Return(true).Times(1) 
	err = repo.Add(input3)

	assert.NotNil(t,err)
}

func TestRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mock_database.NewMockDatabaseService(ctrl)
	
	repo := NewRepository(mockDatabase)

	input1 := "word"
	input2 := "aword"	
	
	key_1 := input1
	key_2 := input2
	
	mockDatabase.EXPECT().Delete(key_1).Times(1)
	mockDatabase.EXPECT().Delete(key_2).Times(1)
	
	err := repo.Remove(input1)
	

	assert.Nil(t,err)
	
	err = repo.Remove(input2)
	
	assert.Nil(t,err)
	
	mockDatabase.EXPECT().Delete(key_1).Return(fmt.Errorf("Error removing data")).Times(1)
	
	err = repo.Remove(input1)

	assert.NotNil(t,err)
}


func TestGetRandN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mock_database.NewMockDatabaseService(ctrl)
	
	repo := NewRepository(mockDatabase)


	mockDatabase.EXPECT().FetchRandom(1).Return(nil,fmt.Errorf("Could not get data from database")).Times(1) 
	
	data, err := repo.GetRandN(1)

	assert.Nil(t,data)
	assert.NotNil(t,err)
	mockDatabase.EXPECT().FetchRandom(1).Return([]string{"url1"}, nil).Times(1)
	
	data, err = repo.GetRandN(1)


	assert.NotNil(t,data)
	assert.Nil(t,err)


}




