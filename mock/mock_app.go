// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/crowdigit/ymm/app (interfaces: Application)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// DownloadPlaylist mocks base method.
func (m *MockApplication) DownloadPlaylist(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DownloadPlaylist", arg0)
}

// DownloadPlaylist indicates an expected call of DownloadPlaylist.
func (mr *MockApplicationMockRecorder) DownloadPlaylist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadPlaylist", reflect.TypeOf((*MockApplication)(nil).DownloadPlaylist), arg0)
}

// DownloadSingle mocks base method.
func (m *MockApplication) DownloadSingle(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DownloadSingle", arg0)
}

// DownloadSingle indicates an expected call of DownloadSingle.
func (mr *MockApplicationMockRecorder) DownloadSingle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadSingle", reflect.TypeOf((*MockApplication)(nil).DownloadSingle), arg0)
}