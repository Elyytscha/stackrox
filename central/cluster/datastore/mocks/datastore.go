// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/datastore.go -source datastore.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	central "github.com/stackrox/rox/generated/internalapi/central"
	storage "github.com/stackrox/rox/generated/storage"
	concurrency "github.com/stackrox/rox/pkg/concurrency"
	effectiveaccessscope "github.com/stackrox/rox/pkg/sac/effectiveaccessscope"
	search "github.com/stackrox/rox/pkg/search"
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

// AddCluster mocks base method.
func (m *MockDataStore) AddCluster(ctx context.Context, cluster *storage.Cluster) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCluster", ctx, cluster)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCluster indicates an expected call of AddCluster.
func (mr *MockDataStoreMockRecorder) AddCluster(ctx, cluster any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCluster", reflect.TypeOf((*MockDataStore)(nil).AddCluster), ctx, cluster)
}

// Count mocks base method.
func (m *MockDataStore) Count(ctx context.Context, q *v1.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, q)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockDataStoreMockRecorder) Count(ctx, q any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDataStore)(nil).Count), ctx, q)
}

// CountClusters mocks base method.
func (m *MockDataStore) CountClusters(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountClusters", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountClusters indicates an expected call of CountClusters.
func (mr *MockDataStoreMockRecorder) CountClusters(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountClusters", reflect.TypeOf((*MockDataStore)(nil).CountClusters), ctx)
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
func (mr *MockDataStoreMockRecorder) Exists(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockDataStore)(nil).Exists), ctx, id)
}

// GetCluster mocks base method.
func (m *MockDataStore) GetCluster(ctx context.Context, id string) (*storage.Cluster, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", ctx, id)
	ret0, _ := ret[0].(*storage.Cluster)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCluster indicates an expected call of GetCluster.
func (mr *MockDataStoreMockRecorder) GetCluster(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockDataStore)(nil).GetCluster), ctx, id)
}

// GetClusterName mocks base method.
func (m *MockDataStore) GetClusterName(ctx context.Context, id string) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterName", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetClusterName indicates an expected call of GetClusterName.
func (mr *MockDataStoreMockRecorder) GetClusterName(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterName", reflect.TypeOf((*MockDataStore)(nil).GetClusterName), ctx, id)
}

