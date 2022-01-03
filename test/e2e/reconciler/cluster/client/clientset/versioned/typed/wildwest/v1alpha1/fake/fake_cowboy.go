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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha1 "github.com/kcp-dev/kcp/test/e2e/reconciler/cluster/apis/wildwest/v1alpha1"
)

// FakeCowboys implements CowboyInterface
type FakeCowboys struct {
	Fake *FakeWildwestV1alpha1
	ns   string
}

var cowboysResource = schema.GroupVersionResource{Group: "wildwest.dev", Version: "v1alpha1", Resource: "cowboys"}

var cowboysKind = schema.GroupVersionKind{Group: "wildwest.dev", Version: "v1alpha1", Kind: "Cowboy"}

// Get takes name of the cowboy, and returns the corresponding cowboy object, and an error if there is any.
func (c *FakeCowboys) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Cowboy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(cowboysResource, c.ns, name), &v1alpha1.Cowboy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Cowboy), err
}

// List takes label and field selectors, and returns the list of Cowboys that match those selectors.
func (c *FakeCowboys) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.CowboyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(cowboysResource, cowboysKind, c.ns, opts), &v1alpha1.CowboyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.CowboyList{ListMeta: obj.(*v1alpha1.CowboyList).ListMeta}
	for _, item := range obj.(*v1alpha1.CowboyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cowboys.
func (c *FakeCowboys) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(cowboysResource, c.ns, opts))

}

// Create takes the representation of a cowboy and creates it.  Returns the server's representation of the cowboy, and an error, if there is any.
func (c *FakeCowboys) Create(ctx context.Context, cowboy *v1alpha1.Cowboy, opts v1.CreateOptions) (result *v1alpha1.Cowboy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(cowboysResource, c.ns, cowboy), &v1alpha1.Cowboy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Cowboy), err
}

// Update takes the representation of a cowboy and updates it. Returns the server's representation of the cowboy, and an error, if there is any.
func (c *FakeCowboys) Update(ctx context.Context, cowboy *v1alpha1.Cowboy, opts v1.UpdateOptions) (result *v1alpha1.Cowboy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(cowboysResource, c.ns, cowboy), &v1alpha1.Cowboy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Cowboy), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCowboys) UpdateStatus(ctx context.Context, cowboy *v1alpha1.Cowboy, opts v1.UpdateOptions) (*v1alpha1.Cowboy, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(cowboysResource, "status", c.ns, cowboy), &v1alpha1.Cowboy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Cowboy), err
}

// Delete takes name of the cowboy and deletes it. Returns an error if one occurs.
func (c *FakeCowboys) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(cowboysResource, c.ns, name), &v1alpha1.Cowboy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCowboys) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(cowboysResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.CowboyList{})
	return err
}

// Patch applies the patch and returns the patched cowboy.
func (c *FakeCowboys) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Cowboy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(cowboysResource, c.ns, name, pt, data, subresources...), &v1alpha1.Cowboy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Cowboy), err
}
