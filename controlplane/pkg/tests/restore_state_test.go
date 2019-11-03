package tests

import (
	"context"
	"testing"
	"time"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/common"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/kernel"

	. "github.com/onsi/gomega"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
)

func TestRestoreConnectionState(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()

	srv.AddFakeForwarder("dp1", "tcp:some_address")

	g.Expect(srv.nsmServer.Manager().WaitForForwarder(context.Background(), 1*time.Millisecond).Error()).To(Equal("failed to wait for NSMD stare restore... timeout 1ms happened"))

	xcons := []*crossconnect.CrossConnect{}
	srv.nsmServer.Manager().RestoreConnections(xcons, "dp1", srv.nsmServer)
	g.Expect(srv.nsmServer.Manager().WaitForForwarder(context.Background(), 1*time.Second)).To(BeNil())
}

func TestRestoreConnectionStateWrongDst(t *testing.T) {
	g := NewWithT(t)

	storage := NewSharedStorage()
	srv := NewNSMDFullServer(Master, storage)
	defer srv.Stop()

	srv.AddFakeForwarder("dp1", "tcp:some_address")
	srv.registerFakeEndpointWithName("ns1", "IP", Worker, "ep2")

	nsmClient := srv.RequestNSM("nsm1")

	requestConnection := &connection.Connection{
		Id:             "1",
		NetworkService: "ns1",
		Mechanism: &connection.Mechanism{
			Parameters: map[string]string{
				common.Workspace: nsmClient.Workspace,
			},
		},
		NetworkServiceManagers: []string{"src"},
	}

	dstConnection := &connection.Connection{
		Id: "2",
		Mechanism: &connection.Mechanism{
			Type: kernel.MECHANISM,
			Parameters: map[string]string{
				kernel.WorkspaceNSEName: "nse1",
			},
		},
		NetworkService:         "ns1",
		NetworkServiceManagers: []string{"src"},
	}
	xcons := []*crossconnect.CrossConnect{
		{
			Source:      requestConnection,
			Destination: dstConnection,
			Id:          "1",
		},
	}
	srv.nsmServer.Manager().RestoreConnections(xcons, "dp1", srv.nsmServer)
	g.Expect(srv.nsmServer.Manager().WaitForForwarder(context.Background(), 1*time.Second)).To(BeNil())
	g.Expect(len(srv.TestModel.GetAllClientConnections())).To(Equal(0))
}
