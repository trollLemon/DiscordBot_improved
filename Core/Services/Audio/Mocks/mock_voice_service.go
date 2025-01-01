// Code generated by MockGen. DO NOT EDIT.
// Source: voice_service.go
//
// Generated by this command:
//
//	mockgen --source=voice_service.go --destination=mock_voice_service.go
//

// Package mock_audio is a generated GoMock package.
package mock_audio

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockVoiceService is a mock of VoiceService interface.
type MockVoiceService struct {
	ctrl     *gomock.Controller
	recorder *MockVoiceServiceMockRecorder
}

// MockVoiceServiceMockRecorder is the mock recorder for MockVoiceService.
type MockVoiceServiceMockRecorder struct {
	mock *MockVoiceService
}

// NewMockVoiceService creates a new mock instance.
func NewMockVoiceService(ctrl *gomock.Controller) *MockVoiceService {
	mock := &MockVoiceService{ctrl: ctrl}
	mock.recorder = &MockVoiceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVoiceService) EXPECT() *MockVoiceServiceMockRecorder {
	return m.recorder
}

// Disconnect mocks base method.
func (m *MockVoiceService) Disconnect() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disconnect")
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockVoiceServiceMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockVoiceService)(nil).Disconnect))
}

// PlayAudioFile mocks base method.
func (m *MockVoiceService) PlayAudioFile(url string, Done chan bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PlayAudioFile", url, Done)
}

// PlayAudioFile indicates an expected call of PlayAudioFile.
func (mr *MockVoiceServiceMockRecorder) PlayAudioFile(url, Done any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayAudioFile", reflect.TypeOf((*MockVoiceService)(nil).PlayAudioFile), url, Done)
}
