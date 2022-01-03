#!/usr/bin/env bash

# Copyright 2022 The KCP Authors.
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

CURRENT_DIR="$(pwd)"
DEMO_ROOT="$(dirname "${BASH_SOURCE}")"
KCP_ROOT="$(cd ${DEMO_ROOT}/../.. && pwd)"
export KCP_DATA_ROOT=${KCP_DATA_ROOT:-$KCP_ROOT}

source ${DEMO_ROOT}/.startUtils
setupTraps $0

KUBECONFIG=${KCP_DATA_ROOT}/.kcp/admin.kubeconfig
export KCP_LISTEN_ADDR="127.0.0.1:6443"

${DEMO_ROOT}/startKcpAndClusterController.sh --auto-publish-apis=true --resources-to-sync deployments.apps &
KCP_PID=$!

wait_command "grep 'Serving securely' ${CURRENT_DIR}/kcp.log"

echo ""
echo "Starting Deployment Splitter"
${KCP_ROOT}/bin/deployment-splitter -kubeconfig=${KUBECONFIG} &> deployment-splitter.log &
SPLIT_PID=$!
echo "Deployment Splitter started: $SPLIT_PID" 

echo ""
echo "Use ctrl-C to stop all components"
echo ""

tail -f deployment-splitter.log  &

wait
