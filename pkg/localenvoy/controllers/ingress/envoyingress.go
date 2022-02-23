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

package ingress

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	clusterLabel     = "kcp.dev/cluster"
	toEnvoyLabel     = "ingress.kcp.dev/envoy"
	ownedByCluster   = "ingress.kcp.dev/owned-by-cluster"
	ownedByIngress   = "ingress.kcp.dev/owned-by-ingress"
	ownedByNamespace = "ingress.kcp.dev/owned-by-namespace"
)

// reconcile is triggered on every change to an ingress resource, or it's associated services (by tracker).
func (c *Controller) reconcile(ctx context.Context, ingress *networkingv1.Ingress) error {
	klog.InfoS("reconciling Ingress", "ClusterName", ingress.ClusterName, "Namespace", ingress.Namespace, "Name", ingress.Name)

	if ingress.Labels[clusterLabel] == "" {
		// this is the root. Ignore.
		return nil
	}

	ingressRootKey := rootIngressKeyFor(ingress)
	obj, exists, err := c.ingressIndexer.GetByKey(ingressRootKey)
	if err != nil {
		klog.Warningf("failed to get root ingress: %v", err)
		return nil
	}
	if !exists {
		// TODO(jmprusi): A leaf without rootIngress? use OwnerRefs to avoid this.
		// TODO(jmprusi): Add user-facing condition to leaf.
		klog.Warningf("root Ingress not found %s", ingressRootKey)
		return nil
	}

	rootIngress := obj.(*networkingv1.Ingress).DeepCopy()
	rootIngress.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{
		Hostname: generateStatusHost(c.domain, rootIngress),
	}}

	// Label the received ingress for envoy, as we want the controlplane to use this leaf
	// for updating the envoy config.
	ingress.Labels[toEnvoyLabel] = "true"

	// Update the rootIngress status with our desired LB.
	if _, err := c.client.Cluster(rootIngress.ClusterName).NetworkingV1().Ingresses(rootIngress.Namespace).UpdateStatus(ctx, rootIngress, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("failed to update root ingress status: %w", err)
	}

	return nil
}

// TODO(jmprusi): Review the hash algorithm.
func domainHashString(s string) string {
	h := fnv.New32a()
	// nolint: errcheck
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum32())
}

// generateStatusHost returns a string that represent the desired status hostname for the ingress.
// If the host is part of the same domain, it will be preserved as the status hostname, if not
// a new one will be generated based on a hash of the ingress name, namespace and clusterName.
func generateStatusHost(domain string, ingress *networkingv1.Ingress) string {
	// TODO(jmprusi): using "contains" is a bad idea as it could be abused by crafting a malicious hostname, but for a PoC it should be good enough?
	allRulesAreDomain := true
	for _, rule := range ingress.Spec.Rules {
		if !strings.Contains(rule.Host, domain) {
			allRulesAreDomain = false
			break
		}
	}

	//TODO(jmprusi): Hardcoded to the first one...
	if allRulesAreDomain {
		return ingress.Spec.Rules[0].Host
	}

	return domainHashString(ingress.Name+ingress.Namespace+ingress.ClusterName) + "." + domain
}
