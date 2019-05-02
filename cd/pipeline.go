package cd

import "time"

type Pipeline struct {
	ID     string  `json:"id"`
	Steps  []Step  `json:"steps"`
	Checks []Check `json:"checks"`
}

// TODO: This is probably mostly what Drone uses. Maybe we should copy that struct :)
type Step struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Commands []string `json:"commands"`
}

type Check struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Duration    time.Duration     `json:"duration"`
	Environment map[string]string `json:"environment"`
}
