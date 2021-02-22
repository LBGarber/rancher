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

type NodePoolHandler func(string, *v3.NodePool) (*v3.NodePool, error)

type NodePoolController interface {
	generic.ControllerMeta
	NodePoolClient

	OnChange(ctx context.Context, name string, sync NodePoolHandler)
	OnRemove(ctx context.Context, name string, sync NodePoolHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() NodePoolCache
}

type NodePoolClient interface {
	Create(*v3.NodePool) (*v3.NodePool, error)
	Update(*v3.NodePool) (*v3.NodePool, error)
	UpdateStatus(*v3.NodePool) (*v3.NodePool, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.NodePool, error)
	List(namespace string, opts metav1.ListOptions) (*v3.NodePoolList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.NodePool, err error)
}

type NodePoolCache interface {
	Get(namespace, name string) (*v3.NodePool, error)
	List(namespace string, selector labels.Selector) ([]*v3.NodePool, error)

	AddIndexer(indexName string, indexer NodePoolIndexer)
	GetByIndex(indexName, key string) ([]*v3.NodePool, error)
}

type NodePoolIndexer func(obj *v3.NodePool) ([]string, error)

type nodePoolController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewNodePoolController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) NodePoolController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &nodePoolController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromNodePoolHandlerToHandler(sync NodePoolHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.NodePool
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.NodePool))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *nodePoolController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.NodePool))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateNodePoolDeepCopyOnChange(client NodePoolClient, obj *v3.NodePool, handler func(obj *v3.NodePool) (*v3.NodePool, error)) (*v3.NodePool, error) {
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

func (c *nodePoolController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *nodePoolController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *nodePoolController) OnChange(ctx context.Context, name string, sync NodePoolHandler) {
	c.AddGenericHandler(ctx, name, FromNodePoolHandlerToHandler(sync))
}

func (c *nodePoolController) OnRemove(ctx context.Context, name string, sync NodePoolHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromNodePoolHandlerToHandler(sync)))
}

func (c *nodePoolController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *nodePoolController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *nodePoolController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *nodePoolController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *nodePoolController) Cache() NodePoolCache {
	return &nodePoolCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *nodePoolController) Create(obj *v3.NodePool) (*v3.NodePool, error) {
	result := &v3.NodePool{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *nodePoolController) Update(obj *v3.NodePool) (*v3.NodePool, error) {
	result := &v3.NodePool{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *nodePoolController) UpdateStatus(obj *v3.NodePool) (*v3.NodePool, error) {
	result := &v3.NodePool{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *nodePoolController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *nodePoolController) Get(namespace, name string, options metav1.GetOptions) (*v3.NodePool, error) {
	result := &v3.NodePool{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *nodePoolController) List(namespace string, opts metav1.ListOptions) (*v3.NodePoolList, error) {
	result := &v3.NodePoolList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *nodePoolController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *nodePoolController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.NodePool, error) {
	result := &v3.NodePool{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type nodePoolCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *nodePoolCache) Get(namespace, name string) (*v3.NodePool, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.NodePool), nil
}

func (c *nodePoolCache) List(namespace string, selector labels.Selector) (ret []*v3.NodePool, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.NodePool))
	})

	return ret, err
}

func (c *nodePoolCache) AddIndexer(indexName string, indexer NodePoolIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.NodePool))
		},
	}))
}

func (c *nodePoolCache) GetByIndex(indexName, key string) (result []*v3.NodePool, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.NodePool, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.NodePool))
	}
	return result, nil
}

type NodePoolStatusHandler func(obj *v3.NodePool, status v3.NodePoolStatus) (v3.NodePoolStatus, error)

type NodePoolGeneratingHandler func(obj *v3.NodePool, status v3.NodePoolStatus) ([]runtime.Object, v3.NodePoolStatus, error)

func RegisterNodePoolStatusHandler(ctx context.Context, controller NodePoolController, condition condition.Cond, name string, handler NodePoolStatusHandler) {
	statusHandler := &nodePoolStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromNodePoolHandlerToHandler(statusHandler.sync))
}

func RegisterNodePoolGeneratingHandler(ctx context.Context, controller NodePoolController, apply apply.Apply,
	condition condition.Cond, name string, handler NodePoolGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &nodePoolGeneratingHandler{
		NodePoolGeneratingHandler: handler,
		apply:                     apply,
		name:                      name,
		gvk:                       controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterNodePoolStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type nodePoolStatusHandler struct {
	client    NodePoolClient
	condition condition.Cond
	handler   NodePoolStatusHandler
}

func (a *nodePoolStatusHandler) sync(key string, obj *v3.NodePool) (*v3.NodePool, error) {
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

type nodePoolGeneratingHandler struct {
	NodePoolGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *nodePoolGeneratingHandler) Remove(key string, obj *v3.NodePool) (*v3.NodePool, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.NodePool{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *nodePoolGeneratingHandler) Handle(obj *v3.NodePool, status v3.NodePoolStatus) (v3.NodePoolStatus, error) {
	objs, newStatus, err := a.NodePoolGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
