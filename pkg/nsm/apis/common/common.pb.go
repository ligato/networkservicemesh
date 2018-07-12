// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type InterfaceType int32

const (
	InterfaceType_DEFAULT_INTF     InterfaceType = 0
	InterfaceType_KERNEL_INTERFACE InterfaceType = 1
	InterfaceType_VHOST_INTERFACE  InterfaceType = 2
	InterfaceType_MEM_INTERFACE    InterfaceType = 3
	InterfaceType_SRIOV_INTERFACE  InterfaceType = 4
	InterfaceType_HW_INTERFACE     InterfaceType = 5
)

var InterfaceType_name = map[int32]string{
	0: "DEFAULT_INTF",
	1: "KERNEL_INTERFACE",
	2: "VHOST_INTERFACE",
	3: "MEM_INTERFACE",
	4: "SRIOV_INTERFACE",
	5: "HW_INTERFACE",
}
var InterfaceType_value = map[string]int32{
	"DEFAULT_INTF":     0,
	"KERNEL_INTERFACE": 1,
	"VHOST_INTERFACE":  2,
	"MEM_INTERFACE":    3,
	"SRIOV_INTERFACE":  4,
	"HW_INTERFACE":     5,
}

func (x InterfaceType) String() string {
	return proto.EnumName(InterfaceType_name, int32(x))
}
func (InterfaceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_common_e2b5cb4e98f7bf11, []int{0}
}

type InterfacePreference int32

const (
	InterfacePreference_DEFAULT_PREF InterfacePreference = 0
	InterfacePreference_FIRST        InterfacePreference = 1
	InterfacePreference_SECOND       InterfacePreference = 2
	InterfacePreference_THIRD        InterfacePreference = 3
	InterfacePreference_FORTH        InterfacePreference = 4
	InterfacePreference_FIFTH        InterfacePreference = 5
)

var InterfacePreference_name = map[int32]string{
	0: "DEFAULT_PREF",
	1: "FIRST",
	2: "SECOND",
	3: "THIRD",
	4: "FORTH",
	5: "FIFTH",
}
var InterfacePreference_value = map[string]int32{
	"DEFAULT_PREF": 0,
	"FIRST":        1,
	"SECOND":       2,
	"THIRD":        3,
	"FORTH":        4,
	"FIFTH":        5,
}

func (x InterfacePreference) String() string {
	return proto.EnumName(InterfacePreference_name, int32(x))
}
func (InterfacePreference) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_common_e2b5cb4e98f7bf11, []int{1}
}

type Label struct {
	Selector             map[string]string `protobuf:"bytes,1,rep,name=selector,proto3" json:"selector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Label) Reset()         { *m = Label{} }
func (m *Label) String() string { return proto.CompactTextString(m) }
func (*Label) ProtoMessage()    {}
func (*Label) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_e2b5cb4e98f7bf11, []int{0}
}
func (m *Label) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Label.Unmarshal(m, b)
}
func (m *Label) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Label.Marshal(b, m, deterministic)
}
func (dst *Label) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Label.Merge(dst, src)
}
func (m *Label) XXX_Size() int {
	return xxx_messageInfo_Label.Size(m)
}
func (m *Label) XXX_DiscardUnknown() {
	xxx_messageInfo_Label.DiscardUnknown(m)
}

var xxx_messageInfo_Label proto.InternalMessageInfo

func (m *Label) GetSelector() map[string]string {
	if m != nil {
		return m.Selector
	}
	return nil
}

type Metadata struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Labels               *Label   `protobuf:"bytes,3,opt,name=labels,proto3" json:"labels,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_e2b5cb4e98f7bf11, []int{1}
}
func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metadata.Unmarshal(m, b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
}
func (dst *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(dst, src)
}
func (m *Metadata) XXX_Size() int {
	return xxx_messageInfo_Metadata.Size(m)
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Metadata) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Metadata) GetLabels() *Label {
	if m != nil {
		return m.Labels
	}
	return nil
}

type InterfaceParameters struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InterfaceParameters) Reset()         { *m = InterfaceParameters{} }
func (m *InterfaceParameters) String() string { return proto.CompactTextString(m) }
func (*InterfaceParameters) ProtoMessage()    {}
func (*InterfaceParameters) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_e2b5cb4e98f7bf11, []int{2}
}
func (m *InterfaceParameters) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InterfaceParameters.Unmarshal(m, b)
}
func (m *InterfaceParameters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InterfaceParameters.Marshal(b, m, deterministic)
}
func (dst *InterfaceParameters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterfaceParameters.Merge(dst, src)
}
func (m *InterfaceParameters) XXX_Size() int {
	return xxx_messageInfo_InterfaceParameters.Size(m)
}
func (m *InterfaceParameters) XXX_DiscardUnknown() {
	xxx_messageInfo_InterfaceParameters.DiscardUnknown(m)
}

