// Code generated by MockGen. DO NOT EDIT.
// Source: notification_service.go
//
// Generated by this command:
//
//	mockgen --source=notification_service.go --destination=./Mocks/mock_notification_service.go
//

// Package mock_audio is a generated GoMock package.
package mock_audio

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockNotificationService is a mock of NotificationService interface.
type MockNotificationService struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceMockRecorder
}

// MockNotificationServiceMockRecorder is the mock recorder for MockNotificationService.
type MockNotificationServiceMockRecorder struct {
	mock *MockNotificationService
}

// NewMockNotificationService creates a new mock instance.
func NewMockNotificationService(ctrl *gomock.Controller) *MockNotificationService {
	mock := &MockNotificationService{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationService) EXPECT() *MockNotificationServiceMockRecorder {
	return m.recorder
}

// SendError mocks base method.
func (m *MockNotificationService) SendError(error string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendError", error)
}

// SendError indicates an expected call of SendError.
func (mr *MockNotificationServiceMockRecorder) SendError(error any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendError", reflect.TypeOf((*MockNotificationService)(nil).SendError), error)
}

// SendNotification mocks base method.
func (m *MockNotificationService) SendNotification(content string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendNotification", content)
}

// SendNotification indicates an expected call of SendNotification.
func (mr *MockNotificationServiceMockRecorder) SendNotification(content any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotification", reflect.TypeOf((*MockNotificationService)(nil).SendNotification), content)
}

// UpdateService mocks base method.
func (m *MockNotificationService) UpdateService(channel string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateService", channel)
}

// UpdateService indicates an expected call of UpdateService.
func (mr *MockNotificationServiceMockRecorder) UpdateService(channel any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateService", reflect.TypeOf((*MockNotificationService)(nil).UpdateService), channel)
}
