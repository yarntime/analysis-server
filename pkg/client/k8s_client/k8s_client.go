package client

import (
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewK8sClint(config *rest.Config) *k8s.Clientset {

	clientSet, err := k8s.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet
}
