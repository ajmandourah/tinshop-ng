// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ajmandourah/tinshop/repository (interfaces: Source)

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	http "net/http"
	reflect "reflect"

	repository "github.com/ajmandourah/tinshop/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockSource) Download(arg0 http.ResponseWriter, arg1 *http.Request, arg2, arg3 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Download", arg0, arg1, arg2, arg3)
}

// Download indicates an expected call of Download.
func (mr *MockSourceMockRecorder) Download(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockSource)(nil).Download), arg0, arg1, arg2, arg3)
}

// GetFiles mocks base method.
func (m *MockSource) GetFiles() []repository.FileDesc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFiles")
	ret0, _ := ret[0].([]repository.FileDesc)
	return ret0
}

// GetFiles indicates an expected call of GetFiles.
func (mr *MockSourceMockRecorder) GetFiles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFiles", reflect.TypeOf((*MockSource)(nil).GetFiles))
}

// Load mocks base method.
func (m *MockSource) Load(arg0 []string, arg1 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Load", arg0, arg1)
}

// Load indicates an expected call of Load.
func (mr *MockSourceMockRecorder) Load(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockSource)(nil).Load), arg0, arg1)
}

// Reset mocks base method.
func (m *MockSource) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockSourceMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockSource)(nil).Reset))
}

// UnWatchAll mocks base method.
func (m *MockSource) UnWatchAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnWatchAll")
}

// UnWatchAll indicates an expected call of UnWatchAll.
func (mr *MockSourceMockRecorder) UnWatchAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnWatchAll", reflect.TypeOf((*MockSource)(nil).UnWatchAll))
}
