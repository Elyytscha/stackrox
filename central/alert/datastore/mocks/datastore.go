// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

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

// AddAlertTags mocks base method.
func (m *MockDataStore) AddAlertTags(ctx context.Context, alertID string, tags []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAlertTags", ctx, alertID, tags)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAlertTags indicates an expected call of AddAlertTags.
func (mr *MockDataStoreMockRecorder) AddAlertTags(ctx, alertID, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAlertTags", reflect.TypeOf((*MockDataStore)(nil).AddAlertTags), ctx, alertID, tags)
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

// CountAlerts mocks base method.
func (m *MockDataStore) CountAlerts(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAlerts", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAlerts indicates an expected call of CountAlerts.
func (mr *MockDataStoreMockRecorder) CountAlerts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAlerts", reflect.TypeOf((*MockDataStore)(nil).CountAlerts), ctx)
}

// DeleteAlerts mocks base method.
func (m *MockDataStore) DeleteAlerts(ctx context.Context, ids ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAlerts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAlerts indicates an expected call of DeleteAlerts.
func (mr *MockDataStoreMockRecorder) DeleteAlerts(ctx interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAlerts", reflect.TypeOf((*MockDataStore)(nil).DeleteAlerts), varargs...)
}

// GetAlert mocks base method.
func (m *MockDataStore) GetAlert(ctx context.Context, id string) (*storage.Alert, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlert", ctx, id)
	ret0, _ := ret[0].(*storage.Alert)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAlert indicates an expected call of GetAlert.
func (mr *MockDataStoreMockRecorder) GetAlert(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlert", reflect.TypeOf((*MockDataStore)(nil).GetAlert), ctx, id)
}

// ListAlerts mocks base method.
func (m *MockDataStore) ListAlerts(ctx context.Context, request *v1.ListAlertsRequest) ([]*storage.ListAlert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAlerts", ctx, request)
	ret0, _ := ret[0].([]*storage.ListAlert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAlerts indicates an expected call of ListAlerts.
func (mr *MockDataStoreMockRecorder) ListAlerts(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAlerts", reflect.TypeOf((*MockDataStore)(nil).ListAlerts), ctx, request)
}

// MarkAlertStale mocks base method.
func (m *MockDataStore) MarkAlertStale(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAlertStale", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAlertStale indicates an expected call of MarkAlertStale.
func (mr *MockDataStoreMockRecorder) MarkAlertStale(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAlertStale", reflect.TypeOf((*MockDataStore)(nil).MarkAlertStale), ctx, id)
}

// RemoveAlertTags mocks base method.
func (m *MockDataStore) RemoveAlertTags(ctx context.Context, alertID string, tags []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAlertTags", ctx, alertID, tags)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAlertTags indicates an expected call of RemoveAlertTags.
func (mr *MockDataStoreMockRecorder) RemoveAlertTags(ctx, alertID, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAlertTags", reflect.TypeOf((*MockDataStore)(nil).RemoveAlertTags), ctx, alertID, tags)
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

// SearchAlerts mocks base method.
func (m *MockDataStore) SearchAlerts(ctx context.Context, q *aux.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAlerts", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchAlerts indicates an expected call of SearchAlerts.
func (mr *MockDataStoreMockRecorder) SearchAlerts(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAlerts", reflect.TypeOf((*MockDataStore)(nil).SearchAlerts), ctx, q)
}

// SearchListAlerts mocks base method.
func (m *MockDataStore) SearchListAlerts(ctx context.Context, q *aux.Query) ([]*storage.ListAlert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchListAlerts", ctx, q)
	ret0, _ := ret[0].([]*storage.ListAlert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchListAlerts indicates an expected call of SearchListAlerts.
func (mr *MockDataStoreMockRecorder) SearchListAlerts(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchListAlerts", reflect.TypeOf((*MockDataStore)(nil).SearchListAlerts), ctx, q)
}

// SearchRawAlerts mocks base method.
func (m *MockDataStore) SearchRawAlerts(ctx context.Context, q *aux.Query) ([]*storage.Alert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawAlerts", ctx, q)
	ret0, _ := ret[0].([]*storage.Alert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawAlerts indicates an expected call of SearchRawAlerts.
func (mr *MockDataStoreMockRecorder) SearchRawAlerts(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawAlerts", reflect.TypeOf((*MockDataStore)(nil).SearchRawAlerts), ctx, q)
}

// UpsertAlert mocks base method.
func (m *MockDataStore) UpsertAlert(ctx context.Context, alert *storage.Alert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertAlert", ctx, alert)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertAlert indicates an expected call of UpsertAlert.
func (mr *MockDataStoreMockRecorder) UpsertAlert(ctx, alert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertAlert", reflect.TypeOf((*MockDataStore)(nil).UpsertAlert), ctx, alert)
}

// UpsertAlerts mocks base method.
func (m *MockDataStore) UpsertAlerts(ctx context.Context, alerts []*storage.Alert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertAlerts", ctx, alerts)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertAlerts indicates an expected call of UpsertAlerts.
func (mr *MockDataStoreMockRecorder) UpsertAlerts(ctx, alerts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertAlerts", reflect.TypeOf((*MockDataStore)(nil).UpsertAlerts), ctx, alerts)
}

// WalkAll mocks base method.
func (m *MockDataStore) WalkAll(ctx context.Context, fn func(*storage.ListAlert) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalkAll", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WalkAll indicates an expected call of WalkAll.
func (mr *MockDataStoreMockRecorder) WalkAll(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalkAll", reflect.TypeOf((*MockDataStore)(nil).WalkAll), ctx, fn)
}
