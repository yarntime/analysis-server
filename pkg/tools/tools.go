package tools

import (
	mata_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientConfig(host string) (*rest.Config, error) {
	if host != "" {
		return clientcmd.BuildConfigFromFlags(host, "")
	}
	return rest.InClusterConfig()
}

func GetKeyOfResource(meta mata_v1.ObjectMeta) string {
	return meta.Namespace + "/" + meta.Name
}
