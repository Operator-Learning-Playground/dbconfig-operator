/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	dbconfigv1alpha1 "github.com/myoperator/dbconfigoperator/pkg/apis/dbconfig/v1alpha1"
	versioned "github.com/myoperator/dbconfigoperator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/myoperator/dbconfigoperator/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/myoperator/dbconfigoperator/pkg/client/listers/dbconfig/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DbConfigInformer provides access to a shared informer and lister for
// DbConfigs.
type DbConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.DbConfigLister
}

type dbConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDbConfigInformer constructs a new informer for DbConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDbConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDbConfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDbConfigInformer constructs a new informer for DbConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDbConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApiV1alpha1().DbConfigs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApiV1alpha1().DbConfigs(namespace).Watch(context.TODO(), options)
			},
		},
		&dbconfigv1alpha1.DbConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *dbConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDbConfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *dbConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&dbconfigv1alpha1.DbConfig{}, f.defaultInformer)
}

func (f *dbConfigInformer) Lister() v1alpha1.DbConfigLister {
	return v1alpha1.NewDbConfigLister(f.Informer().GetIndexer())
}
