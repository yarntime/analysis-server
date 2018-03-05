package main

import (
	"flag"
	"github.com/golang/glog"
	k8sclient "github.com/yarntime/analysis-server/pkg/client/k8s_client"
	mtclient "github.com/yarntime/analysis-server/pkg/client/mtclient"
	c "github.com/yarntime/analysis-server/pkg/controller"
	"github.com/yarntime/analysis-server/pkg/tools"
	"time"
)

var (
	apiserverAddress      string
	concurrentJobHandlers int
	resyncPeriod          time.Duration
	baseImage             string
	jobNamespace          string
)

func init() {
	flag.StringVar(&apiserverAddress, "apiserver_address", "", "Kubernetes apiserver address")
	flag.IntVar(&concurrentJobHandlers, "concurrent_job_handlers", 4, "Concurrent job handlers")
	flag.DurationVar(&resyncPeriod, "resync_period", time.Minute*30, "Resync period")
	flag.StringVar(&baseImage, "base_image", "registry.harbor:5000/sky-firmament/predict:latest", "Job image")
	flag.StringVar(&jobNamespace, "job_namespace", "sky-firmament", "Job namespace")
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
		BaseImage:             baseImage,
		JobNamespace:          jobNamespace,
		K8sClient:             k8sclient.NewK8sClint(restConfig),
		MTClient:              mtclient.NewMTClient(restConfig),
	}

	mtc := c.NewMTController(config)

	glog.Info("run controller.")
	go mtc.Run(stop)

	select {}
}
