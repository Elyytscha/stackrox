// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	reflect "reflect"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// GetNamespace mocks base method
func (m *MockDataStore) GetNamespace(ctx context.Context, id string) (*storage.NamespaceMetadata, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace", ctx, id)
	ret0, _ := ret[0].(*storage.NamespaceMetadata)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNamespace indicates an expected call of GetNamespace
func (mr *MockDataStoreMockRecorder) GetNamespace(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockDataStore)(nil).GetNamespace), ctx, id)
}

// GetNamespaces mocks base method
func (m *MockDataStore) GetNamespaces(ctx context.Context) ([]*storage.NamespaceMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespaces", ctx)
	ret0, _ := ret[0].([]*storage.NamespaceMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespaces indicates an expected call of GetNamespaces
func (mr *MockDataStoreMockRecorder) GetNamespaces(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespaces", reflect.TypeOf((*MockDataStore)(nil).GetNamespaces), ctx)
}

// AddNamespace mocks base method
func (m *MockDataStore) AddNamespace(arg0 context.Context, arg1 *storage.NamespaceMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNamespace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNamespace indicates an expected call of AddNamespace
func (mr *MockDataStoreMockRecorder) AddNamespace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNamespace", reflect.TypeOf((*MockDataStore)(nil).AddNamespace), arg0, arg1)
}

// UpdateNamespace mocks base method
func (m *MockDataStore) UpdateNamespace(arg0 context.Context, arg1 *storage.NamespaceMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNamespace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNamespace indicates an expected call of UpdateNamespace
func (mr *MockDataStoreMockRecorder) UpdateNamespace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNamespace", reflect.TypeOf((*MockDataStore)(nil).UpdateNamespace), arg0, arg1)
}

// RemoveNamespace mocks base method
func (m *MockDataStore) RemoveNamespace(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNamespace", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNamespace indicates an expected call of RemoveNamespace
func (mr *MockDataStoreMockRecorder) RemoveNamespace(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNamespace", reflect.TypeOf((*MockDataStore)(nil).RemoveNamespace), ctx, id)
}

// SearchResults mocks base method
func (m *MockDataStore) SearchResults(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchResults", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchResults indicates an expected call of SearchResults
func (mr *MockDataStoreMockRecorder) SearchResults(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchResults", reflect.TypeOf((*MockDataStore)(nil).SearchResults), ctx, q)
}

// Search mocks base method
func (m *MockDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// SearchNamespaces mocks base method
func (m *MockDataStore) SearchNamespaces(ctx context.Context, q *v1.Query) ([]*storage.NamespaceMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchNamespaces", ctx, q)
	ret0, _ := ret[0].([]*storage.NamespaceMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchNamespaces indicates an expected call of SearchNamespaces
func (mr *MockDataStoreMockRecorder) SearchNamespaces(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchNamespaces", reflect.TypeOf((*MockDataStore)(nil).SearchNamespaces), ctx, q)
}
