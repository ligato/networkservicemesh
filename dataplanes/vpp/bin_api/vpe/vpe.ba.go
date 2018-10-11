// Code generated by GoVPP binapi-generator. DO NOT EDIT.
// source: api/vpe.api.json

/*
Package vpe is a generated VPP binary API of the 'vpe' VPP module.

It is generated from this file:
	vpe.api.json

It contains these VPP binary API objects:
	16 messages
	8 services
*/
package vpe

import "git.fd.io/govpp.git/api"
import "github.com/lunixbochs/struc"
import "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = api.RegisterMessage
var _ = struc.Pack
var _ = bytes.NewBuffer

/* Messages */

// ControlPing represents the VPP binary API message 'control_ping'.
// Generated from 'vpe.api.json', line 4:
//
//            "control_ping",
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
//            {
//                "crc": "0x51077d14"
//            }
//
type ControlPing struct{}

func (*ControlPing) GetMessageName() string {
	return "control_ping"
}
func (*ControlPing) GetCrcString() string {
	return "51077d14"
}
func (*ControlPing) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewControlPing() api.Message {
	return &ControlPing{}
}

// ControlPingReply represents the VPP binary API message 'control_ping_reply'.
// Generated from 'vpe.api.json', line 22:
//
//            "control_ping_reply",
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
//                "client_index"
//            ],
//            [
//                "u32",
//                "vpe_pid"
//            ],
//            {
//                "crc": "0xf6b0b8ca"
//            }
//
type ControlPingReply struct {
	Retval      int32
	ClientIndex uint32
	VpePID      uint32
}

