// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockClusterManager is a mock of ClusterManager interface.
type MockClusterManager struct {
	ctrl     *gomock.Controller
	recorder *MockClusterManagerMockRecorder
}

// MockClusterManagerMockRecorder is the mock recorder for MockClusterManager.
type MockClusterManagerMockRecorder struct {
	mock *MockClusterManager
}

// NewMockClusterManager creates a new mock instance.
func NewMockClusterManager(ctrl *gomock.Controller) *MockClusterManager {
	mock := &MockClusterManager{ctrl: ctrl}
	mock.recorder = &MockClusterManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterManager) EXPECT() *MockClusterManagerMockRecorder {
	return m.recorder
}

// GetCluster mocks base method.
func (m *MockClusterManager) GetCluster(ctx context.Context, id string) (*storage.Cluster, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", ctx, id)
	ret0, _ := ret[0].(*storage.Cluster)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCluster indicates an expected call of GetCluster.
func (mr *MockClusterManagerMockRecorder) GetCluster(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockClusterManager)(nil).GetCluster), ctx, id)
}

// GetClusters mocks base method.
func (m *MockClusterManager) GetClusters(ctx context.Context) ([]*storage.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusters", ctx)
	ret0, _ := ret[0].([]*storage.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusters indicates an expected call of GetClusters.
func (mr *MockClusterManagerMockRecorder) GetClusters(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusters", reflect.TypeOf((*MockClusterManager)(nil).GetClusters), ctx)
}

// UpdateClusterHealth mocks base method.
func (m *MockClusterManager) UpdateClusterHealth(ctx context.Context, id string, status *storage.ClusterHealthStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterHealth", ctx, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterHealth indicates an expected call of UpdateClusterHealth.
func (mr *MockClusterManagerMockRecorder) UpdateClusterHealth(ctx, id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterHealth", reflect.TypeOf((*MockClusterManager)(nil).UpdateClusterHealth), ctx, id, status)
}

// UpdateClusterUpgradeStatus mocks base method.
func (m *MockClusterManager) UpdateClusterUpgradeStatus(ctx context.Context, clusterID string, status *storage.ClusterUpgradeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterUpgradeStatus", ctx, clusterID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterUpgradeStatus indicates an expected call of UpdateClusterUpgradeStatus.
func (mr *MockClusterManagerMockRecorder) UpdateClusterUpgradeStatus(ctx, clusterID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterUpgradeStatus", reflect.TypeOf((*MockClusterManager)(nil).UpdateClusterUpgradeStatus), ctx, clusterID, status)
}

// UpdateSensorDeploymentIdentification mocks base method.
func (m *MockClusterManager) UpdateSensorDeploymentIdentification(ctx context.Context, clusterID string, identification *storage.SensorDeploymentIdentification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSensorDeploymentIdentification", ctx, clusterID, identification)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSensorDeploymentIdentification indicates an expected call of UpdateSensorDeploymentIdentification.
func (mr *MockClusterManagerMockRecorder) UpdateSensorDeploymentIdentification(ctx, clusterID, identification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSensorDeploymentIdentification", reflect.TypeOf((*MockClusterManager)(nil).UpdateSensorDeploymentIdentification), ctx, clusterID, identification)
}

// MockPolicyManager is a mock of PolicyManager interface.
type MockPolicyManager struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyManagerMockRecorder
}

// MockPolicyManagerMockRecorder is the mock recorder for MockPolicyManager.
type MockPolicyManagerMockRecorder struct {
	mock *MockPolicyManager
}

// NewMockPolicyManager creates a new mock instance.
func NewMockPolicyManager(ctrl *gomock.Controller) *MockPolicyManager {
	mock := &MockPolicyManager{ctrl: ctrl}
	mock.recorder = &MockPolicyManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPolicyManager) EXPECT() *MockPolicyManagerMockRecorder {
	return m.recorder
}

// GetAllPolicies mocks base method.
func (m *MockPolicyManager) GetAllPolicies(ctx context.Context) ([]*storage.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPolicies", ctx)
	ret0, _ := ret[0].([]*storage.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPolicies indicates an expected call of GetAllPolicies.
func (mr *MockPolicyManagerMockRecorder) GetAllPolicies(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPolicies", reflect.TypeOf((*MockPolicyManager)(nil).GetAllPolicies), ctx)
}

// MockProcessBaselineManager is a mock of ProcessBaselineManager interface.
type MockProcessBaselineManager struct {
	ctrl     *gomock.Controller
	recorder *MockProcessBaselineManagerMockRecorder
}

// MockProcessBaselineManagerMockRecorder is the mock recorder for MockProcessBaselineManager.
type MockProcessBaselineManagerMockRecorder struct {
	mock *MockProcessBaselineManager
}

// NewMockProcessBaselineManager creates a new mock instance.
func NewMockProcessBaselineManager(ctrl *gomock.Controller) *MockProcessBaselineManager {
	mock := &MockProcessBaselineManager{ctrl: ctrl}
	mock.recorder = &MockProcessBaselineManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessBaselineManager) EXPECT() *MockProcessBaselineManagerMockRecorder {
	return m.recorder
}

// WalkAll mocks base method.
func (m *MockProcessBaselineManager) WalkAll(ctx context.Context, fn func(*storage.ProcessBaseline) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalkAll", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WalkAll indicates an expected call of WalkAll.
func (mr *MockProcessBaselineManagerMockRecorder) WalkAll(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalkAll", reflect.TypeOf((*MockProcessBaselineManager)(nil).WalkAll), ctx, fn)
}

// MockNetworkBaselineManager is a mock of NetworkBaselineManager interface.
type MockNetworkBaselineManager struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkBaselineManagerMockRecorder
}

// MockNetworkBaselineManagerMockRecorder is the mock recorder for MockNetworkBaselineManager.
type MockNetworkBaselineManagerMockRecorder struct {
	mock *MockNetworkBaselineManager
}

// NewMockNetworkBaselineManager creates a new mock instance.
func NewMockNetworkBaselineManager(ctrl *gomock.Controller) *MockNetworkBaselineManager {
	mock := &MockNetworkBaselineManager{ctrl: ctrl}
	mock.recorder = &MockNetworkBaselineManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkBaselineManager) EXPECT() *MockNetworkBaselineManagerMockRecorder {
	return m.recorder
}

// Walk mocks base method.
func (m *MockNetworkBaselineManager) Walk(ctx context.Context, fn func(*storage.NetworkBaseline) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Walk", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Walk indicates an expected call of Walk.
func (mr *MockNetworkBaselineManagerMockRecorder) Walk(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Walk", reflect.TypeOf((*MockNetworkBaselineManager)(nil).Walk), ctx, fn)
}

// MockDelegatedRegistryConfigManager is a mock of DelegatedRegistryConfigManager interface.
type MockDelegatedRegistryConfigManager struct {
	ctrl     *gomock.Controller
	recorder *MockDelegatedRegistryConfigManagerMockRecorder
}

// MockDelegatedRegistryConfigManagerMockRecorder is the mock recorder for MockDelegatedRegistryConfigManager.
type MockDelegatedRegistryConfigManagerMockRecorder struct {
	mock *MockDelegatedRegistryConfigManager
}

// NewMockDelegatedRegistryConfigManager creates a new mock instance.
func NewMockDelegatedRegistryConfigManager(ctrl *gomock.Controller) *MockDelegatedRegistryConfigManager {
	mock := &MockDelegatedRegistryConfigManager{ctrl: ctrl}
	mock.recorder = &MockDelegatedRegistryConfigManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDelegatedRegistryConfigManager) EXPECT() *MockDelegatedRegistryConfigManagerMockRecorder {
	return m.recorder
}

// GetConfig mocks base method.
func (m *MockDelegatedRegistryConfigManager) GetConfig(ctx context.Context) (*storage.DelegatedRegistryConfig, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig", ctx)
	ret0, _ := ret[0].(*storage.DelegatedRegistryConfig)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockDelegatedRegistryConfigManagerMockRecorder) GetConfig(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockDelegatedRegistryConfigManager)(nil).GetConfig), ctx)
}

// MockImageIntegrationManager is a mock of ImageIntegrationManager interface.
type MockImageIntegrationManager struct {
	ctrl     *gomock.Controller
	recorder *MockImageIntegrationManagerMockRecorder
}

// MockImageIntegrationManagerMockRecorder is the mock recorder for MockImageIntegrationManager.
type MockImageIntegrationManagerMockRecorder struct {
	mock *MockImageIntegrationManager
}

// NewMockImageIntegrationManager creates a new mock instance.
func NewMockImageIntegrationManager(ctrl *gomock.Controller) *MockImageIntegrationManager {
	mock := &MockImageIntegrationManager{ctrl: ctrl}
	mock.recorder = &MockImageIntegrationManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageIntegrationManager) EXPECT() *MockImageIntegrationManagerMockRecorder {
	return m.recorder
}

// GetImageIntegrations mocks base method.
func (m *MockImageIntegrationManager) GetImageIntegrations(ctx context.Context, integration *v1.GetImageIntegrationsRequest) ([]*storage.ImageIntegration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageIntegrations", ctx, integration)
	ret0, _ := ret[0].([]*storage.ImageIntegration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageIntegrations indicates an expected call of GetImageIntegrations.
func (mr *MockImageIntegrationManagerMockRecorder) GetImageIntegrations(ctx, integration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageIntegrations", reflect.TypeOf((*MockImageIntegrationManager)(nil).GetImageIntegrations), ctx, integration)
}

// MockNetworkEntityManager is a mock of NetworkEntityManager interface.
type MockNetworkEntityManager struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkEntityManagerMockRecorder
}

// MockNetworkEntityManagerMockRecorder is the mock recorder for MockNetworkEntityManager.
type MockNetworkEntityManagerMockRecorder struct {
	mock *MockNetworkEntityManager
}

// NewMockNetworkEntityManager creates a new mock instance.
func NewMockNetworkEntityManager(ctrl *gomock.Controller) *MockNetworkEntityManager {
	mock := &MockNetworkEntityManager{ctrl: ctrl}
	mock.recorder = &MockNetworkEntityManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkEntityManager) EXPECT() *MockNetworkEntityManagerMockRecorder {
	return m.recorder
}

// GetAllEntitiesForCluster mocks base method.
func (m *MockNetworkEntityManager) GetAllEntitiesForCluster(ctx context.Context, clusterID string) ([]*storage.NetworkEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllEntitiesForCluster", ctx, clusterID)
	ret0, _ := ret[0].([]*storage.NetworkEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllEntitiesForCluster indicates an expected call of GetAllEntitiesForCluster.
func (mr *MockNetworkEntityManagerMockRecorder) GetAllEntitiesForCluster(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllEntitiesForCluster", reflect.TypeOf((*MockNetworkEntityManager)(nil).GetAllEntitiesForCluster), ctx, clusterID)
}
