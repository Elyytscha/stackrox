// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	compliance "github.com/stackrox/rox/central/compliance"
	types "github.com/stackrox/rox/central/compliance/datastore/types"
	storage "github.com/stackrox/rox/generated/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// ClearAggregationResults mocks base method.
func (m *MockStore) ClearAggregationResults(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearAggregationResults", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearAggregationResults indicates an expected call of ClearAggregationResults.
func (mr *MockStoreMockRecorder) ClearAggregationResults(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearAggregationResults", reflect.TypeOf((*MockStore)(nil).ClearAggregationResults), ctx)
}

// GetAggregationResult mocks base method.
func (m *MockStore) GetAggregationResult(ctx context.Context, queryString string, groupBy []storage.ComplianceAggregation_Scope, unit storage.ComplianceAggregation_Scope) ([]*storage.ComplianceAggregation_Result, []*storage.ComplianceAggregation_Source, map[*storage.ComplianceAggregation_Result]*storage.ComplianceDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregationResult", ctx, queryString, groupBy, unit)
	ret0, _ := ret[0].([]*storage.ComplianceAggregation_Result)
	ret1, _ := ret[1].([]*storage.ComplianceAggregation_Source)
	ret2, _ := ret[2].(map[*storage.ComplianceAggregation_Result]*storage.ComplianceDomain)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetAggregationResult indicates an expected call of GetAggregationResult.
func (mr *MockStoreMockRecorder) GetAggregationResult(ctx, queryString, groupBy, unit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregationResult", reflect.TypeOf((*MockStore)(nil).GetAggregationResult), ctx, queryString, groupBy, unit)
}

// GetConfig mocks base method.
func (m *MockStore) GetConfig(ctx context.Context, id string) (*storage.ComplianceConfig, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig", ctx, id)
	ret0, _ := ret[0].(*storage.ComplianceConfig)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockStoreMockRecorder) GetConfig(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockStore)(nil).GetConfig), ctx, id)
}

// GetLatestRunMetadataBatch mocks base method.
func (m *MockStore) GetLatestRunMetadataBatch(ctx context.Context, clusterID string, standardIDs []string) (map[compliance.ClusterStandardPair]types.ComplianceRunsMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRunMetadataBatch", ctx, clusterID, standardIDs)
	ret0, _ := ret[0].(map[compliance.ClusterStandardPair]types.ComplianceRunsMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestRunMetadataBatch indicates an expected call of GetLatestRunMetadataBatch.
func (mr *MockStoreMockRecorder) GetLatestRunMetadataBatch(ctx, clusterID, standardIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRunMetadataBatch", reflect.TypeOf((*MockStore)(nil).GetLatestRunMetadataBatch), ctx, clusterID, standardIDs)
}

// GetLatestRunResults mocks base method.
func (m *MockStore) GetLatestRunResults(ctx context.Context, clusterID, standardID string, flags types.GetFlags) (types.ResultsWithStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRunResults", ctx, clusterID, standardID, flags)
	ret0, _ := ret[0].(types.ResultsWithStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestRunResults indicates an expected call of GetLatestRunResults.
func (mr *MockStoreMockRecorder) GetLatestRunResults(ctx, clusterID, standardID, flags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRunResults", reflect.TypeOf((*MockStore)(nil).GetLatestRunResults), ctx, clusterID, standardID, flags)
}

// GetLatestRunResultsBatch mocks base method.
func (m *MockStore) GetLatestRunResultsBatch(ctx context.Context, clusterIDs, standardIDs []string, flags types.GetFlags) (map[compliance.ClusterStandardPair]types.ResultsWithStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRunResultsBatch", ctx, clusterIDs, standardIDs, flags)
	ret0, _ := ret[0].(map[compliance.ClusterStandardPair]types.ResultsWithStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestRunResultsBatch indicates an expected call of GetLatestRunResultsBatch.
func (mr *MockStoreMockRecorder) GetLatestRunResultsBatch(ctx, clusterIDs, standardIDs, flags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRunResultsBatch", reflect.TypeOf((*MockStore)(nil).GetLatestRunResultsBatch), ctx, clusterIDs, standardIDs, flags)
}

// GetSpecificRunResults mocks base method.
func (m *MockStore) GetSpecificRunResults(ctx context.Context, clusterID, standardID, runID string, flags types.GetFlags) (types.ResultsWithStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpecificRunResults", ctx, clusterID, standardID, runID, flags)
	ret0, _ := ret[0].(types.ResultsWithStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpecificRunResults indicates an expected call of GetSpecificRunResults.
func (mr *MockStoreMockRecorder) GetSpecificRunResults(ctx, clusterID, standardID, runID, flags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpecificRunResults", reflect.TypeOf((*MockStore)(nil).GetSpecificRunResults), ctx, clusterID, standardID, runID, flags)
}

// StoreAggregationResult mocks base method.
func (m *MockStore) StoreAggregationResult(ctx context.Context, queryString string, groupBy []storage.ComplianceAggregation_Scope, unit storage.ComplianceAggregation_Scope, results []*storage.ComplianceAggregation_Result, sources []*storage.ComplianceAggregation_Source, domains map[*storage.ComplianceAggregation_Result]*storage.ComplianceDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreAggregationResult", ctx, queryString, groupBy, unit, results, sources, domains)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreAggregationResult indicates an expected call of StoreAggregationResult.
func (mr *MockStoreMockRecorder) StoreAggregationResult(ctx, queryString, groupBy, unit, results, sources, domains interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreAggregationResult", reflect.TypeOf((*MockStore)(nil).StoreAggregationResult), ctx, queryString, groupBy, unit, results, sources, domains)
}

// StoreComplianceDomain mocks base method.
func (m *MockStore) StoreComplianceDomain(ctx context.Context, domain *storage.ComplianceDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreComplianceDomain", ctx, domain)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreComplianceDomain indicates an expected call of StoreComplianceDomain.
func (mr *MockStoreMockRecorder) StoreComplianceDomain(ctx, domain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreComplianceDomain", reflect.TypeOf((*MockStore)(nil).StoreComplianceDomain), ctx, domain)
}

// StoreFailure mocks base method.
func (m *MockStore) StoreFailure(ctx context.Context, metadata *storage.ComplianceRunMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreFailure", ctx, metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreFailure indicates an expected call of StoreFailure.
func (mr *MockStoreMockRecorder) StoreFailure(ctx, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreFailure", reflect.TypeOf((*MockStore)(nil).StoreFailure), ctx, metadata)
}

// StoreRunResults mocks base method.
func (m *MockStore) StoreRunResults(ctx context.Context, results *storage.ComplianceRunResults) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreRunResults", ctx, results)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreRunResults indicates an expected call of StoreRunResults.
func (mr *MockStoreMockRecorder) StoreRunResults(ctx, results interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreRunResults", reflect.TypeOf((*MockStore)(nil).StoreRunResults), ctx, results)
}

// UpdateConfig mocks base method.
func (m *MockStore) UpdateConfig(ctx context.Context, config *storage.ComplianceConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfig", ctx, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfig indicates an expected call of UpdateConfig.
func (mr *MockStoreMockRecorder) UpdateConfig(ctx, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfig", reflect.TypeOf((*MockStore)(nil).UpdateConfig), ctx, config)
}
