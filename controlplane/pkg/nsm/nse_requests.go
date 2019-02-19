package nsm

import (
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/local/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/nsm"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/registry"
	remote_connection "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/connection"
	remote_networkservice "github.com/networkservicemesh/networkservicemesh/controlplane/pkg/apis/remote/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/model"
)

func (srv *networkServiceManager) createRemoteNSMRequest(endpoint *registry.NSERegistration, requestConnection nsm.NSMConnection, dataplane *model.Dataplane, existingConnection *model.ClientConnection) *remote_networkservice.NetworkServiceRequest {
	// We need to obtain parameters for remote mechanism
	remoteM := []*remote_connection.Mechanism{}
	for _, mechanism := range dataplane.RemoteMechanisms {
		remoteM = append(remoteM, mechanism)
	}
	var message *remote_networkservice.NetworkServiceRequest
	if existingConnection != nil {
		remoteDst := existingConnection.Xcon.GetRemoteDestination()
		message = &remote_networkservice.NetworkServiceRequest{

			Connection: &remote_connection.Connection{
				Id:                                   remoteDst.GetId(),
				NetworkService:                       remoteDst.NetworkService,
				Context:                              remoteDst.GetContext(),
				Labels:                               remoteDst.GetLabels(),
				DestinationNetworkServiceManagerName: endpoint.GetNetworkServiceManager().GetName(),
				SourceNetworkServiceManagerName:      srv.getNetworkServiceManagerName(),
				NetworkServiceEndpointName:           endpoint.GetNetworkserviceEndpoint().GetEndpointName(),
			},
			MechanismPreferences: remoteM,
		}
	} else {
		message = &remote_networkservice.NetworkServiceRequest{
			Connection: &remote_connection.Connection{
				Id:                                   "-",
				NetworkService:                       requestConnection.GetNetworkService(),
				Context:                              requestConnection.GetContext(),
				Labels:                               requestConnection.GetLabels(),
				DestinationNetworkServiceManagerName: endpoint.GetNetworkServiceManager().GetName(),
				SourceNetworkServiceManagerName:      srv.getNetworkServiceManagerName(),
				NetworkServiceEndpointName:           endpoint.GetNetworkserviceEndpoint().GetEndpointName(),
			},
			MechanismPreferences: remoteM,
		}
	}
	return message
}

func (srv *networkServiceManager) createLocalNSERequest(endpoint *registry.NSERegistration, requestConnection nsm.NSMConnection) *networkservice.NetworkServiceRequest {

	message := &networkservice.NetworkServiceRequest{
		Connection: &connection.Connection{
			Id:             srv.createConnectionId(),
			NetworkService: endpoint.GetNetworkService().GetName(),
			Context:        requestConnection.GetContext(),
			Labels:         requestConnection.GetLabels(),
		},
		MechanismPreferences: []*connection.Mechanism{
			{
				Type:       connection.MechanismType_MEM_INTERFACE,
				Parameters: map[string]string{},
			},
			{
				Type:       connection.MechanismType_KERNEL_INTERFACE,
				Parameters: map[string]string{},
			},
		},
	}
	return message
}