// GetClusters mocks base method.
func (m *MockDataStore) GetClusters(ctx context.Context) ([]*storage.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusters", ctx)
	ret0, _ := ret[0].([]*storage.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusters indicates an expected call of GetClusters.
func (mr *MockDataStoreMockRecorder) GetClusters(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusters", reflect.TypeOf((*MockDataStore)(nil).GetClusters), ctx)
}

// GetClustersForSAC mocks base method.
func (m *MockDataStore) GetClustersForSAC(ctx context.Context) ([]effectiveaccessscope.ClusterForSAC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClustersForSAC", ctx)
	ret0, _ := ret[0].([]effectiveaccessscope.ClusterForSAC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClustersForSAC indicates an expected call of GetClustersForSAC.
func (mr *MockDataStoreMockRecorder) GetClustersForSAC(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClustersForSAC", reflect.TypeOf((*MockDataStore)(nil).GetClustersForSAC), ctx)
}

// LookupOrCreateClusterFromConfig mocks base method.
func (m *MockDataStore) LookupOrCreateClusterFromConfig(ctx context.Context, clusterID, bundleID string, hello *central.SensorHello) (*storage.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupOrCreateClusterFromConfig", ctx, clusterID, bundleID, hello)
	ret0, _ := ret[0].(*storage.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupOrCreateClusterFromConfig indicates an expected call of LookupOrCreateClusterFromConfig.
func (mr *MockDataStoreMockRecorder) LookupOrCreateClusterFromConfig(ctx, clusterID, bundleID, hello any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupOrCreateClusterFromConfig", reflect.TypeOf((*MockDataStore)(nil).LookupOrCreateClusterFromConfig), ctx, clusterID, bundleID, hello)
}

// RemoveCluster mocks base method.
func (m *MockDataStore) RemoveCluster(ctx context.Context, id string, done *concurrency.Signal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCluster", ctx, id, done)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCluster indicates an expected call of RemoveCluster.
func (mr *MockDataStoreMockRecorder) RemoveCluster(ctx, id, done any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCluster", reflect.TypeOf((*MockDataStore)(nil).RemoveCluster), ctx, id, done)
}

// Search mocks base method.
func (m *MockDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockDataStoreMockRecorder) Search(ctx, q any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// SearchRawClusters mocks base method.
func (m *MockDataStore) SearchRawClusters(ctx context.Context, q *v1.Query) ([]*storage.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawClusters", ctx, q)
	ret0, _ := ret[0].([]*storage.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawClusters indicates an expected call of SearchRawClusters.
func (mr *MockDataStoreMockRecorder) SearchRawClusters(ctx, q any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawClusters", reflect.TypeOf((*MockDataStore)(nil).SearchRawClusters), ctx, q)
}

// SearchResults mocks base method.
func (m *MockDataStore) SearchResults(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchResults", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchResults indicates an expected call of SearchResults.
func (mr *MockDataStoreMockRecorder) SearchResults(ctx, q any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchResults", reflect.TypeOf((*MockDataStore)(nil).SearchResults), ctx, q)
}

// UpdateAuditLogFileStates mocks base method.
func (m *MockDataStore) UpdateAuditLogFileStates(ctx context.Context, id string, states map[string]*storage.AuditLogFileState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAuditLogFileStates", ctx, id, states)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAuditLogFileStates indicates an expected call of UpdateAuditLogFileStates.
func (mr *MockDataStoreMockRecorder) UpdateAuditLogFileStates(ctx, id, states any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAuditLogFileStates", reflect.TypeOf((*MockDataStore)(nil).UpdateAuditLogFileStates), ctx, id, states)
}

// UpdateCluster mocks base method.
func (m *MockDataStore) UpdateCluster(ctx context.Context, cluster *storage.Cluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCluster", ctx, cluster)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCluster indicates an expected call of UpdateCluster.
func (mr *MockDataStoreMockRecorder) UpdateCluster(ctx, cluster any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCluster", reflect.TypeOf((*MockDataStore)(nil).UpdateCluster), ctx, cluster)
}

// UpdateClusterCertExpiryStatus mocks base method.
func (m *MockDataStore) UpdateClusterCertExpiryStatus(ctx context.Context, id string, clusterCertExpiryStatus *storage.ClusterCertExpiryStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterCertExpiryStatus", ctx, id, clusterCertExpiryStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterCertExpiryStatus indicates an expected call of UpdateClusterCertExpiryStatus.
func (mr *MockDataStoreMockRecorder) UpdateClusterCertExpiryStatus(ctx, id, clusterCertExpiryStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterCertExpiryStatus", reflect.TypeOf((*MockDataStore)(nil).UpdateClusterCertExpiryStatus), ctx, id, clusterCertExpiryStatus)
}

// UpdateClusterHealth mocks base method.
func (m *MockDataStore) UpdateClusterHealth(ctx context.Context, id string, clusterHealthStatus *storage.ClusterHealthStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterHealth", ctx, id, clusterHealthStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterHealth indicates an expected call of UpdateClusterHealth.
func (mr *MockDataStoreMockRecorder) UpdateClusterHealth(ctx, id, clusterHealthStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterHealth", reflect.TypeOf((*MockDataStore)(nil).UpdateClusterHealth), ctx, id, clusterHealthStatus)
}

// UpdateClusterStatus mocks base method.
func (m *MockDataStore) UpdateClusterStatus(ctx context.Context, id string, status *storage.ClusterStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterStatus", ctx, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterStatus indicates an expected call of UpdateClusterStatus.
func (mr *MockDataStoreMockRecorder) UpdateClusterStatus(ctx, id, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterStatus", reflect.TypeOf((*MockDataStore)(nil).UpdateClusterStatus), ctx, id, status)
}

// UpdateClusterUpgradeStatus mocks base method.
func (m *MockDataStore) UpdateClusterUpgradeStatus(ctx context.Context, id string, clusterUpgradeStatus *storage.ClusterUpgradeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterUpgradeStatus", ctx, id, clusterUpgradeStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterUpgradeStatus indicates an expected call of UpdateClusterUpgradeStatus.
func (mr *MockDataStoreMockRecorder) UpdateClusterUpgradeStatus(ctx, id, clusterUpgradeStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterUpgradeStatus", reflect.TypeOf((*MockDataStore)(nil).UpdateClusterUpgradeStatus), ctx, id, clusterUpgradeStatus)
}

// UpdateSensorDeploymentIdentification mocks base method.
func (m *MockDataStore) UpdateSensorDeploymentIdentification(ctx context.Context, id string, identification *storage.SensorDeploymentIdentification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSensorDeploymentIdentification", ctx, id, identification)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSensorDeploymentIdentification indicates an expected call of UpdateSensorDeploymentIdentification.
func (mr *MockDataStoreMockRecorder) UpdateSensorDeploymentIdentification(ctx, id, identification any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSensorDeploymentIdentification", reflect.TypeOf((*MockDataStore)(nil).UpdateSensorDeploymentIdentification), ctx, id, identification)
}
