package signalcd

import "time"

type Pipeline struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Steps  []Step  `json:"steps"`
	Checks []Check `json:"checks"`
}

// TODO: This is probably mostly what Drone uses. Maybe we should copy that struct :)
type Step struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Commands []string `json:"commands"`

	Status *StepStatus `json:"status,omitempty"`
}

type StepStatus struct {
	ExitCode int   `json:"exit_code"`
	Started  int64 `json:"started,omitempty"`
	Stopped  int64 `json:"stopped,omitempty"`
}

type Check struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Duration    time.Duration     `json:"duration"`
	Environment map[string]string `json:"environment"`
}
