// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/crowdigit/ymm/ydl (interfaces: YoutubeDL)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	ydl "github.com/crowdigit/ymm/ydl"
	gomock "github.com/golang/mock/gomock"
)

// MockYoutubeDL is a mock of YoutubeDL interface.
type MockYoutubeDL struct {
	ctrl     *gomock.Controller
	recorder *MockYoutubeDLMockRecorder
}

// MockYoutubeDLMockRecorder is the mock recorder for MockYoutubeDL.
type MockYoutubeDLMockRecorder struct {
	mock *MockYoutubeDL
}

// NewMockYoutubeDL creates a new mock instance.
func NewMockYoutubeDL(ctrl *gomock.Controller) *MockYoutubeDL {
	mock := &MockYoutubeDL{ctrl: ctrl}
	mock.recorder = &MockYoutubeDLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockYoutubeDL) EXPECT() *MockYoutubeDLMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockYoutubeDL) Download(arg0 string, arg1 ydl.VideoMetadata) (ydl.DownloadResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0, arg1)
	ret0, _ := ret[0].(ydl.DownloadResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockYoutubeDLMockRecorder) Download(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockYoutubeDL)(nil).Download), arg0, arg1)
}

// PlaylistMetadata mocks base method.
func (m *MockYoutubeDL) PlaylistMetadata(arg0 string) ([][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaylistMetadata", arg0)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlaylistMetadata indicates an expected call of PlaylistMetadata.
func (mr *MockYoutubeDLMockRecorder) PlaylistMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaylistMetadata", reflect.TypeOf((*MockYoutubeDL)(nil).PlaylistMetadata), arg0)
}

// VideoMetadata mocks base method.
func (m *MockYoutubeDL) VideoMetadata(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VideoMetadata", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VideoMetadata indicates an expected call of VideoMetadata.
func (mr *MockYoutubeDLMockRecorder) VideoMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VideoMetadata", reflect.TypeOf((*MockYoutubeDL)(nil).VideoMetadata), arg0)
}
