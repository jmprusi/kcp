package mutators

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ConfigMapMutator struct {
	fromConfig          *rest.Config
	toConfig            *rest.Config
	registeredWorkspace string
}

func (cmm *ConfigMapMutator) getGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
}

func (cmm *ConfigMapMutator) Register(mutators map[schema.GroupVersionResource]Mutator) {
	if _, ok := mutators[cmm.getGVR()]; !ok {
		mutators[cmm.getGVR()] = cmm
	}
}

func NewConfigMapMutator(fromConfig, toConfig *rest.Config, registeredWorkspace string) *ConfigMapMutator {
	return &ConfigMapMutator{
		fromConfig:          fromConfig,
		toConfig:            toConfig,
		registeredWorkspace: registeredWorkspace,
	}
}

func (cmm *ConfigMapMutator) ApplyDownstreamName(downstreamObj *unstructured.Unstructured) error {
	var configmap corev1.ConfigMap
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(
		downstreamObj.UnstructuredContent(),
		&configmap)
	if err != nil {
		return err
	}

	originalName, ok := configmap.Labels[originalNameLabel]

	// If the original name label is not present, that means we are syncing from kcp to workloadcluster
	// let's set the original name label to the name of the configmap, and do any transformation required.
	if !ok || originalName == "" {
		// Set the original label
		configmap.Labels[originalNameLabel] = configmap.Name

		if configmap.Name == "kube-root-ca.crt" {
			configmap.Name = "kcp-root-ca.crt"
		}

	} else {
		// If the original name label is present, that means we are syncing from workloadcluster to kcp
		// let's restore the original name to the object.
		configmap.Name = originalName
	}

	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&configmap)
	if err != nil {
		return err
	}
	// Set the changes back into the obj.
	downstreamObj.SetUnstructuredContent(unstructured)
	return nil
}

// ApplyStatus makes modifications to the Status of the deployment object.
func (cmm *ConfigMapMutator) ApplyStatus(upstreamObj *unstructured.Unstructured) error {
	// No transformations
	return nil
}

// ApplySpec makes modifications to the Spec of the deployment object.
func (cmm *ConfigMapMutator) ApplySpec(downstreamObj *unstructured.Unstructured) error {
	// No transformations
	return nil
}
