// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dapr/components-contrib/pubsub (interfaces: PubSub)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	metadata "github.com/dapr/components-contrib/metadata"
	pubsub "github.com/dapr/components-contrib/pubsub"
	gomock "github.com/golang/mock/gomock"
)

// MockPubSub is a mock of PubSub interface.
type MockPubSub struct {
	ctrl     *gomock.Controller
	recorder *MockPubSubMockRecorder
}

// MockPubSubMockRecorder is the mock recorder for MockPubSub.
type MockPubSubMockRecorder struct {
	mock *MockPubSub
}

// NewMockPubSub creates a new mock instance.
func NewMockPubSub(ctrl *gomock.Controller) *MockPubSub {
	mock := &MockPubSub{ctrl: ctrl}
	mock.recorder = &MockPubSubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubSub) EXPECT() *MockPubSubMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPubSub) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockPubSubMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPubSub)(nil).Close))
}

// Features mocks base method.
func (m *MockPubSub) Features() []pubsub.Feature {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Features")
	ret0, _ := ret[0].([]pubsub.Feature)
	return ret0
}

// Features indicates an expected call of Features.
func (mr *MockPubSubMockRecorder) Features() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Features", reflect.TypeOf((*MockPubSub)(nil).Features))
}

// GetComponentMetadata mocks base method.
func (m *MockPubSub) GetComponentMetadata() metadata.MetadataMap {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComponentMetadata")
	ret0, _ := ret[0].(metadata.MetadataMap)
	return ret0
}

// GetComponentMetadata indicates an expected call of GetComponentMetadata.
func (mr *MockPubSubMockRecorder) GetComponentMetadata() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComponentMetadata", reflect.TypeOf((*MockPubSub)(nil).GetComponentMetadata))
}

// Init mocks base method.
func (m *MockPubSub) Init(ctx context.Context, metadata pubsub.Metadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx, metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockPubSubMockRecorder) Init(ctx, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockPubSub)(nil).Init), ctx, metadata)
}

// Publish mocks base method.
func (m *MockPubSub) Publish(ctx context.Context, req *pubsub.PublishRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPubSubMockRecorder) Publish(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPubSub)(nil).Publish), ctx, req)
}

// Subscribe mocks base method.
func (m *MockPubSub) Subscribe(ctx context.Context, req pubsub.SubscribeRequest, handler pubsub.Handler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, req, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPubSubMockRecorder) Subscribe(ctx, req, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPubSub)(nil).Subscribe), ctx, req, handler)
}
