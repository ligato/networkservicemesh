// Code generated by protoc-gen-go. DO NOT EDIT.
// source: registry.proto

package registry

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

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

type NetworkServiceEndpoint struct {
	NetworkServiceName        string            `protobuf:"bytes,1,opt,name=network_service_name,json=networkServiceName,proto3" json:"network_service_name,omitempty"`
	Payload                   string            `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	NetworkServiceManagerName string            `protobuf:"bytes,3,opt,name=network_service_manager_name,json=networkServiceManagerName,proto3" json:"network_service_manager_name,omitempty"`
	EndpointName              string            `protobuf:"bytes,4,opt,name=endpoint_name,json=endpointName,proto3" json:"endpoint_name,omitempty"`
	Labels                    map[string]string `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	State                     string            `protobuf:"bytes,6,opt,name=state,proto3" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}          `json:"-"`
	XXX_unrecognized          []byte            `json:"-"`
	XXX_sizecache             int32             `json:"-"`
}

func (m *NetworkServiceEndpoint) Reset()         { *m = NetworkServiceEndpoint{} }
func (m *NetworkServiceEndpoint) String() string { return proto.CompactTextString(m) }
func (*NetworkServiceEndpoint) ProtoMessage()    {}
func (*NetworkServiceEndpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{0}
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

func (m *NetworkServiceEndpoint) GetNetworkServiceName() string {
	if m != nil {
		return m.NetworkServiceName
	}
	return ""
}

func (m *NetworkServiceEndpoint) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *NetworkServiceEndpoint) GetNetworkServiceManagerName() string {
	if m != nil {
		return m.NetworkServiceManagerName
	}
	return ""
}

func (m *NetworkServiceEndpoint) GetEndpointName() string {
	if m != nil {
		return m.EndpointName
	}
	return ""
}

func (m *NetworkServiceEndpoint) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *NetworkServiceEndpoint) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

type NetworkService struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Payload              string   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Matches              []*Match `protobuf:"bytes,3,rep,name=matches,proto3" json:"matches,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkService) Reset()         { *m = NetworkService{} }
func (m *NetworkService) String() string { return proto.CompactTextString(m) }
func (*NetworkService) ProtoMessage()    {}
func (*NetworkService) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{1}
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

func (m *NetworkService) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NetworkService) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *NetworkService) GetMatches() []*Match {
	if m != nil {
		return m.Matches
	}
	return nil
}

type Match struct {
	SourceSelector       map[string]string `protobuf:"bytes,1,rep,name=source_selector,json=sourceSelector,proto3" json:"source_selector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Routes               []*Destination    `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Match) Reset()         { *m = Match{} }
func (m *Match) String() string { return proto.CompactTextString(m) }
func (*Match) ProtoMessage()    {}
func (*Match) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{2}
}
func (m *Match) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Match.Unmarshal(m, b)
}
func (m *Match) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Match.Marshal(b, m, deterministic)
}
func (dst *Match) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Match.Merge(dst, src)
}
func (m *Match) XXX_Size() int {
	return xxx_messageInfo_Match.Size(m)
}
func (m *Match) XXX_DiscardUnknown() {
	xxx_messageInfo_Match.DiscardUnknown(m)
}

var xxx_messageInfo_Match proto.InternalMessageInfo

func (m *Match) GetSourceSelector() map[string]string {
	if m != nil {
		return m.SourceSelector
	}
	return nil
}

func (m *Match) GetRoutes() []*Destination {
	if m != nil {
		return m.Routes
	}
	return nil
}

type Destination struct {
	DestinationSelector  map[string]string `protobuf:"bytes,1,rep,name=destination_selector,json=destinationSelector,proto3" json:"destination_selector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Weight               uint32            `protobuf:"varint,2,opt,name=weight,proto3" json:"weight,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Destination) Reset()         { *m = Destination{} }
func (m *Destination) String() string { return proto.CompactTextString(m) }
func (*Destination) ProtoMessage()    {}
func (*Destination) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{3}
}
func (m *Destination) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Destination.Unmarshal(m, b)
}
func (m *Destination) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Destination.Marshal(b, m, deterministic)
}
func (dst *Destination) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Destination.Merge(dst, src)
}
func (m *Destination) XXX_Size() int {
	return xxx_messageInfo_Destination.Size(m)
}
func (m *Destination) XXX_DiscardUnknown() {
	xxx_messageInfo_Destination.DiscardUnknown(m)
}

