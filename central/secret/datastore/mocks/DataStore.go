// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import v1 "github.com/stackrox/rox/generated/api/v1"

// DataStore is an autogenerated mock type for the DataStore type
type DataStore struct {
	mock.Mock
}

// GetSecret provides a mock function with given fields: id
func (_m *DataStore) GetSecret(id string) (*v1.Secret, bool, error) {
	ret := _m.Called(id)

	var r0 *v1.Secret
	if rf, ok := ret.Get(0).(func(string) *v1.Secret); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Secret)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RemoveSecret provides a mock function with given fields: id
func (_m *DataStore) RemoveSecret(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SearchListSecrets provides a mock function with given fields: q
func (_m *DataStore) SearchListSecrets(q *v1.Query) ([]*v1.ListSecret, error) {
	ret := _m.Called(q)

	var r0 []*v1.ListSecret
	if rf, ok := ret.Get(0).(func(*v1.Query) []*v1.ListSecret); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.ListSecret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.Query) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchSecrets provides a mock function with given fields: q
func (_m *DataStore) SearchSecrets(q *v1.Query) ([]*v1.SearchResult, error) {
	ret := _m.Called(q)

	var r0 []*v1.SearchResult
	if rf, ok := ret.Get(0).(func(*v1.Query) []*v1.SearchResult); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.SearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.Query) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertSecret provides a mock function with given fields: request
func (_m *DataStore) UpsertSecret(request *v1.Secret) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.Secret) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
