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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	v1alpha1 "github.com/networkservicemesh/networkservicemesh/k8s/pkg/apis/networkservice/v1alpha1"
)

// NetworkServiceManagerLister helps list NetworkServiceManagers.
type NetworkServiceManagerLister interface {
	// List lists all NetworkServiceManagers in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.NetworkServiceManager, err error)
	// NetworkServiceManagers returns an object that can list and get NetworkServiceManagers.
	NetworkServiceManagers(namespace string) NetworkServiceManagerNamespaceLister
	NetworkServiceManagerListerExpansion
}

// networkServiceManagerLister implements the NetworkServiceManagerLister interface.
type networkServiceManagerLister struct {
	indexer cache.Indexer
}

// NewNetworkServiceManagerLister returns a new NetworkServiceManagerLister.
func NewNetworkServiceManagerLister(indexer cache.Indexer) NetworkServiceManagerLister {
	return &networkServiceManagerLister{indexer: indexer}
}

// List lists all NetworkServiceManagers in the indexer.
func (s *networkServiceManagerLister) List(selector labels.Selector) (ret []*v1alpha1.NetworkServiceManager, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NetworkServiceManager))
	})
	return ret, err
}

// NetworkServiceManagers returns an object that can list and get NetworkServiceManagers.
func (s *networkServiceManagerLister) NetworkServiceManagers(namespace string) NetworkServiceManagerNamespaceLister {
	return networkServiceManagerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NetworkServiceManagerNamespaceLister helps list and get NetworkServiceManagers.
type NetworkServiceManagerNamespaceLister interface {
	// List lists all NetworkServiceManagers in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.NetworkServiceManager, err error)
	// Get retrieves the NetworkServiceManager from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.NetworkServiceManager, error)
	NetworkServiceManagerNamespaceListerExpansion
}

// networkServiceManagerNamespaceLister implements the NetworkServiceManagerNamespaceLister
// interface.
type networkServiceManagerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all NetworkServiceManagers in the indexer for a given namespace.
func (s networkServiceManagerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.NetworkServiceManager, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NetworkServiceManager))
	})
	return ret, err
}

// Get retrieves the NetworkServiceManager from the indexer for a given namespace and name.
func (s networkServiceManagerNamespaceLister) Get(name string) (*v1alpha1.NetworkServiceManager, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("networkservicemanager"), name)
	}
	return obj.(*v1alpha1.NetworkServiceManager), nil
}
