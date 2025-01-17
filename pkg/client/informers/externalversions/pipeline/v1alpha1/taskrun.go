/*
Copyright 2019 The Tekton Authors

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
package v1alpha1

import (
	time "time"

	pipeline_v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	versioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	internalinterfaces "github.com/tektoncd/pipeline/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/tektoncd/pipeline/pkg/client/listers/pipeline/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TaskRunInformer provides access to a shared informer and lister for
// TaskRuns.
type TaskRunInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.TaskRunLister
}

type taskRunInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTaskRunInformer constructs a new informer for TaskRun type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTaskRunInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTaskRunInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTaskRunInformer constructs a new informer for TaskRun type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTaskRunInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TektonV1alpha1().TaskRuns(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TektonV1alpha1().TaskRuns(namespace).Watch(options)
			},
		},
		&pipeline_v1alpha1.TaskRun{},
		resyncPeriod,
		indexers,
	)
}

func (f *taskRunInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTaskRunInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *taskRunInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&pipeline_v1alpha1.TaskRun{}, f.defaultInformer)
}

func (f *taskRunInformer) Lister() v1alpha1.TaskRunLister {
	return v1alpha1.NewTaskRunLister(f.Informer().GetIndexer())
}
