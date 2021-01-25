// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	common "github.com/stackrox/rox/central/sensor/service/common"
	connection "github.com/stackrox/rox/central/sensor/service/connection"
	pipeline "github.com/stackrox/rox/central/sensor/service/pipeline"
	central "github.com/stackrox/rox/generated/internalapi/central"
	storage "github.com/stackrox/rox/generated/storage"
	concurrency "github.com/stackrox/rox/pkg/concurrency"
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

// Start mocks base method
func (m *MockManager) Start(mgr common.ClusterManager, netEntitiesMgr common.NetworkEntityManager, policyMgr common.PolicyManager, baselineMgr common.ProcessBaselineManager, autoTriggerUpgrades *concurrency.Flag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", mgr, netEntitiesMgr, policyMgr, baselineMgr, autoTriggerUpgrades)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockManagerMockRecorder) Start(mgr, netEntitiesMgr, policyMgr, baselineMgr, autoTriggerUpgrades interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockManager)(nil).Start), mgr, netEntitiesMgr, policyMgr, baselineMgr, autoTriggerUpgrades)
}

// HandleConnection mocks base method
func (m *MockManager) HandleConnection(ctx context.Context, sensorHello *central.SensorHello, cluster *storage.Cluster, eventPipeline pipeline.ClusterPipeline, server central.SensorService_CommunicateServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleConnection", ctx, sensorHello, cluster, eventPipeline, server)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleConnection indicates an expected call of HandleConnection
func (mr *MockManagerMockRecorder) HandleConnection(ctx, sensorHello, cluster, eventPipeline, server interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleConnection", reflect.TypeOf((*MockManager)(nil).HandleConnection), ctx, sensorHello, cluster, eventPipeline, server)
}

// GetConnection mocks base method
func (m *MockManager) GetConnection(clusterID string) connection.SensorConnection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnection", clusterID)
	ret0, _ := ret[0].(connection.SensorConnection)
	return ret0
}

// GetConnection indicates an expected call of GetConnection
func (mr *MockManagerMockRecorder) GetConnection(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockManager)(nil).GetConnection), clusterID)
}

// GetActiveConnections mocks base method
func (m *MockManager) GetActiveConnections() []connection.SensorConnection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveConnections")
	ret0, _ := ret[0].([]connection.SensorConnection)
	return ret0
}

// GetActiveConnections indicates an expected call of GetActiveConnections
func (mr *MockManagerMockRecorder) GetActiveConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveConnections", reflect.TypeOf((*MockManager)(nil).GetActiveConnections))
}

// BroadcastMessage mocks base method
func (m *MockManager) BroadcastMessage(msg *central.MsgToSensor) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BroadcastMessage", msg)
}

// BroadcastMessage indicates an expected call of BroadcastMessage
func (mr *MockManagerMockRecorder) BroadcastMessage(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastMessage", reflect.TypeOf((*MockManager)(nil).BroadcastMessage), msg)
}

// SendMessage mocks base method
func (m *MockManager) SendMessage(clusterID string, msg *central.MsgToSensor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", clusterID, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage
func (mr *MockManagerMockRecorder) SendMessage(clusterID, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockManager)(nil).SendMessage), clusterID, msg)
}

// TriggerUpgrade mocks base method
func (m *MockManager) TriggerUpgrade(ctx context.Context, clusterID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TriggerUpgrade", ctx, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// TriggerUpgrade indicates an expected call of TriggerUpgrade
func (mr *MockManagerMockRecorder) TriggerUpgrade(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TriggerUpgrade", reflect.TypeOf((*MockManager)(nil).TriggerUpgrade), ctx, clusterID)
}

// TriggerCertRotation mocks base method
func (m *MockManager) TriggerCertRotation(ctx context.Context, clusterID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TriggerCertRotation", ctx, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// TriggerCertRotation indicates an expected call of TriggerCertRotation
func (mr *MockManagerMockRecorder) TriggerCertRotation(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TriggerCertRotation", reflect.TypeOf((*MockManager)(nil).TriggerCertRotation), ctx, clusterID)
}

// ProcessCheckInFromUpgrader mocks base method
func (m *MockManager) ProcessCheckInFromUpgrader(ctx context.Context, clusterID string, req *central.UpgradeCheckInFromUpgraderRequest) (*central.UpgradeCheckInFromUpgraderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessCheckInFromUpgrader", ctx, clusterID, req)
	ret0, _ := ret[0].(*central.UpgradeCheckInFromUpgraderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessCheckInFromUpgrader indicates an expected call of ProcessCheckInFromUpgrader
func (mr *MockManagerMockRecorder) ProcessCheckInFromUpgrader(ctx, clusterID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessCheckInFromUpgrader", reflect.TypeOf((*MockManager)(nil).ProcessCheckInFromUpgrader), ctx, clusterID, req)
}

// ProcessUpgradeCheckInFromSensor mocks base method
func (m *MockManager) ProcessUpgradeCheckInFromSensor(ctx context.Context, clusterID string, req *central.UpgradeCheckInFromSensorRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessUpgradeCheckInFromSensor", ctx, clusterID, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessUpgradeCheckInFromSensor indicates an expected call of ProcessUpgradeCheckInFromSensor
func (mr *MockManagerMockRecorder) ProcessUpgradeCheckInFromSensor(ctx, clusterID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessUpgradeCheckInFromSensor", reflect.TypeOf((*MockManager)(nil).ProcessUpgradeCheckInFromSensor), ctx, clusterID, req)
}

// PushExternalNetworkEntitiesToSensor mocks base method
func (m *MockManager) PushExternalNetworkEntitiesToSensor(ctx context.Context, clusterID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushExternalNetworkEntitiesToSensor", ctx, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushExternalNetworkEntitiesToSensor indicates an expected call of PushExternalNetworkEntitiesToSensor
func (mr *MockManagerMockRecorder) PushExternalNetworkEntitiesToSensor(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushExternalNetworkEntitiesToSensor", reflect.TypeOf((*MockManager)(nil).PushExternalNetworkEntitiesToSensor), ctx, clusterID)
}

// PushExternalNetworkEntitiesToAllSensors mocks base method
func (m *MockManager) PushExternalNetworkEntitiesToAllSensors(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushExternalNetworkEntitiesToAllSensors", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushExternalNetworkEntitiesToAllSensors indicates an expected call of PushExternalNetworkEntitiesToAllSensors
func (mr *MockManagerMockRecorder) PushExternalNetworkEntitiesToAllSensors(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushExternalNetworkEntitiesToAllSensors", reflect.TypeOf((*MockManager)(nil).PushExternalNetworkEntitiesToAllSensors), ctx)
}
