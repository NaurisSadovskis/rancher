/*
Copyright 2021 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type RkeAddonHandler func(string, *v3.RkeAddon) (*v3.RkeAddon, error)

type RkeAddonController interface {
	generic.ControllerMeta
	RkeAddonClient

	OnChange(ctx context.Context, name string, sync RkeAddonHandler)
	OnRemove(ctx context.Context, name string, sync RkeAddonHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() RkeAddonCache
}

type RkeAddonClient interface {
	Create(*v3.RkeAddon) (*v3.RkeAddon, error)
	Update(*v3.RkeAddon) (*v3.RkeAddon, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.RkeAddon, error)
	List(namespace string, opts metav1.ListOptions) (*v3.RkeAddonList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.RkeAddon, err error)
}

type RkeAddonCache interface {
	Get(namespace, name string) (*v3.RkeAddon, error)
	List(namespace string, selector labels.Selector) ([]*v3.RkeAddon, error)

	AddIndexer(indexName string, indexer RkeAddonIndexer)
	GetByIndex(indexName, key string) ([]*v3.RkeAddon, error)
}

type RkeAddonIndexer func(obj *v3.RkeAddon) ([]string, error)

type rkeAddonController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewRkeAddonController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) RkeAddonController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &rkeAddonController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromRkeAddonHandlerToHandler(sync RkeAddonHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.RkeAddon
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.RkeAddon))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *rkeAddonController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.RkeAddon))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateRkeAddonDeepCopyOnChange(client RkeAddonClient, obj *v3.RkeAddon, handler func(obj *v3.RkeAddon) (*v3.RkeAddon, error)) (*v3.RkeAddon, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *rkeAddonController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *rkeAddonController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *rkeAddonController) OnChange(ctx context.Context, name string, sync RkeAddonHandler) {
	c.AddGenericHandler(ctx, name, FromRkeAddonHandlerToHandler(sync))
}

func (c *rkeAddonController) OnRemove(ctx context.Context, name string, sync RkeAddonHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromRkeAddonHandlerToHandler(sync)))
}

func (c *rkeAddonController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *rkeAddonController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *rkeAddonController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *rkeAddonController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *rkeAddonController) Cache() RkeAddonCache {
	return &rkeAddonCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *rkeAddonController) Create(obj *v3.RkeAddon) (*v3.RkeAddon, error) {
	result := &v3.RkeAddon{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *rkeAddonController) Update(obj *v3.RkeAddon) (*v3.RkeAddon, error) {
	result := &v3.RkeAddon{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *rkeAddonController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *rkeAddonController) Get(namespace, name string, options metav1.GetOptions) (*v3.RkeAddon, error) {
	result := &v3.RkeAddon{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *rkeAddonController) List(namespace string, opts metav1.ListOptions) (*v3.RkeAddonList, error) {
	result := &v3.RkeAddonList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *rkeAddonController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *rkeAddonController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.RkeAddon, error) {
	result := &v3.RkeAddon{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type rkeAddonCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *rkeAddonCache) Get(namespace, name string) (*v3.RkeAddon, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.RkeAddon), nil
}

func (c *rkeAddonCache) List(namespace string, selector labels.Selector) (ret []*v3.RkeAddon, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.RkeAddon))
	})

	return ret, err
}

func (c *rkeAddonCache) AddIndexer(indexName string, indexer RkeAddonIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.RkeAddon))
		},
	}))
}

func (c *rkeAddonCache) GetByIndex(indexName, key string) (result []*v3.RkeAddon, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.RkeAddon, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.RkeAddon))
	}
	return result, nil
}
