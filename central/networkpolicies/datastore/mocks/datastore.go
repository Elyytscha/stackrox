// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockUndoDataStore is a mock of UndoDataStore interface
type MockUndoDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockUndoDataStoreMockRecorder
}

// MockUndoDataStoreMockRecorder is the mock recorder for MockUndoDataStore
type MockUndoDataStoreMockRecorder struct {
	mock *MockUndoDataStore
}

// NewMockUndoDataStore creates a new mock instance
func NewMockUndoDataStore(ctrl *gomock.Controller) *MockUndoDataStore {
	mock := &MockUndoDataStore{ctrl: ctrl}
	mock.recorder = &MockUndoDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUndoDataStore) EXPECT() *MockUndoDataStoreMockRecorder {
	return m.recorder
}

// GetUndoRecord mocks base method
func (m *MockUndoDataStore) GetUndoRecord(ctx context.Context, clusterID string) (*storage.NetworkPolicyApplicationUndoRecord, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUndoRecord", ctx, clusterID)
	ret0, _ := ret[0].(*storage.NetworkPolicyApplicationUndoRecord)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUndoRecord indicates an expected call of GetUndoRecord
func (mr *MockUndoDataStoreMockRecorder) GetUndoRecord(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUndoRecord", reflect.TypeOf((*MockUndoDataStore)(nil).GetUndoRecord), ctx, clusterID)
}

// UpsertUndoRecord mocks base method
func (m *MockUndoDataStore) UpsertUndoRecord(ctx context.Context, clusterID string, undoRecord *storage.NetworkPolicyApplicationUndoRecord) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertUndoRecord", ctx, clusterID, undoRecord)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertUndoRecord indicates an expected call of UpsertUndoRecord
func (mr *MockUndoDataStoreMockRecorder) UpsertUndoRecord(ctx, clusterID, undoRecord interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUndoRecord", reflect.TypeOf((*MockUndoDataStore)(nil).UpsertUndoRecord), ctx, clusterID, undoRecord)
}

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

// GetNetworkPolicy mocks base method
func (m *MockDataStore) GetNetworkPolicy(ctx context.Context, id string) (*storage.NetworkPolicy, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkPolicy", ctx, id)
	ret0, _ := ret[0].(*storage.NetworkPolicy)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNetworkPolicy indicates an expected call of GetNetworkPolicy
func (mr *MockDataStoreMockRecorder) GetNetworkPolicy(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkPolicy", reflect.TypeOf((*MockDataStore)(nil).GetNetworkPolicy), ctx, id)
}

// GetNetworkPolicies mocks base method
func (m *MockDataStore) GetNetworkPolicies(ctx context.Context, clusterID, namespace string) ([]*storage.NetworkPolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkPolicies", ctx, clusterID, namespace)
	ret0, _ := ret[0].([]*storage.NetworkPolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkPolicies indicates an expected call of GetNetworkPolicies
func (mr *MockDataStoreMockRecorder) GetNetworkPolicies(ctx, clusterID, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkPolicies", reflect.TypeOf((*MockDataStore)(nil).GetNetworkPolicies), ctx, clusterID, namespace)
}

// CountMatchingNetworkPolicies mocks base method
func (m *MockDataStore) CountMatchingNetworkPolicies(ctx context.Context, clusterID, namespace string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountMatchingNetworkPolicies", ctx, clusterID, namespace)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountMatchingNetworkPolicies indicates an expected call of CountMatchingNetworkPolicies
func (mr *MockDataStoreMockRecorder) CountMatchingNetworkPolicies(ctx, clusterID, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountMatchingNetworkPolicies", reflect.TypeOf((*MockDataStore)(nil).CountMatchingNetworkPolicies), ctx, clusterID, namespace)
}

// AddNetworkPolicy mocks base method
func (m *MockDataStore) AddNetworkPolicy(ctx context.Context, np *storage.NetworkPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNetworkPolicy", ctx, np)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNetworkPolicy indicates an expected call of AddNetworkPolicy
func (mr *MockDataStoreMockRecorder) AddNetworkPolicy(ctx, np interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNetworkPolicy", reflect.TypeOf((*MockDataStore)(nil).AddNetworkPolicy), ctx, np)
}

// UpdateNetworkPolicy mocks base method
func (m *MockDataStore) UpdateNetworkPolicy(ctx context.Context, np *storage.NetworkPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNetworkPolicy", ctx, np)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNetworkPolicy indicates an expected call of UpdateNetworkPolicy
func (mr *MockDataStoreMockRecorder) UpdateNetworkPolicy(ctx, np interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNetworkPolicy", reflect.TypeOf((*MockDataStore)(nil).UpdateNetworkPolicy), ctx, np)
}

// RemoveNetworkPolicy mocks base method
func (m *MockDataStore) RemoveNetworkPolicy(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNetworkPolicy", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNetworkPolicy indicates an expected call of RemoveNetworkPolicy
func (mr *MockDataStoreMockRecorder) RemoveNetworkPolicy(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNetworkPolicy", reflect.TypeOf((*MockDataStore)(nil).RemoveNetworkPolicy), ctx, id)
}

// GetUndoRecord mocks base method
func (m *MockDataStore) GetUndoRecord(ctx context.Context, clusterID string) (*storage.NetworkPolicyApplicationUndoRecord, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUndoRecord", ctx, clusterID)
	ret0, _ := ret[0].(*storage.NetworkPolicyApplicationUndoRecord)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUndoRecord indicates an expected call of GetUndoRecord
func (mr *MockDataStoreMockRecorder) GetUndoRecord(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUndoRecord", reflect.TypeOf((*MockDataStore)(nil).GetUndoRecord), ctx, clusterID)
}

// UpsertUndoRecord mocks base method
func (m *MockDataStore) UpsertUndoRecord(ctx context.Context, clusterID string, undoRecord *storage.NetworkPolicyApplicationUndoRecord) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertUndoRecord", ctx, clusterID, undoRecord)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertUndoRecord indicates an expected call of UpsertUndoRecord
func (mr *MockDataStoreMockRecorder) UpsertUndoRecord(ctx, clusterID, undoRecord interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUndoRecord", reflect.TypeOf((*MockDataStore)(nil).UpsertUndoRecord), ctx, clusterID, undoRecord)
}
