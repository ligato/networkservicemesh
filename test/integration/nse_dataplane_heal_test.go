package nsmd_integration_tests

import (
	"fmt"
	"github.com/networkservicemesh/networkservicemesh/test/integration/nsmd_test_utils"
	"github.com/networkservicemesh/networkservicemesh/test/kube_testing/pods"
	"testing"
	"time"

	"github.com/networkservicemesh/networkservicemesh/test/kube_testing"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestDataplaneHealLocal(t *testing.T) {
	RegisterTestingT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	testDataplaneHeal(t, 1)
}

func TestDataplaneHealRemote(t *testing.T) {
	RegisterTestingT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	testDataplaneHeal(t, 2)
}

/**
If passed 1 both will be on same node, if not on different.
*/
func testDataplaneHeal(t *testing.T, nodesCount int) {
	k8s, err := kube_testing.NewK8s()
	defer k8s.Cleanup()

	Expect(err).To(BeNil())

	s1 := time.Now()
	k8s.Prepare("nsmd", "nsc", "nsmd-dataplane", "icmp-responder-nse", "jaeger")
	logrus.Printf("Cleanup done: %v", time.Since(s1))

	// Deploy open tracing to see what happening.
	nodes_setup := nsmd_test_utils.SetupNodes(k8s, nodesCount )

	// Run ICMP on latest node
	_ = nsmd_test_utils.DeployIcmp(k8s, nodes_setup[nodesCount-1].Node, "icmp-responder-nse1")

	nscPodNode := nsmd_test_utils.DeployNsc(k8s, nodes_setup[0].Node, "nsc1")
	var nscInfo *nsmd_test_utils.NSCCheckInfo
	failures := InterceptGomegaFailures(func() {
		nscInfo = nsmd_test_utils.CheckNSC(k8s, t, nscPodNode)
	})
	// Do dumping of container state to dig into what is happened.
	printErrors(failures, k8s, nodes_setup, nscInfo, t)

	logrus.Infof("Delete Selected dataplane")
	k8s.DeletePods(nodes_setup[nodesCount-1].Dataplane)

	logrus.Infof("Wait NSMD is waiting for dataplane recovery")
	k8s.WaitLogsContains(nodes_setup[nodesCount-1].Nsmd, "nsmd", "Waiting for Dataplane to recovery...", 60*time.Second)
	// Now are are in dataplane dead state, and in Heal procedure waiting for dataplane.
	dpName := fmt.Sprintf("nsmd-dataplane-recovered-%d", nodesCount-1)

	logrus.Infof("Starting recovered dataplane...")
	startTime := time.Now()
	nodes_setup[nodesCount-1].Dataplane = k8s.CreatePod(pods.VPPDataplanePod(dpName, nodes_setup[nodesCount-1].Node))
	logrus.Printf("Started new Dataplane: %v on node %s", time.Since(startTime), nodes_setup[nodesCount-1].Node.Name)

	// Check NSMd goint into HEAL state.

	logrus.Infof("Waiting for connection recovery...")
	if nodesCount > 1 {
		k8s.WaitLogsContains(nodes_setup[nodesCount-1].Nsmd, "nsmd", "Healing will be continued on source side...", 60*time.Second)
		k8s.WaitLogsContains(nodes_setup[0].Nsmd, "nsmd", "Heal: Connection recovered:", 60*time.Second)
	} else {
		k8s.WaitLogsContains(nodes_setup[nodesCount-1].Nsmd, "nsmd", "Heal: Connection recovered:", 60*time.Second)
	}
	logrus.Infof("Waiting for connection recovery Done...")

	failures = InterceptGomegaFailures(func() {
		nscInfo = nsmd_test_utils.CheckNSC(k8s, t, nscPodNode)
	})
	printErrors(failures, k8s, nodes_setup, nscInfo, t)
}