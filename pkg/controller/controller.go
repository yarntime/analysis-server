package controller

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/yarntime/analysis-server/pkg/client/mtclient"
	"github.com/yarntime/analysis-server/pkg/tools"
	"github.com/yarntime/analysis-server/pkg/types"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	k8s "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	batch "k8s.io/client-go/pkg/apis/batch/v2alpha1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type MTController struct {
	k8sClient *k8s.Clientset

	mtClient *mtclient.MTClient

	recorder record.EventRecorder

	concurrentJobHandlers int

	resyncPeriod time.Duration

	queue workqueue.RateLimitingInterface
}

func NewMTController(c *Config) *MTController {

	mtController := &MTController{
		k8sClient:             c.K8sClient,
		mtClient:              c.MTClient,
		concurrentJobHandlers: c.ConcurrentJobHandlers,
		resyncPeriod:          c.ResyncPeriod,
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "monitored"),
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(mtController.k8sClient.CoreV1().RESTClient()).Events("")})

	_, mtlw := cache.NewInformer(
		cache.NewListWatchFromClient(mtController.mtClient.RESTClient(), "monitoredtargets", meta_v1.NamespaceAll, fields.Everything()),
		&types.MonitoredTarget{},
		mtController.resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: mtController.enqueueController,
			UpdateFunc: func(old, cur interface{}) {
				mt := cur.(*types.MonitoredTarget)
				mtController.enqueueController(mt)
			},
		},
	)

	_, cjlw := cache.NewInformer(
		cache.NewListWatchFromClient(mtController.k8sClient.BatchV2alpha1Client.RESTClient(), "cronjobs", meta_v1.NamespaceAll, fields.Everything()),
		&batch.CronJob{},
		mtController.resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: mtController.enqueueController,
			UpdateFunc: func(old, cur interface{}) {
				mt := cur.(*types.MonitoredTarget)
				mtController.enqueueController(mt)
			},
		},
	)

	go mtlw.Run(c.StopCh)
	go cjlw.Run(c.StopCh)

	return mtController
}

func (mtc *MTController) Run(stopCh chan struct{}) {
	for i := 0; i < mtc.concurrentJobHandlers; i++ {
		go wait.Until(mtc.startHandler, time.Second, stopCh)
	}

	<-stopCh
}

func (mtc *MTController) enqueueController(obj interface{}) {
	mt := obj.(*types.MonitoredTarget)
	key := tools.GetKeyOfResource(mt.ObjectMeta)
	mtc.queue.Add(key)
}

func (mtc *MTController) startHandler() {
	for mtc.processNextWorkItem() {
	}
}

func (mtc *MTController) processNextWorkItem() bool {
	key, quit := mtc.queue.Get()
	if quit {
		return false
	}
	defer mtc.queue.Done(key)

	mtc.processMonitoredTarget(key.(string))
	return true
}

func (mtc *MTController) processMonitoredTarget(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing monitored target %q (%v)", key, time.Now().Sub(startTime))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	if len(ns) == 0 || len(name) == 0 {
		return fmt.Errorf("invalid monitored target key %q: either namespace or name is missing", key)
	}

	mt, err := mtc.mtClient.MonitoredTargets(ns).Get(name, meta_v1.GetOptions{})
	if err != nil {
		glog.Warningf("Failed get monitored target %q (%v) from kubernetes", key, time.Now().Sub(startTime))
		if errors.IsNotFound(err) {
			glog.V(4).Infof("MonitoredTarget has been deleted: %v", key)
			return nil
		}
		return err
	}

	if mt.Status.StartTime == nil {
		// create cron job
		return nil
	}

	return nil
}
