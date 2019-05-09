package signalcd

import (
	appsv1 "k8s.io/api/apps/v1"
)

type Agent struct {
	Name   string                  `json:"name"`
	Status appsv1.DeploymentStatus `json:"status"`
}
