package service

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const (
	clusterLabel = "kcp.dev/cluster"
	ownedByLabel = "kcp.dev/owned-by"
)

func (c *Controller) reconcile(ctx context.Context, service *corev1.Service) error {
	klog.Infof("reconciling service %q", service.Name)

	if service.Labels == nil || service.Labels[clusterLabel] == "" {
		// This is a root service; get its leafs.
		sel, err := labels.Parse(fmt.Sprintf("%s=%s", ownedByLabel, service.Name))
		if err != nil {
			return err
		}
		leafs, err := c.lister.List(sel)
		if err != nil {
			return err
		}

		if len(leafs) == 0 {
			if err := c.createLeafs(ctx, service); err != nil {
				return err
			}
		}

	} else if service.Labels[ownedByLabel] != "" {
		rootServiceName := service.Labels[ownedByLabel]
		// A leaf service was updated; get others and aggregate status.
		sel, err := labels.Parse(fmt.Sprintf("%s=%s", ownedByLabel, rootServiceName))
		if err != nil {
			return err
		}
		others, err := c.lister.List(sel)
		if err != nil {
			return err
		}

		var rootService *corev1.Service

		rootIf, exists, err := c.indexer.Get(&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:   service.Namespace,
				Name:        rootServiceName,
				ClusterName: service.GetClusterName(),
			},
		})
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("Root Service not found: %s", rootServiceName)
		}

		rootService = rootIf.(*corev1.Service)

		// Aggregate .status from all leafs.

		rootService = rootService.DeepCopy()

		// Cheat and set the root .status.conditions to the first leaf's .status.conditions.
		// TODO: do better.
		if len(others) > 0 {
			rootService.Status.Conditions = others[0].Status.Conditions
		}

		if _, err := c.client.Services(rootService.Namespace).UpdateStatus(ctx, rootService, metav1.UpdateOptions{}); err != nil {
			if errors.IsConflict(err) {
				key, err := cache.MetaNamespaceKeyFunc(service)
				if err != nil {
					return err
				}
				c.queue.AddRateLimited(key)
				return nil
			}
			return err
		}
	}

	return nil
}

func (c *Controller) createLeafs(ctx context.Context, root *corev1.Service) error {
	cls, err := c.clusterLister.List(labels.Everything())
	if err != nil {
		return err
	}

	// No clusters; nothing to do.
	if len(cls) == 0 {
		return nil
	}

	for _, cl := range cls {
		vd := root.DeepCopy()
		vd.Name = fmt.Sprintf("%s--%s", root.Name, cl.Name)

		if vd.Labels == nil {
			vd.Labels = map[string]string{}
		}
		vd.Labels[clusterLabel] = cl.Name
		vd.Labels[ownedByLabel] = root.Name

		// Set OwnerReference so deleting the Deployment deletes all virtual deployments.
		vd.OwnerReferences = []metav1.OwnerReference{{
			APIVersion: "core/v1",
			Kind:       "Service",
			UID:        root.UID,
			Name:       root.Name,
		}}

		// TODO: munge namespace
		vd.SetResourceVersion("")
		if _, err := c.kubeClient.CoreV1().Services(root.Namespace).Create(ctx, vd, metav1.CreateOptions{}); err != nil {
			return err
		}
		klog.Infof("created child service %q", vd.Name)
	}

	return nil
}
