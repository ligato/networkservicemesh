package health

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/onsi/gomega"

	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
)

func TestGrpcHealth(t *testing.T) {
	_ = os.Setenv(tools.InsecureEnv, "true")
	assert := gomega.NewWithT(t)
	_ = os.Remove("soc")
	sock, err := net.Listen("unix", "soc")
	assert.Expect(err).Should(gomega.BeNil())
	server := tools.NewServer(context.Background())
	health := NewGrpcHealth(server, tools.NewAddr("unix", "soc"), time.Second)
	go func() {
		serveErr := server.Serve(sock)
		if serveErr != nil {
			logrus.Error(serveErr.Error())
		}
	}()
	err = health.Check()
	assert.Expect(err).Should(gomega.BeNil())
	_ = os.Remove("soc")
}
