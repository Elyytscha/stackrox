// Code generated by MockGen. DO NOT EDIT.
// Source: search.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
)

// MockSearcher is a mock of Searcher interface.
type MockSearcher struct {
	ctrl     *gomock.Controller
	recorder *MockSearcherMockRecorder
}

// MockSearcherMockRecorder is the mock recorder for MockSearcher.
type MockSearcherMockRecorder struct {
	mock *MockSearcher
}

// NewMockSearcher creates a new mock instance.
func NewMockSearcher(ctrl *gomock.Controller) *MockSearcher {
	mock := &MockSearcher{ctrl: ctrl}
	mock.recorder = &MockSearcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearcher) EXPECT() *MockSearcherMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockSearcher) Count(ctx context.Context, query *aux.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockSearcherMockRecorder) Count(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockSearcher)(nil).Count), ctx, query)
}

// Search mocks base method.
func (m *MockSearcher) Search(ctx context.Context, query *aux.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockSearcherMockRecorder) Search(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockSearcher)(nil).Search), ctx, query)
}

// SearchReportConfigurations mocks base method.
func (m *MockSearcher) SearchReportConfigurations(ctx context.Context, query *aux.Query) ([]*storage.ReportConfiguration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchReportConfigurations", ctx, query)
	ret0, _ := ret[0].([]*storage.ReportConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchReportConfigurations indicates an expected call of SearchReportConfigurations.
func (mr *MockSearcherMockRecorder) SearchReportConfigurations(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchReportConfigurations", reflect.TypeOf((*MockSearcher)(nil).SearchReportConfigurations), ctx, query)
}
