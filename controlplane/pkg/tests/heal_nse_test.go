package tests

import (
	"testing"
	"time"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/nsm"

	. "github.com/onsi/gomega"
	"golang.org/x/net/context"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
)

func TestHealRemoteNSE(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	srv2 := NewNSMDFullServer(Worker, storage)
	defer srv.Stop()
	defer srv2.Stop()

	srv.TestModel.AddDataplane(context.Background(), testDataplane1)
	srv2.TestModel.AddDataplane(context.Background(), testDataplane2)

	// Register in both
	nseReg := srv2.registerFakeEndpointWithName("golden_network", "test", Worker, "ep1")
	nseReg2 := srv2.registerFakeEndpointWithName("golden_network", "test", Worker, "ep2")

	// Add to local endpoints for Server2
	srv2.TestModel.AddEndpoint(context.Background(), nseReg)
	srv2.TestModel.AddEndpoint(context.Background(), nseReg2)

	l1 := newTestConnectionModelListener()
	l2 := newTestConnectionModelListener()

	srv.TestModel.AddListener(l1)
	srv2.TestModel.AddListener(l2)

	// Now we could try to connect via Client API
	nsmClient, conn := srv.requestNSMConnection("nsm-1")
	defer conn.Close()

	request := &networkservice.NetworkServiceRequest{
		Connection: &connection.Connection{
			NetworkService: "golden_network",
			Context: &connectioncontext.ConnectionContext{
				IpContext: &connectioncontext.IPContext{
					DstIpRequired: true,
					SrcIpRequired: true,
				},
			},
			Labels: make(map[string]string),
		},
		MechanismPreferences: []*connection.Mechanism{
			{
				Type: connection.MechanismType_KERNEL_INTERFACE,
				Parameters: map[string]string{
					connection.NetNsInodeKey:    "10",
					connection.InterfaceNameKey: "icmp-responder1",
				},
			},
		},
	}

	timeout := time.Second * 10

	nsmResponse, err := nsmClient.Request(context.Background(), request)
	g.Expect(err).To(BeNil())
	g.Expect(nsmResponse.GetNetworkService()).To(Equal("golden_network"))

	// We need to check for cross connections.
	clientConnection1 := srv.TestModel.GetClientConnection(nsmResponse.GetId())
	g.Expect(clientConnection1.GetID()).To(Equal("1"))

	clientConnection2 := srv2.TestModel.GetClientConnection(clientConnection1.Xcon.GetRemoteDestination().GetId())
	g.Expect(clientConnection2.GetID()).To(Equal("1"))

	// We need to inform cross connection monitor about this connection, since dataplane is fake one.
	l1.WaitAdd(1, timeout, t)

	epName := clientConnection1.Endpoint.GetNetworkServiceEndpoint().GetName()
	_, err = srv.nseRegistry.RemoveNSE(context.Background(), &registry.RemoveNSERequest{
		NetworkServiceEndpointName: epName,
	})
	if err != nil {
		t.Fatal("Err must be nil")
	}

	srv2.TestModel.DeleteEndpoint(context.Background(), epName)

	// Simlate delete
	clientConnection2.Xcon.GetLocalDestination().State = connection.State_DOWN
	srv.manager.GetHealProperties().HealDSTNSEWaitTimeout = time.Second * 1
	srv2.manager.Heal(clientConnection2, nsm.HealStateDstDown)

	// First update, is delete
	// Second update is update
	l1.WaitUpdate(7, timeout, t)

	clientConnection1_1 := srv.TestModel.GetClientConnection(nsmResponse.GetId())
	g.Expect(clientConnection1_1.GetID()).To(Equal("1"))
	g.Expect(clientConnection1_1.Xcon.GetRemoteDestination().GetId()).To(Equal("3"))
}
