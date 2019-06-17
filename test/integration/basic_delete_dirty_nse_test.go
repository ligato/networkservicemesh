// +build basic

package nsmd_integration_tests

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/test/kubetest"
	"github.com/networkservicemesh/networkservicemesh/test/kubetest/pods"
)

func TestDeleteDirtyNSE(t *testing.T) {
	RegisterTestingT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	logrus.Print("Running delete dirty NSE test")

	k8s, err := kubetest.NewK8s(true)
	Expect(err).To(BeNil())
	defer k8s.Cleanup()

	nodesConf, err := kubetest.SetupNodesConfig(k8s, 1, defaultTimeout, []*pods.NSMgrPodConfig{}, k8s.GetK8sNamespace())
	Expect(err).To(BeNil())
	defer kubetest.FailLogger(k8s, nodesConf, t)

	nsePod := kubetest.DeployDirtyICMP(k8s, nodesConf[0].Node, "dirty-icmp-responder-nse", defaultTimeout)

	kubetest.ExpectNSEsCountToBe(k8s, 0, 1)

	k8s.DeletePods(nsePod)

	kubetest.ExpectNSEsCountToBe(k8s, 1, 0)
}

func TestDeleteDirtyNSEWithClient(t *testing.T) {
	RegisterTestingT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	logrus.Print("Running delete dirty NSE with client test")

	k8s, err := kubetest.NewK8s(true)
	Expect(err).To(BeNil())
	defer k8s.Cleanup()

	nodesConf, err := kubetest.SetupNodesConfig(k8s, 1, defaultTimeout, []*pods.NSMgrPodConfig{}, k8s.GetK8sNamespace())
	Expect(err).To(BeNil())
	defer kubetest.FailLogger(k8s, nodesConf, t)

	nsePod := kubetest.DeployDirtyICMP(k8s, nodesConf[0].Node, "dirty-icmp-responder-nse", defaultTimeout)
	kubetest.DeployNSC(k8s, nodesConf[0].Node, "nsc-1", defaultTimeout)

	kubetest.ExpectNSEsCountToBe(k8s, 0, 1)

	k8s.DeletePods(nsePod)

	kubetest.ExpectNSEsCountToBe(k8s, 1, 1)
}
