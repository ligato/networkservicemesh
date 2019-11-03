package networkservice

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
)

// Clone clones request
func (ns *NetworkServiceRequest) Clone() *NetworkServiceRequest {
	return proto.Clone(ns).(*NetworkServiceRequest)
}

// GetRequestConnection returns request connection
func (ns *NetworkServiceRequest) GetRequestConnection() *connection.Connection {
	return ns.GetConnection()
}

// SetRequestConnection sets request connection
func (ns *NetworkServiceRequest) SetRequestConnection(conn *connection.Connection) {
	ns.Connection = conn
}

// GetRequestMechanismPreferences returns request mechanism preferences
func (ns *NetworkServiceRequest) GetRequestMechanismPreferences() []*connection.Mechanism {
	preferences := make([]*connection.Mechanism, 0, len(ns.MechanismPreferences))
	for _, m := range ns.MechanismPreferences {
		preferences = append(preferences, m)
	}

	return preferences
}

// SetRequestMechanismPreferences sets request mechanism preferences
func (ns *NetworkServiceRequest) SetRequestMechanismPreferences(mechanismPreferences []*connection.Mechanism) {
	ns.MechanismPreferences = mechanismPreferences
}

// IsValid returns if request is valid
func (ns *NetworkServiceRequest) IsValid() error {
	if ns == nil {
		return errors.New("request cannot be nil")
	}

	if ns.GetConnection() == nil {
		return errors.Errorf("request.Connection cannot be nil %v", ns)
	}

	if err := ns.GetConnection().IsValid(); err != nil {
		return errors.Errorf("request.Connection is invalid: %s: %v", err, ns)
	}

	if ns.GetMechanismPreferences() == nil {
		return errors.Errorf("request.MechanismPreferences cannot be nil: %v", ns)
	}

	if len(ns.GetMechanismPreferences()) < 1 {
		return errors.Errorf("request.MechanismPreferences must have at least one entry: %v", ns)
	}

	return nil
}
