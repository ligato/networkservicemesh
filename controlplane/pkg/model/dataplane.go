package model

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sync"

	local "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	remote "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/connection"
)

type DataplaneState int8

const (
	DataplaneState_None  DataplaneState = 0 // In case dataplane is not yet configured for connection
	DataplaneState_Ready DataplaneState = 1 // In case dataplane is configured for connection.
)

type Dataplane struct {
	RegisteredName       string
	SocketLocation       string
	LocalMechanisms      []*local.Mechanism
	RemoteMechanisms     []*remote.Mechanism
	MechanismsConfigured bool
}

func (d *Dataplane) Clone() *Dataplane {
	if d == nil {
		return nil
	}

	var lm []*local.Mechanism
	for _, m := range d.LocalMechanisms {
		lm = append(lm, proto.Clone(m).(*local.Mechanism))
	}

	var rm []*remote.Mechanism
	for _, m := range d.RemoteMechanisms {
		rm = append(rm, proto.Clone(m).(*remote.Mechanism))
	}

	return &Dataplane{
		RegisteredName:       d.RegisteredName,
		SocketLocation:       d.SocketLocation,
		LocalMechanisms:      lm,
		RemoteMechanisms:     rm,
		MechanismsConfigured: d.MechanismsConfigured,
	}
}

type dataplaneDomain struct {
	baseDomain
	inner sync.Map
}

func (d *dataplaneDomain) AddDataplane(dp *Dataplane) {
	d.inner.Store(dp.RegisteredName, dp.Clone())
	d.resourceAdded(dp.Clone())
}

func (d *dataplaneDomain) GetDataplane(name string) *Dataplane {
	v, _ := d.inner.Load(name)
	if v != nil {
		return v.(*Dataplane).Clone()
	}
	return nil
}

func (d *dataplaneDomain) DeleteDataplane(name string) {
	v := d.GetDataplane(name)
	if v == nil {
		return
	}
	d.inner.Delete(name)
	d.resourceDeleted(v)
}

func (d *dataplaneDomain) UpdateDataplane(dp *Dataplane) {
	v := d.GetDataplane(dp.RegisteredName)
	if v == nil {
		d.AddDataplane(dp)
		return
	}
	d.inner.Store(dp.RegisteredName, dp.Clone())
	d.resourceUpdated(v, dp.Clone())
}

func (d *dataplaneDomain) SelectDataplane(dataplaneSelector func(dp *Dataplane) bool) (*Dataplane, error) {
	var rv *Dataplane
	d.inner.Range(func(key, value interface{}) bool {
		dp := value.(*Dataplane)

		if dataplaneSelector == nil {
			rv = dp
			return false
		}

		if dataplaneSelector(dp) {
			rv = dp
			return false
		}

		return true
	})

	if rv == nil {
		return nil, fmt.Errorf("no appropriate dataplanes found")
	}

	return rv.Clone(), nil
}

func (d *dataplaneDomain) SetDataplaneModificationHandler(h *ModificationHandler) func() {
	deleteFunc := d.addHandler(h)
	d.inner.Range(func(key, value interface{}) bool {
		d.resourceAdded(value.(*Dataplane).Clone())
		return true
	})
	return deleteFunc
}