var xxx_messageInfo_InterfaceParameters proto.InternalMessageInfo

type Interface struct {
	Type                 InterfaceType        `protobuf:"varint,1,opt,name=type,proto3,enum=common.InterfaceType" json:"type,omitempty"`
	Metadata             *Metadata            `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Preference           InterfacePreference  `protobuf:"varint,3,opt,name=preference,proto3,enum=common.InterfacePreference" json:"preference,omitempty"`
	Parmeters            *InterfaceParameters `protobuf:"bytes,4,opt,name=parmeters,proto3" json:"parmeters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Interface) Reset()         { *m = Interface{} }
func (m *Interface) String() string { return proto.CompactTextString(m) }
func (*Interface) ProtoMessage()    {}
func (*Interface) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_e2b5cb4e98f7bf11, []int{3}
}
func (m *Interface) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Interface.Unmarshal(m, b)
}
func (m *Interface) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Interface.Marshal(b, m, deterministic)
}
func (dst *Interface) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Interface.Merge(dst, src)
}
func (m *Interface) XXX_Size() int {
	return xxx_messageInfo_Interface.Size(m)
}
func (m *Interface) XXX_DiscardUnknown() {
	xxx_messageInfo_Interface.DiscardUnknown(m)
}

var xxx_messageInfo_Interface proto.InternalMessageInfo

func (m *Interface) GetType() InterfaceType {
	if m != nil {
		return m.Type
	}
	return InterfaceType_DEFAULT_INTF
}

func (m *Interface) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Interface) GetPreference() InterfacePreference {
	if m != nil {
		return m.Preference
	}
	return InterfacePreference_DEFAULT_PREF
}

func (m *Interface) GetParmeters() *InterfaceParameters {
	if m != nil {
		return m.Parmeters
	}
	return nil
}

func init() {
	proto.RegisterType((*Label)(nil), "common.Label")
	proto.RegisterMapType((map[string]string)(nil), "common.Label.SelectorEntry")
	proto.RegisterType((*Metadata)(nil), "common.Metadata")
	proto.RegisterType((*InterfaceParameters)(nil), "common.InterfaceParameters")
	proto.RegisterType((*Interface)(nil), "common.Interface")
	proto.RegisterEnum("common.InterfaceType", InterfaceType_name, InterfaceType_value)
	proto.RegisterEnum("common.InterfacePreference", InterfacePreference_name, InterfacePreference_value)
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_common_e2b5cb4e98f7bf11) }

