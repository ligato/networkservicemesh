package model

import (
	"fmt"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/registry"
	. "github.com/onsi/gomega"
	"strconv"
	"testing"
)

func TestAddAndGetСс(t *testing.T) {
	RegisterTestingT(t)

	cc := &ClientConnection{
		ConnectionId: "1",
		Xcon: &crossconnect.CrossConnect{
			Id: "1",
		},
		RemoteNsm: &registry.NetworkServiceManager{
			Name: "master",
			Url:  "1.1.1.1",
		},
		Endpoint: &registry.NSERegistration{
			NetworkService: &registry.NetworkService{
				Name: "ns1",
			},
			NetworkServiceManager: &registry.NetworkServiceManager{
				Name: "worker",
				Url:  "2.2.2.2",
			},
			NetworkserviceEndpoint: &registry.NetworkServiceEndpoint{
				NetworkServiceName: "ns1",
				EndpointName:       "endp1",
			},
		},
		Dataplane: &Dataplane{
			RegisteredName: "dp1",
		},
		ConnectionState: ClientConnection_Healing,
		Request:         nil,
		DataplaneState:  DataplaneState_Ready,
	}

	ccd := clientConnectionDomain{}
	ccd.AddClientConnection(cc)
	getConn := ccd.GetClientConnection("1")

	Expect(getConn.ConnectionId).To(Equal(cc.ConnectionId))
	Expect(getConn.ConnectionState).To(Equal(cc.ConnectionState))
	Expect(getConn.DataplaneState).To(Equal(cc.DataplaneState))
	Expect(getConn.Request).To(BeNil())

	Expect(getConn.GetNetworkService()).To(Equal(cc.GetNetworkService()))
	Expect(getConn.GetId()).To(Equal(cc.GetId()))
	Expect(getConn.GetConnectionSource()).To(Equal(cc.GetConnectionSource()))

	Expect(fmt.Sprintf("%p", getConn.RemoteNsm)).ToNot(Equal(fmt.Sprintf("%p", cc.RemoteNsm)))
	Expect(fmt.Sprintf("%p", getConn.Dataplane)).ToNot(Equal(fmt.Sprintf("%p", cc.Dataplane)))
	Expect(fmt.Sprintf("%p", getConn.Endpoint)).ToNot(Equal(fmt.Sprintf("%p", cc.Endpoint)))
	Expect(fmt.Sprintf("%p", getConn.Endpoint.NetworkServiceManager)).
		ToNot(Equal(fmt.Sprintf("%p", cc.Endpoint.NetworkServiceManager)))
}

func TestGetAllСс(t *testing.T) {
	RegisterTestingT(t)

	ccd := clientConnectionDomain{}
	amount := 5

	for i := 0; i < amount; i++ {
		ccd.AddClientConnection(&ClientConnection{
			ConnectionId: fmt.Sprintf("%d", i),
			Xcon: &crossconnect.CrossConnect{
				Id: "1",
			},
			RemoteNsm: &registry.NetworkServiceManager{
				Name: "master",
				Url:  "1.1.1.1",
			},
			Endpoint: &registry.NSERegistration{
				NetworkService: &registry.NetworkService{
					Name: "ns1",
				},
				NetworkServiceManager: &registry.NetworkServiceManager{
					Name: "worker",
					Url:  "2.2.2.2",
				},
				NetworkserviceEndpoint: &registry.NetworkServiceEndpoint{
					NetworkServiceName: "ns1",
					EndpointName:       "endp1",
				},
			},
			Dataplane: &Dataplane{
				RegisteredName: "dp1",
			},
			ConnectionState: ClientConnection_Healing,
			Request:         nil,
			DataplaneState:  DataplaneState_Ready,
		})
	}

	all := ccd.GetAllClientConnections()
	Expect(len(all)).To(Equal(amount))

	expected := make([]bool, amount)
	for i := 0; i < amount; i++ {
		index, _ := strconv.ParseInt(all[i].ConnectionId, 10, 64)
		expected[index] = true
	}

	for i := 0; i < amount; i++ {
		Expect(expected[i]).To(BeTrue())
	}
}

