// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
	time "time"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetCluster mocks base method
func (m *MockStore) GetCluster(id string) (*storage.Cluster, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", id)
	ret0, _ := ret[0].(*storage.Cluster)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCluster indicates an expected call of GetCluster
func (mr *MockStoreMockRecorder) GetCluster(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockStore)(nil).GetCluster), id)
}

// GetClusters mocks base method
func (m *MockStore) GetClusters() ([]*storage.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusters")
	ret0, _ := ret[0].([]*storage.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusters indicates an expected call of GetClusters
func (mr *MockStoreMockRecorder) GetClusters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusters", reflect.TypeOf((*MockStore)(nil).GetClusters))
}

// CountClusters mocks base method
func (m *MockStore) CountClusters() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountClusters")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountClusters indicates an expected call of CountClusters
func (mr *MockStoreMockRecorder) CountClusters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountClusters", reflect.TypeOf((*MockStore)(nil).CountClusters))
}

// AddCluster mocks base method
func (m *MockStore) AddCluster(cluster *storage.Cluster) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCluster", cluster)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCluster indicates an expected call of AddCluster
func (mr *MockStoreMockRecorder) AddCluster(cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCluster", reflect.TypeOf((*MockStore)(nil).AddCluster), cluster)
}

// UpdateCluster mocks base method
func (m *MockStore) UpdateCluster(cluster *storage.Cluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCluster", cluster)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCluster indicates an expected call of UpdateCluster
func (mr *MockStoreMockRecorder) UpdateCluster(cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCluster", reflect.TypeOf((*MockStore)(nil).UpdateCluster), cluster)
}

// RemoveCluster mocks base method
func (m *MockStore) RemoveCluster(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCluster", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCluster indicates an expected call of RemoveCluster
func (mr *MockStoreMockRecorder) RemoveCluster(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCluster", reflect.TypeOf((*MockStore)(nil).RemoveCluster), id)
}

// UpdateClusterContactTimes mocks base method
func (m *MockStore) UpdateClusterContactTimes(t time.Time, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{t}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateClusterContactTimes", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterContactTimes indicates an expected call of UpdateClusterContactTimes
func (mr *MockStoreMockRecorder) UpdateClusterContactTimes(t interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{t}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterContactTimes", reflect.TypeOf((*MockStore)(nil).UpdateClusterContactTimes), varargs...)
}

// UpdateClusterStatus mocks base method
func (m *MockStore) UpdateClusterStatus(id string, status *storage.ClusterStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterStatus", id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterStatus indicates an expected call of UpdateClusterStatus
func (mr *MockStoreMockRecorder) UpdateClusterStatus(id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterStatus", reflect.TypeOf((*MockStore)(nil).UpdateClusterStatus), id, status)
}

// UpdateClusterUpgradeStatus mocks base method
func (m *MockStore) UpdateClusterUpgradeStatus(id string, status *storage.ClusterUpgradeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterUpgradeStatus", id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterUpgradeStatus indicates an expected call of UpdateClusterUpgradeStatus
func (mr *MockStoreMockRecorder) UpdateClusterUpgradeStatus(id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterUpgradeStatus", reflect.TypeOf((*MockStore)(nil).UpdateClusterUpgradeStatus), id, status)
}
