// Code generated by MockGen. DO NOT EDIT.
// Source: identity.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	authproviders "github.com/stackrox/rox/pkg/auth/authproviders"
	reflect "reflect"
	time "time"
)

// MockIdentity is a mock of Identity interface
type MockIdentity struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityMockRecorder
}

// MockIdentityMockRecorder is the mock recorder for MockIdentity
type MockIdentityMockRecorder struct {
	mock *MockIdentity
}

// NewMockIdentity creates a new mock instance
func NewMockIdentity(ctrl *gomock.Controller) *MockIdentity {
	mock := &MockIdentity{ctrl: ctrl}
	mock.recorder = &MockIdentityMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentity) EXPECT() *MockIdentityMockRecorder {
	return m.recorder
}

// UID mocks base method
func (m *MockIdentity) UID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UID indicates an expected call of UID
func (mr *MockIdentityMockRecorder) UID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UID", reflect.TypeOf((*MockIdentity)(nil).UID))
}

// FriendlyName mocks base method
func (m *MockIdentity) FriendlyName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FriendlyName")
	ret0, _ := ret[0].(string)
	return ret0
}

// FriendlyName indicates an expected call of FriendlyName
func (mr *MockIdentityMockRecorder) FriendlyName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FriendlyName", reflect.TypeOf((*MockIdentity)(nil).FriendlyName))
}

// FullName mocks base method
func (m *MockIdentity) FullName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullName")
	ret0, _ := ret[0].(string)
	return ret0
}

// FullName indicates an expected call of FullName
func (mr *MockIdentityMockRecorder) FullName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullName", reflect.TypeOf((*MockIdentity)(nil).FullName))
}

// User mocks base method
func (m *MockIdentity) User() *storage.UserInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User")
	ret0, _ := ret[0].(*storage.UserInfo)
	return ret0
}

// User indicates an expected call of User
func (mr *MockIdentityMockRecorder) User() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockIdentity)(nil).User))
}

// Role mocks base method
func (m *MockIdentity) Role() *storage.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Role")
	ret0, _ := ret[0].(*storage.Role)
	return ret0
}

// Role indicates an expected call of Role
func (mr *MockIdentityMockRecorder) Role() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Role", reflect.TypeOf((*MockIdentity)(nil).Role))
}

// Service mocks base method
func (m *MockIdentity) Service() *storage.ServiceIdentity {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Service")
	ret0, _ := ret[0].(*storage.ServiceIdentity)
	return ret0
}

// Service indicates an expected call of Service
func (mr *MockIdentityMockRecorder) Service() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Service", reflect.TypeOf((*MockIdentity)(nil).Service))
}

// Attributes mocks base method
func (m *MockIdentity) Attributes() map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes")
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// Attributes indicates an expected call of Attributes
func (mr *MockIdentityMockRecorder) Attributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockIdentity)(nil).Attributes))
}

// Expiry mocks base method
func (m *MockIdentity) Expiry() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expiry")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Expiry indicates an expected call of Expiry
func (mr *MockIdentityMockRecorder) Expiry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expiry", reflect.TypeOf((*MockIdentity)(nil).Expiry))
}

// ExternalAuthProvider mocks base method
func (m *MockIdentity) ExternalAuthProvider() authproviders.Provider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExternalAuthProvider")
	ret0, _ := ret[0].(authproviders.Provider)
	return ret0
}

// ExternalAuthProvider indicates an expected call of ExternalAuthProvider
func (mr *MockIdentityMockRecorder) ExternalAuthProvider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExternalAuthProvider", reflect.TypeOf((*MockIdentity)(nil).ExternalAuthProvider))
}
