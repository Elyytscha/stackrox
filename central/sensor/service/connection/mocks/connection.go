// Code generated by MockGen. DO NOT EDIT.
// Source: connection.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	scrape "github.com/stackrox/rox/central/scrape"
	networkentities "github.com/stackrox/rox/central/sensor/networkentities"
	networkpolicies "github.com/stackrox/rox/central/sensor/networkpolicies"
	telemetry "github.com/stackrox/rox/central/sensor/telemetry"
	central "github.com/stackrox/rox/generated/internalapi/central"
	centralsensor "github.com/stackrox/rox/pkg/centralsensor"
	concurrency "github.com/stackrox/rox/pkg/concurrency"
	gomock "go.uber.org/mock/gomock"
)

// MockSensorConnection is a mock of SensorConnection interface.
type MockSensorConnection struct {
	ctrl     *gomock.Controller
	recorder *MockSensorConnectionMockRecorder
}

// MockSensorConnectionMockRecorder is the mock recorder for MockSensorConnection.
type MockSensorConnectionMockRecorder struct {
	mock *MockSensorConnection
}

// NewMockSensorConnection creates a new mock instance.
func NewMockSensorConnection(ctrl *gomock.Controller) *MockSensorConnection {
	mock := &MockSensorConnection{ctrl: ctrl}
	mock.recorder = &MockSensorConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSensorConnection) EXPECT() *MockSensorConnectionMockRecorder {
	return m.recorder
}

// CheckAutoUpgradeSupport mocks base method.
func (m *MockSensorConnection) CheckAutoUpgradeSupport() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAutoUpgradeSupport")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckAutoUpgradeSupport indicates an expected call of CheckAutoUpgradeSupport.
func (mr *MockSensorConnectionMockRecorder) CheckAutoUpgradeSupport() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAutoUpgradeSupport", reflect.TypeOf((*MockSensorConnection)(nil).CheckAutoUpgradeSupport))
}

// ClusterID mocks base method.
func (m *MockSensorConnection) ClusterID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClusterID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ClusterID indicates an expected call of ClusterID.
func (mr *MockSensorConnectionMockRecorder) ClusterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClusterID", reflect.TypeOf((*MockSensorConnection)(nil).ClusterID))
}

// HasCapability mocks base method.
func (m *MockSensorConnection) HasCapability(capability centralsensor.SensorCapability) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasCapability", capability)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasCapability indicates an expected call of HasCapability.
func (mr *MockSensorConnectionMockRecorder) HasCapability(capability interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasCapability", reflect.TypeOf((*MockSensorConnection)(nil).HasCapability), capability)
}

// InjectMessage mocks base method.
func (m *MockSensorConnection) InjectMessage(ctx concurrency.Waitable, msg *central.MsgToSensor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectMessage", ctx, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectMessage indicates an expected call of InjectMessage.
func (mr *MockSensorConnectionMockRecorder) InjectMessage(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectMessage", reflect.TypeOf((*MockSensorConnection)(nil).InjectMessage), ctx, msg)
}

// InjectMessageIntoQueue mocks base method.
func (m *MockSensorConnection) InjectMessageIntoQueue(msg *central.MsgFromSensor) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InjectMessageIntoQueue", msg)
}

// InjectMessageIntoQueue indicates an expected call of InjectMessageIntoQueue.
func (mr *MockSensorConnectionMockRecorder) InjectMessageIntoQueue(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectMessageIntoQueue", reflect.TypeOf((*MockSensorConnection)(nil).InjectMessageIntoQueue), msg)
}

// NetworkEntities mocks base method.
func (m *MockSensorConnection) NetworkEntities() networkentities.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkEntities")
	ret0, _ := ret[0].(networkentities.Controller)
	return ret0
}

// NetworkEntities indicates an expected call of NetworkEntities.
func (mr *MockSensorConnectionMockRecorder) NetworkEntities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkEntities", reflect.TypeOf((*MockSensorConnection)(nil).NetworkEntities))
}

// NetworkPolicies mocks base method.
func (m *MockSensorConnection) NetworkPolicies() networkpolicies.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkPolicies")
	ret0, _ := ret[0].(networkpolicies.Controller)
	return ret0
}

// NetworkPolicies indicates an expected call of NetworkPolicies.
func (mr *MockSensorConnectionMockRecorder) NetworkPolicies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkPolicies", reflect.TypeOf((*MockSensorConnection)(nil).NetworkPolicies))
}

// ObjectsDeletedByReconciliation mocks base method.
func (m *MockSensorConnection) ObjectsDeletedByReconciliation() (map[string]int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectsDeletedByReconciliation")
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ObjectsDeletedByReconciliation indicates an expected call of ObjectsDeletedByReconciliation.
func (mr *MockSensorConnectionMockRecorder) ObjectsDeletedByReconciliation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectsDeletedByReconciliation", reflect.TypeOf((*MockSensorConnection)(nil).ObjectsDeletedByReconciliation))
}

// Scrapes mocks base method.
func (m *MockSensorConnection) Scrapes() scrape.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scrapes")
	ret0, _ := ret[0].(scrape.Controller)
	return ret0
}

// Scrapes indicates an expected call of Scrapes.
func (mr *MockSensorConnectionMockRecorder) Scrapes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scrapes", reflect.TypeOf((*MockSensorConnection)(nil).Scrapes))
}

// Stopped mocks base method.
func (m *MockSensorConnection) Stopped() concurrency.ReadOnlyErrorSignal {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stopped")
	ret0, _ := ret[0].(concurrency.ReadOnlyErrorSignal)
	return ret0
}

// Stopped indicates an expected call of Stopped.
func (mr *MockSensorConnectionMockRecorder) Stopped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stopped", reflect.TypeOf((*MockSensorConnection)(nil).Stopped))
}

// Telemetry mocks base method.
func (m *MockSensorConnection) Telemetry() telemetry.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Telemetry")
	ret0, _ := ret[0].(telemetry.Controller)
	return ret0
}

// Telemetry indicates an expected call of Telemetry.
func (mr *MockSensorConnectionMockRecorder) Telemetry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Telemetry", reflect.TypeOf((*MockSensorConnection)(nil).Telemetry))
}

// Terminate mocks base method.
func (m *MockSensorConnection) Terminate(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Terminate", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Terminate indicates an expected call of Terminate.
func (mr *MockSensorConnectionMockRecorder) Terminate(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Terminate", reflect.TypeOf((*MockSensorConnection)(nil).Terminate), err)
}
