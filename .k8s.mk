# Copyright (c) 2018 Cisco and/or its affiliates.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at:
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

K8S_CONF_DIR = k8s/conf

# Need nsmdp and icmp-responder-nse here as well, but missing yaml files
DEPLOY_TRACING = jaeger
DEPLOY_NSM = nsmd vppagent-dataplane
DEPLOY_MONITOR = crossconnect-monitor skydive
DEPLOY_INFRA = $(DEPLOY_TRACING) $(DEPLOY_NSM) $(DEPLOY_MONITOR)
DEPLOY_ICMP_KERNEL = icmp-responder-nse nsc
DEPLOY_ICMP_VPP = vppagent-icmp-responder-nse vppagent-nsc
DEPLOY_ICMP = $(DEPLOY_ICMP_KERNEL) $(DEPLOY_ICMP_VPP)
DEPLOY_VPN = secure-intranet-connectivity vppagent-firewall-nse vpn-gateway-nse vpn-gateway-nsc
DEPLOYS = $(DEPLOY_INFRA) $(DEPLOY_ICMP) $(DEPLOY_ICMP_VPP) $(DEPLOY_VPN)

CLUSTER_CONFIGS = cluster-role-admin cluster-role-binding cluster-role-view

# All of the rules that use vagrant are intentionally written in such a way
# That you could set the CLUSTER_RULES_PREFIX different and introduce
# a new platform to run on with k8s by adding a new include ${method}.mk
# and setting CLUSTER_RULES_PREFIX to a different value
ifeq ($(CLUSTER_RULES_PREFIX),)
CLUSTER_RULES_PREFIX := vagrant
endif
include .vagrant.mk
include .packet.mk

# .null.mk allows you to skip the vagrant machinery with:
# export CLUSTER_RULES_PREFIX=null
# before running make
include .null.mk

# Pull in docker targets
CONTAINER_BUILD_PREFIX = docker
include .docker.mk

