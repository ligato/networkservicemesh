// Copyright (c) 2020 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build interdomain

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/networkservicemesh/networkservicemesh/applications/nsmrs/pkg/serviceregistryserver"
	"github.com/networkservicemesh/networkservicemesh/controlplane/pkg/nsmd"

	. "github.com/onsi/gomega"

	"github.com/networkservicemesh/networkservicemesh/k8s/pkg/proxyregistryserver"
	"github.com/networkservicemesh/networkservicemesh/test/kubetest"
	"github.com/networkservicemesh/networkservicemesh/test/kubetest/pods"
)

func TestFloatingInterdomainFewClientsWithDelay(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	g := NewWithT(t)
	clustersCount := 2
	k8ss := []*kubetest.ExtK8s{}
	for i := 0; i < clustersCount; i++ {
		kubeconfig := os.Getenv(fmt.Sprintf("KUBECONFIG_CLUSTER_%d", i+1))
		g.Expect(len(kubeconfig)).ToNot(Equal(0))

		k8s, err := kubetest.NewK8sForConfig(g, true, kubeconfig)
		g.Expect(err).To(BeNil())

		defer k8s.Cleanup()
		defer k8s.SaveTestArtifacts(t)

		k8ss = append(k8ss, &kubetest.ExtK8s{
			K8s:        k8s,
			NodesSetup: nil,
		})
	}

	nsmrsNode := &k8ss[clustersCount-1].K8s.GetNodesWait(2, defaultTimeout)[1]
	nsmrsPod := kubetest.DeployNSMRS(k8ss[clustersCount-1].K8s, nsmrsNode, "nsmrs", defaultTimeout, map[string]string{
		serviceregistryserver.NSEExpirationTimeoutEnv.Name(): "30s",
	})

	nsmrsExternalIP, err := kubetest.GetNodeExternalIP(nsmrsNode)
	if err != nil {
		nsmrsExternalIP, err = kubetest.GetNodeInternalIP(nsmrsNode)
		g.Expect(err).To(BeNil())
	}
	nsmrsInternalIP, err := kubetest.GetNodeInternalIP(nsmrsNode)
	g.Expect(err).To(BeNil())

	for i := 0; i < clustersCount; i++ {
		k8s := k8ss[i].K8s

		nodesSetup, err := kubetest.SetupNodesConfig(k8s, 1, defaultTimeout, []*pods.NSMgrPodConfig{
			{
				Variables: map[string]string{
					nsmd.NSETrackingIntervalSecondsEnv.Name(): "10s",
					nsmd.PublicAPIAddressEnv.Name():           "127.0.0.1:5001",
				},
				Namespace:          k8s.GetK8sNamespace(),
				ForwarderVariables: kubetest.DefaultForwarderVariables(k8s.GetForwardingPlane()),
			},
		}, k8s.GetK8sNamespace())
		g.Expect(err).To(BeNil())

		k8ss[i].NodesSetup = nodesSetup

		pnsmdName := fmt.Sprintf("pnsmgr-%s", nodesSetup[0].Node.Name)
		proxyNSMgrConfig := &pods.NSMgrPodConfig{
			Variables: pods.DefaultProxyNSMD(),
			Namespace: k8s.GetK8sNamespace(),
		}
		proxyNSMgrConfig.Variables[proxyregistryserver.NSMRSAddressEnv] = nsmrsInternalIP + ":80"
		proxyNSMgrConfig.Variables[nsmd.PublicAPIAddressEnv.Name()] = "127.0.0.1:5001"
		_, err = kubetest.DeployProxyNSMgrWithConfig(k8s, nodesSetup[0].Node, pnsmdName, defaultTimeout, proxyNSMgrConfig)
		g.Expect(err).To(BeNil())

		serviceCleanup := kubetest.RunProxyNSMgrService(k8s)
		defer serviceCleanup()
	}

	_ = kubetest.DeployICMP(k8ss[clustersCount-1].K8s, k8ss[clustersCount-1].NodesSetup[0].Node, "icmp-responder-nse-1", defaultTimeout)
	k8ss[clustersCount-1].K8s.WaitLogsContains(nsmrsPod, "nsmrs", "Registered NSE entry", defaultTimeout)

	nsc1 := kubetest.DeployNSCWithEnv(k8ss[0].K8s, k8ss[0].NodesSetup[0].Node, "nsc-1", defaultTimeout, map[string]string{
		"CLIENT_LABELS":          "app=icmp",
		"CLIENT_NETWORK_SERVICE": fmt.Sprintf("icmp-responder@%s", nsmrsExternalIP),
	})

	kubetest.CheckNSC(k8ss[0].K8s, nsc1)
	<-time.After(time.Second * 35)

	nsc2 := kubetest.DeployNSCWithEnv(k8ss[0].K8s, k8ss[0].NodesSetup[0].Node, "nsc-2", defaultTimeout, map[string]string{
		"CLIENT_LABELS":          "app=icmp",
		"CLIENT_NETWORK_SERVICE": fmt.Sprintf("icmp-responder@%s", nsmrsExternalIP),
	})
	kubetest.CheckNSC(k8ss[0].K8s, nsc2)
}
