The simplest possible case for Network Service Mesh is to have is connecting a Client via a vWire to another Pod that is providing a Network Service.
Network Service Mesh allows flexibility in the choice of mechanisms used to provide that vWire to a workload.
[the icmp responder example](icmp-responder.md) does this with kernel interfaces.  The vpp-icmp-responder provides and
consumes the same 'icmp-responder' Network Service, but has Client's and Endpoint's that use a [memif](https://www.youtube.com/watch?v=6aVr32WgY0Q) high speed
memory interfaces to achieve performance unavailable via kernel interfaces.


![vpp-icmp-responder-example](../images/vpp-icmp-responder-example.svg)

## Deploy
Utilize the [Run](../guide-quickstart.md) instructions to install the NSM infrastructure, and then type:

```bash
make helm-install-vpp-icmp-responder
```

## What it does

This will install two Deployments:

Name | Description 
:--------|:--------
**vpp-icmp-responder-nsc** | The Clients, four replicas
**vpp-icmp-responder-nse** | The Endpoints, two replicas

And cause each Client to get a vWire connecting it to one of the Endpoints.  Network Service Mesh handles the
Network Service Discovery and Routing, as well as the vWire 'Connection Handling' for setting all of this up.

![vpp-icmp-responder-example-2](../images/vpp-icmp-responder-example-2.svg)

In order to make this case more interesting, Endpoint1 and Endpoint2 are deployed on two separate Nodes using
PodAntiAffinity, so that the Network Service Mesh has to demonstrate the ability to string vWires between Clients and
Endpoints on the same Node and Clients and Endpoints on different Nodes.

## Verify

First verify that the vpp-icmp-responder example Pods are all up and running:

```bash
kubectl get pods | grep vpp-icmp-responder
```

To see the vpp-icmp-responder example in action, you can run:

```bash
make k8s-icmp-check
```
