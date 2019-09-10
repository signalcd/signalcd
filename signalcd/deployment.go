package signalcd

import "time"

//DeploymentPhase is the current state of a Deployment
type DeploymentPhase string

const (
	//Unknown is the state of a Deployment that is currently unknown
	Unknown DeploymentPhase = "unknown"
	//Success is the state of a Deployment that was successfully executed
	Success DeploymentPhase = "success"
	//Failure is the state of a Deployment that failed to execute
	Failure DeploymentPhase = "failure"
	//Progress is the state of a Deployment that is currently running
	Progress DeploymentPhase = "progress"
	//Pending is the state of a Deployment that is schedule but not yet progressing
	Pending DeploymentPhase = "pending"
	//Killed is the state of a Deployment that was killed during execution
	Killed DeploymentPhase = "killed"
)

//Deployment is a specific execution of a Pipeline with more meta data
type Deployment struct {
	Number   int64
	Created  time.Time
	Started  time.Time
	Finished time.Time
	Status   DeploymentStatus

	Pipeline Pipeline
}

//DeploymentStatus is the status of a Deployment
type DeploymentStatus struct {
	Phase DeploymentPhase
}
