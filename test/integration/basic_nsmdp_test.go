// +build basic_suite

package integration

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/networkservicemesh/networkservicemesh/test/kubetest"
	"github.com/networkservicemesh/networkservicemesh/test/kubetest/pods"
)

func TestNSMDDP(t *testing.T) {
	g := NewWithT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	k8s, err := kubetest.NewK8s(g, kubetest.ReuseNSMResources)
	defer k8s.Cleanup(t)

	g.Expect(err).To(BeNil())
	defer k8s.SaveTestArtifacts(t)
	nodes, err := kubetest.SetupNodes(k8s, 1, defaultTimeout)
	g.Expect(err).To(BeNil())
	icmpPod := kubetest.DeployICMP(k8s, nodes[0].Node, "icmp-responder-nse-1", defaultTimeout)

	nsmdName := nodes[0].Nsmd.Name
	k8s.DeletePods(nodes[0].Nsmd, icmpPod)
	nodes[0].Nsmd = k8s.CreatePod(pods.NSMgrPod(nsmdName, nodes[0].Node, k8s.GetK8sNamespace())) // Recovery NSEs
	// Wait for NSMgr to be deployed, to not get admission error
	kubetest.WaitNSMgrDeployed(k8s, nodes[0].Nsmd, defaultTimeout)
	icmpPod = kubetest.DeployICMP(k8s, nodes[0].Node, "icmp-responder-nse-2", defaultTimeout)
	g.Expect(icmpPod).ToNot(BeNil())
}
