// Code generated by protoc-gen-go. DO NOT EDIT.
// source: netmesh.proto

package netmesh

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/ligato/networkservicemesh/pkg/nsm/apis/common"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NetworkServiceChannel struct {
	Metadata             *common.Metadata    `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	NetworkServiceName   string              `protobuf:"bytes,2,opt,name=network_service_name,json=networkServiceName,proto3" json:"network_service_name,omitempty"`
	Payload              string              `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	SocketLocation       string              `protobuf:"bytes,4,opt,name=socket_location,json=socketLocation,proto3" json:"socket_location,omitempty"`
	Interface            []*common.Interface `protobuf:"bytes,5,rep,name=interface,proto3" json:"interface,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *NetworkServiceChannel) Reset()         { *m = NetworkServiceChannel{} }
func (m *NetworkServiceChannel) String() string { return proto.CompactTextString(m) }
func (*NetworkServiceChannel) ProtoMessage()    {}
func (*NetworkServiceChannel) Descriptor() ([]byte, []int) {
	return fileDescriptor_netmesh_8944ebf4e75e6af3, []int{0}
}
func (m *NetworkServiceChannel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkServiceChannel.Unmarshal(m, b)
}
func (m *NetworkServiceChannel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkServiceChannel.Marshal(b, m, deterministic)
}
func (dst *NetworkServiceChannel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkServiceChannel.Merge(dst, src)
}
func (m *NetworkServiceChannel) XXX_Size() int {
	return xxx_messageInfo_NetworkServiceChannel.Size(m)
}
func (m *NetworkServiceChannel) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkServiceChannel.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkServiceChannel proto.InternalMessageInfo

func (m *NetworkServiceChannel) GetMetadata() *common.Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *NetworkServiceChannel) GetNetworkServiceName() string {
	if m != nil {
		return m.NetworkServiceName
	}
	return ""
}

func (m *NetworkServiceChannel) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *NetworkServiceChannel) GetSocketLocation() string {
	if m != nil {
		return m.SocketLocation
	}
	return ""
}

func (m *NetworkServiceChannel) GetInterface() []*common.Interface {
	if m != nil {
		return m.Interface
	}
	return nil
}

type NetworkServiceEndpoint struct {
	Metadata             *common.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *NetworkServiceEndpoint) Reset()         { *m = NetworkServiceEndpoint{} }
func (m *NetworkServiceEndpoint) String() string { return proto.CompactTextString(m) }
func (*NetworkServiceEndpoint) ProtoMessage()    {}
func (*NetworkServiceEndpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_netmesh_8944ebf4e75e6af3, []int{1}
}
func (m *NetworkServiceEndpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkServiceEndpoint.Unmarshal(m, b)
}
func (m *NetworkServiceEndpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkServiceEndpoint.Marshal(b, m, deterministic)
}
func (dst *NetworkServiceEndpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkServiceEndpoint.Merge(dst, src)
}
func (m *NetworkServiceEndpoint) XXX_Size() int {
	return xxx_messageInfo_NetworkServiceEndpoint.Size(m)
}
func (m *NetworkServiceEndpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkServiceEndpoint.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkServiceEndpoint proto.InternalMessageInfo

func (m *NetworkServiceEndpoint) GetMetadata() *common.Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type NetworkService struct {
	Metadata             *common.Metadata         `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Channel              []*NetworkServiceChannel `protobuf:"bytes,2,rep,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *NetworkService) Reset()         { *m = NetworkService{} }
func (m *NetworkService) String() string { return proto.CompactTextString(m) }
func (*NetworkService) ProtoMessage()    {}
func (*NetworkService) Descriptor() ([]byte, []int) {
	return fileDescriptor_netmesh_8944ebf4e75e6af3, []int{2}
}
func (m *NetworkService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkService.Unmarshal(m, b)
}
func (m *NetworkService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkService.Marshal(b, m, deterministic)
}
func (dst *NetworkService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkService.Merge(dst, src)
}
func (m *NetworkService) XXX_Size() int {
	return xxx_messageInfo_NetworkService.Size(m)
}
func (m *NetworkService) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkService.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkService proto.InternalMessageInfo

func (m *NetworkService) GetMetadata() *common.Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *NetworkService) GetChannel() []*NetworkServiceChannel {
	if m != nil {
		return m.Channel
	}
	return nil
}

func init() {
	proto.RegisterType((*NetworkServiceChannel)(nil), "netmesh.NetworkServiceChannel")
	proto.RegisterType((*NetworkServiceEndpoint)(nil), "netmesh.NetworkServiceEndpoint")
	proto.RegisterType((*NetworkService)(nil), "netmesh.NetworkService")
}

func init() { proto.RegisterFile("netmesh.proto", fileDescriptor_netmesh_8944ebf4e75e6af3) }

var fileDescriptor_netmesh_8944ebf4e75e6af3 = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0x49, 0xab, 0xd6, 0x6e, 0xb1, 0xea, 0xa2, 0xb2, 0x78, 0x90, 0xd2, 0x8b, 0x3d, 0x48,
	0x56, 0xea, 0xc5, 0xbb, 0x54, 0x10, 0xb4, 0x87, 0xf8, 0x00, 0x65, 0xba, 0x1d, 0xdb, 0x25, 0xd9,
	0x99, 0x90, 0xac, 0xff, 0x9e, 0xd9, 0x97, 0x10, 0x93, 0x8d, 0x12, 0xf0, 0xd2, 0xd3, 0x32, 0xb3,
	0xbf, 0x99, 0x6f, 0x66, 0x3e, 0x71, 0x40, 0xe8, 0x1d, 0x96, 0x9b, 0x38, 0x2f, 0xd8, 0xb3, 0xec,
	0x85, 0xf0, 0x7c, 0xb6, 0xb6, 0x7e, 0xf3, 0xba, 0x8c, 0x0d, 0x3b, 0x9d, 0xd9, 0x35, 0x78, 0xd6,
	0x84, 0xfe, 0x9d, 0x8b, 0xb4, 0xc4, 0xe2, 0xcd, 0x1a, 0xfc, 0xa1, 0x74, 0x9e, 0xae, 0x35, 0x95,
	0x4e, 0x43, 0x6e, 0x4b, 0x6d, 0xd8, 0x39, 0xa6, 0xf0, 0xd4, 0xfd, 0xc6, 0x5f, 0x91, 0x38, 0x9d,
	0xd7, 0x75, 0xcf, 0x75, 0xdd, 0xdd, 0x06, 0x88, 0x30, 0x93, 0x57, 0x62, 0xdf, 0xa1, 0x87, 0x15,
	0x78, 0x50, 0xd1, 0x28, 0x9a, 0x0c, 0xa6, 0x47, 0x71, 0x28, 0x7d, 0x0a, 0xf9, 0xe4, 0x97, 0x90,
	0xd7, 0xe2, 0x24, 0xc8, 0x2f, 0x82, 0xfe, 0x82, 0xc0, 0xa1, 0xea, 0x8c, 0xa2, 0x49, 0x3f, 0x91,
	0xd4, 0x92, 0x98, 0x83, 0x43, 0xa9, 0x44, 0x2f, 0x87, 0xcf, 0x8c, 0x61, 0xa5, 0xba, 0x15, 0xd4,
	0x84, 0xf2, 0x52, 0x1c, 0x96, 0x6c, 0x52, 0xf4, 0x8b, 0x8c, 0x0d, 0x78, 0xcb, 0xa4, 0x76, 0x2a,
	0x62, 0x58, 0xa7, 0x1f, 0x43, 0x56, 0x6a, 0xd1, 0xb7, 0xe4, 0xb1, 0x78, 0x01, 0x83, 0x6a, 0x77,
	0xd4, 0x9d, 0x0c, 0xa6, 0xc7, 0xcd, 0x8c, 0x0f, 0xcd, 0x47, 0xf2, 0xc7, 0x8c, 0xef, 0xc5, 0x59,
	0x7b, 0xd9, 0x19, 0xad, 0x72, 0xb6, 0xe4, 0xb7, 0xdb, 0x76, 0xfc, 0x21, 0x86, 0xed, 0x3e, 0x5b,
	0x5e, 0xeb, 0x56, 0xf4, 0x4c, 0x7d, 0x66, 0xd5, 0xa9, 0xc6, 0xbe, 0x88, 0x1b, 0x9b, 0xff, 0x35,
	0x23, 0x69, 0xf0, 0xe5, 0x5e, 0x65, 0xdb, 0xcd, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x5b,
	0xe8, 0x82, 0x17, 0x02, 0x00, 0x00,
}
