// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/networkflow/store (interfaces: FlowStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	types "github.com/gogo/protobuf/types"
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	timestamp "github.com/stackrox/rox/pkg/timestamp"
	reflect "reflect"
)

// MockFlowStore is a mock of FlowStore interface
type MockFlowStore struct {
	ctrl     *gomock.Controller
	recorder *MockFlowStoreMockRecorder
}

// MockFlowStoreMockRecorder is the mock recorder for MockFlowStore
type MockFlowStoreMockRecorder struct {
	mock *MockFlowStore
}

// NewMockFlowStore creates a new mock instance
func NewMockFlowStore(ctrl *gomock.Controller) *MockFlowStore {
	mock := &MockFlowStore{ctrl: ctrl}
	mock.recorder = &MockFlowStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFlowStore) EXPECT() *MockFlowStoreMockRecorder {
	return m.recorder
}

// GetAllFlows mocks base method
func (m *MockFlowStore) GetAllFlows(arg0 *types.Timestamp) ([]*storage.NetworkFlow, types.Timestamp, error) {
	ret := m.ctrl.Call(m, "GetAllFlows", arg0)
	ret0, _ := ret[0].([]*storage.NetworkFlow)
	ret1, _ := ret[1].(types.Timestamp)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllFlows indicates an expected call of GetAllFlows
func (mr *MockFlowStoreMockRecorder) GetAllFlows(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFlows", reflect.TypeOf((*MockFlowStore)(nil).GetAllFlows), arg0)
}

// GetFlow mocks base method
func (m *MockFlowStore) GetFlow(arg0 *storage.NetworkFlowProperties) (*storage.NetworkFlow, error) {
	ret := m.ctrl.Call(m, "GetFlow", arg0)
	ret0, _ := ret[0].(*storage.NetworkFlow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlow indicates an expected call of GetFlow
func (mr *MockFlowStoreMockRecorder) GetFlow(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlow", reflect.TypeOf((*MockFlowStore)(nil).GetFlow), arg0)
}

// RemoveFlow mocks base method
func (m *MockFlowStore) RemoveFlow(arg0 *storage.NetworkFlowProperties) error {
	ret := m.ctrl.Call(m, "RemoveFlow", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFlow indicates an expected call of RemoveFlow
func (mr *MockFlowStoreMockRecorder) RemoveFlow(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFlow", reflect.TypeOf((*MockFlowStore)(nil).RemoveFlow), arg0)
}

// UpsertFlows mocks base method
func (m *MockFlowStore) UpsertFlows(arg0 []*storage.NetworkFlow, arg1 timestamp.MicroTS) error {
	ret := m.ctrl.Call(m, "UpsertFlows", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertFlows indicates an expected call of UpsertFlows
func (mr *MockFlowStoreMockRecorder) UpsertFlows(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertFlows", reflect.TypeOf((*MockFlowStore)(nil).UpsertFlows), arg0, arg1)
}
