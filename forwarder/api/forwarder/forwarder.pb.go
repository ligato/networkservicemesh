// Code generated by protoc-gen-go. DO NOT EDIT.
// source: forwarder.proto

package forwarder

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	networkservice "github.com/networkservicemesh/api/pkg/api/networkservice"
	crossconnect "github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
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

// Message sent by forwarder module informing NSM of any changes in its
// operations parameters or constraints
type MechanismUpdate struct {
	RemoteMechanisms     []*networkservice.Mechanism `protobuf:"bytes,1,rep,name=remote_mechanisms,json=remoteMechanisms,proto3" json:"remote_mechanisms,omitempty"`
	LocalMechanisms      []*networkservice.Mechanism `protobuf:"bytes,2,rep,name=local_mechanisms,json=localMechanisms,proto3" json:"local_mechanisms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *MechanismUpdate) Reset()         { *m = MechanismUpdate{} }
func (m *MechanismUpdate) String() string { return proto.CompactTextString(m) }
func (*MechanismUpdate) ProtoMessage()    {}
func (*MechanismUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_19bff53f4d11db23, []int{0}
}

func (m *MechanismUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MechanismUpdate.Unmarshal(m, b)
}
func (m *MechanismUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MechanismUpdate.Marshal(b, m, deterministic)
}
func (m *MechanismUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MechanismUpdate.Merge(m, src)
}
func (m *MechanismUpdate) XXX_Size() int {
	return xxx_messageInfo_MechanismUpdate.Size(m)
}
func (m *MechanismUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_MechanismUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_MechanismUpdate proto.InternalMessageInfo

func (m *MechanismUpdate) GetRemoteMechanisms() []*networkservice.Mechanism {
	if m != nil {
		return m.RemoteMechanisms
	}
	return nil
}

func (m *MechanismUpdate) GetLocalMechanisms() []*networkservice.Mechanism {
	if m != nil {
		return m.LocalMechanisms
	}
	return nil
}

func init() {
	proto.RegisterType((*MechanismUpdate)(nil), "forwarder.MechanismUpdate")
}

func init() { proto.RegisterFile("forwarder.proto", fileDescriptor_19bff53f4d11db23) }

var fileDescriptor_19bff53f4d11db23 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0xb1, 0x4e, 0xc3, 0x30,
	0x14, 0x54, 0x40, 0x80, 0xfa, 0x18, 0x92, 0x58, 0x02, 0x55, 0x9e, 0x10, 0x13, 0x93, 0x41, 0x61,
	0x64, 0x01, 0x22, 0x90, 0x18, 0xba, 0x54, 0x62, 0x05, 0xa5, 0xee, 0x6b, 0x89, 0x94, 0xf8, 0x19,
	0xfb, 0x55, 0xa8, 0xdf, 0xc0, 0x07, 0xf0, 0xbb, 0xa8, 0x75, 0x9a, 0x04, 0xa4, 0x66, 0xb1, 0xfc,
	0xee, 0xce, 0x67, 0xdf, 0x19, 0xe2, 0x05, 0xb9, 0xaf, 0xc2, 0xcd, 0xd1, 0x29, 0xeb, 0x88, 0x49,
	0x8c, 0x5a, 0x40, 0x26, 0x9a, 0x8c, 0x41, 0xcd, 0x25, 0x99, 0x40, 0x4a, 0xa1, 0x1d, 0x79, 0xdf,
	0xc0, 0x0d, 0x36, 0xb6, 0xbc, 0xb6, 0xe8, 0xaf, 0xb1, 0xb6, 0xbc, 0x0e, 0x6b, 0x60, 0x2e, 0x7f,
	0x22, 0x88, 0x27, 0xa8, 0x3f, 0x0a, 0x53, 0xfa, 0xfa, 0xd5, 0xce, 0x0b, 0x46, 0xf1, 0x08, 0xa9,
	0xc3, 0x9a, 0x18, 0xdf, 0xeb, 0x1d, 0xe3, 0xc7, 0xd1, 0xc5, 0xe1, 0xd5, 0x69, 0x76, 0xa6, 0x7a,
	0xf7, 0xb5, 0xe7, 0xa6, 0x49, 0xd0, 0xb7, 0x80, 0x17, 0xf7, 0x90, 0x54, 0xa4, 0x8b, 0xaa, 0x6f,
	0x71, 0x30, 0x64, 0x11, 0x6f, 0xe5, 0x9d, 0x43, 0xf6, 0x1d, 0xc1, 0xe8, 0x79, 0x97, 0x53, 0x3c,
	0xc0, 0xc9, 0x14, 0x3f, 0x57, 0xe8, 0x59, 0x48, 0xf5, 0x27, 0x61, 0xbe, 0x19, 0xf2, 0x30, 0xc8,
	0x01, 0x4e, 0xdc, 0xc1, 0x51, 0x5e, 0x91, 0xc7, 0x41, 0x83, 0x73, 0xb5, 0x24, 0x5a, 0x56, 0x18,
	0xea, 0x99, 0xad, 0x16, 0xea, 0x69, 0xd3, 0x56, 0xf6, 0x06, 0x69, 0xf7, 0xb6, 0x09, 0x99, 0x92,
	0xc9, 0x89, 0x17, 0x48, 0x9b, 0x6d, 0x2f, 0xf9, 0x1e, 0x07, 0x29, 0x55, 0xf7, 0x8d, 0xff, 0x1a,
	0xbf, 0x89, 0x66, 0xc7, 0x5b, 0xf5, 0xed, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x82, 0x8d, 0x40,
	0x63, 0xec, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ForwarderClient is the client API for Forwarder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ForwarderClient interface {
	Request(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*crossconnect.CrossConnect, error)
	Close(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*empty.Empty, error)
}

type forwarderClient struct {
	cc grpc.ClientConnInterface
}

func NewForwarderClient(cc grpc.ClientConnInterface) ForwarderClient {
	return &forwarderClient{cc}
}

func (c *forwarderClient) Request(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*crossconnect.CrossConnect, error) {
	out := new(crossconnect.CrossConnect)
	err := c.cc.Invoke(ctx, "/forwarder.Forwarder/Request", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forwarderClient) Close(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/forwarder.Forwarder/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForwarderServer is the server API for Forwarder service.
type ForwarderServer interface {
	Request(context.Context, *crossconnect.CrossConnect) (*crossconnect.CrossConnect, error)
	Close(context.Context, *crossconnect.CrossConnect) (*empty.Empty, error)
}

// UnimplementedForwarderServer can be embedded to have forward compatible implementations.
type UnimplementedForwarderServer struct {
}

func (*UnimplementedForwarderServer) Request(ctx context.Context, req *crossconnect.CrossConnect) (*crossconnect.CrossConnect, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Request not implemented")
}
func (*UnimplementedForwarderServer) Close(ctx context.Context, req *crossconnect.CrossConnect) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}

func RegisterForwarderServer(s *grpc.Server, srv ForwarderServer) {
	s.RegisterService(&_Forwarder_serviceDesc, srv)
}

func _Forwarder_Request_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(crossconnect.CrossConnect)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForwarderServer).Request(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forwarder.Forwarder/Request",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForwarderServer).Request(ctx, req.(*crossconnect.CrossConnect))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forwarder_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(crossconnect.CrossConnect)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForwarderServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forwarder.Forwarder/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForwarderServer).Close(ctx, req.(*crossconnect.CrossConnect))
	}
	return interceptor(ctx, in, info, handler)
}

var _Forwarder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "forwarder.Forwarder",
	HandlerType: (*ForwarderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Request",
			Handler:    _Forwarder_Request_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Forwarder_Close_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "forwarder.proto",
}

// MechanismsMonitorClient is the client API for MechanismsMonitor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MechanismsMonitorClient interface {
	MonitorMechanisms(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (MechanismsMonitor_MonitorMechanismsClient, error)
}

type mechanismsMonitorClient struct {
	cc grpc.ClientConnInterface
}

func NewMechanismsMonitorClient(cc grpc.ClientConnInterface) MechanismsMonitorClient {
	return &mechanismsMonitorClient{cc}
}

func (c *mechanismsMonitorClient) MonitorMechanisms(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (MechanismsMonitor_MonitorMechanismsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MechanismsMonitor_serviceDesc.Streams[0], "/forwarder.MechanismsMonitor/MonitorMechanisms", opts...)
	if err != nil {
		return nil, err
	}
	x := &mechanismsMonitorMonitorMechanismsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MechanismsMonitor_MonitorMechanismsClient interface {
	Recv() (*MechanismUpdate, error)
	grpc.ClientStream
}

type mechanismsMonitorMonitorMechanismsClient struct {
	grpc.ClientStream
}

func (x *mechanismsMonitorMonitorMechanismsClient) Recv() (*MechanismUpdate, error) {
	m := new(MechanismUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MechanismsMonitorServer is the server API for MechanismsMonitor service.
type MechanismsMonitorServer interface {
	MonitorMechanisms(*empty.Empty, MechanismsMonitor_MonitorMechanismsServer) error
}

// UnimplementedMechanismsMonitorServer can be embedded to have forward compatible implementations.
type UnimplementedMechanismsMonitorServer struct {
}

func (*UnimplementedMechanismsMonitorServer) MonitorMechanisms(req *empty.Empty, srv MechanismsMonitor_MonitorMechanismsServer) error {
	return status.Errorf(codes.Unimplemented, "method MonitorMechanisms not implemented")
}

func RegisterMechanismsMonitorServer(s *grpc.Server, srv MechanismsMonitorServer) {
	s.RegisterService(&_MechanismsMonitor_serviceDesc, srv)
}

func _MechanismsMonitor_MonitorMechanisms_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MechanismsMonitorServer).MonitorMechanisms(m, &mechanismsMonitorMonitorMechanismsServer{stream})
}

type MechanismsMonitor_MonitorMechanismsServer interface {
	Send(*MechanismUpdate) error
	grpc.ServerStream
}

type mechanismsMonitorMonitorMechanismsServer struct {
	grpc.ServerStream
}

func (x *mechanismsMonitorMonitorMechanismsServer) Send(m *MechanismUpdate) error {
	return x.ServerStream.SendMsg(m)
}

var _MechanismsMonitor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "forwarder.MechanismsMonitor",
	HandlerType: (*MechanismsMonitorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MonitorMechanisms",
			Handler:       _MechanismsMonitor_MonitorMechanisms_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "forwarder.proto",
}
