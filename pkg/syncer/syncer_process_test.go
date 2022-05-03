/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package syncer

//
//import (
//	"github.com/kcp-dev/apimachinery/pkg/logicalcluster"
//	workloadv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/workload/v1alpha1"
//	"github.com/stretchr/testify/require"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
//	"k8s.io/apimachinery/pkg/runtime/schema"
//	"testing"
//	"time"
//)
//
//func TestSyncerProcess(t *testing.T) {
//	tests := map[string]struct {
//		fromNamespaces []*corev1.Namespace
//		fromResources  map[schema.GroupVersionResource][]*unstructured.Unstructured
//
//		direction SyncDirection
//		workloadClusterName string
//
//		listLocationsError        error
//		listAPIBindingsError      error
//		listWorkloadClustersError error
//		patchNamespaceError       error
//
//		wantError           bool
//	}{
//		"SpecSyncer without cluster state label": {
//			fromNamespaces: []*corev1.Namespace{
//				{
//					ObjectMeta: metav1.ObjectMeta{
//						Name:        "test",
//						ClusterName: "root:org:ws",
//						Labels: map[string]string{
//							"state.internal.workloads.kcp.dev/us-west1": "Sync",
//						},
//					},
//				},
//			},
//			fromResources: map[schema.GroupVersionResource][]*unstructured.Unstructured{
//				schema.GroupVersionResource{ Group: "apps", Version: "v1", Resource: "deployments"}: {
//					{
//						Object: map[string]interface{}{
//							"metadata": map[string]interface{}{
//
//							},
//							"spec": map[string]interface{}{
//
//							},
//						},
//					},
//				},
//,			},
//			wantError: false,
//		},
//		"SpecSyncer with empty cluster state label": {
//			fromNamespaces: []*corev1.Namespace{
//				{
//					ObjectMeta: metav1.ObjectMeta{
//						Name:        "test",
//						ClusterName: "root:org:ws",
//						Labels: map[string]string{
//							"state.internal.workloads.kcp.dev/us-west1": "Sync",
//						},
//					},
//				},
//			},
//			wantError: false,
//		},
//	}
//
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//			controller := Controller{
//				name:                string(direction) + "--" + kcpClusterName.String() + "--" + pcluster
//				pcluster:            "",
//				queue:               nil,
//				fromInformers:       nil,
//				fromClient:          nil,
//				toClient:            nil,
//				upsertFn:            nil,
//				deleteFn:            nil,
//				direction:           "",
//				upstreamClusterName: logicalcluster.LogicalCluster{},
//				mutators:            nil,
//			}
//
//			status, err := r.reconcile(context.Background(), tc.namespace)
//			if tc.wantError {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//			}
//
//			require.Equal(t, status, tc.wantReconcileStatus)
//			require.Equal(t, tc.wantRequeue, requeuedAfter)
//			require.Equal(t, tc.wantPatch, gotPatch)
//		})
//	}
//}
//}