var xxx_messageInfo_Destination proto.InternalMessageInfo

func (m *Destination) GetDestinationSelector() map[string]string {
	if m != nil {
		return m.DestinationSelector
	}
	return nil
}

func (m *Destination) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

type NetworkServiceManager struct {
	Name                 string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url                  string               `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	LastSeen             *timestamp.Timestamp `protobuf:"bytes,3,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
	State                string               `protobuf:"bytes,4,opt,name=state,proto3" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *NetworkServiceManager) Reset()         { *m = NetworkServiceManager{} }
func (m *NetworkServiceManager) String() string { return proto.CompactTextString(m) }
func (*NetworkServiceManager) ProtoMessage()    {}
func (*NetworkServiceManager) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{4}
}
func (m *NetworkServiceManager) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkServiceManager.Unmarshal(m, b)
}
func (m *NetworkServiceManager) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkServiceManager.Marshal(b, m, deterministic)
}
func (dst *NetworkServiceManager) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkServiceManager.Merge(dst, src)
}
func (m *NetworkServiceManager) XXX_Size() int {
	return xxx_messageInfo_NetworkServiceManager.Size(m)
}
func (m *NetworkServiceManager) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkServiceManager.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkServiceManager proto.InternalMessageInfo

func (m *NetworkServiceManager) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NetworkServiceManager) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *NetworkServiceManager) GetLastSeen() *timestamp.Timestamp {
	if m != nil {
		return m.LastSeen
	}
	return nil
}

func (m *NetworkServiceManager) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

type RemoveNSERequest struct {
	EndpointName         string   `protobuf:"bytes,1,opt,name=endpoint_name,json=endpointName,proto3" json:"endpoint_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveNSERequest) Reset()         { *m = RemoveNSERequest{} }
func (m *RemoveNSERequest) String() string { return proto.CompactTextString(m) }
func (*RemoveNSERequest) ProtoMessage()    {}
func (*RemoveNSERequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{5}
}
func (m *RemoveNSERequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveNSERequest.Unmarshal(m, b)
}
func (m *RemoveNSERequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveNSERequest.Marshal(b, m, deterministic)
}
func (dst *RemoveNSERequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveNSERequest.Merge(dst, src)
}
func (m *RemoveNSERequest) XXX_Size() int {
	return xxx_messageInfo_RemoveNSERequest.Size(m)
}
func (m *RemoveNSERequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveNSERequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveNSERequest proto.InternalMessageInfo

func (m *RemoveNSERequest) GetEndpointName() string {
	if m != nil {
		return m.EndpointName
	}
	return ""
}

type FindNetworkServiceRequest struct {
	NetworkServiceName   string   `protobuf:"bytes,1,opt,name=network_service_name,json=networkServiceName,proto3" json:"network_service_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindNetworkServiceRequest) Reset()         { *m = FindNetworkServiceRequest{} }
func (m *FindNetworkServiceRequest) String() string { return proto.CompactTextString(m) }
func (*FindNetworkServiceRequest) ProtoMessage()    {}
func (*FindNetworkServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{6}
}
func (m *FindNetworkServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindNetworkServiceRequest.Unmarshal(m, b)
}
func (m *FindNetworkServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindNetworkServiceRequest.Marshal(b, m, deterministic)
}
func (dst *FindNetworkServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindNetworkServiceRequest.Merge(dst, src)
}
func (m *FindNetworkServiceRequest) XXX_Size() int {
	return xxx_messageInfo_FindNetworkServiceRequest.Size(m)
}
func (m *FindNetworkServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindNetworkServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindNetworkServiceRequest proto.InternalMessageInfo

func (m *FindNetworkServiceRequest) GetNetworkServiceName() string {
	if m != nil {
		return m.NetworkServiceName
	}
	return ""
}

type FindNetworkServiceResponse struct {
	Payload                 string                            `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	NetworkService          *NetworkService                   `protobuf:"bytes,2,opt,name=network_service,json=networkService,proto3" json:"network_service,omitempty"`
	NetworkServiceManagers  map[string]*NetworkServiceManager `protobuf:"bytes,3,rep,name=network_service_managers,json=networkServiceManagers,proto3" json:"network_service_managers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NetworkServiceEndpoints []*NetworkServiceEndpoint         `protobuf:"bytes,4,rep,name=network_service_endpoints,json=networkServiceEndpoints,proto3" json:"network_service_endpoints,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}                          `json:"-"`
	XXX_unrecognized        []byte                            `json:"-"`
	XXX_sizecache           int32                             `json:"-"`
}

func (m *FindNetworkServiceResponse) Reset()         { *m = FindNetworkServiceResponse{} }
func (m *FindNetworkServiceResponse) String() string { return proto.CompactTextString(m) }
func (*FindNetworkServiceResponse) ProtoMessage()    {}
func (*FindNetworkServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{7}
}
func (m *FindNetworkServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindNetworkServiceResponse.Unmarshal(m, b)
}
func (m *FindNetworkServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindNetworkServiceResponse.Marshal(b, m, deterministic)
}
func (dst *FindNetworkServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindNetworkServiceResponse.Merge(dst, src)
}
func (m *FindNetworkServiceResponse) XXX_Size() int {
	return xxx_messageInfo_FindNetworkServiceResponse.Size(m)
}
func (m *FindNetworkServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FindNetworkServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FindNetworkServiceResponse proto.InternalMessageInfo

func (m *FindNetworkServiceResponse) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *FindNetworkServiceResponse) GetNetworkService() *NetworkService {
	if m != nil {
		return m.NetworkService
	}
	return nil
}

func (m *FindNetworkServiceResponse) GetNetworkServiceManagers() map[string]*NetworkServiceManager {
	if m != nil {
		return m.NetworkServiceManagers
	}
	return nil
}

func (m *FindNetworkServiceResponse) GetNetworkServiceEndpoints() []*NetworkServiceEndpoint {
	if m != nil {
		return m.NetworkServiceEndpoints
	}
	return nil
}

type NSERegistration struct {
	NetworkService         *NetworkService         `protobuf:"bytes,1,opt,name=network_service,json=networkService,proto3" json:"network_service,omitempty"`
	NetworkServiceManager  *NetworkServiceManager  `protobuf:"bytes,2,opt,name=network_service_manager,json=networkServiceManager,proto3" json:"network_service_manager,omitempty"`
	NetworkserviceEndpoint *NetworkServiceEndpoint `protobuf:"bytes,3,opt,name=networkservice_endpoint,json=networkserviceEndpoint,proto3" json:"networkservice_endpoint,omitempty"`
	XXX_NoUnkeyedLiteral   struct{}                `json:"-"`
	XXX_unrecognized       []byte                  `json:"-"`
	XXX_sizecache          int32                   `json:"-"`
}

func (m *NSERegistration) Reset()         { *m = NSERegistration{} }
func (m *NSERegistration) String() string { return proto.CompactTextString(m) }
func (*NSERegistration) ProtoMessage()    {}
func (*NSERegistration) Descriptor() ([]byte, []int) {
	return fileDescriptor_registry_2ea353b676a489cf, []int{8}
}
func (m *NSERegistration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NSERegistration.Unmarshal(m, b)
}
func (m *NSERegistration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NSERegistration.Marshal(b, m, deterministic)
}
func (dst *NSERegistration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NSERegistration.Merge(dst, src)
}
func (m *NSERegistration) XXX_Size() int {
	return xxx_messageInfo_NSERegistration.Size(m)
}
func (m *NSERegistration) XXX_DiscardUnknown() {
	xxx_messageInfo_NSERegistration.DiscardUnknown(m)
}

var xxx_messageInfo_NSERegistration proto.InternalMessageInfo

func (m *NSERegistration) GetNetworkService() *NetworkService {
	if m != nil {
		return m.NetworkService
	}
	return nil
}

func (m *NSERegistration) GetNetworkServiceManager() *NetworkServiceManager {
	if m != nil {
		return m.NetworkServiceManager
	}
	return nil
}

func (m *NSERegistration) GetNetworkserviceEndpoint() *NetworkServiceEndpoint {
	if m != nil {
		return m.NetworkserviceEndpoint
	}
	return nil
}

func init() {
	proto.RegisterType((*NetworkServiceEndpoint)(nil), "registry.NetworkServiceEndpoint")
	proto.RegisterMapType((map[string]string)(nil), "registry.NetworkServiceEndpoint.LabelsEntry")
	proto.RegisterType((*NetworkService)(nil), "registry.NetworkService")
	proto.RegisterType((*Match)(nil), "registry.Match")
	proto.RegisterMapType((map[string]string)(nil), "registry.Match.SourceSelectorEntry")
	proto.RegisterType((*Destination)(nil), "registry.Destination")
	proto.RegisterMapType((map[string]string)(nil), "registry.Destination.DestinationSelectorEntry")
	proto.RegisterType((*NetworkServiceManager)(nil), "registry.NetworkServiceManager")
	proto.RegisterType((*RemoveNSERequest)(nil), "registry.RemoveNSERequest")
	proto.RegisterType((*FindNetworkServiceRequest)(nil), "registry.FindNetworkServiceRequest")
	proto.RegisterType((*FindNetworkServiceResponse)(nil), "registry.FindNetworkServiceResponse")
	proto.RegisterMapType((map[string]*NetworkServiceManager)(nil), "registry.FindNetworkServiceResponse.NetworkServiceManagersEntry")
	proto.RegisterType((*NSERegistration)(nil), "registry.NSERegistration")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NetworkServiceRegistryClient is the client API for NetworkServiceRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkServiceRegistryClient interface {
	RegisterNSE(ctx context.Context, in *NSERegistration, opts ...grpc.CallOption) (*NSERegistration, error)
	RemoveNSE(ctx context.Context, in *RemoveNSERequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type networkServiceRegistryClient struct {
	cc *grpc.ClientConn
}

func NewNetworkServiceRegistryClient(cc *grpc.ClientConn) NetworkServiceRegistryClient {
	return &networkServiceRegistryClient{cc}
}

func (c *networkServiceRegistryClient) RegisterNSE(ctx context.Context, in *NSERegistration, opts ...grpc.CallOption) (*NSERegistration, error) {
	out := new(NSERegistration)
	err := c.cc.Invoke(ctx, "/registry.NetworkServiceRegistry/RegisterNSE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceRegistryClient) RemoveNSE(ctx context.Context, in *RemoveNSERequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/registry.NetworkServiceRegistry/RemoveNSE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServiceRegistryServer is the server API for NetworkServiceRegistry service.
type NetworkServiceRegistryServer interface {
	RegisterNSE(context.Context, *NSERegistration) (*NSERegistration, error)
	RemoveNSE(context.Context, *RemoveNSERequest) (*empty.Empty, error)
}

func RegisterNetworkServiceRegistryServer(s *grpc.Server, srv NetworkServiceRegistryServer) {
	s.RegisterService(&_NetworkServiceRegistry_serviceDesc, srv)
}

func _NetworkServiceRegistry_RegisterNSE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NSERegistration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceRegistryServer).RegisterNSE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.NetworkServiceRegistry/RegisterNSE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceRegistryServer).RegisterNSE(ctx, req.(*NSERegistration))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkServiceRegistry_RemoveNSE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveNSERequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceRegistryServer).RemoveNSE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.NetworkServiceRegistry/RemoveNSE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceRegistryServer).RemoveNSE(ctx, req.(*RemoveNSERequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NetworkServiceRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "registry.NetworkServiceRegistry",
	HandlerType: (*NetworkServiceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNSE",
			Handler:    _NetworkServiceRegistry_RegisterNSE_Handler,
		},
		{
			MethodName: "RemoveNSE",
			Handler:    _NetworkServiceRegistry_RemoveNSE_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registry.proto",
}

// NetworkServiceDiscoveryClient is the client API for NetworkServiceDiscovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkServiceDiscoveryClient interface {
	FindNetworkService(ctx context.Context, in *FindNetworkServiceRequest, opts ...grpc.CallOption) (*FindNetworkServiceResponse, error)
}

type networkServiceDiscoveryClient struct {
	cc *grpc.ClientConn
}

func NewNetworkServiceDiscoveryClient(cc *grpc.ClientConn) NetworkServiceDiscoveryClient {
	return &networkServiceDiscoveryClient{cc}
}

func (c *networkServiceDiscoveryClient) FindNetworkService(ctx context.Context, in *FindNetworkServiceRequest, opts ...grpc.CallOption) (*FindNetworkServiceResponse, error) {
	out := new(FindNetworkServiceResponse)
	err := c.cc.Invoke(ctx, "/registry.NetworkServiceDiscovery/FindNetworkService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServiceDiscoveryServer is the server API for NetworkServiceDiscovery service.
type NetworkServiceDiscoveryServer interface {
	FindNetworkService(context.Context, *FindNetworkServiceRequest) (*FindNetworkServiceResponse, error)
}

func RegisterNetworkServiceDiscoveryServer(s *grpc.Server, srv NetworkServiceDiscoveryServer) {
	s.RegisterService(&_NetworkServiceDiscovery_serviceDesc, srv)
}

func _NetworkServiceDiscovery_FindNetworkService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindNetworkServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceDiscoveryServer).FindNetworkService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.NetworkServiceDiscovery/FindNetworkService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceDiscoveryServer).FindNetworkService(ctx, req.(*FindNetworkServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NetworkServiceDiscovery_serviceDesc = grpc.ServiceDesc{
	ServiceName: "registry.NetworkServiceDiscovery",
	HandlerType: (*NetworkServiceDiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindNetworkService",
			Handler:    _NetworkServiceDiscovery_FindNetworkService_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registry.proto",
}

func init() { proto.RegisterFile("registry.proto", fileDescriptor_registry_2ea353b676a489cf) }

var fileDescriptor_registry_2ea353b676a489cf = []byte{
	// 770 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0xcf, 0x4e, 0xdb, 0x4e,
	0x10, 0x96, 0x93, 0x10, 0x60, 0xf2, 0x23, 0x41, 0x0b, 0x04, 0xc7, 0xbf, 0x4a, 0x45, 0xa1, 0x07,
	0x2a, 0xb5, 0x4e, 0x15, 0x54, 0xd1, 0x3f, 0x07, 0x8a, 0x4a, 0x38, 0x41, 0x0e, 0x4e, 0xa5, 0xaa,
	0x52, 0xa5, 0xc8, 0x49, 0xa6, 0xc6, 0xc5, 0xde, 0x75, 0xbd, 0x9b, 0x20, 0xf3, 0x04, 0x3d, 0xf4,
	0x19, 0xfa, 0x2a, 0x3d, 0xf6, 0xda, 0x57, 0xe8, 0x9b, 0x54, 0x59, 0xdb, 0xc4, 0x76, 0x1c, 0x20,
	0x97, 0x68, 0xff, 0xcc, 0x7c, 0x33, 0xf3, 0xed, 0x7c, 0x13, 0x43, 0xd5, 0x47, 0xcb, 0xe6, 0xc2,
	0x0f, 0x74, 0xcf, 0x67, 0x82, 0x91, 0xb5, 0x78, 0xaf, 0x1d, 0x5a, 0xb6, 0xb8, 0x1c, 0x0f, 0xf4,
	0x21, 0x73, 0x5b, 0x16, 0x73, 0x4c, 0x6a, 0xb5, 0xa4, 0xc9, 0x60, 0xfc, 0xa5, 0xe5, 0x89, 0xc0,
	0x43, 0xde, 0x42, 0xd7, 0x13, 0x41, 0xf8, 0x1b, 0xba, 0x6b, 0x6f, 0xef, 0x77, 0x12, 0xb6, 0x8b,
	0x5c, 0x98, 0xae, 0x37, 0x5b, 0x85, 0xce, 0xcd, 0xbf, 0x05, 0xa8, 0x77, 0x51, 0x5c, 0x33, 0xff,
	0xaa, 0x87, 0xfe, 0xc4, 0x1e, 0x62, 0x87, 0x8e, 0x3c, 0x66, 0x53, 0x41, 0x5e, 0xc0, 0x36, 0x0d,
	0x6f, 0xfa, 0x3c, 0xbc, 0xea, 0x53, 0xd3, 0x45, 0x55, 0xd9, 0x53, 0x0e, 0xd6, 0x0d, 0x42, 0x53,
	0x5e, 0x5d, 0xd3, 0x45, 0xa2, 0xc2, 0xaa, 0x67, 0x06, 0x0e, 0x33, 0x47, 0x6a, 0x41, 0x1a, 0xc5,
	0x5b, 0x72, 0x0c, 0x8f, 0xb2, 0x58, 0xae, 0x49, 0x4d, 0x0b, 0xfd, 0x10, 0xb3, 0x28, 0xcd, 0x1b,
	0x69, 0xcc, 0x8b, 0xd0, 0x42, 0x42, 0xef, 0xc3, 0x06, 0x46, 0x89, 0x85, 0x1e, 0x25, 0xe9, 0xf1,
	0x5f, 0x7c, 0x28, 0x8d, 0x4e, 0xa1, 0xec, 0x98, 0x03, 0x74, 0xb8, 0xba, 0xb2, 0x57, 0x3c, 0xa8,
	0xb4, 0x9f, 0xe9, 0xb7, 0x4c, 0xe7, 0xd7, 0xa8, 0x9f, 0x4b, 0xf3, 0x0e, 0x15, 0x7e, 0x60, 0x44,
	0xbe, 0x64, 0x1b, 0x56, 0xb8, 0x30, 0x05, 0xaa, 0x65, 0x19, 0x22, 0xdc, 0x68, 0xaf, 0xa1, 0x92,
	0x30, 0x26, 0x9b, 0x50, 0xbc, 0xc2, 0x20, 0xe2, 0x62, 0xba, 0x9c, 0xba, 0x4d, 0x4c, 0x67, 0x8c,
	0x51, 0xe9, 0xe1, 0xe6, 0x4d, 0xe1, 0x95, 0xd2, 0xb4, 0xa1, 0x9a, 0x0e, 0x4f, 0x08, 0x94, 0x12,
	0x54, 0xca, 0xf5, 0x1d, 0xe4, 0x3d, 0x85, 0x55, 0xd7, 0x14, 0xc3, 0x4b, 0xe4, 0x6a, 0x51, 0xd6,
	0x55, 0x9b, 0xd5, 0x75, 0x31, 0xbd, 0x30, 0xe2, 0xfb, 0xe6, 0x6f, 0x05, 0x56, 0xe4, 0x11, 0x39,
	0x87, 0x1a, 0x67, 0x63, 0x7f, 0x88, 0x7d, 0x8e, 0x0e, 0x0e, 0x05, 0xf3, 0x55, 0x45, 0x3a, 0xef,
	0x67, 0x9c, 0xf5, 0x9e, 0x34, 0xeb, 0x45, 0x56, 0x21, 0x17, 0x55, 0x9e, 0x3a, 0x24, 0xcf, 0xa1,
	0xec, 0xb3, 0xb1, 0x40, 0xae, 0x16, 0x24, 0xc8, 0xce, 0x0c, 0xe4, 0x14, 0xb9, 0xb0, 0xa9, 0x29,
	0x6c, 0x46, 0x8d, 0xc8, 0x48, 0x3b, 0x81, 0xad, 0x1c, 0xd4, 0xa5, 0x48, 0xfb, 0xa3, 0x40, 0x25,
	0x01, 0x4d, 0x4c, 0xd8, 0x1e, 0xcd, 0xb6, 0xd9, 0xa2, 0xf4, 0xdc, 0x7c, 0x92, 0xeb, 0x74, 0x7d,
	0x5b, 0xa3, 0xf9, 0x1b, 0x52, 0x87, 0xf2, 0x35, 0xda, 0xd6, 0xa5, 0x90, 0xd9, 0x6c, 0x18, 0xd1,
	0x4e, 0x3b, 0x03, 0x75, 0x11, 0xd0, 0x52, 0x25, 0xfd, 0x50, 0x60, 0xa7, 0x9b, 0xd7, 0xe1, 0xb9,
	0xfd, 0xb0, 0x09, 0xc5, 0xb1, 0xef, 0x44, 0x28, 0xd3, 0x25, 0x39, 0x82, 0x75, 0xc7, 0xe4, 0xa2,
	0xcf, 0x11, 0xa9, 0x54, 0x4c, 0xa5, 0xad, 0xe9, 0x16, 0x63, 0x96, 0x83, 0x7a, 0xac, 0x78, 0xfd,
	0x43, 0x2c, 0x70, 0x63, 0x6d, 0x6a, 0xdc, 0x43, 0xa4, 0xb3, 0x8e, 0x2e, 0x25, 0x3a, 0xba, 0x79,
	0x04, 0x9b, 0x06, 0xba, 0x6c, 0x82, 0xdd, 0x5e, 0xc7, 0xc0, 0x6f, 0x63, 0xe4, 0x62, 0x5e, 0x66,
	0xca, 0xbc, 0xcc, 0x9a, 0x17, 0xd0, 0x38, 0xb3, 0xe9, 0x28, 0x5d, 0x4a, 0x8c, 0xb0, 0xf4, 0xd4,
	0x68, 0xfe, 0x2a, 0x82, 0x96, 0x87, 0xc7, 0x3d, 0x46, 0x79, 0x4a, 0x17, 0x4a, 0x5a, 0x17, 0x27,
	0x50, 0xcb, 0x84, 0x92, 0x6c, 0x55, 0xda, 0xea, 0x22, 0xdd, 0x1b, 0xd5, 0x74, 0x7c, 0x72, 0x03,
	0xea, 0x82, 0xb9, 0x14, 0x6b, 0xed, 0xdd, 0x0c, 0x6b, 0x71, 0x92, 0x7a, 0xee, 0xb3, 0x46, 0x73,
	0xa5, 0x9e, 0x3b, 0xd5, 0x38, 0xf9, 0x0c, 0x8d, 0x6c, 0xec, 0x98, 0x66, 0xae, 0x96, 0x64, 0xf0,
	0xbd, 0xfb, 0x06, 0x98, 0xb1, 0x4b, 0x73, 0xcf, 0xb9, 0xf6, 0x15, 0xfe, 0xbf, 0x23, 0xa9, 0x9c,
	0xbe, 0x7d, 0x99, 0xec, 0xdb, 0x4a, 0xfb, 0xf1, 0xa2, 0xd0, 0x11, 0x4e, 0xb2, 0xb1, 0xbf, 0x17,
	0xa0, 0x26, 0x9b, 0x48, 0x3a, 0x84, 0x7a, 0xcd, 0x79, 0x1c, 0x65, 0xc9, 0xc7, 0xf9, 0x08, 0xbb,
	0x0b, 0x1e, 0xe7, 0xa1, 0x39, 0xee, 0xe4, 0x52, 0x4f, 0x3e, 0xdd, 0x02, 0x67, 0x89, 0x8f, 0x64,
	0x75, 0x3f, 0xef, 0xf5, 0x34, 0x40, 0x7c, 0xde, 0xfe, 0xa9, 0x64, 0xff, 0x4f, 0x23, 0x56, 0x02,
	0xf2, 0x1e, 0x2a, 0xe1, 0x1a, 0xfd, 0x6e, 0xaf, 0x43, 0x1a, 0x89, 0x18, 0x69, 0xee, 0xb4, 0xc5,
	0x57, 0xe4, 0x18, 0xd6, 0x6f, 0x45, 0x4b, 0xb4, 0x99, 0x5d, 0x56, 0xc9, 0x5a, 0x7d, 0x6e, 0x32,
	0x74, 0xa6, 0xdf, 0x0c, 0xed, 0x1b, 0xd8, 0x4d, 0xe7, 0x77, 0x6a, 0xf3, 0x21, 0x9b, 0xa0, 0x1f,
	0x90, 0x3e, 0x90, 0xf9, 0x16, 0x27, 0xfb, 0x77, 0x0b, 0x20, 0x8c, 0xf6, 0xe4, 0x21, 0x2a, 0x19,
	0x94, 0x65, 0x2e, 0x87, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa4, 0xae, 0x18, 0x5a, 0x01, 0x09,
	0x00, 0x00,
}
