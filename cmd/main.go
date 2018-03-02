package main

import (
	"flag"
	k8sclient "github.com/yarntime/analysis-server/pkg/client/k8s_client"
	mtclient "github.com/yarntime/analysis-server/pkg/client/mtclient"
	c "github.com/yarntime/analysis-server/pkg/controller"
	"github.com/yarntime/analysis-server/pkg/tools"
	"github.com/golang/glog"
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
	flag.Set("alsologtostderr", "true")
	flag.Parse()
}

func main() {
	stop := make(chan struct{})

	restConfig, err := tools.GetClientConfig(apiserverAddress)
	if err != nil {
		panic(err.Error())
	}

	glog.Info("register monitored target.")
	err = mtclient.RegisterMonitoredTarget(restConfig)
	if err != nil {
		panic(err.Error())
	}

	config := &c.Config{
		Address:               apiserverAddress,
		ConcurrentJobHandlers: concurrentJobHandlers,
		ResyncPeriod:          resyncPeriod,
		StopCh:                stop,
		K8sClient:             k8sclient.NewK8sClint(restConfig),
		MTClient:              mtclient.NewMTClient(restConfig),
	}

	mtc := c.NewMTController(config)

	glog.Info("run controller.")
	go mtc.Run(stop)

	select {}
}
