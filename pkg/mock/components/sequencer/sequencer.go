// Code generated by MockGen. DO NOT EDIT.
// Source: components/sequencer/store.go

// Package mock_sequencer is a generated GoMock package.
package mock_sequencer

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sequencer "mosn.io/layotto/components/sequencer"
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

// GetNextId mocks base method.
func (m *MockStore) GetNextId(arg0 *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextId", arg0)
	ret0, _ := ret[0].(*sequencer.GetNextIdResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextId indicates an expected call of GetNextId.
func (mr *MockStoreMockRecorder) GetNextId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextId", reflect.TypeOf((*MockStore)(nil).GetNextId), arg0)
}

// GetSegment mocks base method.
func (m *MockStore) GetSegment(arg0 *sequencer.GetSegmentRequest) (bool, *sequencer.GetSegmentResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegment", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*sequencer.GetSegmentResponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSegment indicates an expected call of GetSegment.
func (mr *MockStoreMockRecorder) GetSegment(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegment", reflect.TypeOf((*MockStore)(nil).GetSegment), arg0)
}

// Init mocks base method.
func (m *MockStore) Init(config sequencer.Configuration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockStoreMockRecorder) Init(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockStore)(nil).Init), config)
}
