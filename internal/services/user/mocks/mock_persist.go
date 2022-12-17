// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user (interfaces: Persistent)

// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	reflect "reflect"

	models "github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockPersistent is a mock of Persistent interface.
type MockPersistent struct {
	ctrl     *gomock.Controller
	recorder *MockPersistentMockRecorder
}

// MockPersistentMockRecorder is the mock recorder for MockPersistent.
type MockPersistentMockRecorder struct {
	mock *MockPersistent
}

// NewMockPersistent creates a new mock instance.
func NewMockPersistent(ctrl *gomock.Controller) *MockPersistent {
	mock := &MockPersistent{ctrl: ctrl}
	mock.recorder = &MockPersistentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersistent) EXPECT() *MockPersistentMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPersistent) Create(arg0 context.Context, arg1 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPersistentMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPersistent)(nil).Create), arg0, arg1)
}

// FindByEmail mocks base method.
func (m *MockPersistent) FindByEmail(arg0 context.Context, arg1 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockPersistentMockRecorder) FindByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockPersistent)(nil).FindByEmail), arg0, arg1)
}
