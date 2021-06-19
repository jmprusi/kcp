#!/bin/bash
#
# Copyright 2021 Red Hat, Inc.
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
#
set -euo pipefail

export GOROOT=$(go env GOROOT)
export KUADRANT_NAMESPACE="kuadrant-system"
export KIND_BIN="./bin/kind"
TEMP_DIR="./tmp"

# TODO(jmprusi): Hardcoded, perhaps better to make this configurable.
KIND_CLUSTER_A="kcp-cluster-a"
KIND_CLUSTER_B="kcp-cluster-b"

mkdir -p ${TEMP_DIR}


# TODO(jmprusi): Split this setup into up/clean actions.
echo "Deleting any previous kind clusters."
{
  ${KIND_BIN} delete cluster --name ${KIND_CLUSTER_A}
  ${KIND_BIN} delete cluster --name ${KIND_CLUSTER_B}
} &> /dev/null

echo "Deploying two kind k8s clusters locally."
{
  ${KIND_BIN} create cluster --name ${KIND_CLUSTER_A}
  ${KIND_BIN} create cluster --name ${KIND_CLUSTER_B}
} &>/dev/null


echo "Exporting KUBECONFIG=.kcp/data/admin.kubeconfig"
export KUBECONFIG=.kcp/data/admin.kubeconfig

echo "Creating Cluster objects for each of the k8s cluster."

${KIND_BIN} get kubeconfig --name=${KIND_CLUSTER_A} | sed -e 's/^/    /' | cat contrib/examples/cluster.yaml - | sed -e "s/name: local/name: ${KIND_CLUSTER_A}/" > ${TEMP_DIR}/${KIND_CLUSTER_A}.yaml
${KIND_BIN} get kubeconfig --name=${KIND_CLUSTER_B} | sed -e 's/^/    /' | cat contrib/examples/cluster.yaml - | sed -e "s/name: local/name: ${KIND_CLUSTER_B}/" > ${TEMP_DIR}/${KIND_CLUSTER_B}.yaml

echo "The clusters are ready, and you can find the cluster objects in ./tmp/"
echo "To run KCP:"
echo ""
echo "./bin/kcp start --push_mode --install_cluster_controller"
echo ""
echo "And then you can register the clusters, remember to set the proper KUBECONFIG:"
echo ""
echo "export KUBECONFIG=.kcp/data/admin.kubeconfig"
echo "kubectl apply -f ./tmp"
echo ""












