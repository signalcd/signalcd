package signalcd

import "time"

//DeploymentPhase is the current state of a Deployment
type DeploymentPhase string

const (
	//Unknown is the state of a Deployment that is currently unknown
	Unknown DeploymentPhase = "unknown"
	//Success is the state of a Deployment that was successfully executed
	Success DeploymentPhase = "success"
	//Failed is the state of a Deployment that failed to execute
	Failed DeploymentPhase = "failure"
	//Progress is the state of a Deployment that is currently running
	Progress DeploymentPhase = "progress"
)

//Deployment is a specific execution of a Pipeline with more meta data
type Deployment struct {
	Number  int64
	Created time.Time
	Status  DeploymentStatus

	Pipeline Pipeline
}

//DeploymentStatus is the status of a Deployment
type DeploymentStatus struct {
	Phase DeploymentPhase
}
