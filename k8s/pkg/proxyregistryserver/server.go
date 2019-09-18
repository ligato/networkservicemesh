package proxyregistryserver

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/clusterinfo"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
	nsmClientset "github.com/networkservicemesh/networkservicemesh/k8s/pkg/networkservice/clientset/versioned"
	"github.com/networkservicemesh/networkservicemesh/k8s/pkg/registryserver"
)

// New starts proxy Network Service Discovery Server and Cluster Info Server
func New(clientset *nsmClientset.Clientset, clusterInfoService clusterinfo.ClusterInfoServer) *grpc.Server {
	tracer := opentracing.GlobalTracer()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer)))

	cache := registryserver.NewRegistryCache(clientset, &registryserver.ResourceFilterConfig{})

	discovery := newDiscoveryService(cache, clusterInfoService)

	registry.RegisterNetworkServiceDiscoveryServer(server, discovery)
	clusterinfo.RegisterClusterInfoServer(server, clusterInfoService)

	if err := cache.Start(); err != nil {
		logrus.Error(err)
	}
	logrus.Info("RegistryCache started")

	return server
}
