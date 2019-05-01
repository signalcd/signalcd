package cd

type Pipeline struct {
	ID       string   `json:"id"`
	Artifact Artifact `json:"artifact"`
	Steps    []Step   `json:"steps"`
}

// TODO: This is probably mostly what Drone uses. Maybe we should copy that struct :)
type Step struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Commands []string `json:"commands"`
}

type Artifact struct {
	URL string `json:"url"`
}
