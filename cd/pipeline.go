package cd

type Pipeline struct {
	ID    string
	Steps []Step `json:"steps"`
}

// TODO: This is probably mostly what Drone uses. Maybe we should copy that struct :)
type Step struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Commands []string `json:"commands"`
}
