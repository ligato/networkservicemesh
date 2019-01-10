#!/bin/bash

# Copyright (c) 2016-2017 Bitnami
# Copyright (c) 2018 Cisco and/or its affiliates.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -xe

KUBECTL_VERSION=v1.11.3

# Install kubectl
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/"${KUBECTL_VERSION}"/bin/linux/amd64/kubectl && \
 	chmod +x "kubectl" && sudo mv "kubectl" /usr/local/bin/

kubectl get nodes -o wide
kubectl version
kubectl api-versions
kubectl label --overwrite --all=true nodes app=nsmd-ds

kubectl apply -f k8s/conf/cluster-role-admin.yaml
kubectl apply -f k8s/conf/cluster-role-binding.yaml

# vim: sw=4 ts=4 et si
