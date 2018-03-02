package main

import (
	"flag"
	k8sclient "github.com/yarntime/analysis-server/pkg/client/k8s_client"
	mtclient "github.com/yarntime/analysis-server/pkg/client/mtclient"
	c "github.com/yarntime/analysis-server/pkg/controller"
	"time"
)

var (
	apiserverAddress      string
	concurrentJobHandlers int
	resyncPeriod          time.Duration
)

func init() {
	flag.StringVar(&apiserverAddress, "apiserver_address", "", "Kubernetes apiserver address")
	flag.IntVar(&concurrentJobHandlers, "concurrent_job_handlers", 4, "Concurrent job handlers")
	flag.DurationVar(&resyncPeriod, "resync_period", time.Minute*30, "Resync period")
	flag.Parse()
}

func main() {
	stop := make(chan struct{})
	config := &c.Config{
		Address:               apiserverAddress,
		ConcurrentJobHandlers: concurrentJobHandlers,
		ResyncPeriod:          resyncPeriod,
		StopCh:                stop,
		K8sClient:             k8sclient.NewK8sClint(apiserverAddress),
		MTClient:              mtclient.NewMTClient(apiserverAddress),
	}

	mtc := c.NewMTController(config)

	go mtc.Run()

	select {}
}
