package types

import (
	"github.com/Azure/go-autorest/autorest/date"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type TargetPhase string

const (
	Creating TargetPhase = "creating"
	Finished TargetPhase = "finished"
)

type MonitoredTarget struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata,omitempty"`
	Spec               TargetSpec   `json:"spec"`
	Status             TargetStatus `json:"status"`
}

type TargetSpec struct {
	ResourceType      string        `json:"type"`
	ResourceNamespace string        `json:"namespace"`
	ResourceName      string        `json:"name"`
	MetricName        string        `json:"metric"`
	Cron              string        `json:"cron"`
	AbnormalDetection bool          `json:"abnormaldetection"`
	Period            time.Duration `json:"period"`
}

type TargetStatus struct {
	Phase     TargetPhase   `json:"phase,omitempty"`
	StartTime *meta_v1.Time `json:"startTime,omitempty"`
	RunTimes  int           `json:"runtimes"`
	CalDate   date.Date     `json:"date"`
}

type MonitoredTargetList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []MonitoredTarget `json:"items"`
}
