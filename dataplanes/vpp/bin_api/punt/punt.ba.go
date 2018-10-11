// Code generated by GoVPP binapi-generator. DO NOT EDIT.
// source: api/punt.api.json

/*
Package punt is a generated VPP binary API of the 'punt' VPP module.

It is generated from this file:
	punt.api.json

It contains these VPP binary API objects:
	6 messages
	3 services
*/
package punt

import "git.fd.io/govpp.git/api"
import "github.com/lunixbochs/struc"
import "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = api.RegisterMessage
var _ = struc.Pack
var _ = bytes.NewBuffer

/* Messages */

// Punt represents the VPP binary API message 'punt'.
// Generated from 'punt.api.json', line 4:
//
//            "punt",
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
//                "ipv"
//            ],
//            [
//                "u8",
//                "l4_protocol"
//            ],
//            [
//                "u16",
//                "l4_port"
//            ],
//            {
//                "crc": "0x37760008"
//            }
//
type Punt struct {
	IsAdd      uint8
	IPv        uint8
	L4Protocol uint8
	L4Port     uint16
}

func (*Punt) GetMessageName() string {
	return "punt"
}
func (*Punt) GetCrcString() string {
	return "37760008"
}
func (*Punt) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewPunt() api.Message {
	return &Punt{}
}

// PuntReply represents the VPP binary API message 'punt_reply'.
// Generated from 'punt.api.json', line 38:
//
//            "punt_reply",
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
type PuntReply struct {
	Retval int32
}

func (*PuntReply) GetMessageName() string {
	return "punt_reply"
}
func (*PuntReply) GetCrcString() string {
	return "e8d4e804"
}
func (*PuntReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewPuntReply() api.Message {
	return &PuntReply{}
}

// PuntSocketRegister represents the VPP binary API message 'punt_socket_register'.
// Generated from 'punt.api.json', line 56:
//
//            "punt_socket_register",
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
//                "header_version"
//            ],
//            [
//                "u8",
//                "is_ip4"
//            ],
//            [
//                "u8",
//                "l4_protocol"
//            ],
//            [
//                "u16",
//                "l4_port"
//            ],
//            [
//                "u8",
//                "pathname",
//                108
//            ],
//            {
//                "crc": "0xc163b363"
//            }
//
type PuntSocketRegister struct {
	HeaderVersion uint32
	IsIP4         uint8
	L4Protocol    uint8
	L4Port        uint16
	Pathname      []byte `struc:"[108]byte"`
}

func (*PuntSocketRegister) GetMessageName() string {
	return "punt_socket_register"
}
func (*PuntSocketRegister) GetCrcString() string {
	return "c163b363"
}
func (*PuntSocketRegister) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewPuntSocketRegister() api.Message {
	return &PuntSocketRegister{}
}

// PuntSocketRegisterReply represents the VPP binary API message 'punt_socket_register_reply'.
// Generated from 'punt.api.json', line 95:
//
//            "punt_socket_register_reply",
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
//                "u8",
//                "pathname",
//                64
//            ],
//            {
//                "crc": "0x42dc0ee6"
//            }
//
type PuntSocketRegisterReply struct {
	Retval   int32
	Pathname []byte `struc:"[64]byte"`
}

func (*PuntSocketRegisterReply) GetMessageName() string {
	return "punt_socket_register_reply"
}
func (*PuntSocketRegisterReply) GetCrcString() string {
	return "42dc0ee6"
}
func (*PuntSocketRegisterReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewPuntSocketRegisterReply() api.Message {
	return &PuntSocketRegisterReply{}
}

// PuntSocketDeregister represents the VPP binary API message 'punt_socket_deregister'.
// Generated from 'punt.api.json', line 118:
//
//            "punt_socket_deregister",
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
//                "is_ip4"
//            ],
//            [
//                "u8",
//                "l4_protocol"
//            ],
//            [
//                "u16",
//                "l4_port"
//            ],
//            {
//                "crc": "0x9846a4cc"
//            }
//
type PuntSocketDeregister struct {
	IsIP4      uint8
	L4Protocol uint8
	L4Port     uint16
}

func (*PuntSocketDeregister) GetMessageName() string {
	return "punt_socket_deregister"
}
func (*PuntSocketDeregister) GetCrcString() string {
	return "9846a4cc"
}
func (*PuntSocketDeregister) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewPuntSocketDeregister() api.Message {
	return &PuntSocketDeregister{}
}

// PuntSocketDeregisterReply represents the VPP binary API message 'punt_socket_deregister_reply'.
// Generated from 'punt.api.json', line 148:
//
//            "punt_socket_deregister_reply",
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
type PuntSocketDeregisterReply struct {
	Retval int32
}

func (*PuntSocketDeregisterReply) GetMessageName() string {
	return "punt_socket_deregister_reply"
}
func (*PuntSocketDeregisterReply) GetCrcString() string {
	return "e8d4e804"
}
func (*PuntSocketDeregisterReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewPuntSocketDeregisterReply() api.Message {
	return &PuntSocketDeregisterReply{}
}

/* Services */

type Services interface {
	Punt(*Punt) (*PuntReply, error)
	PuntSocketDeregister(*PuntSocketDeregister) (*PuntSocketDeregisterReply, error)
	PuntSocketRegister(*PuntSocketRegister) (*PuntSocketRegisterReply, error)
}

func init() {
	api.RegisterMessage((*Punt)(nil), "punt.Punt")
	api.RegisterMessage((*PuntReply)(nil), "punt.PuntReply")
	api.RegisterMessage((*PuntSocketRegister)(nil), "punt.PuntSocketRegister")
	api.RegisterMessage((*PuntSocketRegisterReply)(nil), "punt.PuntSocketRegisterReply")
	api.RegisterMessage((*PuntSocketDeregister)(nil), "punt.PuntSocketDeregister")
	api.RegisterMessage((*PuntSocketDeregisterReply)(nil), "punt.PuntSocketDeregisterReply")
}
