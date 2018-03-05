package mtclient

import (
	"github.com/yarntime/analysis-server/pkg/types"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

type MonitoredTargetGetter interface {
	MonitoredTargets(namespace string) MonitoredTargetInterface
}

type MonitoredTargetInterface interface {
	Create(*types.MonitoredTarget) (*types.MonitoredTarget, error)
	Update(*types.MonitoredTarget) (*types.MonitoredTarget, error)
	UpdateStatus(*types.MonitoredTarget) (*types.MonitoredTarget, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*types.MonitoredTarget, error)
	List(opts meta_v1.ListOptions) (*types.MonitoredTargetList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
}

type monitoredTargets struct {
	client rest.Interface
	ns     string
}

func newMonitoredTargets(c *MTClient, namespace string) *monitoredTargets {
	return &monitoredTargets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

func (c *monitoredTargets) Create(monitoredTarget *types.MonitoredTarget) (result *types.MonitoredTarget, err error) {
	result = &types.MonitoredTarget{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(ResourceKind).
		Body(monitoredTarget).
		Do().
		Into(result)
	return
}

func (c *monitoredTargets) Update(monitoredTarget *types.MonitoredTarget) (result *types.MonitoredTarget, err error) {
	result = &types.MonitoredTarget{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(monitoredTarget.Name).
		Body(monitoredTarget).
		Do().
		Into(result)
	return
}

func (c *monitoredTargets) UpdateStatus(monitoredTarget *types.MonitoredTarget) (result *types.MonitoredTarget, err error) {
	result = &types.MonitoredTarget{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(monitoredTarget.Name).
		SubResource("status").
		Body(monitoredTarget).
		Do().
		Into(result)
	return
}

func (c *monitoredTargets) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(name).
		Body(options).
		Do().
		Error()
}

func (c *monitoredTargets) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(ResourceKind).
		VersionedParams(&listOptions, ParameterCodec).
		Body(options).
		Do().
		Error()
}

func (c *monitoredTargets) Get(name string, options meta_v1.GetOptions) (result *types.MonitoredTarget, err error) {
	result = &types.MonitoredTarget{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(name).
		VersionedParams(&options, ParameterCodec).
		Do().
		Into(result)
	return
}

func (c *monitoredTargets) List(opts meta_v1.ListOptions) (result *types.MonitoredTargetList, err error) {
	result = &types.MonitoredTargetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(ResourceKind).
		VersionedParams(&opts, ParameterCodec).
		Do().
		Into(result)
	return
}

func (c *monitoredTargets) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource(ResourceKind).
		VersionedParams(&opts, ParameterCodec).
		Watch()
}
