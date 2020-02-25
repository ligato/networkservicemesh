// Code generated by protoc-gen-go. DO NOT EDIT.
// source: update.proto

package update

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	networkservice "github.com/networkservicemesh/api/pkg/api/networkservice"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type AddDNSContextMessage struct {
	Context              *networkservice.DNSContext `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	ConnectionID         string                     `protobuf:"bytes,2,opt,name=connectionID,proto3" json:"connectionID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *AddDNSContextMessage) Reset()         { *m = AddDNSContextMessage{} }
func (m *AddDNSContextMessage) String() string { return proto.CompactTextString(m) }
func (*AddDNSContextMessage) ProtoMessage()    {}
func (*AddDNSContextMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0fa214029f1c21, []int{0}
}

func (m *AddDNSContextMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddDNSContextMessage.Unmarshal(m, b)
}
func (m *AddDNSContextMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddDNSContextMessage.Marshal(b, m, deterministic)
}
func (m *AddDNSContextMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddDNSContextMessage.Merge(m, src)
}
func (m *AddDNSContextMessage) XXX_Size() int {
	return xxx_messageInfo_AddDNSContextMessage.Size(m)
}
func (m *AddDNSContextMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_AddDNSContextMessage.DiscardUnknown(m)
}

var xxx_messageInfo_AddDNSContextMessage proto.InternalMessageInfo

func (m *AddDNSContextMessage) GetContext() *networkservice.DNSContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *AddDNSContextMessage) GetConnectionID() string {
	if m != nil {
		return m.ConnectionID
	}
	return ""
}

