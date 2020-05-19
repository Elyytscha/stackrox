// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	analystnotes "github.com/stackrox/rox/central/analystnotes"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	reflect "reflect"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// Search mocks base method
func (m *MockDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// SearchDeployments mocks base method
func (m *MockDataStore) SearchDeployments(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchDeployments", ctx, q)
	ret0, _ := ret[0].([]*v1.SearchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchDeployments indicates an expected call of SearchDeployments
func (mr *MockDataStoreMockRecorder) SearchDeployments(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchDeployments", reflect.TypeOf((*MockDataStore)(nil).SearchDeployments), ctx, q)
}

// SearchRawDeployments mocks base method
func (m *MockDataStore) SearchRawDeployments(ctx context.Context, q *v1.Query) ([]*storage.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawDeployments", ctx, q)
	ret0, _ := ret[0].([]*storage.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawDeployments indicates an expected call of SearchRawDeployments
func (mr *MockDataStoreMockRecorder) SearchRawDeployments(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawDeployments", reflect.TypeOf((*MockDataStore)(nil).SearchRawDeployments), ctx, q)
}

// SearchListDeployments mocks base method
func (m *MockDataStore) SearchListDeployments(ctx context.Context, q *v1.Query) ([]*storage.ListDeployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchListDeployments", ctx, q)
	ret0, _ := ret[0].([]*storage.ListDeployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchListDeployments indicates an expected call of SearchListDeployments
func (mr *MockDataStoreMockRecorder) SearchListDeployments(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchListDeployments", reflect.TypeOf((*MockDataStore)(nil).SearchListDeployments), ctx, q)
}

// ListDeployment mocks base method
func (m *MockDataStore) ListDeployment(ctx context.Context, id string) (*storage.ListDeployment, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDeployment", ctx, id)
	ret0, _ := ret[0].(*storage.ListDeployment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListDeployment indicates an expected call of ListDeployment
func (mr *MockDataStoreMockRecorder) ListDeployment(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeployment", reflect.TypeOf((*MockDataStore)(nil).ListDeployment), ctx, id)
}

// GetDeployment mocks base method
func (m *MockDataStore) GetDeployment(ctx context.Context, id string) (*storage.Deployment, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeployment", ctx, id)
	ret0, _ := ret[0].(*storage.Deployment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDeployment indicates an expected call of GetDeployment
func (mr *MockDataStoreMockRecorder) GetDeployment(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeployment", reflect.TypeOf((*MockDataStore)(nil).GetDeployment), ctx, id)
}

// GetDeployments mocks base method
func (m *MockDataStore) GetDeployments(ctx context.Context, ids []string) ([]*storage.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeployments", ctx, ids)
	ret0, _ := ret[0].([]*storage.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeployments indicates an expected call of GetDeployments
func (mr *MockDataStoreMockRecorder) GetDeployments(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeployments", reflect.TypeOf((*MockDataStore)(nil).GetDeployments), ctx, ids)
}

// CountDeployments mocks base method
func (m *MockDataStore) CountDeployments(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountDeployments", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountDeployments indicates an expected call of CountDeployments
func (mr *MockDataStoreMockRecorder) CountDeployments(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountDeployments", reflect.TypeOf((*MockDataStore)(nil).CountDeployments), ctx)
}

// UpsertDeployment mocks base method
func (m *MockDataStore) UpsertDeployment(ctx context.Context, deployment *storage.Deployment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertDeployment", ctx, deployment)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertDeployment indicates an expected call of UpsertDeployment
func (mr *MockDataStoreMockRecorder) UpsertDeployment(ctx, deployment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertDeployment", reflect.TypeOf((*MockDataStore)(nil).UpsertDeployment), ctx, deployment)
}

// AddTagsToProcessKey mocks base method
func (m *MockDataStore) AddTagsToProcessKey(ctx context.Context, key *analystnotes.ProcessNoteKey, tags []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTagsToProcessKey", ctx, key, tags)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTagsToProcessKey indicates an expected call of AddTagsToProcessKey
func (mr *MockDataStoreMockRecorder) AddTagsToProcessKey(ctx, key, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTagsToProcessKey", reflect.TypeOf((*MockDataStore)(nil).AddTagsToProcessKey), ctx, key, tags)
}

// RemoveTagsFromProcessKey mocks base method
func (m *MockDataStore) RemoveTagsFromProcessKey(ctx context.Context, key *analystnotes.ProcessNoteKey, tags []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveTagsFromProcessKey", ctx, key, tags)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveTagsFromProcessKey indicates an expected call of RemoveTagsFromProcessKey
func (mr *MockDataStoreMockRecorder) RemoveTagsFromProcessKey(ctx, key, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveTagsFromProcessKey", reflect.TypeOf((*MockDataStore)(nil).RemoveTagsFromProcessKey), ctx, key, tags)
}

// GetTagsForProcessKey mocks base method
func (m *MockDataStore) GetTagsForProcessKey(ctx context.Context, key *analystnotes.ProcessNoteKey) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagsForProcessKey", ctx, key)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagsForProcessKey indicates an expected call of GetTagsForProcessKey
func (mr *MockDataStoreMockRecorder) GetTagsForProcessKey(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagsForProcessKey", reflect.TypeOf((*MockDataStore)(nil).GetTagsForProcessKey), ctx, key)
}

// RemoveDeployment mocks base method
func (m *MockDataStore) RemoveDeployment(ctx context.Context, clusterID, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDeployment", ctx, clusterID, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDeployment indicates an expected call of RemoveDeployment
func (mr *MockDataStoreMockRecorder) RemoveDeployment(ctx, clusterID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDeployment", reflect.TypeOf((*MockDataStore)(nil).RemoveDeployment), ctx, clusterID, id)
}

// GetImagesForDeployment mocks base method
func (m *MockDataStore) GetImagesForDeployment(ctx context.Context, deployment *storage.Deployment) ([]*storage.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImagesForDeployment", ctx, deployment)
	ret0, _ := ret[0].([]*storage.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImagesForDeployment indicates an expected call of GetImagesForDeployment
func (mr *MockDataStoreMockRecorder) GetImagesForDeployment(ctx, deployment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImagesForDeployment", reflect.TypeOf((*MockDataStore)(nil).GetImagesForDeployment), ctx, deployment)
}

// GetDeploymentIDs mocks base method
func (m *MockDataStore) GetDeploymentIDs() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeploymentIDs")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeploymentIDs indicates an expected call of GetDeploymentIDs
func (mr *MockDataStoreMockRecorder) GetDeploymentIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeploymentIDs", reflect.TypeOf((*MockDataStore)(nil).GetDeploymentIDs))
}
