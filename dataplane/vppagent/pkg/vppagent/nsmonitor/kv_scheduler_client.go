package nsmonitor

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	kvSchedulerPort   = 9191
	kvSchedulerPrefix = "/scheduler"
	downStreamResync  = "/downstream-resync"
)

// KVSchedulerClient - Client to vpp KVScheduler server.
type KVSchedulerClient interface {
	//DownstreamResync - Calls downstream-resync in KVScheduler
	DownstreamResync()
}

type kvSchedulerClient struct {
	httpClient          http.Client
	kvSchedulerEndpoint string
}

//NewKVSchedulerClient - Creates a new client for KVScheduelr. Can return an error if vppAgentEndpoint has an incorrect format.
func NewKVSchedulerClient(vppAgentEndpoint string) (KVSchedulerClient, error) {
	kvSchedulerEndpoint, err := buildKvSchedulerDownStreamPath(vppAgentEndpoint)
	if err != nil {
		return nil, err
	}
	return &kvSchedulerClient{
		kvSchedulerEndpoint: kvSchedulerEndpoint,
	}, nil
}

func (c *kvSchedulerClient) DownstreamResync() {
	downSteamResyncPath := c.kvSchedulerEndpoint + kvSchedulerPrefix + downStreamResync
	request, err := http.NewRequest("POST", downSteamResyncPath, nil)
	if err != nil {
		logrus.Errorf("kvSchedulerClient:, can't create request %v", err)
	}
	resp, err := c.httpClient.Do(request)
	if err != nil {
		logrus.Errorf("kvSchedulerClient:, can't do request %v, error: %v", resp, err)
	}
	err = resp.Body.Close()
	if err != nil {
		logrus.Errorf("kvSchedulerClient:, can't close response body: %v", err)
	}
	logrus.Infof("kvSchedulerClient: response %v from %v", resp, downSteamResyncPath)
}

func buildKvSchedulerDownStreamPath(vppAgentEndpoint string) (string, error) {
	parts := strings.Split(vppAgentEndpoint, ":")
	serverURL := fmt.Sprintf("http://%v:%v", parts[0], kvSchedulerPort)
	_, err := url.Parse(vppAgentEndpoint)
	if err != nil {
		return "", err
	}
	return serverURL, nil
}
