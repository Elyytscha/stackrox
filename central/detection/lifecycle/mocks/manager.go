// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/stackrox/rox/central/sensor/service/common"
	storage "github.com/stackrox/rox/generated/storage"
	enricher "github.com/stackrox/rox/pkg/images/enricher"
	reflect "reflect"
)

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

// IndicatorAdded mocks base method
func (m *MockManager) IndicatorAdded(indicator *storage.ProcessIndicator, injector common.MessageInjector) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndicatorAdded", indicator, injector)
	ret0, _ := ret[0].(error)
	return ret0
}

// IndicatorAdded indicates an expected call of IndicatorAdded
func (mr *MockManagerMockRecorder) IndicatorAdded(indicator, injector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndicatorAdded", reflect.TypeOf((*MockManager)(nil).IndicatorAdded), indicator, injector)
}

// DeploymentUpdated mocks base method
func (m *MockManager) DeploymentUpdated(ctx enricher.EnrichmentContext, deployment *storage.Deployment, injector common.MessageInjector) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeploymentUpdated", ctx, deployment, injector)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeploymentUpdated indicates an expected call of DeploymentUpdated
func (mr *MockManagerMockRecorder) DeploymentUpdated(ctx, deployment, injector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeploymentUpdated", reflect.TypeOf((*MockManager)(nil).DeploymentUpdated), ctx, deployment, injector)
}

// UpsertPolicy mocks base method
func (m *MockManager) UpsertPolicy(policy *storage.Policy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertPolicy", policy)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertPolicy indicates an expected call of UpsertPolicy
func (mr *MockManagerMockRecorder) UpsertPolicy(policy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertPolicy", reflect.TypeOf((*MockManager)(nil).UpsertPolicy), policy)
}

// RecompilePolicy mocks base method
func (m *MockManager) RecompilePolicy(policy *storage.Policy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecompilePolicy", policy)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecompilePolicy indicates an expected call of RecompilePolicy
func (mr *MockManagerMockRecorder) RecompilePolicy(policy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecompilePolicy", reflect.TypeOf((*MockManager)(nil).RecompilePolicy), policy)
}

// DeploymentRemoved mocks base method
func (m *MockManager) DeploymentRemoved(deployment *storage.Deployment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeploymentRemoved", deployment)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeploymentRemoved indicates an expected call of DeploymentRemoved
func (mr *MockManagerMockRecorder) DeploymentRemoved(deployment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeploymentRemoved", reflect.TypeOf((*MockManager)(nil).DeploymentRemoved), deployment)
}

// RemovePolicy mocks base method
func (m *MockManager) RemovePolicy(policyID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePolicy", policyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePolicy indicates an expected call of RemovePolicy
func (mr *MockManagerMockRecorder) RemovePolicy(policyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePolicy", reflect.TypeOf((*MockManager)(nil).RemovePolicy), policyID)
}