.PHONY: k8s-deploy
k8s-deploy: k8s-delete $(addsuffix -deploy,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-infra-deploy
k8s-infra-deploy: k8s-infra-delete $(addsuffix -deploy,$(addprefix k8s-,$(DEPLOY_INFRA)))

.PHONY: k8s-icmp-deploy
k8s-icmp-deploy: k8s-icmp-delete $(addsuffix -deploy,$(addprefix k8s-,$(DEPLOY_ICMP)))

.PHONY: k8s-vpn-deploy
k8s-vpn-deploy: k8s-vpn-delete $(addsuffix -deploy,$(addprefix k8s-,$(DEPLOY_VPN)))

.PHONY: k8s-redeploy
k8s-redeploy: k8s-delete $(addsuffix -deployonly,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-deployonly
k8s-deployonly: $(addsuffix -deployonly,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-jaeger-deploy
k8s-jaeger-deploy:  k8s-start k8s-config k8s-jaeger-delete
	@until ! $$(kubectl get pods | grep -q ^jaeger ); do echo "Wait for jaeger to terminate"; sleep 1; done
	@sed "s;\(image:[ \t]*networkservicemesh/[^:]*\).*;\1$${COMMIT/$${COMMIT}/:$${COMMIT}};" ${K8S_CONF_DIR}/jaeger.yaml | kubectl apply -f -

.PHONY: k8s-%-deploy
k8s-%-deploy:  k8s-start k8s-config k8s-%-delete k8s-%-load-images
	@until ! $$(kubectl get pods | grep -q ^$* ); do echo "Wait for $* to terminate"; sleep 1; done
	@sed "s;\(image:[ \t]*networkservicemesh/[^:]*\).*;\1$${COMMIT/$${COMMIT}/:$${COMMIT}};" ${K8S_CONF_DIR}/$*.yaml | kubectl apply -f -


.PHONY: k8s-%-deployonly
k8s-%-deployonly:
	@until ! $$(kubectl get pods | grep -q ^$* ); do echo "Wait for $* to terminate"; sleep 1; done
	@sed "s;\(image:[ \t]*networkservicemesh/[^:]*\).*;\1$${COMMIT/$${COMMIT}/:$${COMMIT}};" ${K8S_CONF_DIR}/$*.yaml | kubectl apply -f -

.PHONY: k8s-delete
k8s-delete: $(addsuffix -delete,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-infra-delete
k8s-infra-delete: $(addsuffix -delete,$(addprefix k8s-,$(DEPLOY_INFRA)))

.PHONY: k8s-icmp-delete
k8s-icmp-delete: $(addsuffix -delete,$(addprefix k8s-,$(DEPLOY_ICMP)))

.PHONY: k8s-vpn-delete
k8s-vpn-delete: $(addsuffix -delete,$(addprefix k8s-,$(DEPLOY_VPN)))

.PHONY: k8s-%-delete
k8s-%-delete:
	@echo "Deleting ${K8S_CONF_DIR}/$*.yaml"
	@kubectl delete -f ${K8S_CONF_DIR}/$*.yaml > /dev/null 2>&1 || echo "$* does not exist and thus cannot be deleted"

.PHONY: k8s-load-images
k8s-load-images: $(addsuffix -load-images,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-%-load-images
k8s-%-load-images:  k8s-start $(CLUSTER_RULES_PREFIX)-%-load-images
	@echo "Delegated to $(CLUSTER_RULES_PREFIX)-$*-load-images"

.PHONY: k8s-%-config
k8s-%-config:  k8s-start
	@kubectl apply -f ${K8S_CONF_DIR}/$*.yaml

.PHONY: k8s-config
k8s-config: $(addsuffix -config,$(addprefix k8s-,$(CLUSTER_CONFIGS)))

.PHONY: k8s-start
k8s-start: $(CLUSTER_RULES_PREFIX)-start

.PHONY: k8s-start
k8s-restart: $(CLUSTER_RULES_PREFIX)-restart

.PHONY: k8s-build
k8s-build: $(addsuffix -build,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-jaeger-build
k8s-jaeger-build:

.PHONY: k8s-jaeger-save
k8s-jaeger-save:

.PHONY: k8s-jaeger-load-images
k8s-jaeger-load-images:

.PHONY: k8s-save
k8s-save: $(addsuffix -save,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-save-deploy
k8s-save-deploy: k8s-delete $(addsuffix -save-deploy,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-%-save-deploy
k8s-%-save-deploy:  k8s-start k8s-config k8s-%-save  k8s-%-load-images
	sed "s;\(image:[ \t]*networkservicemesh/[^:]*\).*;\1$${COMMIT/$${COMMIT}/:$${COMMIT}};" ${K8S_CONF_DIR}/$*.yaml | kubectl apply -f -

NSMD_CONTAINERS = nsmd nsmdp nsmd-k8s
.PHONY: k8s-nsmd-build
k8s-nsmd-build:  $(addsuffix -build,$(addprefix ${CONTAINER_BUILD_PREFIX}-,$(NSMD_CONTAINERS)))

.PHONY: k8s-nsmd-save
k8s-nsmd-save:  $(addsuffix -save,$(addprefix ${CONTAINER_BUILD_PREFIX}-,$(NSMD_CONTAINERS)))

.PHONY: k8s-nsmd-load-images
k8s-nsmd-load-images:  k8s-start $(addsuffix -load-images,$(addprefix ${CLUSTER_RULES_PREFIX}-,$(NSMD_CONTAINERS)))

VPPAGENT_DATAPLANE_CONTAINERS = vppagent-dataplane
.PHONY: k8s-vppagent-dataplane-build
k8s-vppagent-dataplane-build:  $(addsuffix -build,$(addprefix ${CONTAINER_BUILD_PREFIX}-,$(VPPAGENT_DATAPLANE_CONTAINERS)))
 .PHONY: k8s-vppagent-dataplane-save
k8s-vppagent-dataplane-save:  $(addsuffix -save,$(addprefix ${CONTAINER_BUILD_PREFIX}-,$(VPPAGENT_DATAPLANE_CONTAINERS)))
 .PHONY: k8s-vppagent-dataplane-load-images
k8s-vppagent-dataplane-load-images:  k8s-start $(addsuffix -load-images,$(addprefix ${CLUSTER_RULES_PREFIX}-,$(VPPAGENT_DATAPLANE_CONTAINERS)))

.PHONY: k8s-secure-intranet-connectivity-build
k8s-secure-intranet-connectivity-build:

.PHONY: k8s-secure-intranet-connectivity-save
k8s-secure-intranet-connectivity-save:

.PHONY: k8s-secure-intranet-connectivity-load-images
k8s-secure-intranet-connectivity-load-images:
	@echo "Wait for nsmd to register the resources"
	@sleep 10

.PHONY: k8s-skydive-build
k8s-skydive-build:

.PHONY: k8s-skydive-save
k8s-skydive-save: k8s-skydive-build

.PHONY: k8s-skydive-load-images
k8s-skydive-load-images:

.PHONY: k8s-vpn-gateway-nse-build
k8s-vpn-gateway-nse-build:

.PHONY: k8s-vpn-gateway-nse-save
k8s-vpn-gateway-nse-save:

.PHONY: k8s-vpn-gateway-nse-load-images
k8s-vpn-gateway-nse-load-images: k8s-icmp-responder-nse-load-images

.PHONY: k8s-vpn-gateway-nsc-build
k8s-vpn-gateway-nsc-build:

.PHONY: k8s-vpn-gateway-nsc-save
k8s-vpn-gateway-nsc-save:

.PHONY: k8s-vpn-gateway-nsc-load-images
k8s-vpn-gateway-nsc-load-images: k8s-nsc-load-images

.PHONY: k8s-nsc-build
k8s-nsc-build:  ${CONTAINER_BUILD_PREFIX}-nsc-build

.PHONY: k8s-nsc-save
k8s-nsc-save:  ${CONTAINER_BUILD_PREFIX}-nsc-save


.PHONY: k8s-icmp-responder-nse-build
k8s-icmp-responder-nse-build:  ${CONTAINER_BUILD_PREFIX}-icmp-responder-nse-build

.PHONY: k8s-icmp-responder-nse-save
k8s-icmp-responder-nse-save:  ${CONTAINER_BUILD_PREFIX}-icmp-responder-nse-save

.PHONY: k8s-vppagent-icmp-responder-nse-build
k8s-vppagent-icmp-responder-nse-build:  ${CONTAINER_BUILD_PREFIX}-vppagent-icmp-responder-nse-build

.PHONY: k8s-vppagent-icmp-responder-nse-save
k8s-vppagent-icmp-responder-nse-save:  ${CONTAINER_BUILD_PREFIX}-vppagent-icmp-responder-nse-save

.PHONY: k8s-vppagent-firewall-nse-build
k8s-vppagent-firewall-nse-build:  ${CONTAINER_BUILD_PREFIX}-vppagent-firewall-nse-build

.PHONY: k8s-vppagent-firewall-nse-save
k8s-vppagent-firewall-nse-save:  ${CONTAINER_BUILD_PREFIX}-vppagent-firewall-nse-save

.PHONY: k8s-vppagent-nsc-build
k8s-vppagent-nsc-build:  ${CONTAINER_BUILD_PREFIX}-vppagent-nsc-build

.PHONY: k8s-vppagent-nsc-save
k8s-vppagent-nsc-save:  ${CONTAINER_BUILD_PREFIX}-vppagent-nsc-save


.PHONY: k8s-crossconnect-monitor-build
k8s-crossconnect-monitor-build: ${CONTAINER_BUILD_PREFIX}-crossconnect-monitor-build

.PHONY: k8s-crossconnect-monitor-save
k8s-crossconnect-monitor-save: ${CONTAINER_BUILD_PREFIX}-crossconnect-monitor-save

.PHONY: k8s-crossconnect-load-images
k8s-crossconnect-monitor-load-images:  k8s-start $(addsuffix -load-images,$(addprefix ${CLUSTER_RULES_PREFIX}-,crossconnect-monitor))


.PHONY: k8s-skydive-build
k8s-skydive-build:

.PHONY: k8s-skydive-save
k8s-skydive-save: k8s-skydive-build

.PHONY: k8s-skydive-load-images
k8s-skydive-load-images:

.PHONY: k8s-vpn-gateway-nse-build
k8s-vpn-gateway-nse-build:

.PHONY: k8s-vpn-gateway-nse-save
k8s-vpn-gateway-nse-save:

.PHONY: k8s-vpn-gateway-nse-load-images
k8s-vpn-gateway-nse-load-images: k8s-icmp-responder-nse-load-images

.PHONY: k8s-vpn-gateway-nsc-build
k8s-vpn-gateway-nsc-build:

.PHONY: k8s-vpn-gateway-nsc-save
k8s-vpn-gateway-nsc-save:

.PHONY: k8s-vpn-gateway-nsc-load-images
k8s-vpn-gateway-nsc-load-images: k8s-nsc-load-images

# TODO add k8s-%-logs and k8s-logs to capture all the logs from k8s

.PHONY: k8s-logs
k8s-logs: $(addsuffix -logs,$(addprefix k8s-,$(DEPLOYS)))

.PHONY: k8s-%logs
k8s-%-logs:
	@echo "K8s logs for $*"
	@for pod in $$(kubectl get pods --all-namespaces | grep $* | awk '{print $$2}');do \
		echo '******************************************'; \
		echo "Logs: $${pod}:"; \
		kubectl logs $${pod} || true; \
		kubectl logs -p $${pod} || true; \
		echo '>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>'; \
		echo "Network information for $${pod}"; \
		kubectl exec -ti $${pod} ip addr; \
		kubectl exec -ti $${pod} ip neigh; \
		if [[ "$${pod}" == *"vppagent"* ]]; then \
			echo "vpp information for $${pod}"; \
			kubectl exec -it $${pod} vppctl show int; \
			kubectl exec -it $${pod} vppctl show int addr; \
			kubectl exec -it $${pod} vppctl show vxlan tunnel; \
			kubectl exec -it $${pod} vppctl show memif; \
		fi; \
	done

.PHONY: k8s-nsmd-logs
k8s-nsmd-logs:
	@echo "K8s logs for nsmds"
	@echo '******************************************'
	@for pod in $$(kubectl get pods --all-namespaces | grep nsmd | awk '{print $$2}'); do \
		for container in nsmd nsmdp nsmd-k8s; do \
			echo '------------------------------------------'; \
			echo "K8s logs for $${pod} container $${container}"; \
			kubectl logs $${pod} --container $${container} || true; \
			kubectl logs -p $${pod} --container $${container} || true ;\
			echo '>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>'; \
			echo 'NSMD Network information'; \
			kubectl exec -ti $${pod} --container $${container} ip addr; \
		done \
	done

.PHONY: k8s-%-debug
k8s-%-debug:
	@echo "Debugging $*"
	@kubectl exec -ti $$(kubectl get pods | grep $*- | cut -d \  -f1) /go/src/github.com/networkservicemesh/networkservicemesh/scripts/debug.sh $*

.PHONY: k8s-nsmd-debug
k8s-nsmd-debug:
	@kubectl exec -ti $(pod) -c nsmd /go/src/github.com/networkservicemesh/networkservicemesh/scripts/debug.sh nsmd

.PHONY: k8s-forward
k8s-forward:
	@echo "Forwarding local $(port1) to $(port2) for $(pod)"
	@kubectl port-forward $$(kubectl get pods | grep $(pod) | cut -d \  -f1) $(port1):$(port2)

.PHONY: k8s-check
k8s-check:
	./scripts/nsc_ping_all.sh
	./scripts/verify_vpn_gateway.sh

.PHONY: k8s-terminating-cleanup
k8s-terminating-cleanup:
	@kubectl get pods -o wide |grep Terminating | cut -d \  -f 1 | xargs kubectl delete pods --force --grace-period 0 {}

.PHONE: k8s-kublet-restart
k8s-kublet-restart: vagrant-kublet-restart

.PHONE: k8s-pods
k8s-pods:
	@kubectl get pods -o wide

.PHONY: k8s-nsmd-master-tlogs
k8s-nsmd-master-tlogs:
	@kubectl logs -f $$(kubectl get pods -o wide | grep kube-master | grep nsmd | cut -d\  -f1) -c nsmd

.PHONY: k8s-nsmd-worker-tlogs
k8s-nsmd-worker-tlogs:
	@kubectl logs -f $$(kubectl get pods -o wide | grep kube-worker | grep nsmd | cut -d\  -f1) -c nsmd





