// Code generated by MockGen. DO NOT EDIT.
// Source: dispatcher.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	central "github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/sensor/kubernetes/eventpipeline/message"
	resources "github.com/stackrox/rox/sensor/kubernetes/listener/resources"
)

// MockDispatcher is a mock of Dispatcher interface.
type MockDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockDispatcherMockRecorder
}

// MockDispatcherMockRecorder is the mock recorder for MockDispatcher.
type MockDispatcherMockRecorder struct {
	mock *MockDispatcher
}

// NewMockDispatcher creates a new mock instance.
func NewMockDispatcher(ctrl *gomock.Controller) *MockDispatcher {
	mock := &MockDispatcher{ctrl: ctrl}
	mock.recorder = &MockDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDispatcher) EXPECT() *MockDispatcherMockRecorder {
	return m.recorder
}

// ProcessEvent mocks base method.
func (m *MockDispatcher) ProcessEvent(obj, oldObj interface{}, action central.ResourceAction) *message.ResourceEvent {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessEvent", obj, oldObj, action)
	ret0, _ := ret[0].(*message.ResourceEvent)
	return ret0
}

// ProcessEvent indicates an expected call of ProcessEvent.
func (mr *MockDispatcherMockRecorder) ProcessEvent(obj, oldObj, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessEvent", reflect.TypeOf((*MockDispatcher)(nil).ProcessEvent), obj, oldObj, action)
}

// MockDispatcherRegistry is a mock of DispatcherRegistry interface.
type MockDispatcherRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockDispatcherRegistryMockRecorder
}

// MockDispatcherRegistryMockRecorder is the mock recorder for MockDispatcherRegistry.
type MockDispatcherRegistryMockRecorder struct {
	mock *MockDispatcherRegistry
}

// NewMockDispatcherRegistry creates a new mock instance.
func NewMockDispatcherRegistry(ctrl *gomock.Controller) *MockDispatcherRegistry {
	mock := &MockDispatcherRegistry{ctrl: ctrl}
	mock.recorder = &MockDispatcherRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDispatcherRegistry) EXPECT() *MockDispatcherRegistryMockRecorder {
	return m.recorder
}

// ForClusterOperators mocks base method.
func (m *MockDispatcherRegistry) ForClusterOperators() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForClusterOperators")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForClusterOperators indicates an expected call of ForClusterOperators.
func (mr *MockDispatcherRegistryMockRecorder) ForClusterOperators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForClusterOperators", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForClusterOperators))
}

// ForComplianceOperatorProfiles mocks base method.
func (m *MockDispatcherRegistry) ForComplianceOperatorProfiles() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForComplianceOperatorProfiles")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForComplianceOperatorProfiles indicates an expected call of ForComplianceOperatorProfiles.
func (mr *MockDispatcherRegistryMockRecorder) ForComplianceOperatorProfiles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForComplianceOperatorProfiles", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForComplianceOperatorProfiles))
}

// ForComplianceOperatorResults mocks base method.
func (m *MockDispatcherRegistry) ForComplianceOperatorResults() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForComplianceOperatorResults")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForComplianceOperatorResults indicates an expected call of ForComplianceOperatorResults.
func (mr *MockDispatcherRegistryMockRecorder) ForComplianceOperatorResults() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForComplianceOperatorResults", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForComplianceOperatorResults))
}

// ForComplianceOperatorRules mocks base method.
func (m *MockDispatcherRegistry) ForComplianceOperatorRules() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForComplianceOperatorRules")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForComplianceOperatorRules indicates an expected call of ForComplianceOperatorRules.
func (mr *MockDispatcherRegistryMockRecorder) ForComplianceOperatorRules() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForComplianceOperatorRules", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForComplianceOperatorRules))
}

// ForComplianceOperatorScanSettingBindings mocks base method.
func (m *MockDispatcherRegistry) ForComplianceOperatorScanSettingBindings() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForComplianceOperatorScanSettingBindings")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForComplianceOperatorScanSettingBindings indicates an expected call of ForComplianceOperatorScanSettingBindings.
func (mr *MockDispatcherRegistryMockRecorder) ForComplianceOperatorScanSettingBindings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForComplianceOperatorScanSettingBindings", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForComplianceOperatorScanSettingBindings))
}

