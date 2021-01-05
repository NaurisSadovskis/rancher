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
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
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

type ClusterLoggingHandler func(string, *v3.ClusterLogging) (*v3.ClusterLogging, error)

type ClusterLoggingController interface {
	generic.ControllerMeta
	ClusterLoggingClient

	OnChange(ctx context.Context, name string, sync ClusterLoggingHandler)
	OnRemove(ctx context.Context, name string, sync ClusterLoggingHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterLoggingCache
}

type ClusterLoggingClient interface {
	Create(*v3.ClusterLogging) (*v3.ClusterLogging, error)
	Update(*v3.ClusterLogging) (*v3.ClusterLogging, error)
	UpdateStatus(*v3.ClusterLogging) (*v3.ClusterLogging, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterLogging, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ClusterLoggingList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ClusterLogging, err error)
}

type ClusterLoggingCache interface {
	Get(namespace, name string) (*v3.ClusterLogging, error)
	List(namespace string, selector labels.Selector) ([]*v3.ClusterLogging, error)

	AddIndexer(indexName string, indexer ClusterLoggingIndexer)
	GetByIndex(indexName, key string) ([]*v3.ClusterLogging, error)
}

type ClusterLoggingIndexer func(obj *v3.ClusterLogging) ([]string, error)

type clusterLoggingController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterLoggingController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterLoggingController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterLoggingController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterLoggingHandlerToHandler(sync ClusterLoggingHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ClusterLogging
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ClusterLogging))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterLoggingController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ClusterLogging))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterLoggingDeepCopyOnChange(client ClusterLoggingClient, obj *v3.ClusterLogging, handler func(obj *v3.ClusterLogging) (*v3.ClusterLogging, error)) (*v3.ClusterLogging, error) {
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

func (c *clusterLoggingController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterLoggingController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterLoggingController) OnChange(ctx context.Context, name string, sync ClusterLoggingHandler) {
	c.AddGenericHandler(ctx, name, FromClusterLoggingHandlerToHandler(sync))
}

func (c *clusterLoggingController) OnRemove(ctx context.Context, name string, sync ClusterLoggingHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterLoggingHandlerToHandler(sync)))
}

func (c *clusterLoggingController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *clusterLoggingController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *clusterLoggingController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterLoggingController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterLoggingController) Cache() ClusterLoggingCache {
	return &clusterLoggingCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterLoggingController) Create(obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
	result := &v3.ClusterLogging{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *clusterLoggingController) Update(obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
	result := &v3.ClusterLogging{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterLoggingController) UpdateStatus(obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
	result := &v3.ClusterLogging{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterLoggingController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *clusterLoggingController) Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterLogging, error) {
	result := &v3.ClusterLogging{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *clusterLoggingController) List(namespace string, opts metav1.ListOptions) (*v3.ClusterLoggingList, error) {
	result := &v3.ClusterLoggingList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *clusterLoggingController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *clusterLoggingController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ClusterLogging, error) {
	result := &v3.ClusterLogging{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterLoggingCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterLoggingCache) Get(namespace, name string) (*v3.ClusterLogging, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ClusterLogging), nil
}

func (c *clusterLoggingCache) List(namespace string, selector labels.Selector) (ret []*v3.ClusterLogging, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ClusterLogging))
	})

	return ret, err
}

func (c *clusterLoggingCache) AddIndexer(indexName string, indexer ClusterLoggingIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ClusterLogging))
		},
	}))
}

func (c *clusterLoggingCache) GetByIndex(indexName, key string) (result []*v3.ClusterLogging, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ClusterLogging, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ClusterLogging))
	}
	return result, nil
}

type ClusterLoggingStatusHandler func(obj *v3.ClusterLogging, status v3.ClusterLoggingStatus) (v3.ClusterLoggingStatus, error)

type ClusterLoggingGeneratingHandler func(obj *v3.ClusterLogging, status v3.ClusterLoggingStatus) ([]runtime.Object, v3.ClusterLoggingStatus, error)

func RegisterClusterLoggingStatusHandler(ctx context.Context, controller ClusterLoggingController, condition condition.Cond, name string, handler ClusterLoggingStatusHandler) {
	statusHandler := &clusterLoggingStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterLoggingHandlerToHandler(statusHandler.sync))
}

func RegisterClusterLoggingGeneratingHandler(ctx context.Context, controller ClusterLoggingController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterLoggingGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterLoggingGeneratingHandler{
		ClusterLoggingGeneratingHandler: handler,
		apply:                           apply,
		name:                            name,
		gvk:                             controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterLoggingStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterLoggingStatusHandler struct {
	client    ClusterLoggingClient
	condition condition.Cond
	handler   ClusterLoggingStatusHandler
}

func (a *clusterLoggingStatusHandler) sync(key string, obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type clusterLoggingGeneratingHandler struct {
	ClusterLoggingGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterLoggingGeneratingHandler) Remove(key string, obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ClusterLogging{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterLoggingGeneratingHandler) Handle(obj *v3.ClusterLogging, status v3.ClusterLoggingStatus) (v3.ClusterLoggingStatus, error) {
	objs, newStatus, err := a.ClusterLoggingGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
