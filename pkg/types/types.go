package types

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
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
	Period            time.Duration `json:"period"`
}

type TargetStatus struct {
	RunTimes int `json:"runtimes"`
}

type MonitoredTargetList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []MonitoredTarget `json:"items"`
}
