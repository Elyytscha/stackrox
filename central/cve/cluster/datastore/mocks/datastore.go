// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "github.com/gogo/protobuf/types"
	gomock "github.com/golang/mock/gomock"
	converter "github.com/stackrox/rox/central/cve/converter/v2"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
)

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockDataStore) Count(ctx context.Context, q *aux.Query) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, q)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockDataStoreMockRecorder) Count(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDataStore)(nil).Count), ctx, q)
}

// DeleteClusterCVEsInternal mocks base method.
func (m *MockDataStore) DeleteClusterCVEsInternal(ctx context.Context, clusterID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteClusterCVEsInternal", ctx, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteClusterCVEsInternal indicates an expected call of DeleteClusterCVEsInternal.
func (mr *MockDataStoreMockRecorder) DeleteClusterCVEsInternal(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteClusterCVEsInternal", reflect.TypeOf((*MockDataStore)(nil).DeleteClusterCVEsInternal), ctx, clusterID)
}

// Exists mocks base method.
func (m *MockDataStore) Exists(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockDataStoreMockRecorder) Exists(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockDataStore)(nil).Exists), ctx, id)
}

// Get mocks base method.
func (m *MockDataStore) Get(ctx context.Context, id string) (*storage.ClusterCVE, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*storage.ClusterCVE)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockDataStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDataStore)(nil).Get), ctx, id)
}

// GetBatch mocks base method.
func (m *MockDataStore) GetBatch(ctx context.Context, id []string) ([]*storage.ClusterCVE, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatch", ctx, id)
	ret0, _ := ret[0].([]*storage.ClusterCVE)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatch indicates an expected call of GetBatch.
func (mr *MockDataStoreMockRecorder) GetBatch(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatch", reflect.TypeOf((*MockDataStore)(nil).GetBatch), ctx, id)
}

// Search mocks base method.
func (m *MockDataStore) Search(ctx context.Context, q *aux.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// SearchCVEs mocks base method.
func (m *MockDataStore) SearchCVEs(ctx context.Context, q *aux.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCVEs", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchCVEs indicates an expected call of SearchCVEs.
func (mr *MockDataStoreMockRecorder) SearchCVEs(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCVEs", reflect.TypeOf((*MockDataStore)(nil).SearchCVEs), ctx, q)
}

// SearchRawCVEs mocks base method.
func (m *MockDataStore) SearchRawCVEs(ctx context.Context, q *aux.Query) ([]*storage.ClusterCVE, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawCVEs", ctx, q)
	ret0, _ := ret[0].([]*storage.ClusterCVE)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawCVEs indicates an expected call of SearchRawCVEs.
func (mr *MockDataStoreMockRecorder) SearchRawCVEs(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawCVEs", reflect.TypeOf((*MockDataStore)(nil).SearchRawCVEs), ctx, q)
}

// Suppress mocks base method.
func (m *MockDataStore) Suppress(ctx context.Context, start *types.Timestamp, duration *types.Duration, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, start, duration}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Suppress", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Suppress indicates an expected call of Suppress.
func (mr *MockDataStoreMockRecorder) Suppress(ctx, start, duration interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, start, duration}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Suppress", reflect.TypeOf((*MockDataStore)(nil).Suppress), varargs...)
}

// Unsuppress mocks base method.
func (m *MockDataStore) Unsuppress(ctx context.Context, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Unsuppress", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsuppress indicates an expected call of Unsuppress.
func (mr *MockDataStoreMockRecorder) Unsuppress(ctx interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsuppress", reflect.TypeOf((*MockDataStore)(nil).Unsuppress), varargs...)
}

// UpsertClusterCVEsInternal mocks base method.
func (m *MockDataStore) UpsertClusterCVEsInternal(ctx context.Context, cveType storage.CVE_CVEType, cveParts ...converter.ClusterCVEParts) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, cveType}
	for _, a := range cveParts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertClusterCVEsInternal", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertClusterCVEsInternal indicates an expected call of UpsertClusterCVEsInternal.
func (mr *MockDataStoreMockRecorder) UpsertClusterCVEsInternal(ctx, cveType interface{}, cveParts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, cveType}, cveParts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertClusterCVEsInternal", reflect.TypeOf((*MockDataStore)(nil).UpsertClusterCVEsInternal), varargs...)
}
