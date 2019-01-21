// Code generated by protoc-gen-go. DO NOT EDIT.
// source: connection.proto

package connection

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import connectioncontext "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/connectioncontext"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MechanismType int32

const (
	MechanismType_NONE          MechanismType = 0
	MechanismType_VXLAN         MechanismType = 1
	MechanismType_VXLAN_GPE     MechanismType = 2
	MechanismType_GRE           MechanismType = 3
	MechanismType_SRV6          MechanismType = 4
	MechanismType_MPLSoEthernet MechanismType = 5
	MechanismType_MPLSoGRE      MechanismType = 6
	MechanismType_MPLSoUDP      MechanismType = 7
)

var MechanismType_name = map[int32]string{
	0: "NONE",
	1: "VXLAN",
	2: "VXLAN_GPE",
	3: "GRE",
	4: "SRV6",
	5: "MPLSoEthernet",
	6: "MPLSoGRE",
	7: "MPLSoUDP",
}
var MechanismType_value = map[string]int32{
	"NONE":          0,
	"VXLAN":         1,
	"VXLAN_GPE":     2,
	"GRE":           3,
	"SRV6":          4,
	"MPLSoEthernet": 5,
	"MPLSoGRE":      6,
	"MPLSoUDP":      7,
}

func (x MechanismType) String() string {
	return proto.EnumName(MechanismType_name, int32(x))
}
func (MechanismType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{0}
}

type State int32

const (
	State_UP   State = 0
	State_DOWN State = 1
)

var State_name = map[int32]string{
	0: "UP",
	1: "DOWN",
}
var State_value = map[string]int32{
	"UP":   0,
	"DOWN": 1,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{1}
}

type ConnectionEventType int32

const (
	ConnectionEventType_INITIAL_STATE_TRANSFER ConnectionEventType = 0
	ConnectionEventType_UPDATE                 ConnectionEventType = 1
	ConnectionEventType_DELETE                 ConnectionEventType = 2
)

var ConnectionEventType_name = map[int32]string{
	0: "INITIAL_STATE_TRANSFER",
	1: "UPDATE",
	2: "DELETE",
}
var ConnectionEventType_value = map[string]int32{
	"INITIAL_STATE_TRANSFER": 0,
	"UPDATE":                 1,
	"DELETE":                 2,
}

func (x ConnectionEventType) String() string {
	return proto.EnumName(ConnectionEventType_name, int32(x))
}
func (ConnectionEventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{2}
}

