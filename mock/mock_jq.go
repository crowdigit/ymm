// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/crowdigit/ymm/jq (interfaces: Jq)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJq is a mock of Jq interface.
type MockJq struct {
	ctrl     *gomock.Controller
	recorder *MockJqMockRecorder
}

// MockJqMockRecorder is the mock recorder for MockJq.
type MockJqMockRecorder struct {
	mock *MockJq
}

// NewMockJq creates a new mock instance.
func NewMockJq(ctrl *gomock.Controller) *MockJq {
	mock := &MockJq{ctrl: ctrl}
	mock.recorder = &MockJqMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJq) EXPECT() *MockJqMockRecorder {
	return m.recorder
}

// Slurp mocks base method.
func (m *MockJq) Slurp(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Slurp", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Slurp indicates an expected call of Slurp.
func (mr *MockJqMockRecorder) Slurp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Slurp", reflect.TypeOf((*MockJq)(nil).Slurp), arg0)
}