package vppagent

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ligato/vpp-agent/api/configurator"
	"github.com/ligato/vpp-agent/api/models/vpp"
	acl "github.com/ligato/vpp-agent/api/models/vpp/acl"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/networkservice"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
	"github.com/sirupsen/logrus"
)

const (
	action     = "action"     // DENY, PERMIT, REFLECT
	dstNet     = "dstnet"     // IPv4 or IPv6 CIDR
	srcNet     = "srcnet"     // IPv4 or IPv6 CIDR
	icmpType   = "icmptype"   // 8-bit unsigned integer
	tcpLowPort = "tcplowport" // 16-bit unsigned integer
	tcpUpPort  = "tcpupport"  // 16-bit unsigned integer
	udpLowPort = "udplowport" // 16-bit unsigned integer
	udpUpPort  = "udpupport"  // 16-bit unsigned integer
)

// ACL is a VPP Agent ACL composite
type ACL struct {
	endpoint.BaseCompositeEndpoint
	Rules       map[string]string
	Connections map[string]*ConnectionData
}

// Request implements the request handler
func (a *ACL) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	if a.GetNext() == nil {
		err := fmt.Errorf("composite requires that there is Next set")
		return nil, err
	}

	incomingConnection, err := a.GetNext().Request(ctx, request)
	if err != nil {
		return nil, err
	}

	connectionData, err := getConnectionData(a.GetNext(), incomingConnection, false)
	if err != nil {
		return nil, err
	}

	if connectionData.InConnName == "" {
		err = fmt.Errorf("found empty incoming connection name")
		return nil, err
	}

	connectionData.DataChange, err = a.appendDataChange(connectionData.DataChange, connectionData.InConnName)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	a.Connections[incomingConnection.GetId()] = connectionData

	return incomingConnection, nil
}

// Close implements the close handler
func (a *ACL) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	if a.GetNext() != nil {
		return a.GetNext().Close(ctx, connection)
	}
	return &empty.Empty{}, nil
}

// GetOpaque will return the corresponding connection data
func (a *ACL) GetOpaque(incoming interface{}) interface{} {
	incomingConnection := incoming.(*connection.Connection)
	if connectionData, ok := a.Connections[incomingConnection.GetId()]; ok {
		return connectionData
	}
	logrus.Errorf("GetOpaque outgoing not found for %v", incomingConnection)
	return nil
}

// Name returns the composite name
func (a *ACL) Name() string {
	return "acl"
}

// NewACL creates an ACL
func NewACL(configuration *common.NSConfiguration, rules map[string]string) *ACL {
	// ensure the env variables are processed
	if configuration == nil {
		configuration = &common.NSConfiguration{}
	}
	configuration.CompleteNSConfiguration()

	return &ACL{
		Rules:       rules,
		Connections: map[string]*ConnectionData{},
	}
}

func (a *ACL) appendDataChange(rv *configurator.Config, ingressInterface string) (*configurator.Config, error) {
	if rv == nil {
		rv = &configurator.Config{}
	}
	if rv.VppConfig == nil {
		rv.VppConfig = &vpp.ConfigData{}
	}

	rules := []*acl.ACL_Rule{}

	for _, rule := range a.Rules {
		parsed := tools.ParseKVStringToMap(rule, ",", "=")

		action, err := getAction(parsed)
		if err != nil {
			logrus.Errorf("Parsing rule %s failed with %v", rule, err)
			return nil, err
		}

		match, err := getMatch(parsed)
		if err != nil {
			logrus.Errorf("Parsing rule %s failed with %v", rule, err)
			return nil, err
		}
		match.Action = action

		rules = append(rules, match)

		rv.VppConfig.Acls = append(rv.VppConfig.Acls, &acl.ACL{
			Name:  "IngressACL",
			Rules: rules,
			Interfaces: &acl.ACL_Interfaces{
				Egress:  []string{},
				Ingress: []string{ingressInterface},
			},
		})
	}

	return rv, nil
}

func getAction(parsed map[string]string) (acl.ACL_Rule_Action, error) {
	actionName, ok := parsed[action]
	if !ok {
		return acl.ACL_Rule_Action(0), fmt.Errorf("rule should have 'action' set")
	}
	action, ok := acl.ACL_Rule_Action_value[strings.ToUpper(actionName)]
	if !ok {
		return acl.ACL_Rule_Action(0), fmt.Errorf("rule should have a valid 'action'")
	}
	return acl.ACL_Rule_Action(action), nil
}

