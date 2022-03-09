package mutators

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ServiceAccountMutator struct {
	fromConfig          *rest.Config
	toConfig            *rest.Config
	registeredWorkspace string
}

func (sam *ServiceAccountMutator) getGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "serviceaccounts",
	}
}

func (sam *ServiceAccountMutator) Register(mutators map[schema.GroupVersionResource]Mutator) {
	if _, ok := mutators[sam.getGVR()]; !ok {
		mutators[sam.getGVR()] = sam
	}
}

func NewServiceAccountMutator(fromConfig, toConfig *rest.Config, registeredWorkspace string) *ServiceAccountMutator {
	return &ServiceAccountMutator{
		fromConfig:          fromConfig,
		toConfig:            toConfig,
		registeredWorkspace: registeredWorkspace,
	}
}

func (sam *ServiceAccountMutator) ApplyDownstreamName(upstreamObj *unstructured.Unstructured) error {
	var serviceAccount corev1.ServiceAccount
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(
		upstreamObj.UnstructuredContent(),
		&serviceAccount)

	if err != nil {
		return err
	}

	originalName, ok := serviceAccount.Labels[originalNameLabel]

	// If the original name label is not present, that means we are syncing from kcp to workloadcluster
	// let's set the original name label to the name of the serviceaccount, and do any transformation required.
	if !ok || originalName == "" {
		// Set the original label
		serviceAccount.Labels[originalNameLabel] = serviceAccount.Name

		if serviceAccount.Name == "default" {
			serviceAccount.Name = "kcp-default"
		}
	} else {
		// If the original name label is present, that means we are syncing from workloadcluster to kcp
		// let's restore the original name to the object.
		serviceAccount.Name = originalName
	}

	// The default service account should be translated to kcp-default to avoid it clashing with the cluster default sa.

	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&serviceAccount)
	if err != nil {
		return err
	}

	// Set the changes back into the obj.
	upstreamObj.SetUnstructuredContent(unstructured)
	return nil
}

// ApplyStatus makes modifications to the Status of the deployment object.
func (sam *ServiceAccountMutator) ApplyStatus(unstrob *unstructured.Unstructured) error {
	// No transformations
	return nil
}

// ApplySpec makes modifications to the Spec of the deployment object.
func (sam *ServiceAccountMutator) ApplySpec(unstrob *unstructured.Unstructured) error {
	// No transformations
	return nil
}
