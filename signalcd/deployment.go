package signalcd

import "time"

type DeploymentPhase string

const (
	Unknown  DeploymentPhase = "unknown"
	Success  DeploymentPhase = "success"
	Failed   DeploymentPhase = "failure"
	Progress DeploymentPhase = "progress"
)

type Deployment struct {
	Number  int64
	Created time.Time
	Status  DeploymentStatus

	Pipeline Pipeline
}

type DeploymentStatus struct {
	Phase DeploymentPhase
}
