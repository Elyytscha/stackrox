// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	storage "github.com/stackrox/rox/generated/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// AddNotifier mocks base method.
func (m *MockDataStore) AddNotifier(ctx context.Context, notifier *storage.Notifier) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNotifier", ctx, notifier)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNotifier indicates an expected call of AddNotifier.
func (mr *MockDataStoreMockRecorder) AddNotifier(ctx, notifier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNotifier", reflect.TypeOf((*MockDataStore)(nil).AddNotifier), ctx, notifier)
}

// Exists mocks base method.
func (m *MockDataStore) Exists(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockDataStoreMockRecorder) Exists(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockDataStore)(nil).Exists), ctx, id)
}

// GetNotifier mocks base method.
func (m *MockDataStore) GetNotifier(ctx context.Context, id string) (*storage.Notifier, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifier", ctx, id)
	ret0, _ := ret[0].(*storage.Notifier)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNotifier indicates an expected call of GetNotifier.
func (mr *MockDataStoreMockRecorder) GetNotifier(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifier", reflect.TypeOf((*MockDataStore)(nil).GetNotifier), ctx, id)
}

// GetNotifiers mocks base method.
func (m *MockDataStore) GetNotifiers(ctx context.Context) ([]*storage.Notifier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifiers", ctx)
	ret0, _ := ret[0].([]*storage.Notifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotifiers indicates an expected call of GetNotifiers.
func (mr *MockDataStoreMockRecorder) GetNotifiers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifiers", reflect.TypeOf((*MockDataStore)(nil).GetNotifiers), ctx)
}

// GetNotifiersFiltered mocks base method.
func (m *MockDataStore) GetNotifiersFiltered(ctx context.Context, filter func(*storage.Notifier) bool) ([]*storage.Notifier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifiersFiltered", ctx, filter)
	ret0, _ := ret[0].([]*storage.Notifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotifiersFiltered indicates an expected call of GetNotifiersFiltered.
func (mr *MockDataStoreMockRecorder) GetNotifiersFiltered(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifiersFiltered", reflect.TypeOf((*MockDataStore)(nil).GetNotifiersFiltered), ctx, filter)
}

// GetScrubbedNotifier mocks base method.
func (m *MockDataStore) GetScrubbedNotifier(ctx context.Context, id string) (*storage.Notifier, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScrubbedNotifier", ctx, id)
	ret0, _ := ret[0].(*storage.Notifier)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetScrubbedNotifier indicates an expected call of GetScrubbedNotifier.
func (mr *MockDataStoreMockRecorder) GetScrubbedNotifier(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScrubbedNotifier", reflect.TypeOf((*MockDataStore)(nil).GetScrubbedNotifier), ctx, id)
}

// GetScrubbedNotifiers mocks base method.
func (m *MockDataStore) GetScrubbedNotifiers(ctx context.Context) ([]*storage.Notifier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScrubbedNotifiers", ctx)
	ret0, _ := ret[0].([]*storage.Notifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScrubbedNotifiers indicates an expected call of GetScrubbedNotifiers.
func (mr *MockDataStoreMockRecorder) GetScrubbedNotifiers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScrubbedNotifiers", reflect.TypeOf((*MockDataStore)(nil).GetScrubbedNotifiers), ctx)
}

// RemoveNotifier mocks base method.
func (m *MockDataStore) RemoveNotifier(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNotifier", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNotifier indicates an expected call of RemoveNotifier.
func (mr *MockDataStoreMockRecorder) RemoveNotifier(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNotifier", reflect.TypeOf((*MockDataStore)(nil).RemoveNotifier), ctx, id)
}

// UpdateNotifier mocks base method.
func (m *MockDataStore) UpdateNotifier(ctx context.Context, notifier *storage.Notifier) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotifier", ctx, notifier)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNotifier indicates an expected call of UpdateNotifier.
func (mr *MockDataStoreMockRecorder) UpdateNotifier(ctx, notifier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotifier", reflect.TypeOf((*MockDataStore)(nil).UpdateNotifier), ctx, notifier)
}

// UpsertNotifier mocks base method.
func (m *MockDataStore) UpsertNotifier(ctx context.Context, notifier *storage.Notifier) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertNotifier", ctx, notifier)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertNotifier indicates an expected call of UpsertNotifier.
func (mr *MockDataStoreMockRecorder) UpsertNotifier(ctx, notifier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertNotifier", reflect.TypeOf((*MockDataStore)(nil).UpsertNotifier), ctx, notifier)
}
