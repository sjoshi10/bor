// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	bor "github.com/ethereum/go-ethereum/consensus/bor"
	mock "github.com/stretchr/testify/mock"
)

// IHeimdallClient is an autogenerated mock type for the IHeimdallClient type
type IHeimdallClient struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: path, query
func (_m *IHeimdallClient) Fetch(path string, query string) (*bor.ResponseWithHeight, error) {
	ret := _m.Called(path, query)

	var r0 *bor.ResponseWithHeight
	if rf, ok := ret.Get(0).(func(string, string) *bor.ResponseWithHeight); ok {
		r0 = rf(path, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bor.ResponseWithHeight)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(path, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchStateSyncEvents provides a mock function with given fields: fromID, to
func (_m *IHeimdallClient) FetchStateSyncEvents(fromID uint64, to int64) ([]*bor.EventRecordWithTime, error) {
	ret := _m.Called(fromID, to)

	var r0 []*bor.EventRecordWithTime
	if rf, ok := ret.Get(0).(func(uint64, int64) []*bor.EventRecordWithTime); ok {
		r0 = rf(fromID, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*bor.EventRecordWithTime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, int64) error); ok {
		r1 = rf(fromID, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchWithRetry provides a mock function with given fields: path, query
func (_m *IHeimdallClient) FetchWithRetry(path string, query string) (*bor.ResponseWithHeight, error) {
	ret := _m.Called(path, query)

	var r0 *bor.ResponseWithHeight
	if rf, ok := ret.Get(0).(func(string, string) *bor.ResponseWithHeight); ok {
		r0 = rf(path, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bor.ResponseWithHeight)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(path, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
