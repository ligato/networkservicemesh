// Copyright (c) 2018 Cisco and/or its affiliates.
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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	networkservicemesh_io_v1 "github.com/ligato/networkservicemesh/pkg/apis/networkservicemesh.io/v1"
	versioned "github.com/ligato/networkservicemesh/pkg/client/clientset/versioned"
	internalinterfaces "github.com/ligato/networkservicemesh/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/ligato/networkservicemesh/pkg/client/listers/networkservicemesh.io/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// NetworkServiceChannelInformer provides access to a shared informer and lister for
// NetworkServiceChannels.
type NetworkServiceChannelInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.NetworkServiceChannelLister
}

type networkServiceChannelInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewNetworkServiceChannelInformer constructs a new informer for NetworkServiceChannel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewNetworkServiceChannelInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredNetworkServiceChannelInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredNetworkServiceChannelInformer constructs a new informer for NetworkServiceChannel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredNetworkServiceChannelInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkserviceV1().NetworkServiceChannels(namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkserviceV1().NetworkServiceChannels(namespace).Watch(options)
			},
		},
		&networkservicemesh_io_v1.NetworkServiceChannel{},
		resyncPeriod,
		indexers,
	)
}

func (f *networkServiceChannelInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredNetworkServiceChannelInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *networkServiceChannelInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&networkservicemesh_io_v1.NetworkServiceChannel{}, f.defaultInformer)
}

func (f *networkServiceChannelInformer) Lister() v1.NetworkServiceChannelLister {
	return v1.NewNetworkServiceChannelLister(f.Informer().GetIndexer())
}