// ForComplianceOperatorScans mocks base method.
func (m *MockDispatcherRegistry) ForComplianceOperatorScans() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForComplianceOperatorScans")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForComplianceOperatorScans indicates an expected call of ForComplianceOperatorScans.
func (mr *MockDispatcherRegistryMockRecorder) ForComplianceOperatorScans() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForComplianceOperatorScans", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForComplianceOperatorScans))
}

// ForComplianceOperatorTailoredProfiles mocks base method.
func (m *MockDispatcherRegistry) ForComplianceOperatorTailoredProfiles() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForComplianceOperatorTailoredProfiles")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForComplianceOperatorTailoredProfiles indicates an expected call of ForComplianceOperatorTailoredProfiles.
func (mr *MockDispatcherRegistryMockRecorder) ForComplianceOperatorTailoredProfiles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForComplianceOperatorTailoredProfiles", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForComplianceOperatorTailoredProfiles))
}

// ForDeployments mocks base method.
func (m *MockDispatcherRegistry) ForDeployments(deploymentType string) resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForDeployments", deploymentType)
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForDeployments indicates an expected call of ForDeployments.
func (mr *MockDispatcherRegistryMockRecorder) ForDeployments(deploymentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForDeployments", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForDeployments), deploymentType)
}

// ForJobs mocks base method.
func (m *MockDispatcherRegistry) ForJobs() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForJobs")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForJobs indicates an expected call of ForJobs.
func (mr *MockDispatcherRegistryMockRecorder) ForJobs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForJobs", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForJobs))
}

// ForNamespaces mocks base method.
func (m *MockDispatcherRegistry) ForNamespaces() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForNamespaces")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForNamespaces indicates an expected call of ForNamespaces.
func (mr *MockDispatcherRegistryMockRecorder) ForNamespaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForNamespaces", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForNamespaces))
}

// ForNetworkPolicies mocks base method.
func (m *MockDispatcherRegistry) ForNetworkPolicies() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForNetworkPolicies")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForNetworkPolicies indicates an expected call of ForNetworkPolicies.
func (mr *MockDispatcherRegistryMockRecorder) ForNetworkPolicies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForNetworkPolicies", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForNetworkPolicies))
}

// ForNodes mocks base method.
func (m *MockDispatcherRegistry) ForNodes() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForNodes")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForNodes indicates an expected call of ForNodes.
func (mr *MockDispatcherRegistryMockRecorder) ForNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForNodes", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForNodes))
}

// ForOpenshiftRoutes mocks base method.
func (m *MockDispatcherRegistry) ForOpenshiftRoutes() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForOpenshiftRoutes")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForOpenshiftRoutes indicates an expected call of ForOpenshiftRoutes.
func (mr *MockDispatcherRegistryMockRecorder) ForOpenshiftRoutes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForOpenshiftRoutes", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForOpenshiftRoutes))
}

// ForRBAC mocks base method.
func (m *MockDispatcherRegistry) ForRBAC() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForRBAC")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForRBAC indicates an expected call of ForRBAC.
func (mr *MockDispatcherRegistryMockRecorder) ForRBAC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForRBAC", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForRBAC))
}

// ForSecrets mocks base method.
func (m *MockDispatcherRegistry) ForSecrets() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForSecrets")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForSecrets indicates an expected call of ForSecrets.
func (mr *MockDispatcherRegistryMockRecorder) ForSecrets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForSecrets", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForSecrets))
}

// ForServiceAccounts mocks base method.
func (m *MockDispatcherRegistry) ForServiceAccounts() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForServiceAccounts")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForServiceAccounts indicates an expected call of ForServiceAccounts.
func (mr *MockDispatcherRegistryMockRecorder) ForServiceAccounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForServiceAccounts", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForServiceAccounts))
}

// ForServices mocks base method.
func (m *MockDispatcherRegistry) ForServices() resources.Dispatcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForServices")
	ret0, _ := ret[0].(resources.Dispatcher)
	return ret0
}

// ForServices indicates an expected call of ForServices.
func (mr *MockDispatcherRegistryMockRecorder) ForServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForServices", reflect.TypeOf((*MockDispatcherRegistry)(nil).ForServices))
}
