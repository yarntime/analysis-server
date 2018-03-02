package mtclient

import (
	"testing"
	"github.com/yarntime/analysis-server/pkg/tools"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRegisterResource(t *testing.T) {
	restConfig, err := tools.GetClientConfig("192.168.254.45:8080")
	if err != nil {
		panic("Failed to create rest config.")
	}

	clientset, err := apiextcs.NewForConfig(restConfig)
	if err != nil {
		panic(err.Error())
	}

	RegisterMonitoredTarget(restConfig)

	_, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(FullCRDName, meta_v1.GetOptions{})

	if err != nil {
		t.Errorf("Failed to register monitored target.")
	}

}