type Mechanism struct {
	Type                 MechanismType     `protobuf:"varint,1,opt,name=type,proto3,enum=remote.connection.MechanismType" json:"type,omitempty"`
	Parameters           map[string]string `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Mechanism) Reset()         { *m = Mechanism{} }
func (m *Mechanism) String() string { return proto.CompactTextString(m) }
func (*Mechanism) ProtoMessage()    {}
func (*Mechanism) Descriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{0}
}
func (m *Mechanism) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Mechanism.Unmarshal(m, b)
}
func (m *Mechanism) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Mechanism.Marshal(b, m, deterministic)
}
func (dst *Mechanism) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Mechanism.Merge(dst, src)
}
func (m *Mechanism) XXX_Size() int {
	return xxx_messageInfo_Mechanism.Size(m)
}
func (m *Mechanism) XXX_DiscardUnknown() {
	xxx_messageInfo_Mechanism.DiscardUnknown(m)
}

var xxx_messageInfo_Mechanism proto.InternalMessageInfo

func (m *Mechanism) GetType() MechanismType {
	if m != nil {
		return m.Type
	}
	return MechanismType_NONE
}

func (m *Mechanism) GetParameters() map[string]string {
	if m != nil {
		return m.Parameters
	}
	return nil
}

type Connection struct {
	Id                                   string                               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	NetworkService                       string                               `protobuf:"bytes,2,opt,name=network_service,json=networkService,proto3" json:"network_service,omitempty"`
	Mechanism                            *Mechanism                           `protobuf:"bytes,3,opt,name=mechanism,proto3" json:"mechanism,omitempty"`
	Context                              *connectioncontext.ConnectionContext `protobuf:"bytes,4,opt,name=context,proto3" json:"context,omitempty"`
	Labels                               map[string]string                    `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	SourceNetworkServiceManagerName      string                               `protobuf:"bytes,6,opt,name=source_network_service_manager_name,json=sourceNetworkServiceManagerName,proto3" json:"source_network_service_manager_name,omitempty"`
	DestinationNetworkServiceManagerName string                               `protobuf:"bytes,7,opt,name=destination_network_service_manager_name,json=destinationNetworkServiceManagerName,proto3" json:"destination_network_service_manager_name,omitempty"`
	NetworkServiceEndpointName           string                               `protobuf:"bytes,8,opt,name=network_service_endpoint_name,json=networkServiceEndpointName,proto3" json:"network_service_endpoint_name,omitempty"`
	State                                State                                `protobuf:"varint,9,opt,name=state,proto3,enum=remote.connection.State" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral                 struct{}                             `json:"-"`
	XXX_unrecognized                     []byte                               `json:"-"`
	XXX_sizecache                        int32                                `json:"-"`
}

func (m *Connection) Reset()         { *m = Connection{} }
func (m *Connection) String() string { return proto.CompactTextString(m) }
func (*Connection) ProtoMessage()    {}
func (*Connection) Descriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{1}
}
func (m *Connection) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Connection.Unmarshal(m, b)
}
func (m *Connection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Connection.Marshal(b, m, deterministic)
}
func (dst *Connection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Connection.Merge(dst, src)
}
func (m *Connection) XXX_Size() int {
	return xxx_messageInfo_Connection.Size(m)
}
func (m *Connection) XXX_DiscardUnknown() {
	xxx_messageInfo_Connection.DiscardUnknown(m)
}

var xxx_messageInfo_Connection proto.InternalMessageInfo

func (m *Connection) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Connection) GetNetworkService() string {
	if m != nil {
		return m.NetworkService
	}
	return ""
}

func (m *Connection) GetMechanism() *Mechanism {
	if m != nil {
		return m.Mechanism
	}
	return nil
}

func (m *Connection) GetContext() *connectioncontext.ConnectionContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *Connection) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Connection) GetSourceNetworkServiceManagerName() string {
	if m != nil {
		return m.SourceNetworkServiceManagerName
	}
	return ""
}

func (m *Connection) GetDestinationNetworkServiceManagerName() string {
	if m != nil {
		return m.DestinationNetworkServiceManagerName
	}
	return ""
}

func (m *Connection) GetNetworkServiceEndpointName() string {
	if m != nil {
		return m.NetworkServiceEndpointName
	}
	return ""
}

func (m *Connection) GetState() State {
	if m != nil {
		return m.State
	}
	return State_UP
}

type ConnectionEvent struct {
	Type                 ConnectionEventType    `protobuf:"varint,1,opt,name=type,proto3,enum=remote.connection.ConnectionEventType" json:"type,omitempty"`
	Connections          map[string]*Connection `protobuf:"bytes,2,rep,name=connections,proto3" json:"connections,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ConnectionEvent) Reset()         { *m = ConnectionEvent{} }
func (m *ConnectionEvent) String() string { return proto.CompactTextString(m) }
func (*ConnectionEvent) ProtoMessage()    {}
func (*ConnectionEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{2}
}
func (m *ConnectionEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectionEvent.Unmarshal(m, b)
}
func (m *ConnectionEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectionEvent.Marshal(b, m, deterministic)
}
func (dst *ConnectionEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectionEvent.Merge(dst, src)
}
func (m *ConnectionEvent) XXX_Size() int {
	return xxx_messageInfo_ConnectionEvent.Size(m)
}
func (m *ConnectionEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectionEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectionEvent proto.InternalMessageInfo

func (m *ConnectionEvent) GetType() ConnectionEventType {
	if m != nil {
		return m.Type
	}
	return ConnectionEventType_INITIAL_STATE_TRANSFER
}

func (m *ConnectionEvent) GetConnections() map[string]*Connection {
	if m != nil {
		return m.Connections
	}
	return nil
}

type MonitorScopeSelector struct {
	NetworkServiceManagerName string   `protobuf:"bytes,1,opt,name=network_service_manager_name,json=networkServiceManagerName,proto3" json:"network_service_manager_name,omitempty"`
	XXX_NoUnkeyedLiteral      struct{} `json:"-"`
	XXX_unrecognized          []byte   `json:"-"`
	XXX_sizecache             int32    `json:"-"`
}

func (m *MonitorScopeSelector) Reset()         { *m = MonitorScopeSelector{} }
func (m *MonitorScopeSelector) String() string { return proto.CompactTextString(m) }
func (*MonitorScopeSelector) ProtoMessage()    {}
func (*MonitorScopeSelector) Descriptor() ([]byte, []int) {
	return fileDescriptor_connection_aa2f93ec523324fd, []int{3}
}
func (m *MonitorScopeSelector) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorScopeSelector.Unmarshal(m, b)
}
func (m *MonitorScopeSelector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorScopeSelector.Marshal(b, m, deterministic)
}
func (dst *MonitorScopeSelector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorScopeSelector.Merge(dst, src)
}
func (m *MonitorScopeSelector) XXX_Size() int {
	return xxx_messageInfo_MonitorScopeSelector.Size(m)
}
func (m *MonitorScopeSelector) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorScopeSelector.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorScopeSelector proto.InternalMessageInfo

func (m *MonitorScopeSelector) GetNetworkServiceManagerName() string {
	if m != nil {
		return m.NetworkServiceManagerName
	}
	return ""
}

func init() {
	proto.RegisterType((*Mechanism)(nil), "remote.connection.Mechanism")
	proto.RegisterMapType((map[string]string)(nil), "remote.connection.Mechanism.ParametersEntry")
	proto.RegisterType((*Connection)(nil), "remote.connection.Connection")
	proto.RegisterMapType((map[string]string)(nil), "remote.connection.Connection.LabelsEntry")
	proto.RegisterType((*ConnectionEvent)(nil), "remote.connection.ConnectionEvent")
	proto.RegisterMapType((map[string]*Connection)(nil), "remote.connection.ConnectionEvent.ConnectionsEntry")
	proto.RegisterType((*MonitorScopeSelector)(nil), "remote.connection.MonitorScopeSelector")
	proto.RegisterEnum("remote.connection.MechanismType", MechanismType_name, MechanismType_value)
	proto.RegisterEnum("remote.connection.State", State_name, State_value)
	proto.RegisterEnum("remote.connection.ConnectionEventType", ConnectionEventType_name, ConnectionEventType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MonitorConnectionClient is the client API for MonitorConnection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MonitorConnectionClient interface {
	MonitorConnections(ctx context.Context, in *MonitorScopeSelector, opts ...grpc.CallOption) (MonitorConnection_MonitorConnectionsClient, error)
}

type monitorConnectionClient struct {
	cc *grpc.ClientConn
}

func NewMonitorConnectionClient(cc *grpc.ClientConn) MonitorConnectionClient {
	return &monitorConnectionClient{cc}
}

func (c *monitorConnectionClient) MonitorConnections(ctx context.Context, in *MonitorScopeSelector, opts ...grpc.CallOption) (MonitorConnection_MonitorConnectionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MonitorConnection_serviceDesc.Streams[0], "/remote.connection.MonitorConnection/MonitorConnections", opts...)
	if err != nil {
		return nil, err
	}
	x := &monitorConnectionMonitorConnectionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MonitorConnection_MonitorConnectionsClient interface {
	Recv() (*ConnectionEvent, error)
	grpc.ClientStream
}

type monitorConnectionMonitorConnectionsClient struct {
	grpc.ClientStream
}

func (x *monitorConnectionMonitorConnectionsClient) Recv() (*ConnectionEvent, error) {
	m := new(ConnectionEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MonitorConnectionServer is the server API for MonitorConnection service.
type MonitorConnectionServer interface {
	MonitorConnections(*MonitorScopeSelector, MonitorConnection_MonitorConnectionsServer) error
}

func RegisterMonitorConnectionServer(s *grpc.Server, srv MonitorConnectionServer) {
	s.RegisterService(&_MonitorConnection_serviceDesc, srv)
}

func _MonitorConnection_MonitorConnections_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MonitorScopeSelector)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MonitorConnectionServer).MonitorConnections(m, &monitorConnectionMonitorConnectionsServer{stream})
}

type MonitorConnection_MonitorConnectionsServer interface {
	Send(*ConnectionEvent) error
	grpc.ServerStream
}

type monitorConnectionMonitorConnectionsServer struct {
	grpc.ServerStream
}

func (x *monitorConnectionMonitorConnectionsServer) Send(m *ConnectionEvent) error {
	return x.ServerStream.SendMsg(m)
}

var _MonitorConnection_serviceDesc = grpc.ServiceDesc{
	ServiceName: "remote.connection.MonitorConnection",
	HandlerType: (*MonitorConnectionServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MonitorConnections",
			Handler:       _MonitorConnection_MonitorConnections_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "connection.proto",
}

func init() { proto.RegisterFile("connection.proto", fileDescriptor_connection_aa2f93ec523324fd) }

var fileDescriptor_connection_aa2f93ec523324fd = []byte{
	// 717 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x5d, 0x6b, 0xdb, 0x30,
	0x14, 0xad, 0x9d, 0xaf, 0xe6, 0xa6, 0x1f, 0xae, 0x56, 0x86, 0x1b, 0x5a, 0x16, 0xba, 0xb2, 0x66,
	0x65, 0x38, 0x23, 0x1d, 0x63, 0x2b, 0x6c, 0x23, 0x6b, 0xbc, 0x12, 0x48, 0xdc, 0x60, 0x27, 0xed,
	0x18, 0x8c, 0xe0, 0x3a, 0x22, 0x31, 0x8d, 0x25, 0x63, 0x2b, 0x5d, 0xf3, 0xbc, 0xff, 0xb7, 0xe7,
	0xfd, 0x9c, 0x61, 0xd9, 0x8d, 0x9d, 0x8f, 0xa6, 0xec, 0x4d, 0xba, 0xf7, 0x9c, 0x73, 0x75, 0x8f,
	0xae, 0x04, 0x92, 0x45, 0x09, 0xc1, 0x16, 0xb3, 0x29, 0x51, 0x5c, 0x8f, 0x32, 0x8a, 0x76, 0x3c,
	0xec, 0x50, 0x86, 0x95, 0x38, 0x51, 0xec, 0x0f, 0x6c, 0x36, 0x1c, 0xdf, 0x28, 0x16, 0x75, 0x2a,
	0x23, 0x7b, 0x60, 0x32, 0x5a, 0x21, 0x98, 0xfd, 0xa2, 0xde, 0xad, 0x8f, 0xbd, 0x3b, 0xdb, 0xc2,
	0x0e, 0xf6, 0x87, 0x15, 0x8b, 0x12, 0xe6, 0xd1, 0x91, 0x3b, 0x32, 0x09, 0xae, 0xb8, 0xb7, 0x83,
	0x8a, 0xe9, 0xda, 0x7e, 0x25, 0x96, 0x09, 0xf2, 0xf8, 0x9e, 0x2d, 0x46, 0xc2, 0xc2, 0x87, 0x7f,
	0x04, 0xc8, 0xb7, 0xb0, 0x35, 0x34, 0x89, 0xed, 0x3b, 0xe8, 0x1d, 0xa4, 0xd9, 0xc4, 0xc5, 0xb2,
	0x50, 0x12, 0xca, 0x5b, 0xd5, 0x92, 0xb2, 0x70, 0x2a, 0x65, 0x8a, 0xed, 0x4c, 0x5c, 0xac, 0x73,
	0x34, 0x6a, 0x02, 0xb8, 0xa6, 0x67, 0x3a, 0x98, 0x61, 0xcf, 0x97, 0xc5, 0x52, 0xaa, 0x5c, 0xa8,
	0xbe, 0x59, 0xc5, 0x55, 0xda, 0x53, 0xb8, 0x4a, 0x98, 0x37, 0xd1, 0x13, 0xfc, 0xe2, 0x27, 0xd8,
	0x9e, 0x4b, 0x23, 0x09, 0x52, 0xb7, 0x78, 0xc2, 0x4f, 0x95, 0xd7, 0x83, 0x25, 0xda, 0x85, 0xcc,
	0x9d, 0x39, 0x1a, 0x63, 0x59, 0xe4, 0xb1, 0x70, 0x73, 0x26, 0x7e, 0x10, 0x0e, 0xff, 0xa6, 0x01,
	0xce, 0xa7, 0x35, 0xd1, 0x16, 0x88, 0x76, 0x3f, 0x62, 0x8a, 0x76, 0x1f, 0x1d, 0xc3, 0x76, 0xe4,
	0x62, 0x2f, 0xb2, 0x31, 0x92, 0xd8, 0x8a, 0xc2, 0x46, 0x18, 0x45, 0x67, 0x90, 0x77, 0x1e, 0xce,
	0x2b, 0xa7, 0x4a, 0x42, 0xb9, 0x50, 0xdd, 0x5f, 0xd5, 0x93, 0x1e, 0xc3, 0xd1, 0x67, 0xc8, 0x45,
	0x2e, 0xcb, 0x69, 0xce, 0x3c, 0x52, 0x16, 0xfd, 0x8f, 0x0f, 0x79, 0x1e, 0x46, 0xf4, 0x07, 0x12,
	0xaa, 0x41, 0x76, 0x64, 0xde, 0xe0, 0x91, 0x2f, 0x67, 0xb8, 0x99, 0xaf, 0x97, 0x14, 0x8e, 0xe9,
	0x4a, 0x93, 0x63, 0x43, 0x27, 0x23, 0x22, 0x6a, 0xc2, 0x4b, 0x9f, 0x8e, 0x3d, 0x0b, 0xf7, 0xe6,
	0xda, 0xed, 0x39, 0x26, 0x31, 0x07, 0xd8, 0xeb, 0x11, 0xd3, 0xc1, 0x72, 0x96, 0xf7, 0xfe, 0x22,
	0x84, 0x6a, 0x33, 0x0e, 0xb4, 0x42, 0x9c, 0x66, 0x3a, 0x18, 0x5d, 0x41, 0xb9, 0x8f, 0x7d, 0x66,
	0x13, 0x33, 0x28, 0xb8, 0x5a, 0x32, 0xc7, 0x25, 0x8f, 0x12, 0xf8, 0xc7, 0x75, 0x6b, 0x70, 0x30,
	0xaf, 0x85, 0x49, 0xdf, 0xa5, 0x36, 0x61, 0xa1, 0xd8, 0x3a, 0x17, 0x2b, 0xce, 0xde, 0x8d, 0x1a,
	0x41, 0xb8, 0x84, 0x02, 0x19, 0x9f, 0x99, 0x0c, 0xcb, 0x79, 0x3e, 0xb3, 0xf2, 0x12, 0xab, 0x8c,
	0x20, 0xaf, 0x87, 0xb0, 0xe2, 0x47, 0x28, 0x24, 0xfc, 0xfa, 0xaf, 0xd1, 0xfa, 0x2d, 0xc2, 0x76,
	0x6c, 0xbb, 0x7a, 0x87, 0x09, 0x43, 0x67, 0x33, 0x2f, 0xe6, 0xd5, 0xca, 0x8b, 0xe2, 0x8c, 0xc4,
	0xbb, 0xe9, 0x42, 0x21, 0xc6, 0x3d, 0x3c, 0x9c, 0xd3, 0xa7, 0x25, 0x12, 0xfb, 0xe8, 0xd6, 0x93,
	0x3a, 0xc5, 0x9f, 0x20, 0xcd, 0x03, 0x96, 0xb4, 0x79, 0x9a, 0x6c, 0xb3, 0x50, 0x3d, 0x58, 0x59,
	0x36, 0xe9, 0xc2, 0x35, 0xec, 0xb6, 0x28, 0xb1, 0x19, 0xf5, 0x0c, 0x8b, 0xba, 0xd8, 0xc0, 0x23,
	0x6c, 0x31, 0xea, 0xa1, 0x2f, 0xb0, 0xbf, 0x72, 0x2e, 0xc2, 0xda, 0x7b, 0xe4, 0xb1, 0x61, 0x38,
	0x19, 0xc3, 0xe6, 0xcc, 0xef, 0x82, 0xd6, 0x21, 0xad, 0x5d, 0x6a, 0xaa, 0xb4, 0x86, 0xf2, 0x90,
	0xb9, 0xfa, 0xde, 0xac, 0x69, 0x92, 0x80, 0x36, 0x21, 0xcf, 0x97, 0xbd, 0x8b, 0xb6, 0x2a, 0x89,
	0x28, 0x07, 0xa9, 0x0b, 0x5d, 0x95, 0x52, 0x01, 0xd8, 0xd0, 0xaf, 0xde, 0x4b, 0x69, 0xb4, 0x03,
	0x9b, 0xad, 0x76, 0xd3, 0xa0, 0x2a, 0x1b, 0x62, 0x8f, 0x60, 0x26, 0x65, 0xd0, 0x06, 0xac, 0xf3,
	0x50, 0x00, 0xcd, 0x4e, 0x77, 0xdd, 0x7a, 0x5b, 0xca, 0x9d, 0xec, 0x41, 0x86, 0x0f, 0x08, 0xca,
	0x82, 0xd8, 0x6d, 0x4b, 0x6b, 0x81, 0x52, 0xfd, 0xf2, 0x5a, 0x93, 0x84, 0x93, 0x06, 0x3c, 0x5b,
	0x72, 0x7b, 0xa8, 0x08, 0xcf, 0x1b, 0x5a, 0xa3, 0xd3, 0xa8, 0x35, 0x7b, 0x46, 0xa7, 0xd6, 0x51,
	0x7b, 0x1d, 0xbd, 0xa6, 0x19, 0xdf, 0x54, 0x5d, 0x5a, 0x43, 0x00, 0xd9, 0x6e, 0xbb, 0x5e, 0xeb,
	0xa8, 0x92, 0x10, 0xac, 0xeb, 0x6a, 0x53, 0xed, 0xa8, 0x92, 0x58, 0xbd, 0x87, 0x9d, 0xc8, 0xb5,
	0xc4, 0xe7, 0x64, 0x01, 0x5a, 0x08, 0xfa, 0xe8, 0x78, 0xd9, 0x37, 0xb3, 0xc4, 0xf1, 0xe2, 0xe1,
	0xd3, 0xa3, 0xf2, 0x56, 0xf8, 0xba, 0xf1, 0x03, 0xe2, 0xfc, 0x4d, 0x96, 0x7f, 0xfb, 0xa7, 0xff,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x53, 0x3c, 0xe3, 0x1b, 0x83, 0x06, 0x00, 0x00,
}
