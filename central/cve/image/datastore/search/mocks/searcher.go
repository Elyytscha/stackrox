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

// SearchImageCVEs mocks base method.
func (m *MockSearcher) SearchImageCVEs(arg0 context.Context, arg1 *aux.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchImageCVEs", arg0, arg1)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchImageCVEs indicates an expected call of SearchImageCVEs.
func (mr *MockSearcherMockRecorder) SearchImageCVEs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchImageCVEs", reflect.TypeOf((*MockSearcher)(nil).SearchImageCVEs), arg0, arg1)
}

// SearchRawImageCVEs mocks base method.
func (m *MockSearcher) SearchRawImageCVEs(ctx context.Context, query *aux.Query) ([]*storage.ImageCVE, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawImageCVEs", ctx, query)
	ret0, _ := ret[0].([]*storage.ImageCVE)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawImageCVEs indicates an expected call of SearchRawImageCVEs.
func (mr *MockSearcherMockRecorder) SearchRawImageCVEs(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawImageCVEs", reflect.TypeOf((*MockSearcher)(nil).SearchRawImageCVEs), ctx, query)
}
