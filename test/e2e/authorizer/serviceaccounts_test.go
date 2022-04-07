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

package authorizer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"

	"github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1/helper"
	"github.com/kcp-dev/kcp/test/e2e/framework"
)

// Testing ... testing
func TestKindServiceAccounts(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	cfg := framework.KindConfig{
		Name: t.Name(),
	}
	kindClusters := framework.NewKindFixture(t, cfg)

	kindClusterConfig := kindClusters.Servers[cfg.Name].DefaultConfig(t)
	kubeClusterClient := kubernetes.NewForConfigOrDie(kindClusterConfig)

	namespacesList, err := kubeClusterClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	require.NoError(t, err)

	for _, namespace := range namespacesList.Items {
		t.Log(namespace.Name)
	}
}

func TestServiceAccounts(t *testing.T) {
	t.Parallel()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	server := framework.SharedKcpServer(t)
	orgClusterName := framework.NewOrganizationFixture(t, server)
	clusterName := framework.NewWorkspaceFixture(t, server, orgClusterName, "Universal")

	cfg := server.DefaultConfig(t)
	kubeClusterClient, err := kubernetes.NewClusterForConfig(cfg)
	require.NoError(t, err)

	kubeClient := kubeClusterClient.Cluster(clusterName)

	t.Log("Creating namespace")
	namespace, err := kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "e2e-sa-",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err, "failed to create namespace")

	t.Log("Creating role to access configmaps")
	_, err = kubeClient.RbacV1().Roles(namespace.Name).Create(ctx, &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sa-access-configmap",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err, "failed to create role")

	t.Log("Creating role binding to access configmaps")
	_, err = kubeClient.RbacV1().RoleBindings(namespace.Name).Create(ctx, &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sa-access-configmap",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: namespace.Name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "sa-access-configmap",
			APIGroup: rbacv1.GroupName,
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err, "failed to create role")

	t.Log("Waiting for service account to be created")
	require.Eventually(t, func() bool {
		_, err := kubeClient.CoreV1().ServiceAccounts(namespace.Name).Get(ctx, "default", metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return false
		} else if err != nil {
			t.Fatalf("unexpected error retrieving service account: %v", err)
		}
		return true
	}, wait.ForeverTestTimeout, time.Millisecond*100, "\"default\" service account not created in namespace %s",
		helper.QualifiedObjectName(namespace),
	)

	t.Log("Waiting for service account secret to be created")
	var tokenSecret corev1.Secret
	require.Eventually(t, func() bool {
		secrets, err := kubeClient.CoreV1().Secrets(namespace.Name).List(ctx, metav1.ListOptions{})
		require.NoError(t, err, "failed to list secrets")

		for _, secret := range secrets.Items {
			if secret.Annotations[corev1.ServiceAccountNameKey] == "default" {
				tokenSecret = secret
				return true
			}
		}
		return false
	}, wait.ForeverTestTimeout, time.Millisecond*100, "token secret for default service account not created")

	testCases := []struct {
		name  string
		token func(t *testing.T) string
	}{
		{"Legacy token", func(t *testing.T) string {
			return string(tokenSecret.Data["token"])
		}},
		{"Bound service token", func(t *testing.T) string {
			t.Log("Creating service account bound token")
			boundToken, err := kubeClient.CoreV1().ServiceAccounts(namespace.Name).CreateToken(ctx, "default", &authenticationv1.TokenRequest{
				Spec: authenticationv1.TokenRequestSpec{
					Audiences:         []string{"https://kcp.default.svc"},
					ExpirationSeconds: pointer.Int64Ptr(3600),
					BoundObjectRef: &authenticationv1.BoundObjectReference{
						APIVersion: "v1",
						Kind:       "Secret",
						Name:       tokenSecret.Name,
						UID:        tokenSecret.UID,
					},
				},
			}, metav1.CreateOptions{})
			require.NoError(t, err, "failed to create token")
			return boundToken.Status.Token
		}},
	}
	for _, ttc := range testCases {
		t.Run(ttc.name, func(t *testing.T) {
			saRestConfig := server.DefaultConfig(t)
			saRestConfig.BearerToken = ttc.token(t)
			saKubeClusterClient, err := kubernetes.NewClusterForConfig(saRestConfig)
			require.NoError(t, err)

			t.Run("Access workspace with the service account", func(t *testing.T) {
				_, err := saKubeClusterClient.Cluster(clusterName).CoreV1().ConfigMaps(namespace.Name).List(ctx, metav1.ListOptions{})
				require.NoError(t, err)
			})

			t.Run("Access another workspace in the same org", func(t *testing.T) {
				t.Log("Create namespace with the same name ")
				otherClusterName := framework.NewWorkspaceFixture(t, server, orgClusterName, "Universal")
				_, err := kubeClusterClient.Cluster(otherClusterName).CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: namespace.Name,
					},
				}, metav1.CreateOptions{})
				require.NoError(t, err, "failed to create namespace in other workspace")

				t.Log("Accessing workspace with the service account")
				obj, err := saKubeClusterClient.Cluster(otherClusterName).CoreV1().ConfigMaps(namespace.Name).List(ctx, metav1.ListOptions{})
				require.Error(t, err, fmt.Sprintf("expected error accessing workspace with the service account, got: %v", obj))
			})

			t.Run("Access an equally named workspace in another org", func(t *testing.T) {
				t.Log("Create namespace with the same name")
				otherOrgClusterName := framework.NewOrganizationFixture(t, server)
				otherClusterName := framework.NewWorkspaceFixture(t, server, otherOrgClusterName, "Universal")
				_, err := kubeClusterClient.Cluster(otherClusterName).CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: namespace.Name,
					},
				}, metav1.CreateOptions{})
				require.NoError(t, err, "failed to create namespace in other workspace")

				t.Log("Accessing workspace with the service account")
				obj, err := saKubeClusterClient.Cluster(otherClusterName).CoreV1().ConfigMaps(namespace.Name).List(ctx, metav1.ListOptions{})
				require.Error(t, err, fmt.Sprintf("expected error accessing workspace with the service account, got: %v", obj))
			})
		})
	}
}
