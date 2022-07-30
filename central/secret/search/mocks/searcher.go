// Code generated by MockGen. DO NOT EDIT.
// Source: searcher.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
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

// SearchListSecrets mocks base method.
func (m *MockSearcher) SearchListSecrets(ctx context.Context, query *aux.Query) ([]*storage.ListSecret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchListSecrets", ctx, query)
	ret0, _ := ret[0].([]*storage.ListSecret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchListSecrets indicates an expected call of SearchListSecrets.
func (mr *MockSearcherMockRecorder) SearchListSecrets(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchListSecrets", reflect.TypeOf((*MockSearcher)(nil).SearchListSecrets), ctx, query)
}

// SearchRawSecrets mocks base method.
func (m *MockSearcher) SearchRawSecrets(ctx context.Context, query *aux.Query) ([]*storage.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawSecrets", ctx, query)
	ret0, _ := ret[0].([]*storage.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawSecrets indicates an expected call of SearchRawSecrets.
func (mr *MockSearcherMockRecorder) SearchRawSecrets(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawSecrets", reflect.TypeOf((*MockSearcher)(nil).SearchRawSecrets), ctx, query)
}

// SearchSecrets mocks base method.
func (m *MockSearcher) SearchSecrets(arg0 context.Context, arg1 *aux.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchSecrets", arg0, arg1)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchSecrets indicates an expected call of SearchSecrets.
func (mr *MockSearcherMockRecorder) SearchSecrets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchSecrets", reflect.TypeOf((*MockSearcher)(nil).SearchSecrets), arg0, arg1)
}
