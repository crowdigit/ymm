// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/crowdigit/ymm/command (interfaces: CommandProvider)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	command "github.com/crowdigit/ymm/command"
	gomock "github.com/golang/mock/gomock"
)

// MockCommandProvider is a mock of CommandProvider interface.
type MockCommandProvider struct {
	ctrl     *gomock.Controller
	recorder *MockCommandProviderMockRecorder
}

// MockCommandProviderMockRecorder is the mock recorder for MockCommandProvider.
type MockCommandProviderMockRecorder struct {
	mock *MockCommandProvider
}

// NewMockCommandProvider creates a new mock instance.
func NewMockCommandProvider(ctrl *gomock.Controller) *MockCommandProvider {
	mock := &MockCommandProvider{ctrl: ctrl}
	mock.recorder = &MockCommandProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandProvider) EXPECT() *MockCommandProviderMockRecorder {
	return m.recorder
}

// NewCommand mocks base method.
func (m *MockCommandProvider) NewCommand(arg0 string, arg1 ...string) command.Command {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewCommand", varargs...)
	ret0, _ := ret[0].(command.Command)
	return ret0
}

// NewCommand indicates an expected call of NewCommand.
func (mr *MockCommandProviderMockRecorder) NewCommand(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCommand", reflect.TypeOf((*MockCommandProvider)(nil).NewCommand), varargs...)
}
