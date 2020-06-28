package signalcd

import "time"

//Phase is the current state of a Deployment Step or Check
type Phase string

const (
	//Unknown is the state of a Deployment that is currently unknown
	Unknown Phase = "unknown"
	//Success is the state of a Deployment that was successfully executed
	Success Phase = "success"
	//Failure is the state of a Deployment that failed to execute
	Failure Phase = "failure"
	//Progress is the state of a Deployment that is currently running
	Progress Phase = "progress"
	//Pending is the state of a Deployment that is schedule but not yet progressing
	Pending Phase = "pending"
	//Killed is the state of a Deployment that was killed during execution
	Killed Phase = "killed"
)

//Deployment is a specific execution of a Pipeline with more meta data
type Deployment struct {
	Number   int64
	Created  time.Time
	Pipeline Pipeline

	Status map[string]*Status
}

type Status struct {
	Steps []StepStatus
}

type StepStatus struct {
	Phase    Phase
	ExitCode int
	Started  time.Time
	Stopped  *time.Time
}
