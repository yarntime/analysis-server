package controller

import (
	"github.com/yarntime/analysis-server/pkg/client/mtclient"
	k8s "k8s.io/client-go/kubernetes"
	"time"
)

type Config struct {
	Address               string
	ConcurrentJobHandlers int
	StopCh                chan struct{}
	ResyncPeriod          time.Duration
	K8sClient             *k8s.Clientset
	MTClient              *mtclient.MTClient
}
