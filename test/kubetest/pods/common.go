package pods

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

const (
	//DefaultKubeletAdminClusterRole default kubelet api admin cluster role name
	DefaultKubeletAdminClusterRole = "system:kubelet-api-admin"
	//DefaultAccount creates on namespace creating
	DefaultAccount = "default"
	// EnvForwardingPlane is the environment variable for configuring the forwarding plane
	EnvForwardingPlane = "FORWARDING_PLANE"
	// EnvForwardingPlaneVPP is the VPP forwarding plane
	EnvForwardingPlaneVPP = "vpp"
	// EnvForwardingPlaneKernel is the Kernel forwarding plane
	EnvForwardingPlaneKernel = "kernel-forwarder"
	// EnvForwardingPlaneDefault is the default forwarding plane
	EnvForwardingPlaneDefault = EnvForwardingPlaneVPP
	// NSEServiceAccount service account for Network Service Endpoints
	NSEServiceAccount = "nse-acc"
	// NSCServiceAccount service account for Network Service Clients
	NSCServiceAccount = "nsc-acc"
	// NSMgrServiceAccount service account for Network Service Managers
	NSMgrServiceAccount = "nsmgr-acc"
	// NSMRSServiceAccount service account for Network Service Manager Registry Server
	NSMRSServiceAccount = "nsmrs-acc"
	// ForwardPlaneServiceAccount service account for Forwarding Plane
	ForwardPlaneServiceAccount = "forward-plane-acc"

	// Keep in sync with ../../../k8s/cmd/admission-webhook/const.go.
	// Limits for 'nsm-monitor' container.
	nsmMonitorCPULimit    = "100m"
	nsmMonitorMemoryLimit = "15Mi"
	// Limits for 'coredns' container.
	corednsCPULimit    = "100m"
	corednsMemoryLimit = "15Mi"
	// Limits for 'nsm-init' container.
	nsmInitCPULimit    = "200m"
	nsmInitMemoryLimit = "20Mi"
	// Limits for 'nsm-dns-init' container.
	nsmDNSInitCPULimit    = "200m"
	nsmDNSInitMemoryLimit = "15Mi"
)

// ForwardingPlane - Wrapper for getting a forwarding plane pod
func ForwardingPlane(name string, node *v1.Node, plane string) *v1.Pod {
	if plane == EnvForwardingPlaneDefault {
		return VPPForwarderPod(name, node)
	} else if plane == EnvForwardingPlaneKernel {
		return KernelForwarderPod(name, node)
	}
	logrus.Error("Forwarding plane error: Unknown forwarder")
	return nil
}

// ForwardingPlaneWithConfig - Wrapper for getting a forwarding plane pod
func ForwardingPlaneWithConfig(name string, node *v1.Node, variables map[string]string, plane string) *v1.Pod {
	if plane == EnvForwardingPlaneDefault {
		return VPPForwarderPodConfig(name, node, variables)
	} else if plane == EnvForwardingPlaneKernel {
		return KernelForwarderPodConfig(name, node, variables)
	}
	logrus.Error("Forwarding plane error: Unknown forwarder")
	return nil
}

// ForwardingPlaneWithLiveCheck - Wrapper for getting a forwarding plane pod with liveness/readiness probes
func ForwardingPlaneWithLiveCheck(name string, node *v1.Node, plane string) *v1.Pod {
	if plane == EnvForwardingPlaneDefault {
		return VPPForwarderPodLiveCheck(name, node)
	} else if plane == EnvForwardingPlaneKernel {
		return KernelForwarderPodLiveCheck(name, node)
	}
	logrus.Error("Forwarding plane error: Unknown forwarder")
	return nil
}
