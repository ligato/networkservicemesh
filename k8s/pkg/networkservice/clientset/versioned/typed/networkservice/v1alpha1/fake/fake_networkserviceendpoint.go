// Copyright (c) 2019 Cisco and/or its affiliates.
// Copyright (c) 2019 Red Hat Inc. and/or its affiliates.
// Copyright (c) 2019 VMware, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	v1alpha1 "github.com/networkservicemesh/networkservicemesh/k8s/pkg/apis/networkservice/v1alpha1"
)

// FakeNetworkServiceEndpoints implements NetworkServiceEndpointInterface
type FakeNetworkServiceEndpoints struct {
	Fake *FakeNetworkserviceV1alpha1
	ns   string
}

var networkserviceendpointsResource = schema.GroupVersionResource{Group: "networkservice", Version: "v1alpha1", Resource: "networkserviceendpoints"}

var networkserviceendpointsKind = schema.GroupVersionKind{Group: "networkservice", Version: "v1alpha1", Kind: "NetworkServiceEndpoint"}

// Get takes name of the networkServiceEndpoint, and returns the corresponding networkServiceEndpoint object, and an error if there is any.
func (c *FakeNetworkServiceEndpoints) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NetworkServiceEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(networkserviceendpointsResource, c.ns, name), &v1alpha1.NetworkServiceEndpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NetworkServiceEndpoint), err
}

// List takes label and field selectors, and returns the list of NetworkServiceEndpoints that match those selectors.
func (c *FakeNetworkServiceEndpoints) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NetworkServiceEndpointList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(networkserviceendpointsResource, networkserviceendpointsKind, c.ns, opts), &v1alpha1.NetworkServiceEndpointList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.NetworkServiceEndpointList{ListMeta: obj.(*v1alpha1.NetworkServiceEndpointList).ListMeta}
	for _, item := range obj.(*v1alpha1.NetworkServiceEndpointList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested networkServiceEndpoints.
func (c *FakeNetworkServiceEndpoints) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(networkserviceendpointsResource, c.ns, opts))

}

// Create takes the representation of a networkServiceEndpoint and creates it.  Returns the server's representation of the networkServiceEndpoint, and an error, if there is any.
func (c *FakeNetworkServiceEndpoints) Create(ctx context.Context, networkServiceEndpoint *v1alpha1.NetworkServiceEndpoint, opts v1.CreateOptions) (result *v1alpha1.NetworkServiceEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(networkserviceendpointsResource, c.ns, networkServiceEndpoint), &v1alpha1.NetworkServiceEndpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NetworkServiceEndpoint), err
}

// Update takes the representation of a networkServiceEndpoint and updates it. Returns the server's representation of the networkServiceEndpoint, and an error, if there is any.
func (c *FakeNetworkServiceEndpoints) Update(ctx context.Context, networkServiceEndpoint *v1alpha1.NetworkServiceEndpoint, opts v1.UpdateOptions) (result *v1alpha1.NetworkServiceEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(networkserviceendpointsResource, c.ns, networkServiceEndpoint), &v1alpha1.NetworkServiceEndpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NetworkServiceEndpoint), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNetworkServiceEndpoints) UpdateStatus(ctx context.Context, networkServiceEndpoint *v1alpha1.NetworkServiceEndpoint, opts v1.UpdateOptions) (*v1alpha1.NetworkServiceEndpoint, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(networkserviceendpointsResource, "status", c.ns, networkServiceEndpoint), &v1alpha1.NetworkServiceEndpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NetworkServiceEndpoint), err
}

// Delete takes name of the networkServiceEndpoint and deletes it. Returns an error if one occurs.
func (c *FakeNetworkServiceEndpoints) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(networkserviceendpointsResource, c.ns, name), &v1alpha1.NetworkServiceEndpoint{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNetworkServiceEndpoints) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(networkserviceendpointsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.NetworkServiceEndpointList{})
	return err
}

// Patch applies the patch and returns the patched networkServiceEndpoint.
func (c *FakeNetworkServiceEndpoints) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NetworkServiceEndpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(networkserviceendpointsResource, c.ns, name, pt, data, subresources...), &v1alpha1.NetworkServiceEndpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NetworkServiceEndpoint), err
}
