package mtclient

import (
	"k8s.io/apimachinery/pkg/runtime"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type MTInterface interface {
	RESTClient() rest.Interface
	MonitoredTargetGetter
}

type MTClient struct {
	restClient rest.Interface
}

func (c *MTClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *MTClient) MonitoredTargets(namespace string) MonitoredTargetInterface {
	return newMonitoredTargets(c, namespace)
}

func NewMTClient(config *rest.Config) *MTClient {
	clientSet, err := newForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet
}

func newForConfig(c *rest.Config) (*MTClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &MTClient{client}, nil
}

func setConfigDefaults(config *rest.Config) error {
	gv := SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

func New(c rest.Interface) *MTClient {
	return &MTClient{c}
}
