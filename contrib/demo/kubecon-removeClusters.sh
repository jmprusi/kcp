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

DEMO_ROOT="$(dirname "${BASH_SOURCE}")"
CLUSTERS_DIR=${DEMO_ROOT}/clusters/kind

kubectl --kubeconfig=${CLUSTERS_DIR}/us-west1.kubeconfig get deployment my-deployment--us-west1 -o yaml -n demo >>us-west1.log 2>&1
kubectl --kubeconfig=${CLUSTERS_DIR}/us-west1.kubeconfig get pods -o wide -n demo >>us-west1.log 2>&1
kubectl --kubeconfig=${CLUSTERS_DIR}/us-east1.kubeconfig get deployment my-deployment--us-east1 -o yaml -n demo >>us-east1.log 2>&1
kubectl --kubeconfig=${CLUSTERS_DIR}/us-east1.kubeconfig get pods -o wide -n demo >>us-east1.log 2>&1

kind delete clusters us-west1 us-west1 us-east1 > /dev/null || true
