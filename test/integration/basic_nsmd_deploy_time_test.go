// +build basic

package nsmd_integration_tests

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/test/kubetest"
)

func TestNSMDDeploy(t *testing.T) {
	g := NewWithT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	logrus.Print("Running NSMgr Deploy test")

	k8s, err := kubetest.NewK8s(g, true)
	g.Expect(err).To(BeNil())

	// Warmup

	_, err = kubetest.SetupNodes(k8s, 1, defaultTimeout)
	g.Expect(err).To(BeNil())
	k8s.Cleanup()

	st := time.Now()
	k8s, err = kubetest.NewK8s(g, true)
	defer k8s.Cleanup()
	defer kubetest.MakeLogsSnapshot(k8s, t)
	_, err = kubetest.SetupNodes(k8s, 1, defaultTimeout)
	g.Expect(err).To(BeNil())
	deploy := time.Now()
	k8s.Cleanup()
	destroy := time.Now()

	logrus.Infof("Pods Start time: %v", deploy.Sub(st))
	g.Expect(deploy.Sub(st) < time.Second*30).To(Equal(true))
	logrus.Infof("Pods Cleanup time: %v", destroy.Sub(deploy))
	g.Expect(destroy.Sub(deploy) < time.Second*30).To(Equal(true))
}
