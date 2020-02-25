package converter

import (
	"strings"

	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/common"
	"github.com/networkservicemesh/api/pkg/api/networkservice/mechanisms/kernel"

	"github.com/ligato/vpp-agent/api/configurator"
	"github.com/ligato/vpp-agent/api/models/linux"
	linux_interfaces "github.com/ligato/vpp-agent/api/models/linux/interfaces"
	linux_l3 "github.com/ligato/vpp-agent/api/models/linux/l3"
	linux_namespace "github.com/ligato/vpp-agent/api/models/linux/namespace"
	"github.com/ligato/vpp-agent/api/models/vpp"
	vpp_interfaces "github.com/ligato/vpp-agent/api/models/vpp/interfaces"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/api/pkg/api/networkservice"
)

type KernelConnectionConverter struct {
	*networkservice.Connection
	conversionParameters *ConnectionConversionParameters
}

// NewKernelConnectionConverter - creates a kernel connection converter.
func NewKernelConnectionConverter(c *networkservice.Connection, conversionParameters *ConnectionConversionParameters) *KernelConnectionConverter {
	return &KernelConnectionConverter{
		Connection:           c,
		conversionParameters: conversionParameters,
	}
}

func (c *KernelConnectionConverter) ToDataRequest(rv *configurator.Config, connect bool) (*configurator.Config, error) {
	if c == nil {
		return rv, errors.New("LocalConnectionConverter cannot be nil")
	}
	if err := c.IsComplete(); err != nil {
		return rv, err
	}
	if c.GetMechanism().GetType() != kernel.MECHANISM {
		return rv, errors.Errorf("KernelConnectionConverter cannot be used on Connection.Mechanism.Type %s", c.GetMechanism().GetType())
	}
	if rv == nil {
		rv = &configurator.Config{
			LinuxConfig: &linux.ConfigData{},
		}
	}
	if rv.VppConfig == nil {
		rv.VppConfig = &vpp.ConfigData{}
	}
	if rv.LinuxConfig == nil {
		rv.LinuxConfig = &linux.ConfigData{}
	}

	m := kernel.ToMechanism(c.GetMechanism())
	filepath, err := netNsFileName(c.GetMechanism())
	if err != nil && connect {
		return nil, err
	}
	var ipAddresses []string
	var mac string
	if c.conversionParameters.Side == DESTINATION {
		ipAddresses = []string{c.Connection.GetContext().GetIpContext().GetDstIpAddr()}
		if !c.GetContext().IsEthernetContextEmtpy() {
			mac = c.GetContext().EthernetContext.DstMac
		}
	}
	if c.conversionParameters.Side == SOURCE {
		ipAddresses = []string{c.Connection.GetContext().GetIpContext().GetSrcIpAddr()}
		if !c.GetContext().IsEthernetContextEmtpy() {
			mac = c.GetContext().EthernetContext.SrcMac
		}
	}

	logrus.Infof("m.GetParameters()[%s]: %s", common.InterfaceNameKey, m.GetParameters()[common.InterfaceNameKey])

	// If we have access to /dev/vhost-net, we can use tapv2.  Otherwise fall back to
	// veth pairs
	if useVHostNet() {
		// We append an Interfaces.  Interfaces creates the vpp side of an interface.
		//   In this case, a Tapv2 interface that has one side in vpp, and the other
		//   as a Linux kernel interface
		rv.VppConfig.Interfaces = append(rv.VppConfig.Interfaces, &vpp_interfaces.Interface{
			Name:    c.conversionParameters.Name,
			Type:    vpp_interfaces.Interface_TAP,
			Enabled: true,
			Link: &vpp_interfaces.Interface_Tap{
				Tap: &vpp_interfaces.TapLink{
					Version: 2,
				},
			},
		})
		logrus.Info("Found /dev/vhost-net - using tapv2")
		// We apply configuration to LinuxInterfaces
		// Important details:
		//    - LinuxInterfaces.HostIfName - must be no longer than 15 chars (linux limitation)
		rv.LinuxConfig.Interfaces = append(rv.LinuxConfig.Interfaces, &linux.Interface{
			Name:        c.conversionParameters.Name,
			Type:        linux_interfaces.Interface_TAP_TO_VPP,
			Enabled:     true,
			IpAddresses: ipAddresses,
			PhysAddress: mac,
			HostIfName:  m.GetParameters()[common.InterfaceNameKey],
			Namespace: &linux_namespace.NetNamespace{
				Type:      linux_namespace.NetNamespace_FD,
				Reference: filepath,
			},
			Link: &linux_interfaces.Interface_Tap{
				Tap: &linux_interfaces.TapLink{
					VppTapIfName: c.conversionParameters.Name,
				},
			},
		})
	} else {
		logrus.Info("Did Not Find /dev/vhost-net - using veth pairs")
		rv.LinuxConfig.Interfaces = append(rv.LinuxConfig.Interfaces, &linux_interfaces.Interface{
			Name:       c.conversionParameters.Name + "-veth",
			Type:       linux_interfaces.Interface_VETH,
			Enabled:    true,
			HostIfName: c.conversionParameters.Name + "-veth",
			Link: &linux_interfaces.Interface_Veth{
				Veth: &linux_interfaces.VethLink{
					PeerIfName: c.conversionParameters.Name,
				},
			},
		})
		rv.LinuxConfig.Interfaces = append(rv.LinuxConfig.Interfaces, &linux_interfaces.Interface{
			Name:        c.conversionParameters.Name,
			Type:        linux_interfaces.Interface_VETH,
			Enabled:     true,
			IpAddresses: ipAddresses,
			PhysAddress: mac,
			HostIfName:  m.GetParameters()[common.InterfaceNameKey],
			Namespace: &linux_namespace.NetNamespace{
				Type:      linux_namespace.NetNamespace_FD,
				Reference: filepath,
			},
			Link: &linux_interfaces.Interface_Veth{
				Veth: &linux_interfaces.VethLink{
					PeerIfName: c.conversionParameters.Name + "-veth",
				},
			},
		})
		rv.VppConfig.Interfaces = append(rv.VppConfig.Interfaces, &vpp_interfaces.Interface{
			Name:    c.conversionParameters.Name,
			Type:    vpp_interfaces.Interface_AF_PACKET,
			Enabled: true,
			Link: &vpp_interfaces.Interface_Afpacket{
				Afpacket: &vpp_interfaces.AfpacketLink{
					HostIfName: c.conversionParameters.Name + "-veth",
				},
			},
		})
	}

	// Process static routes
	var routes []*networkservice.Route
	switch c.conversionParameters.Side {
	case SOURCE:
		routes = c.Connection.GetContext().GetIpContext().GetDstRoutes()
	case DESTINATION:
		routes = c.Connection.GetContext().GetIpContext().GetSrcRoutes()
	}

	duplicatedPrefixes := make(map[string]bool)
	for _, route := range routes {
		if _, ok := duplicatedPrefixes[route.Prefix]; !ok {
			duplicatedPrefixes[route.Prefix] = true
			rv.LinuxConfig.Routes = append(rv.LinuxConfig.Routes, &linux.Route{
				DstNetwork:        route.Prefix,
				OutgoingInterface: c.conversionParameters.Name,
				Scope:             linux_l3.Route_GLOBAL,
				GwAddr:            extractCleanIPAddress(c.Connection.GetContext().GetIpContext().GetDstIpAddr()),
			})
		}
	}

	// Process IP Neighbor entries
	if c.conversionParameters.Side == SOURCE {
		for _, neightbour := range c.Connection.GetContext().GetIpContext().GetIpNeighbors() {
			rv.LinuxConfig.ArpEntries = append(rv.LinuxConfig.ArpEntries, &linux.ARPEntry{
				IpAddress: neightbour.Ip,
				Interface: c.conversionParameters.Name,
				HwAddress: neightbour.HardwareAddress,
			})
		}
		if c.GetContext().EthernetContext != nil && c.GetContext().EthernetContext.DstMac != "" {
			logrus.Infof("set arp for: %v", c.GetContext().String())
			rv.LinuxConfig.ArpEntries = append(rv.LinuxConfig.ArpEntries, &linux.ARPEntry{
				IpAddress: strings.Split(c.GetContext().IpContext.DstIpAddr, "/")[0],
				Interface: c.conversionParameters.Name,
				HwAddress: c.GetContext().EthernetContext.DstMac,
			})
		}
	}
	return rv, nil
}
