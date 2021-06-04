// Code generated by protoc-gen-go. DO NOT EDIT.
// source: appcallback.proto

package runtime

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// TopicEventResponseStatus allows apps to have finer control over handling of the message.
type TopicEventResponse_TopicEventResponseStatus int32

const (
	// SUCCESS is the default behavior: message is acknowledged and not retried or logged.
	TopicEventResponse_SUCCESS TopicEventResponse_TopicEventResponseStatus = 0
	// RETRY status signals runtime to retry the message as part of an expected scenario (no warning is logged).
	TopicEventResponse_RETRY TopicEventResponse_TopicEventResponseStatus = 1
	// DROP status signals runtime to drop the message as part of an unexpected scenario (warning is logged).
	TopicEventResponse_DROP TopicEventResponse_TopicEventResponseStatus = 2
)

var TopicEventResponse_TopicEventResponseStatus_name = map[int32]string{
	0: "SUCCESS",
	1: "RETRY",
	2: "DROP",
}

var TopicEventResponse_TopicEventResponseStatus_value = map[string]int32{
	"SUCCESS": 0,
	"RETRY":   1,
	"DROP":    2,
}

func (x TopicEventResponse_TopicEventResponseStatus) String() string {
	return proto.EnumName(TopicEventResponse_TopicEventResponseStatus_name, int32(x))
}

func (TopicEventResponse_TopicEventResponseStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f48007b62b0c167c, []int{1, 0}
}

