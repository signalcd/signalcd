package signalcd

import "time"

// Pipeline is the definition on how to run steps and long running checks
type Pipeline struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Steps   []Step    `json:"steps"`
	Checks  []Check   `json:"checks"`
	Created time.Time `json:"created"`
}

// Step is a synchronous execution step in a Pipeline
// TODO: This is probably mostly what Drone uses. Maybe we should copy that struct :)
type Step struct {
	Name             string   `json:"name"`
	Image            string   `json:"image"`
	ImagePullSecrets []string `json:"imagePullSecrets"`
	Commands         []string `json:"commands"`

	Status *Status `json:"status,omitempty"`
}

// Status is the state of a Step in a Pipeline
type Status struct {
	ExitCode int   `json:"exit_code"`
	Started  int64 `json:"started,omitempty"`
	Stopped  int64 `json:"stopped,omitempty"`

	Logs []byte `json:"logs"`
}

// Check is an asynchronous long running Check after the Pipeline was successfully executed
type Check struct {
	Name             string            `json:"name"`
	Image            string            `json:"image"`
	ImagePullSecrets []string          `json:"imagePullSecrets"`
	Duration         time.Duration     `json:"duration"`
	Environment      map[string]string `json:"environment"`

	Status *Status `json:"status,omitempty"`
}