var fileDescriptor_common_e2b5cb4e98f7bf11 = []byte{
	// 428 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x4b, 0x6f, 0xd3, 0x40,
	0x10, 0xc7, 0xd9, 0xd8, 0x8e, 0xe2, 0x49, 0x53, 0x96, 0x69, 0x2b, 0x59, 0xc0, 0x21, 0x8a, 0x84,
	0x14, 0x2a, 0x94, 0x43, 0x38, 0xf0, 0xe8, 0xa9, 0x6a, 0xd6, 0xb2, 0x45, 0x1e, 0xd5, 0xda, 0x94,
	0x63, 0xb5, 0x35, 0x9b, 0x0b, 0xf1, 0x43, 0x1b, 0x83, 0x94, 0x03, 0x57, 0x3e, 0x25, 0x1f, 0x06,
	0xed, 0xda, 0xb1, 0x1d, 0x89, 0x53, 0x76, 0xff, 0xf3, 0x9b, 0xcc, 0x6f, 0x56, 0x86, 0xb3, 0x24,
	0x4f, 0xd3, 0x3c, 0x9b, 0x15, 0x2a, 0x2f, 0x73, 0xec, 0x57, 0xb7, 0xc9, 0x6f, 0x70, 0x96, 0xe2,
	0x49, 0xee, 0xf0, 0x03, 0x0c, 0xf6, 0x72, 0x27, 0x93, 0x32, 0x57, 0x1e, 0x19, 0x5b, 0xd3, 0xe1,
	0xfc, 0xd5, 0xac, 0xee, 0x30, 0xc0, 0x2c, 0xaa, 0xab, 0x2c, 0x2b, 0xd5, 0x81, 0x37, 0xf0, 0xcb,
	0x1b, 0x18, 0x9d, 0x94, 0x90, 0x82, 0xf5, 0x43, 0x1e, 0x3c, 0x32, 0x26, 0x53, 0x97, 0xeb, 0x23,
	0x5e, 0x82, 0xf3, 0x4b, 0xec, 0x7e, 0x4a, 0xaf, 0x67, 0xb2, 0xea, 0xf2, 0xb9, 0xf7, 0x91, 0x4c,
	0x12, 0x18, 0xac, 0x64, 0x29, 0xbe, 0x8b, 0x52, 0x20, 0x82, 0x9d, 0x89, 0x54, 0xd6, 0x8d, 0xe6,
	0x8c, 0xaf, 0xc1, 0xd5, 0xbf, 0xfb, 0x42, 0x24, 0xc7, 0xee, 0x36, 0xc0, 0x37, 0xd0, 0xdf, 0x69,
	0xb7, 0xbd, 0x67, 0x8d, 0xc9, 0x74, 0x38, 0x1f, 0x9d, 0x18, 0xf3, 0xba, 0x38, 0xb9, 0x82, 0x8b,
	0x30, 0x2b, 0xa5, 0xda, 0x8a, 0x44, 0xde, 0x0b, 0x25, 0x52, 0x59, 0x4a, 0xb5, 0x9f, 0xfc, 0x25,
	0xe0, 0x36, 0x39, 0xbe, 0x05, 0xbb, 0x3c, 0x14, 0xd5, 0xf4, 0xf3, 0xf9, 0xd5, 0xf1, 0x9f, 0x1a,
	0x20, 0x3e, 0x14, 0x92, 0x1b, 0x04, 0xdf, 0xc1, 0x20, 0xad, 0xa5, 0x8d, 0xd3, 0x70, 0x4e, 0x8f,
	0xf8, 0x71, 0x19, 0xde, 0x10, 0x78, 0x03, 0x50, 0x28, 0xb9, 0x95, 0x4a, 0x66, 0x89, 0x34, 0xa2,
	0xe7, 0xed, 0xd3, 0xb6, 0x5e, 0x0d, 0xc2, 0x3b, 0x38, 0x7e, 0x02, 0xb7, 0x10, 0xaa, 0x12, 0xf6,
	0x6c, 0x33, 0xeb, 0x3f, 0xbd, 0xcd, 0x4e, 0xbc, 0xa5, 0xaf, 0xff, 0x10, 0x18, 0x9d, 0xd8, 0x23,
	0x85, 0xb3, 0x05, 0xf3, 0x6f, 0xbf, 0x2e, 0xe3, 0xc7, 0x70, 0x1d, 0xfb, 0xf4, 0x19, 0x5e, 0x02,
	0xfd, 0xc2, 0xf8, 0x9a, 0x2d, 0x75, 0xc0, 0xb8, 0x7f, 0x7b, 0xc7, 0x28, 0xc1, 0x0b, 0x78, 0xfe,
	0x10, 0x6c, 0xa2, 0xb8, 0x13, 0xf6, 0xf0, 0x05, 0x8c, 0x56, 0x6c, 0xd5, 0x89, 0x2c, 0xcd, 0x45,
	0x3c, 0xdc, 0x3c, 0x74, 0x42, 0x5b, 0x0f, 0x09, 0xbe, 0x75, 0x12, 0xe7, 0xfa, 0xb1, 0xfb, 0xfc,
	0xed, 0x6a, 0x1d, 0x9b, 0x7b, 0xce, 0xb4, 0x8d, 0x0b, 0x8e, 0x1f, 0xf2, 0x28, 0xa6, 0x04, 0x01,
	0xfa, 0x11, 0xbb, 0xdb, 0xac, 0x17, 0xb4, 0xa7, 0xe3, 0x38, 0x08, 0xf9, 0x82, 0x5a, 0x86, 0xd8,
	0xf0, 0x38, 0xa0, 0x76, 0x05, 0xfb, 0x71, 0x40, 0x9d, 0xa7, 0xbe, 0xf9, 0xa4, 0xdf, 0xff, 0x0b,
	0x00, 0x00, 0xff, 0xff, 0xc5, 0x88, 0x62, 0xe2, 0xe2, 0x02, 0x00, 0x00,
}
