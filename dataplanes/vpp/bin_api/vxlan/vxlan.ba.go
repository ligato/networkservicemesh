// Code generated by GoVPP binapi-generator. DO NOT EDIT.
// source: api/vxlan.api.json

/*
Package vxlan is a generated VPP binary API of the 'vxlan' VPP module.

It is generated from this file:
	vxlan.api.json

It contains these VPP binary API objects:
	8 messages
	4 services
*/
package vxlan

import "git.fd.io/govpp.git/api"
import "github.com/lunixbochs/struc"
import "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = api.RegisterMessage
var _ = struc.Pack
var _ = bytes.NewBuffer

/* Messages */

// VxlanAddDelTunnel represents the VPP binary API message 'vxlan_add_del_tunnel'.
// Generated from 'vxlan.api.json', line 4:
//
//            "vxlan_add_del_tunnel",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "client_index"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "u8",
//                "is_add"
//            ],
//            [
//                "u8",
//                "is_ipv6"
//            ],
//            [
//                "u32",
//                "instance"
//            ],
//            [
//                "u8",
//                "src_address",
//                16
//            ],
//            [
//                "u8",
//                "dst_address",
//                16
//            ],
//            [
//                "u32",
//                "mcast_sw_if_index"
//            ],
//            [
//                "u32",
//                "encap_vrf_id"
//            ],
//            [
//                "u32",
//                "decap_next_index"
//            ],
//            [
//                "u32",
//                "vni"
//            ],
//            {
//                "crc": "0x00f4bdd0"
//            }
//
type VxlanAddDelTunnel struct {
	IsAdd          uint8
	IsIPv6         uint8
	Instance       uint32
	SrcAddress     []byte `struc:"[16]byte"`
	DstAddress     []byte `struc:"[16]byte"`
	McastSwIfIndex uint32
	EncapVrfID     uint32
	DecapNextIndex uint32
	Vni            uint32
}

func (*VxlanAddDelTunnel) GetMessageName() string {
	return "vxlan_add_del_tunnel"
}
func (*VxlanAddDelTunnel) GetCrcString() string {
	return "00f4bdd0"
}
func (*VxlanAddDelTunnel) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewVxlanAddDelTunnel() api.Message {
	return &VxlanAddDelTunnel{}
}

// VxlanAddDelTunnelReply represents the VPP binary API message 'vxlan_add_del_tunnel_reply'.
// Generated from 'vxlan.api.json', line 60:
//
//            "vxlan_add_del_tunnel_reply",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "i32",
//                "retval"
//            ],
//            [
//                "u32",
//                "sw_if_index"
//            ],
//            {
//                "crc": "0xfda5941f"
//            }
//
type VxlanAddDelTunnelReply struct {
	Retval    int32
	SwIfIndex uint32
}

