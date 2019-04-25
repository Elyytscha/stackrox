// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/searchbasedpolicies/matcher (interfaces: Builder)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	searchbasedpolicies "github.com/stackrox/rox/central/searchbasedpolicies"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockBuilder is a mock of Builder interface
type MockBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderMockRecorder
}

// MockBuilderMockRecorder is the mock recorder for MockBuilder
type MockBuilderMockRecorder struct {
	mock *MockBuilder
}

// NewMockBuilder creates a new mock instance
func NewMockBuilder(ctrl *gomock.Controller) *MockBuilder {
	mock := &MockBuilder{ctrl: ctrl}
	mock.recorder = &MockBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBuilder) EXPECT() *MockBuilderMockRecorder {
	return m.recorder
}

// ForPolicy mocks base method
func (m *MockBuilder) ForPolicy(arg0 *storage.Policy) (searchbasedpolicies.Matcher, error) {
	ret := m.ctrl.Call(m, "ForPolicy", arg0)
	ret0, _ := ret[0].(searchbasedpolicies.Matcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForPolicy indicates an expected call of ForPolicy
func (mr *MockBuilderMockRecorder) ForPolicy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForPolicy", reflect.TypeOf((*MockBuilder)(nil).ForPolicy), arg0)
}