// TopicEventRequest message is compatible with CloudEvent spec v1.0
// https://github.com/cloudevents/spec/blob/v1.0/spec.md
type TopicEventRequest struct {
	// id identifies the event. Producers MUST ensure that source + id
	// is unique for each distinct event. If a duplicate event is re-sent
	// (e.g. due to a network error) it MAY have the same id.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// source identifies the context in which an event happened.
	// Often this will include information such as the type of the
	// event source, the organization publishing the event or the process
	// that produced the event. The exact syntax and semantics behind
	// the data encoded in the URI is defined by the event producer.
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	// The type of event related to the originating occurrence.
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	// The version of the CloudEvents specification.
	SpecVersion string `protobuf:"bytes,4,opt,name=spec_version,json=specVersion,proto3" json:"spec_version,omitempty"`
	// The content type of data value.
	DataContentType string `protobuf:"bytes,5,opt,name=data_content_type,json=dataContentType,proto3" json:"data_content_type,omitempty"`
	// The content of the event.
	Data []byte `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	// The pubsub topic which publisher sent to.
	Topic string `protobuf:"bytes,6,opt,name=topic,proto3" json:"topic,omitempty"`
	// The name of the pubsub the publisher sent to.
	PubsubName           string   `protobuf:"bytes,8,opt,name=pubsub_name,json=pubsubName,proto3" json:"pubsub_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TopicEventRequest) Reset()         { *m = TopicEventRequest{} }
func (m *TopicEventRequest) String() string { return proto.CompactTextString(m) }
func (*TopicEventRequest) ProtoMessage()    {}
func (*TopicEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f48007b62b0c167c, []int{0}
}

func (m *TopicEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TopicEventRequest.Unmarshal(m, b)
}
func (m *TopicEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TopicEventRequest.Marshal(b, m, deterministic)
}
func (m *TopicEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TopicEventRequest.Merge(m, src)
}
func (m *TopicEventRequest) XXX_Size() int {
	return xxx_messageInfo_TopicEventRequest.Size(m)
}
func (m *TopicEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TopicEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TopicEventRequest proto.InternalMessageInfo

func (m *TopicEventRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TopicEventRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *TopicEventRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *TopicEventRequest) GetSpecVersion() string {
	if m != nil {
		return m.SpecVersion
	}
	return ""
}

func (m *TopicEventRequest) GetDataContentType() string {
	if m != nil {
		return m.DataContentType
	}
	return ""
}

func (m *TopicEventRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *TopicEventRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *TopicEventRequest) GetPubsubName() string {
	if m != nil {
		return m.PubsubName
	}
	return ""
}

// TopicEventResponse is response from app on published message
type TopicEventResponse struct {
	// The list of output bindings.
	Status               TopicEventResponse_TopicEventResponseStatus `protobuf:"varint,1,opt,name=status,proto3,enum=spec.proto.runtime.v1.TopicEventResponse_TopicEventResponseStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *TopicEventResponse) Reset()         { *m = TopicEventResponse{} }
func (m *TopicEventResponse) String() string { return proto.CompactTextString(m) }
func (*TopicEventResponse) ProtoMessage()    {}
func (*TopicEventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f48007b62b0c167c, []int{1}
}

func (m *TopicEventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TopicEventResponse.Unmarshal(m, b)
}
func (m *TopicEventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TopicEventResponse.Marshal(b, m, deterministic)
}
func (m *TopicEventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TopicEventResponse.Merge(m, src)
}
func (m *TopicEventResponse) XXX_Size() int {
	return xxx_messageInfo_TopicEventResponse.Size(m)
}
func (m *TopicEventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TopicEventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TopicEventResponse proto.InternalMessageInfo

func (m *TopicEventResponse) GetStatus() TopicEventResponse_TopicEventResponseStatus {
	if m != nil {
		return m.Status
	}
	return TopicEventResponse_SUCCESS
}

// ListTopicSubscriptionsResponse is the message including the list of the subscribing topics.
type ListTopicSubscriptionsResponse struct {
	// The list of topics.
	Subscriptions        []*TopicSubscription `protobuf:"bytes,1,rep,name=subscriptions,proto3" json:"subscriptions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ListTopicSubscriptionsResponse) Reset()         { *m = ListTopicSubscriptionsResponse{} }
func (m *ListTopicSubscriptionsResponse) String() string { return proto.CompactTextString(m) }
func (*ListTopicSubscriptionsResponse) ProtoMessage()    {}
func (*ListTopicSubscriptionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f48007b62b0c167c, []int{2}
}

func (m *ListTopicSubscriptionsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListTopicSubscriptionsResponse.Unmarshal(m, b)
}
func (m *ListTopicSubscriptionsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListTopicSubscriptionsResponse.Marshal(b, m, deterministic)
}
func (m *ListTopicSubscriptionsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListTopicSubscriptionsResponse.Merge(m, src)
}
func (m *ListTopicSubscriptionsResponse) XXX_Size() int {
	return xxx_messageInfo_ListTopicSubscriptionsResponse.Size(m)
}
func (m *ListTopicSubscriptionsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListTopicSubscriptionsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListTopicSubscriptionsResponse proto.InternalMessageInfo

func (m *ListTopicSubscriptionsResponse) GetSubscriptions() []*TopicSubscription {
	if m != nil {
		return m.Subscriptions
	}
	return nil
}

// TopicSubscription represents topic and metadata.
type TopicSubscription struct {
	// Required. The name of the pubsub containing the topic below to subscribe to.
	PubsubName string `protobuf:"bytes,1,opt,name=pubsub_name,json=pubsubName,proto3" json:"pubsub_name,omitempty"`
	// Required. The name of topic which will be subscribed
	Topic string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	// The optional properties used for this topic's subscription e.g. session id
	Metadata             map[string]string `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TopicSubscription) Reset()         { *m = TopicSubscription{} }
func (m *TopicSubscription) String() string { return proto.CompactTextString(m) }
func (*TopicSubscription) ProtoMessage()    {}
func (*TopicSubscription) Descriptor() ([]byte, []int) {
	return fileDescriptor_f48007b62b0c167c, []int{3}
}

func (m *TopicSubscription) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TopicSubscription.Unmarshal(m, b)
}
func (m *TopicSubscription) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TopicSubscription.Marshal(b, m, deterministic)
}
func (m *TopicSubscription) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TopicSubscription.Merge(m, src)
}
func (m *TopicSubscription) XXX_Size() int {
	return xxx_messageInfo_TopicSubscription.Size(m)
}
func (m *TopicSubscription) XXX_DiscardUnknown() {
	xxx_messageInfo_TopicSubscription.DiscardUnknown(m)
}

var xxx_messageInfo_TopicSubscription proto.InternalMessageInfo

func (m *TopicSubscription) GetPubsubName() string {
	if m != nil {
		return m.PubsubName
	}
	return ""
}

func (m *TopicSubscription) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *TopicSubscription) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterEnum("spec.proto.runtime.v1.TopicEventResponse_TopicEventResponseStatus", TopicEventResponse_TopicEventResponseStatus_name, TopicEventResponse_TopicEventResponseStatus_value)
	proto.RegisterType((*TopicEventRequest)(nil), "spec.proto.runtime.v1.TopicEventRequest")
	proto.RegisterType((*TopicEventResponse)(nil), "spec.proto.runtime.v1.TopicEventResponse")
	proto.RegisterType((*ListTopicSubscriptionsResponse)(nil), "spec.proto.runtime.v1.ListTopicSubscriptionsResponse")
	proto.RegisterType((*TopicSubscription)(nil), "spec.proto.runtime.v1.TopicSubscription")
	proto.RegisterMapType((map[string]string)(nil), "spec.proto.runtime.v1.TopicSubscription.MetadataEntry")
}

func init() {
	proto.RegisterFile("appcallback.proto", fileDescriptor_f48007b62b0c167c)
}

var fileDescriptor_f48007b62b0c167c = []byte{
	// 532 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xae, 0x9d, 0x9f, 0xa6, 0x93, 0xb4, 0x24, 0x23, 0x88, 0xac, 0x20, 0x41, 0xf0, 0x29, 0x70,
	0x70, 0xd4, 0x20, 0x50, 0xd5, 0x72, 0xa1, 0x21, 0x37, 0x68, 0x91, 0x13, 0x90, 0xe8, 0x25, 0x5a,
	0x3b, 0x4b, 0xb0, 0x1a, 0x7b, 0x17, 0xef, 0x6e, 0xa4, 0xbc, 0x14, 0x2f, 0xc4, 0x81, 0x33, 0x6f,
	0x81, 0x76, 0xed, 0xa4, 0x6e, 0x49, 0xab, 0x9c, 0x3c, 0xfb, 0xcd, 0xcc, 0xa7, 0x6f, 0xfe, 0x0c,
	0x2d, 0xc2, 0x79, 0x48, 0x16, 0x8b, 0x80, 0x84, 0xd7, 0x1e, 0x4f, 0x99, 0x64, 0xf8, 0x44, 0x70,
	0x1a, 0x66, 0xb6, 0x97, 0xaa, 0x44, 0x46, 0x31, 0xf5, 0x96, 0xc7, 0x9d, 0xa7, 0x73, 0xc6, 0xe6,
	0x0b, 0xda, 0x37, 0x8e, 0x40, 0x7d, 0xef, 0xd3, 0x98, 0xcb, 0x55, 0x16, 0xe7, 0xfe, 0xb5, 0xa0,
	0x35, 0x61, 0x3c, 0x0a, 0x47, 0x4b, 0x9a, 0x48, 0x9f, 0xfe, 0x54, 0x54, 0x48, 0x3c, 0x02, 0x3b,
	0x9a, 0x39, 0x56, 0xd7, 0xea, 0x1d, 0xf8, 0x76, 0x34, 0xc3, 0x36, 0x54, 0x05, 0x53, 0x69, 0x48,
	0x1d, 0xdb, 0x60, 0xf9, 0x0b, 0x11, 0xca, 0x72, 0xc5, 0xa9, 0x53, 0x32, 0xa8, 0xb1, 0xf1, 0x05,
	0x34, 0xb4, 0x8e, 0xe9, 0x92, 0xa6, 0x22, 0x62, 0x89, 0x53, 0x36, 0xbe, 0xba, 0xc6, 0xbe, 0x66,
	0x10, 0xbe, 0x82, 0xd6, 0x8c, 0x48, 0x32, 0x0d, 0x59, 0x22, 0x69, 0x22, 0xa7, 0x86, 0xa3, 0x62,
	0xe2, 0x1e, 0x69, 0xc7, 0x30, 0xc3, 0x27, 0x9a, 0x0e, 0xa1, 0xac, 0x21, 0x67, 0xbf, 0x6b, 0xf5,
	0x1a, 0xbe, 0xb1, 0xf1, 0x31, 0x54, 0xa4, 0xd6, 0xec, 0x54, 0x4d, 0x4e, 0xf6, 0xc0, 0xe7, 0x50,
	0xe7, 0x2a, 0x10, 0x2a, 0x98, 0x26, 0x24, 0xa6, 0x4e, 0xcd, 0xf8, 0x20, 0x83, 0x2e, 0x48, 0x4c,
	0xdd, 0x5f, 0x16, 0x60, 0xb1, 0x56, 0xc1, 0x59, 0x22, 0x28, 0x5e, 0x41, 0x55, 0x48, 0x22, 0x95,
	0x30, 0x05, 0x1f, 0x0d, 0xce, 0xbd, 0xad, 0x7d, 0xf4, 0xfe, 0x4f, 0xdd, 0x02, 0x8d, 0x0d, 0x93,
	0x9f, 0x33, 0xba, 0xef, 0xc0, 0xb9, 0x2f, 0x06, 0xeb, 0xb0, 0x3f, 0xfe, 0x32, 0x1c, 0x8e, 0xc6,
	0xe3, 0xe6, 0x1e, 0x1e, 0x40, 0xc5, 0x1f, 0x4d, 0xfc, 0x6f, 0x4d, 0x0b, 0x6b, 0x50, 0xfe, 0xe0,
	0x5f, 0x7e, 0x6e, 0xda, 0x2e, 0x87, 0x67, 0x1f, 0x23, 0x21, 0x0d, 0xc3, 0x58, 0x05, 0x22, 0x4c,
	0x23, 0x2e, 0x23, 0x96, 0x88, 0x8d, 0xf6, 0x0b, 0x38, 0x14, 0x45, 0x87, 0x63, 0x75, 0x4b, 0xbd,
	0xfa, 0xa0, 0xf7, 0x50, 0x09, 0x45, 0x26, 0xff, 0x76, 0xba, 0xfb, 0x7b, 0xbd, 0x0e, 0xc5, 0xa0,
	0xbb, 0x9d, 0xb5, 0xee, 0x76, 0xf6, 0x66, 0x20, 0x76, 0x71, 0x20, 0x3e, 0xd4, 0x62, 0x2a, 0x89,
	0x19, 0x5f, 0xc9, 0xe8, 0x7a, 0xbb, 0xab, 0x2e, 0xef, 0x53, 0x9e, 0x38, 0x4a, 0x64, 0xba, 0xf2,
	0x37, 0x3c, 0x9d, 0x33, 0x38, 0xbc, 0xe5, 0xc2, 0x26, 0x94, 0xae, 0xe9, 0x2a, 0xd7, 0xa4, 0x4d,
	0x2d, 0x66, 0x49, 0x16, 0x6a, 0xbd, 0xab, 0xd9, 0xe3, 0xd4, 0x3e, 0xb1, 0x06, 0x7f, 0x2c, 0xa8,
	0xbf, 0xe7, 0x7c, 0x98, 0x9f, 0x0d, 0x46, 0xd0, 0xde, 0xde, 0x5f, 0x6c, 0x7b, 0xd9, 0xd1, 0x78,
	0xeb, 0xa3, 0xf1, 0x46, 0xfa, 0x68, 0x3a, 0x6f, 0xee, 0x29, 0xe0, 0xe1, 0x31, 0xb9, 0x7b, 0x48,
	0xa1, 0x71, 0x99, 0xdc, 0xac, 0x02, 0xf6, 0x76, 0x58, 0x32, 0x73, 0x8b, 0x9d, 0x97, 0x3b, 0xaf,
	0xa3, 0xbb, 0x77, 0x7e, 0x7a, 0x75, 0x32, 0x8f, 0xe4, 0x0f, 0x15, 0x78, 0x21, 0x8b, 0xfb, 0x0b,
	0xb2, 0x62, 0x52, 0xb2, 0xcd, 0x57, 0x13, 0x65, 0xbf, 0x81, 0x7e, 0x4e, 0xd4, 0x5f, 0x1e, 0x9f,
	0xe5, 0x66, 0x50, 0x35, 0x9e, 0xd7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x4c, 0xd2, 0xcb,
	0x5a, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AppCallbackClient is the client API for AppCallback service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AppCallbackClient interface {
	// Lists all topics subscribed by this app.
	ListTopicSubscriptions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTopicSubscriptionsResponse, error)
	// Subscribes events from Pubsub
	OnTopicEvent(ctx context.Context, in *TopicEventRequest, opts ...grpc.CallOption) (*TopicEventResponse, error)
}

type appCallbackClient struct {
	cc grpc.ClientConnInterface
}

func NewAppCallbackClient(cc grpc.ClientConnInterface) AppCallbackClient {
	return &appCallbackClient{cc}
}

func (c *appCallbackClient) ListTopicSubscriptions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTopicSubscriptionsResponse, error) {
	out := new(ListTopicSubscriptionsResponse)
	err := c.cc.Invoke(ctx, "/spec.proto.runtime.v1.AppCallback/ListTopicSubscriptions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appCallbackClient) OnTopicEvent(ctx context.Context, in *TopicEventRequest, opts ...grpc.CallOption) (*TopicEventResponse, error) {
	out := new(TopicEventResponse)
	err := c.cc.Invoke(ctx, "/spec.proto.runtime.v1.AppCallback/OnTopicEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppCallbackServer is the server API for AppCallback service.
type AppCallbackServer interface {
	// Lists all topics subscribed by this app.
	ListTopicSubscriptions(context.Context, *emptypb.Empty) (*ListTopicSubscriptionsResponse, error)
	// Subscribes events from Pubsub
	OnTopicEvent(context.Context, *TopicEventRequest) (*TopicEventResponse, error)
}

// UnimplementedAppCallbackServer can be embedded to have forward compatible implementations.
type UnimplementedAppCallbackServer struct {
}

func (*UnimplementedAppCallbackServer) ListTopicSubscriptions(ctx context.Context, req *emptypb.Empty) (*ListTopicSubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopicSubscriptions not implemented")
}
func (*UnimplementedAppCallbackServer) OnTopicEvent(ctx context.Context, req *TopicEventRequest) (*TopicEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnTopicEvent not implemented")
}

func RegisterAppCallbackServer(s *grpc.Server, srv AppCallbackServer) {
	s.RegisterService(&_AppCallback_serviceDesc, srv)
}

func _AppCallback_ListTopicSubscriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).ListTopicSubscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spec.proto.runtime.v1.AppCallback/ListTopicSubscriptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).ListTopicSubscriptions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppCallback_OnTopicEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).OnTopicEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spec.proto.runtime.v1.AppCallback/OnTopicEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).OnTopicEvent(ctx, req.(*TopicEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AppCallback_serviceDesc = grpc.ServiceDesc{
	ServiceName: "spec.proto.runtime.v1.AppCallback",
	HandlerType: (*AppCallbackServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTopicSubscriptions",
			Handler:    _AppCallback_ListTopicSubscriptions_Handler,
		},
		{
			MethodName: "OnTopicEvent",
			Handler:    _AppCallback_OnTopicEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "appcallback.proto",
}
