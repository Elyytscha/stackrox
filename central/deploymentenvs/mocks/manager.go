// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	deploymentenvs "github.com/stackrox/rox/central/deploymentenvs"
	reflect "reflect"
)

// MockListener is a mock of Listener interface
type MockListener struct {
	ctrl     *gomock.Controller
	recorder *MockListenerMockRecorder
}

// MockListenerMockRecorder is the mock recorder for MockListener
type MockListenerMockRecorder struct {
	mock *MockListener
}

// NewMockListener creates a new mock instance
func NewMockListener(ctrl *gomock.Controller) *MockListener {
	mock := &MockListener{ctrl: ctrl}
	mock.recorder = &MockListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockListener) EXPECT() *MockListenerMockRecorder {
	return m.recorder
}

// OnUpdate mocks base method
func (m *MockListener) OnUpdate(clusterID string, deploymentEnvs []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUpdate", clusterID, deploymentEnvs)
}

// OnUpdate indicates an expected call of OnUpdate
func (mr *MockListenerMockRecorder) OnUpdate(clusterID, deploymentEnvs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUpdate", reflect.TypeOf((*MockListener)(nil).OnUpdate), clusterID, deploymentEnvs)
}

// OnClusterMarkedInactive mocks base method
func (m *MockListener) OnClusterMarkedInactive(clusterID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnClusterMarkedInactive", clusterID)
}

// OnClusterMarkedInactive indicates an expected call of OnClusterMarkedInactive
func (mr *MockListenerMockRecorder) OnClusterMarkedInactive(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnClusterMarkedInactive", reflect.TypeOf((*MockListener)(nil).OnClusterMarkedInactive), clusterID)
}

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// UpdateDeploymentEnvironments mocks base method
func (m *MockManager) UpdateDeploymentEnvironments(clusterID string, deploymentEnvs []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateDeploymentEnvironments", clusterID, deploymentEnvs)
}

// UpdateDeploymentEnvironments indicates an expected call of UpdateDeploymentEnvironments
func (mr *MockManagerMockRecorder) UpdateDeploymentEnvironments(clusterID, deploymentEnvs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeploymentEnvironments", reflect.TypeOf((*MockManager)(nil).UpdateDeploymentEnvironments), clusterID, deploymentEnvs)
}

// MarkClusterInactive mocks base method
func (m *MockManager) MarkClusterInactive(clusterID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkClusterInactive", clusterID)
}

// MarkClusterInactive indicates an expected call of MarkClusterInactive
func (mr *MockManagerMockRecorder) MarkClusterInactive(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkClusterInactive", reflect.TypeOf((*MockManager)(nil).MarkClusterInactive), clusterID)
}

// RegisterListener mocks base method
func (m *MockManager) RegisterListener(listener deploymentenvs.Listener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterListener", listener)
}

// RegisterListener indicates an expected call of RegisterListener
func (mr *MockManagerMockRecorder) RegisterListener(listener interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterListener", reflect.TypeOf((*MockManager)(nil).RegisterListener), listener)
}

// UnregisterListener mocks base method
func (m *MockManager) UnregisterListener(listener deploymentenvs.Listener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnregisterListener", listener)
}

// UnregisterListener indicates an expected call of UnregisterListener
func (mr *MockManagerMockRecorder) UnregisterListener(listener interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterListener", reflect.TypeOf((*MockManager)(nil).UnregisterListener), listener)
}

// GetDeploymentEnvironmentsByClusterID mocks base method
func (m *MockManager) GetDeploymentEnvironmentsByClusterID(block bool) map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeploymentEnvironmentsByClusterID", block)
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// GetDeploymentEnvironmentsByClusterID indicates an expected call of GetDeploymentEnvironmentsByClusterID
func (mr *MockManagerMockRecorder) GetDeploymentEnvironmentsByClusterID(block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeploymentEnvironmentsByClusterID", reflect.TypeOf((*MockManager)(nil).GetDeploymentEnvironmentsByClusterID), block)
}