func TestDeleteСс(t *testing.T) {
	RegisterTestingT(t)

	ccd := clientConnectionDomain{}
	ccd.AddClientConnection(&ClientConnection{
		ConnectionId: "1",
		Xcon: &crossconnect.CrossConnect{
			Id: "1",
		},
		RemoteNsm: &registry.NetworkServiceManager{
			Name: "master",
			Url:  "1.1.1.1",
		},
		Endpoint: &registry.NSERegistration{
			NetworkService: &registry.NetworkService{
				Name: "ns1",
			},
			NetworkServiceManager: &registry.NetworkServiceManager{
				Name: "worker",
				Url:  "2.2.2.2",
			},
			NetworkserviceEndpoint: &registry.NetworkServiceEndpoint{
				NetworkServiceName: "ns1",
				EndpointName:       "endp1",
			},
		},
		Dataplane: &Dataplane{
			RegisteredName: "dp1",
		},
		ConnectionState: ClientConnection_Healing,
		Request:         nil,
		DataplaneState:  DataplaneState_Ready,
	})

	cc := ccd.GetClientConnection("1")
	Expect(cc).ToNot(BeNil())

	ccd.DeleteClientConnection("1")

	ccDel := ccd.GetClientConnection("1")
	Expect(ccDel).To(BeNil())

	ccd.DeleteClientConnection("NotExistingId")

}

func TestUpdateExistingСс(t *testing.T) {
	RegisterTestingT(t)

	cc := &ClientConnection{
		ConnectionId: "1",
		Xcon: &crossconnect.CrossConnect{
			Id: "1",
		},
		RemoteNsm: &registry.NetworkServiceManager{
			Name: "master",
			Url:  "1.1.1.1",
		},
		Endpoint: &registry.NSERegistration{
			NetworkService: &registry.NetworkService{
				Name: "ns1",
			},
			NetworkServiceManager: &registry.NetworkServiceManager{
				Name: "worker",
				Url:  "2.2.2.2",
			},
			NetworkserviceEndpoint: &registry.NetworkServiceEndpoint{
				NetworkServiceName: "ns1",
				EndpointName:       "endp1",
			},
		},
		Dataplane: &Dataplane{
			RegisteredName: "dp1",
		},
		ConnectionState: ClientConnection_Healing,
		Request:         nil,
		DataplaneState:  DataplaneState_Ready,
	}

	ccd := clientConnectionDomain{}
	ccd.AddClientConnection(cc)

	newUrl := "3.3.3.3"
	newDpName := "updatedName"
	cc.Endpoint.NetworkServiceManager.Url = newUrl
	cc.Dataplane.RegisteredName = newDpName

	notUpdated := ccd.GetClientConnection("1")
	Expect(notUpdated.Endpoint.NetworkServiceManager.Url).ToNot(Equal(newUrl))
	Expect(notUpdated.Dataplane.RegisteredName).ToNot(Equal(newDpName))

	ccd.UpdateClientConnection(cc)
	updated := ccd.GetClientConnection("1")
	Expect(updated.Endpoint.NetworkServiceManager.Url).To(Equal(newUrl))
	Expect(updated.Dataplane.RegisteredName).To(Equal(newDpName))
}

func TestUpdateNotExistingСс(t *testing.T) {
	RegisterTestingT(t)

	cc := &ClientConnection{
		ConnectionId: "1",
		Xcon: &crossconnect.CrossConnect{
			Id: "1",
		},
		RemoteNsm: &registry.NetworkServiceManager{
			Name: "master",
			Url:  "1.1.1.1",
		},
		Endpoint: &registry.NSERegistration{
			NetworkService: &registry.NetworkService{
				Name: "ns1",
			},
			NetworkServiceManager: &registry.NetworkServiceManager{
				Name: "worker",
				Url:  "2.2.2.2",
			},
			NetworkserviceEndpoint: &registry.NetworkServiceEndpoint{
				NetworkServiceName: "ns1",
				EndpointName:       "endp1",
			},
		},
		Dataplane: &Dataplane{
			RegisteredName: "dp1",
		},
		ConnectionState: ClientConnection_Healing,
		Request:         nil,
		DataplaneState:  DataplaneState_Ready,
	}

	ccd := clientConnectionDomain{}

	ccd.UpdateClientConnection(cc)
	updated := ccd.GetClientConnection("1")
	Expect(updated.Endpoint.NetworkServiceManager.Url).To(Equal("2.2.2.2"))
	Expect(updated.Dataplane.RegisteredName).To(Equal("dp1"))
}