func (*ControlPingReply) GetMessageName() string {
	return "control_ping_reply"
}
func (*ControlPingReply) GetCrcString() string {
	return "f6b0b8ca"
}
func (*ControlPingReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewControlPingReply() api.Message {
	return &ControlPingReply{}
}

// Cli represents the VPP binary API message 'cli'.
// Generated from 'vpe.api.json', line 48:
//
//            "cli",
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
//                "u64",
//                "cmd_in_shmem"
//            ],
//            {
//                "crc": "0x23bfbfff"
//            }
//
type Cli struct {
	CmdInShmem uint64
}

func (*Cli) GetMessageName() string {
	return "cli"
}
func (*Cli) GetCrcString() string {
	return "23bfbfff"
}
func (*Cli) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewCli() api.Message {
	return &Cli{}
}

// CliInband represents the VPP binary API message 'cli_inband'.
// Generated from 'vpe.api.json', line 70:
//
//            "cli_inband",
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
//                "length"
//            ],
//            [
//                "u8",
//                "cmd",
//                0,
//                "length"
//            ],
//            {
//                "crc": "0x74e00a49"
//            }
//
type CliInband struct {
	Length uint32 `struc:"sizeof=Cmd"`
	Cmd    []byte
}

func (*CliInband) GetMessageName() string {
	return "cli_inband"
}
func (*CliInband) GetCrcString() string {
	return "74e00a49"
}
func (*CliInband) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewCliInband() api.Message {
	return &CliInband{}
}

// CliReply represents the VPP binary API message 'cli_reply'.
// Generated from 'vpe.api.json', line 98:
//
//            "cli_reply",
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
//                "u64",
//                "reply_in_shmem"
//            ],
//            {
//                "crc": "0x06d68297"
//            }
//
type CliReply struct {
	Retval       int32
	ReplyInShmem uint64
}

func (*CliReply) GetMessageName() string {
	return "cli_reply"
}
func (*CliReply) GetCrcString() string {
	return "06d68297"
}
func (*CliReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewCliReply() api.Message {
	return &CliReply{}
}

// CliInbandReply represents the VPP binary API message 'cli_inband_reply'.
// Generated from 'vpe.api.json', line 120:
//
//            "cli_inband_reply",
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
//                "length"
//            ],
//            [
//                "u8",
//                "reply",
//                0,
//                "length"
//            ],
//            {
//                "crc": "0x1f22bbb8"
//            }
//
type CliInbandReply struct {
	Retval int32
	Length uint32 `struc:"sizeof=Reply"`
	Reply  []byte
}

func (*CliInbandReply) GetMessageName() string {
	return "cli_inband_reply"
}
func (*CliInbandReply) GetCrcString() string {
	return "1f22bbb8"
}
func (*CliInbandReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewCliInbandReply() api.Message {
	return &CliInbandReply{}
}

// GetNodeIndex represents the VPP binary API message 'get_node_index'.
// Generated from 'vpe.api.json', line 148:
//
//            "get_node_index",
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
//                "node_name",
//                64
//            ],
//            {
//                "crc": "0x6c9a495d"
//            }
//
type GetNodeIndex struct {
	NodeName []byte `struc:"[64]byte"`
}

func (*GetNodeIndex) GetMessageName() string {
	return "get_node_index"
}
func (*GetNodeIndex) GetCrcString() string {
	return "6c9a495d"
}
func (*GetNodeIndex) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewGetNodeIndex() api.Message {
	return &GetNodeIndex{}
}

// GetNodeIndexReply represents the VPP binary API message 'get_node_index_reply'.
// Generated from 'vpe.api.json', line 171:
//
//            "get_node_index_reply",
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
//                "node_index"
//            ],
//            {
//                "crc": "0xa8600b89"
//            }
//
type GetNodeIndexReply struct {
	Retval    int32
	NodeIndex uint32
}

func (*GetNodeIndexReply) GetMessageName() string {
	return "get_node_index_reply"
}
func (*GetNodeIndexReply) GetCrcString() string {
	return "a8600b89"
}
func (*GetNodeIndexReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewGetNodeIndexReply() api.Message {
	return &GetNodeIndexReply{}
}

// AddNodeNext represents the VPP binary API message 'add_node_next'.
// Generated from 'vpe.api.json', line 193:
//
//            "add_node_next",
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
//                "node_name",
//                64
//            ],
//            [
//                "u8",
//                "next_name",
//                64
//            ],
//            {
//                "crc": "0x9ab92f7a"
//            }
//
type AddNodeNext struct {
	NodeName []byte `struc:"[64]byte"`
	NextName []byte `struc:"[64]byte"`
}

func (*AddNodeNext) GetMessageName() string {
	return "add_node_next"
}
func (*AddNodeNext) GetCrcString() string {
	return "9ab92f7a"
}
func (*AddNodeNext) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewAddNodeNext() api.Message {
	return &AddNodeNext{}
}

// AddNodeNextReply represents the VPP binary API message 'add_node_next_reply'.
// Generated from 'vpe.api.json', line 221:
//
//            "add_node_next_reply",
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
//                "next_index"
//            ],
//            {
//                "crc": "0x2ed75f32"
//            }
//
type AddNodeNextReply struct {
	Retval    int32
	NextIndex uint32
}

func (*AddNodeNextReply) GetMessageName() string {
	return "add_node_next_reply"
}
func (*AddNodeNextReply) GetCrcString() string {
	return "2ed75f32"
}
func (*AddNodeNextReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewAddNodeNextReply() api.Message {
	return &AddNodeNextReply{}
}

// ShowVersion represents the VPP binary API message 'show_version'.
// Generated from 'vpe.api.json', line 243:
//
//            "show_version",
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
//            {
//                "crc": "0x51077d14"
//            }
//
type ShowVersion struct{}

func (*ShowVersion) GetMessageName() string {
	return "show_version"
}
func (*ShowVersion) GetCrcString() string {
	return "51077d14"
}
func (*ShowVersion) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewShowVersion() api.Message {
	return &ShowVersion{}
}

// ShowVersionReply represents the VPP binary API message 'show_version_reply'.
// Generated from 'vpe.api.json', line 261:
//
//            "show_version_reply",
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
//                "program",
//                32
//            ],
//            [
//                "u8",
//                "version",
//                32
//            ],
//            [
//                "u8",
//                "build_date",
//                32
//            ],
//            [
//                "u8",
//                "build_directory",
//                256
//            ],
//            {
//                "crc": "0x8b5a13b4"
//            }
//
type ShowVersionReply struct {
	Retval         int32
	Program        []byte `struc:"[32]byte"`
	Version        []byte `struc:"[32]byte"`
	BuildDate      []byte `struc:"[32]byte"`
	BuildDirectory []byte `struc:"[256]byte"`
}

func (*ShowVersionReply) GetMessageName() string {
	return "show_version_reply"
}
func (*ShowVersionReply) GetCrcString() string {
	return "8b5a13b4"
}
func (*ShowVersionReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewShowVersionReply() api.Message {
	return &ShowVersionReply{}
}

// GetNodeGraph represents the VPP binary API message 'get_node_graph'.
// Generated from 'vpe.api.json', line 299:
//
//            "get_node_graph",
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
//            {
//                "crc": "0x51077d14"
//            }
//
type GetNodeGraph struct{}

func (*GetNodeGraph) GetMessageName() string {
	return "get_node_graph"
}
func (*GetNodeGraph) GetCrcString() string {
	return "51077d14"
}
func (*GetNodeGraph) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewGetNodeGraph() api.Message {
	return &GetNodeGraph{}
}

// GetNodeGraphReply represents the VPP binary API message 'get_node_graph_reply'.
// Generated from 'vpe.api.json', line 317:
//
//            "get_node_graph_reply",
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
//                "u64",
//                "reply_in_shmem"
//            ],
//            {
//                "crc": "0x06d68297"
//            }
//
type GetNodeGraphReply struct {
	Retval       int32
	ReplyInShmem uint64
}

func (*GetNodeGraphReply) GetMessageName() string {
	return "get_node_graph_reply"
}
func (*GetNodeGraphReply) GetCrcString() string {
	return "06d68297"
}
func (*GetNodeGraphReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewGetNodeGraphReply() api.Message {
	return &GetNodeGraphReply{}
}

// GetNextIndex represents the VPP binary API message 'get_next_index'.
// Generated from 'vpe.api.json', line 339:
//
//            "get_next_index",
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
//                "node_name",
//                64
//            ],
//            [
//                "u8",
//                "next_name",
//                64
//            ],
//            {
//                "crc": "0x9ab92f7a"
//            }
//
type GetNextIndex struct {
	NodeName []byte `struc:"[64]byte"`
	NextName []byte `struc:"[64]byte"`
}

func (*GetNextIndex) GetMessageName() string {
	return "get_next_index"
}
func (*GetNextIndex) GetCrcString() string {
	return "9ab92f7a"
}
func (*GetNextIndex) GetMessageType() api.MessageType {
	return api.RequestMessage
}
func NewGetNextIndex() api.Message {
	return &GetNextIndex{}
}

// GetNextIndexReply represents the VPP binary API message 'get_next_index_reply'.
// Generated from 'vpe.api.json', line 367:
//
//            "get_next_index_reply",
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
//                "next_index"
//            ],
//            {
//                "crc": "0x2ed75f32"
//            }
//
type GetNextIndexReply struct {
	Retval    int32
	NextIndex uint32
}

func (*GetNextIndexReply) GetMessageName() string {
	return "get_next_index_reply"
}
func (*GetNextIndexReply) GetCrcString() string {
	return "2ed75f32"
}
func (*GetNextIndexReply) GetMessageType() api.MessageType {
	return api.ReplyMessage
}
func NewGetNextIndexReply() api.Message {
	return &GetNextIndexReply{}
}

/* Services */

type Services interface {
	AddNodeNext(*AddNodeNext) (*AddNodeNextReply, error)
	Cli(*Cli) (*CliReply, error)
	CliInband(*CliInband) (*CliInbandReply, error)
	ControlPing(*ControlPing) (*ControlPingReply, error)
	GetNextIndex(*GetNextIndex) (*GetNextIndexReply, error)
	GetNodeGraph(*GetNodeGraph) (*GetNodeGraphReply, error)
	GetNodeIndex(*GetNodeIndex) (*GetNodeIndexReply, error)
	ShowVersion(*ShowVersion) (*ShowVersionReply, error)
}

func init() {
	api.RegisterMessage((*ControlPing)(nil), "vpe.ControlPing")
	api.RegisterMessage((*ControlPingReply)(nil), "vpe.ControlPingReply")
	api.RegisterMessage((*Cli)(nil), "vpe.Cli")
	api.RegisterMessage((*CliInband)(nil), "vpe.CliInband")
	api.RegisterMessage((*CliReply)(nil), "vpe.CliReply")
	api.RegisterMessage((*CliInbandReply)(nil), "vpe.CliInbandReply")
	api.RegisterMessage((*GetNodeIndex)(nil), "vpe.GetNodeIndex")
	api.RegisterMessage((*GetNodeIndexReply)(nil), "vpe.GetNodeIndexReply")
	api.RegisterMessage((*AddNodeNext)(nil), "vpe.AddNodeNext")
	api.RegisterMessage((*AddNodeNextReply)(nil), "vpe.AddNodeNextReply")
	api.RegisterMessage((*ShowVersion)(nil), "vpe.ShowVersion")
	api.RegisterMessage((*ShowVersionReply)(nil), "vpe.ShowVersionReply")
	api.RegisterMessage((*GetNodeGraph)(nil), "vpe.GetNodeGraph")
	api.RegisterMessage((*GetNodeGraphReply)(nil), "vpe.GetNodeGraphReply")
	api.RegisterMessage((*GetNextIndex)(nil), "vpe.GetNextIndex")
	api.RegisterMessage((*GetNextIndexReply)(nil), "vpe.GetNextIndexReply")
}