func (*VxlanAddDelTunnelReply) GetMessageName() string {
	return "vxlan_add_del_tunnel_reply"
}
func (*VxlanAddDelTunnelReply) GetCrcString() string {
	return "fda5941f"
}
func (*VxlanAddDelTunnelReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewVxlanAddDelTunnelReply() api.Message {
	return &VxlanAddDelTunnelReply{}
}

// VxlanTunnelDump represents the VPP binary API message 'vxlan_tunnel_dump'.
// Generated from 'vxlan.api.json', line 82:
//
//            "vxlan_tunnel_dump",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "client_index"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "u32",
//                "sw_if_index"
//            ],
//            {
//                "crc": "0x529cb13f"
//            }
//
type VxlanTunnelDump struct {
	SwIfIndex uint32
}

func (*VxlanTunnelDump) GetMessageName() string {
	return "vxlan_tunnel_dump"
}
func (*VxlanTunnelDump) GetCrcString() string {
	return "529cb13f"
}
func (*VxlanTunnelDump) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewVxlanTunnelDump() api.Message {
	return &VxlanTunnelDump{}
}

// VxlanTunnelDetails represents the VPP binary API message 'vxlan_tunnel_details'.
// Generated from 'vxlan.api.json', line 104:
//
//            "vxlan_tunnel_details",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "u32",
//                "sw_if_index"
//            ],
//            [
//                "u32",
//                "instance"
//            ],
//            [
//                "u8",
//                "src_address",
//                16
//            ],
//            [
//                "u8",
//                "dst_address",
//                16
//            ],
//            [
//                "u32",
//                "mcast_sw_if_index"
//            ],
//            [
//                "u32",
//                "encap_vrf_id"
//            ],
//            [
//                "u32",
//                "decap_next_index"
//            ],
//            [
//                "u32",
//                "vni"
//            ],
//            [
//                "u8",
//                "is_ipv6"
//            ],
//            {
//                "crc": "0xce38e127"
//            }
//
type VxlanTunnelDetails struct {
	SwIfIndex      uint32
	Instance       uint32
	SrcAddress     []byte `struc:"[16]byte"`
	DstAddress     []byte `struc:"[16]byte"`
	McastSwIfIndex uint32
	EncapVrfID     uint32
	DecapNextIndex uint32
	Vni            uint32
	IsIPv6         uint8
}

func (*VxlanTunnelDetails) GetMessageName() string {
	return "vxlan_tunnel_details"
}
func (*VxlanTunnelDetails) GetCrcString() string {
	return "ce38e127"
}
func (*VxlanTunnelDetails) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewVxlanTunnelDetails() api.Message {
	return &VxlanTunnelDetails{}
}

// SwInterfaceSetVxlanBypass represents the VPP binary API message 'sw_interface_set_vxlan_bypass'.
// Generated from 'vxlan.api.json', line 156:
//
//            "sw_interface_set_vxlan_bypass",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "client_index"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "u32",
//                "sw_if_index"
//            ],
//            [
//                "u8",
//                "is_ipv6"
//            ],
//            [
//                "u8",
//                "enable"
//            ],
//            {
//                "crc": "0xe74ca095"
//            }
//
type SwInterfaceSetVxlanBypass struct {
	SwIfIndex uint32
	IsIPv6    uint8
	Enable    uint8
}

func (*SwInterfaceSetVxlanBypass) GetMessageName() string {
	return "sw_interface_set_vxlan_bypass"
}
func (*SwInterfaceSetVxlanBypass) GetCrcString() string {
	return "e74ca095"
}
func (*SwInterfaceSetVxlanBypass) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewSwInterfaceSetVxlanBypass() api.Message {
	return &SwInterfaceSetVxlanBypass{}
}

// SwInterfaceSetVxlanBypassReply represents the VPP binary API message 'sw_interface_set_vxlan_bypass_reply'.
// Generated from 'vxlan.api.json', line 186:
//
//            "sw_interface_set_vxlan_bypass_reply",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "i32",
//                "retval"
//            ],
//            {
//                "crc": "0xe8d4e804"
//            }
//
type SwInterfaceSetVxlanBypassReply struct {
	Retval int32
}

func (*SwInterfaceSetVxlanBypassReply) GetMessageName() string {
	return "sw_interface_set_vxlan_bypass_reply"
}
func (*SwInterfaceSetVxlanBypassReply) GetCrcString() string {
	return "e8d4e804"
}
func (*SwInterfaceSetVxlanBypassReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewSwInterfaceSetVxlanBypassReply() api.Message {
	return &SwInterfaceSetVxlanBypassReply{}
}

// VxlanOffloadRx represents the VPP binary API message 'vxlan_offload_rx'.
// Generated from 'vxlan.api.json', line 204:
//
//            "vxlan_offload_rx",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "client_index"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "u32",
//                "hw_if_index"
//            ],
//            [
//                "u32",
//                "sw_if_index"
//            ],
//            [
//                "u8",
//                "enable"
//            ],
//            {
//                "crc": "0xf0b08786"
//            }
//
type VxlanOffloadRx struct {
	HwIfIndex uint32
	SwIfIndex uint32
	Enable    uint8
}

func (*VxlanOffloadRx) GetMessageName() string {
	return "vxlan_offload_rx"
}
func (*VxlanOffloadRx) GetCrcString() string {
	return "f0b08786"
}
func (*VxlanOffloadRx) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewVxlanOffloadRx() api.Message {
	return &VxlanOffloadRx{}
}

// VxlanOffloadRxReply represents the VPP binary API message 'vxlan_offload_rx_reply'.
// Generated from 'vxlan.api.json', line 234:
//
//            "vxlan_offload_rx_reply",
//            [
//                "u16",
//                "_vl_msg_id"
//            ],
//            [
//                "u32",
//                "context"
//            ],
//            [
//                "i32",
//                "retval"
//            ],
//            {
//                "crc": "0xe8d4e804"
//            }
//
type VxlanOffloadRxReply struct {
	Retval int32
}

func (*VxlanOffloadRxReply) GetMessageName() string {
	return "vxlan_offload_rx_reply"
}
func (*VxlanOffloadRxReply) GetCrcString() string {
	return "e8d4e804"
}
func (*VxlanOffloadRxReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewVxlanOffloadRxReply() api.Message {
	return &VxlanOffloadRxReply{}
}

/* Services */

type Services interface {
	DumpVxlanTunnel(*VxlanTunnelDump) (*VxlanTunnelDetails, error)
	SwInterfaceSetVxlanBypass(*SwInterfaceSetVxlanBypass) (*SwInterfaceSetVxlanBypassReply, error)
	VxlanAddDelTunnel(*VxlanAddDelTunnel) (*VxlanAddDelTunnelReply, error)
	VxlanOffloadRx(*VxlanOffloadRx) (*VxlanOffloadRxReply, error)
}

func init() {
	api.RegisterMessage((*VxlanAddDelTunnel)(nil), "vxlan.VxlanAddDelTunnel")
	api.RegisterMessage((*VxlanAddDelTunnelReply)(nil), "vxlan.VxlanAddDelTunnelReply")
	api.RegisterMessage((*VxlanTunnelDump)(nil), "vxlan.VxlanTunnelDump")
	api.RegisterMessage((*VxlanTunnelDetails)(nil), "vxlan.VxlanTunnelDetails")
	api.RegisterMessage((*SwInterfaceSetVxlanBypass)(nil), "vxlan.SwInterfaceSetVxlanBypass")
	api.RegisterMessage((*SwInterfaceSetVxlanBypassReply)(nil), "vxlan.SwInterfaceSetVxlanBypassReply")
	api.RegisterMessage((*VxlanOffloadRx)(nil), "vxlan.VxlanOffloadRx")
	api.RegisterMessage((*VxlanOffloadRxReply)(nil), "vxlan.VxlanOffloadRxReply")
}
