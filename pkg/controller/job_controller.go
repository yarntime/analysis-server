package controller

import (
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	batch "k8s.io/client-go/pkg/apis/batch/v2alpha1"
	"k8s.io/client-go/pkg/util"
)

const (
	ContainerNamePrefix = "predict-job"
)

type JobController struct {
	K8sClient    *k8s.Clientset
	JobNamespace string
	BaseImage    string
}

func NewJobController(c *Config) *JobController {
	return &JobController{
		K8sClient:    c.K8sClient,
		JobNamespace: c.JobNamespace,
		BaseImage:    c.BaseImage,
	}
}

func componentCronJob(container v1.Container, namespace string) *batch.CronJob {
	return &batch.CronJob{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      container.Name,
			Namespace: namespace,
			Labels:    map[string]string{"component": container.Name, "tier": "training-job"},
		},
		Spec: batch.CronJobSpec{
			Schedule:                   "",
			ConcurrencyPolicy:          batch.ForbidConcurrent,
			SuccessfulJobsHistoryLimit: util.Int32Ptr(0),
			FailedJobsHistoryLimit:     util.Int32Ptr(10),
			JobTemplate: batch.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Parallelism: util.Int32Ptr(1),
					Completions: util.Int32Ptr(1),
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers:    []v1.Container{container},
							RestartPolicy: v1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
}

func componentResources(cpu string) v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceName(v1.ResourceCPU): resource.MustParse(cpu),
		},
	}
}

func (jc *JobController) StartTrainingJob(node string) {
	job := componentCronJob(v1.Container{
		Name:      ContainerNamePrefix,
		Image:     jc.BaseImage,
		Command:   []string{"/training"},
		Args:      []string{},
		Resources: componentResources("500m")}, jc.JobNamespace)
	_, err := jc.K8sClient.BatchV2alpha1().CronJobs(jc.JobNamespace).Create(job)
	if err != nil {
		glog.Errorf("Failed to create training job: %s/%s", job.Namespace, job.Name)
	}
}
