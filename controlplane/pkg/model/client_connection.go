package model

import (
	"context"

	"github.com/networkservicemesh/networkservicemesh/sdk/compat"

	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"

	"github.com/networkservicemesh/networkservicemesh/sdk/monitor"

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

	// ClientConnectionRequesting means connection - is trying to be established for first time.
	ClientConnectionRequesting ClientConnectionState = 1

	// ClientConnectionBroken means connection failed requesting
	ClientConnectionBroken ClientConnectionState = 2

	// ClientConnectionHealingBegin means connection - is trying to be updated or to be healed.
	ClientConnectionHealingBegin ClientConnectionState = 3

	// ClientConnectionHealing means connection - is in progress of being updated or to be healed.
	ClientConnectionHealing ClientConnectionState = 4

	// ClientConnectionClosing means connection is started closing process
	ClientConnectionClosing ClientConnectionState = 5
)

// ClientConnection struct in model that describes cross connect between NetworkServiceClient and NetworkServiceEndpoint
type ClientConnection struct {
	ConnectionID            string
	Request                 networkservice.Request
	Xcon                    *crossconnect.CrossConnect
	RemoteNsm               *registry.NetworkServiceManager
	Endpoint                *registry.NSERegistration
	ForwarderRegisteredName string
	ConnectionState         ClientConnectionState
	ForwarderState          ForwarderState
	Span                    opentracing.Span
	Monitor                 monitor.Server
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
	if cc.Xcon == nil {
		return nil
	}
	return compat.ConnectionUnifiedToNSM(cc.Xcon.GetSource())
}

// GetConnectionDestination returns destination part of connection
func (cc *ClientConnection) GetConnectionDestination() connection.Connection {
	return compat.ConnectionUnifiedToNSM(cc.Xcon.GetDestination())
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
		ForwarderRegisteredName: cc.ForwarderRegisteredName,
		Request:                 request,
		ConnectionState:         cc.ConnectionState,
		ForwarderState:          cc.ForwarderState,
		Span:                    cc.Span,
		Monitor:                 cc.Monitor,
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

func (d *clientConnectionDomain) DeleteClientConnection(ctx context.Context, connectionID string) {
	d.delete(ctx, connectionID)
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
