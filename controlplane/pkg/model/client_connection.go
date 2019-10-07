package model

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/nsm/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/nsm/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
)

// ClientConnectionState describes state of ClientConnection
type ClientConnectionState int8

const (
	// ClientConnectionReady means connection is in state 'ready'
	ClientConnectionReady ClientConnectionState = 0

	// ClientConnectionRequesting means connection waits answer from NSE or Dp
	ClientConnectionRequesting ClientConnectionState = 1

	// ClientConnectionBroken means connection failed requesting
	ClientConnectionBroken ClientConnectionState = 2

	// ClientConnectionHealing means connection is in 'healing' state
	ClientConnectionHealing ClientConnectionState = 3

	// ClientConnectionClosing means connection is started closing process
	ClientConnectionClosing ClientConnectionState = 4
)

// ClientConnection struct in model that describes cross connect between NetworkServiceClient and NetworkServiceEndpoint
type ClientConnection struct {
	ConnectionID            string
	Request                 networkservice.Request
	Xcon                    *crossconnect.CrossConnect
	RemoteNsm               *registry.NetworkServiceManager
	Endpoint                *registry.NSERegistration
	DataplaneRegisteredName string
	ConnectionState         ClientConnectionState
	DataplaneState          DataplaneState
}

// GetID returns id of clientConnection
func (cc *ClientConnection) GetID() string {
	if cc == nil {
		return ""
	}
	return cc.ConnectionID
}

// GetNetworkService returns name of networkService of clientConnection
func (cc *ClientConnection) GetNetworkService() string {
	if cc == nil {
		return ""
	}
	return cc.Endpoint.GetNetworkService().GetName()
}

// GetConnectionSource returns source part of connection
func (cc *ClientConnection) GetConnectionSource() connection.Connection {
	return cc.Xcon.GetSourceConnection()
}

// GetConnectionDestination returns destination part of connection
func (cc *ClientConnection) GetConnectionDestination() connection.Connection {
	return cc.Xcon.GetDestinationConnection()
}

// Clone return pointer to copy of ClientConnection
func (cc *ClientConnection) clone() cloneable {
	if cc == nil {
		return nil
	}

	var xcon *crossconnect.CrossConnect
	if cc.Xcon != nil {
		xcon = proto.Clone(cc.Xcon).(*crossconnect.CrossConnect)
	}

	var remoteNsm *registry.NetworkServiceManager
	if cc.RemoteNsm != nil {
		remoteNsm = proto.Clone(cc.RemoteNsm).(*registry.NetworkServiceManager)
	}

	var endpoint *registry.NSERegistration
	if cc.Endpoint != nil {
		endpoint = proto.Clone(cc.Endpoint).(*registry.NSERegistration)
	}

	var request networkservice.Request
	if cc.Request != nil {
		request = cc.Request.Clone()
	}

	return &ClientConnection{
		ConnectionID:            cc.ConnectionID,
		Xcon:                    xcon,
		RemoteNsm:               remoteNsm,
		Endpoint:                endpoint,
		DataplaneRegisteredName: cc.DataplaneRegisteredName,
		Request:                 request,
		ConnectionState:         cc.ConnectionState,
		DataplaneState:          cc.DataplaneState,
	}
}

type clientConnectionDomain struct {
	baseDomain
}

func newClientConnectionDomain() clientConnectionDomain {
	return clientConnectionDomain{
		baseDomain: newBase(),
	}
}

func (d *clientConnectionDomain) AddClientConnection(ctx context.Context, cc *ClientConnection) {
	d.store(ctx, cc.ConnectionID, cc)
}

func (d *clientConnectionDomain) GetClientConnection(id string) *ClientConnection {
	v, _ := d.load(id)
	if v != nil {
		return v.(*ClientConnection)
	}
	return nil
}

func (d *clientConnectionDomain) GetAllClientConnections() []*ClientConnection {
	var rv []*ClientConnection
	d.kvRange(func(_ string, value interface{}) bool {
		rv = append(rv, value.(*ClientConnection))
		return true
	})
	return rv
}

func (d *clientConnectionDomain) DeleteClientConnection(ctx context.Context, id string) {
	d.delete(ctx, id)
}

func (d *clientConnectionDomain) UpdateClientConnection(ctx context.Context, cc *ClientConnection) {
	d.store(ctx, cc.ConnectionID, cc)
}

func (d *clientConnectionDomain) ApplyClientConnectionChanges(ctx context.Context, id string, f func(*ClientConnection)) *ClientConnection {
	upd := d.applyChanges(ctx, id, func(v interface{}) { f(v.(*ClientConnection)) })
	if upd != nil {
		return upd.(*ClientConnection)
	}
	return nil
}

func (d *clientConnectionDomain) SetClientConnectionModificationHandler(h *ModificationHandler) func() {
	return d.addHandler(h)
}