func getIP(parsed map[string]string) (*acl.ACL_Rule_IpRule_Ip, error) {
	dstNet, dstNetOk := parsed[dstNet]
	srcNet, srcNetOk := parsed[srcNet]
	if dstNetOk {
		_, _, err := net.ParseCIDR(dstNet)
		if err != nil {
			return nil, fmt.Errorf("dstnet is not a valid CIDR [%v]. Failed with: %v", dstNet, err)
		}
	} else {
		dstNet = ""
	}

	if srcNetOk {
		_, _, err := net.ParseCIDR(srcNet)
		if err != nil {
			return nil, fmt.Errorf("srcnet is not a valid CIDR [%v]. Failed with: %v", srcNet, err)
		}
	} else {
		srcNet = ""
	}

	if dstNetOk || srcNetOk {
		return &acl.ACL_Rule_IpRule_Ip{
			DestinationNetwork: dstNet,
			SourceNetwork:      srcNet,
		}, nil
	}
	return nil, nil
}

func getICMP(parsed map[string]string) (*acl.ACL_Rule_IpRule_Icmp, error) {
	icmpType, ok := parsed[icmpType]
	if !ok {
		return nil, nil
	}
	icmpType8, err := strconv.ParseUint(icmpType, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("failed parsing icmptype [%v] with: %v", icmpType, err)
	}
	return &acl.ACL_Rule_IpRule_Icmp{
		Icmpv6: false,
		IcmpCodeRange: &acl.ACL_Rule_IpRule_Icmp_Range{
			First: uint32(0),
			Last:  uint32(65535),
		},
		IcmpTypeRange: &acl.ACL_Rule_IpRule_Icmp_Range{
			First: uint32(icmpType8),
			Last:  uint32(icmpType8),
		},
	}, nil
}

func getPort(name string, parsed map[string]string) (uint16, bool, error) {
	port, ok := parsed[name]
	if !ok {
		return 0, false, nil
	}
	port16, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, true, fmt.Errorf("failed parsing %s [%v] with: %v", name, port, err)
	}

	return uint16(port16), true, nil
}

func getTCP(parsed map[string]string) (*acl.ACL_Rule_IpRule_Tcp, error) {
	lowerPort, lpFound, lpErr := getPort(tcpLowPort, parsed)
	if !lpFound {
		return nil, nil
	} else if lpErr != nil {
		return nil, lpErr
	}

	upperPort, upFound, upErr := getPort(tcpUpPort, parsed)
	if !upFound {
		return nil, nil
	} else if upErr != nil {
		return nil, lpErr
	}

	return &acl.ACL_Rule_IpRule_Tcp{
		DestinationPortRange: &acl.ACL_Rule_IpRule_PortRange{
			LowerPort: uint32(lowerPort),
			UpperPort: uint32(upperPort),
		},
		SourcePortRange: &acl.ACL_Rule_IpRule_PortRange{
			LowerPort: uint32(0),
			UpperPort: uint32(65535),
		},
		TcpFlagsMask:  0,
		TcpFlagsValue: 0,
	}, nil
}

func getUDP(parsed map[string]string) (*acl.ACL_Rule_IpRule_Udp, error) {
	lowerPort, lpFound, lpErr := getPort(udpLowPort, parsed)
	if !lpFound {
		return nil, nil
	} else if lpErr != nil {
		return nil, lpErr
	}

	upperPort, upFound, upErr := getPort(udpUpPort, parsed)
	if !upFound {
		return nil, nil
	} else if upErr != nil {
		return nil, lpErr
	}

	return &acl.ACL_Rule_IpRule_Udp{
		DestinationPortRange: &acl.ACL_Rule_IpRule_PortRange{
			LowerPort: uint32(lowerPort),
			UpperPort: uint32(upperPort),
		},
		SourcePortRange: &acl.ACL_Rule_IpRule_PortRange{
			LowerPort: uint32(0),
			UpperPort: uint32(65535),
		},
	}, nil
}

func getIPRule(parsed map[string]string) (*acl.ACL_Rule_IpRule, error) {
	ip, err := getIP(parsed)
	if err != nil {
		return nil, err
	}

	icmp, err := getICMP(parsed)
	if err != nil {
		return nil, err
	}

	tcp, err := getTCP(parsed)
	if err != nil {
		return nil, err
	}

	udp, err := getUDP(parsed)
	if err != nil {
		return nil, err
	}

	return &acl.ACL_Rule_IpRule{
		Ip:   ip,
		Icmp: icmp,
		Tcp:  tcp,
		Udp:  udp,
	}, nil
}

func getMatch(parsed map[string]string) (*acl.ACL_Rule, error) {
	ipRule, err := getIPRule(parsed)
	if err != nil {
		return nil, err
	}

	return &acl.ACL_Rule{
		IpRule:    ipRule,
		MacipRule: nil,
	}, nil
}