type RemoveDNSContextMessage struct {
	ConnectionID         string   `protobuf:"bytes,1,opt,name=connectionID,proto3" json:"connectionID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveDNSContextMessage) Reset()         { *m = RemoveDNSContextMessage{} }
func (m *RemoveDNSContextMessage) String() string { return proto.CompactTextString(m) }
func (*RemoveDNSContextMessage) ProtoMessage()    {}
func (*RemoveDNSContextMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0fa214029f1c21, []int{1}
}

func (m *RemoveDNSContextMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveDNSContextMessage.Unmarshal(m, b)
}
func (m *RemoveDNSContextMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveDNSContextMessage.Marshal(b, m, deterministic)
}
func (m *RemoveDNSContextMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveDNSContextMessage.Merge(m, src)
}
func (m *RemoveDNSContextMessage) XXX_Size() int {
	return xxx_messageInfo_RemoveDNSContextMessage.Size(m)
}
func (m *RemoveDNSContextMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveDNSContextMessage.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveDNSContextMessage proto.InternalMessageInfo

func (m *RemoveDNSContextMessage) GetConnectionID() string {
	if m != nil {
		return m.ConnectionID
	}
	return ""
}

func init() {
	proto.RegisterType((*AddDNSContextMessage)(nil), "update.AddDNSContextMessage")
	proto.RegisterType((*RemoveDNSContextMessage)(nil), "update.RemoveDNSContextMessage")
}

func init() { proto.RegisterFile("update.proto", fileDescriptor_3f0fa214029f1c21) }

var fileDescriptor_3f0fa214029f1c21 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x2d, 0x48, 0x49,
	0x2c, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0xa4, 0x24, 0x0a, 0x4a,
	0x2a, 0x0b, 0x52, 0x8b, 0xf5, 0x53, 0x73, 0x0b, 0x4a, 0x2a, 0x21, 0x24, 0x44, 0x85, 0x94, 0x78,
	0x72, 0x7e, 0x5e, 0x5e, 0x6a, 0x72, 0x49, 0x66, 0x7e, 0x5e, 0x72, 0x7e, 0x5e, 0x49, 0x6a, 0x45,
	0x09, 0x44, 0x42, 0xa9, 0x98, 0x4b, 0xc4, 0x31, 0x25, 0xc5, 0xc5, 0x2f, 0xd8, 0x19, 0x22, 0xec,
	0x9b, 0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0x2a, 0x64, 0xce, 0xc5, 0x0e, 0x55, 0x28, 0xc1, 0xa8, 0xc0,
	0xa8, 0xc1, 0x6d, 0x24, 0xab, 0x87, 0x69, 0x04, 0x42, 0x5b, 0x10, 0x4c, 0xb5, 0x90, 0x12, 0x17,
	0x0f, 0x42, 0xa1, 0xa7, 0x8b, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8a, 0x98, 0x92, 0x2d,
	0x97, 0x78, 0x50, 0x6a, 0x6e, 0x7e, 0x59, 0x2a, 0xa6, 0xbd, 0xe8, 0xda, 0x19, 0x31, 0xb5, 0x1b,
	0x2d, 0x63, 0xe4, 0x12, 0x80, 0xe8, 0x4c, 0xcb, 0x4c, 0x0f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e,
	0x15, 0x72, 0xe5, 0xe2, 0x45, 0xf1, 0x88, 0x90, 0x8c, 0x1e, 0x34, 0x8c, 0xb0, 0xf9, 0x4f, 0x4a,
	0x4c, 0x2f, 0x3d, 0x3f, 0x3f, 0x3d, 0x07, 0x1a, 0x82, 0x49, 0xa5, 0x69, 0x7a, 0xae, 0xa0, 0xe0,
	0x12, 0xf2, 0xe6, 0x12, 0x40, 0x77, 0x9a, 0x90, 0x3c, 0xcc, 0x24, 0x1c, 0x8e, 0xc6, 0x65, 0x58,
	0x12, 0x1b, 0x98, 0x6f, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xce, 0xad, 0x50, 0xe6, 0xae, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DNSConfigServiceClient is the client API for DNSConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DNSConfigServiceClient interface {
	AddDNSContext(ctx context.Context, in *AddDNSContextMessage, opts ...grpc.CallOption) (*empty.Empty, error)
	RemoveDNSContext(ctx context.Context, in *RemoveDNSContextMessage, opts ...grpc.CallOption) (*empty.Empty, error)
}

type dNSConfigServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDNSConfigServiceClient(cc grpc.ClientConnInterface) DNSConfigServiceClient {
	return &dNSConfigServiceClient{cc}
}

func (c *dNSConfigServiceClient) AddDNSContext(ctx context.Context, in *AddDNSContextMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/update.DNSConfigService/AddDNSContext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dNSConfigServiceClient) RemoveDNSContext(ctx context.Context, in *RemoveDNSContextMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/update.DNSConfigService/RemoveDNSContext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DNSConfigServiceServer is the server API for DNSConfigService service.
type DNSConfigServiceServer interface {
	AddDNSContext(context.Context, *AddDNSContextMessage) (*empty.Empty, error)
	RemoveDNSContext(context.Context, *RemoveDNSContextMessage) (*empty.Empty, error)
}

// UnimplementedDNSConfigServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDNSConfigServiceServer struct {
}

func (*UnimplementedDNSConfigServiceServer) AddDNSContext(ctx context.Context, req *AddDNSContextMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDNSContext not implemented")
}
func (*UnimplementedDNSConfigServiceServer) RemoveDNSContext(ctx context.Context, req *RemoveDNSContextMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveDNSContext not implemented")
}

func RegisterDNSConfigServiceServer(s *grpc.Server, srv DNSConfigServiceServer) {
	s.RegisterService(&_DNSConfigService_serviceDesc, srv)
}

func _DNSConfigService_AddDNSContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDNSContextMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DNSConfigServiceServer).AddDNSContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/update.DNSConfigService/AddDNSContext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DNSConfigServiceServer).AddDNSContext(ctx, req.(*AddDNSContextMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _DNSConfigService_RemoveDNSContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveDNSContextMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DNSConfigServiceServer).RemoveDNSContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/update.DNSConfigService/RemoveDNSContext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DNSConfigServiceServer).RemoveDNSContext(ctx, req.(*RemoveDNSContextMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _DNSConfigService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "update.DNSConfigService",
	HandlerType: (*DNSConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddDNSContext",
			Handler:    _DNSConfigService_AddDNSContext_Handler,
		},
		{
			MethodName: "RemoveDNSContext",
			Handler:    _DNSConfigService_RemoveDNSContext_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "update.proto",
}
