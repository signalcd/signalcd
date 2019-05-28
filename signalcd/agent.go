package signalcd

import (
	"encoding/json"
	"time"
)

// Agent is a daemon responsible for running the pipelines in namespaces
type Agent struct {
	Name     string   `json:"name"`
	Pipeline Pipeline `json:"pipeline"`
}

// AgentServer is the Agent plus the timestamp of a last heartbeat
type AgentServer struct {
	Agent `json:",inline"`

	Heartbeat time.Time `json:"-"`
}

// Ready returns if an Agent is ready.
// If ready an agent has reported to the API in the past 15s.
func (as AgentServer) Ready() bool {
	return time.Since(as.Heartbeat) < 15*time.Second
}

// MarshalJSON implements the Marshaler interface to return
// ready boolean from last heartbeat
func (as AgentServer) MarshalJSON() ([]byte, error) {
	s := struct {
		Agent `json:",inline"`
		Ready bool `json:"ready"`
	}{
		Agent: as.Agent,
		Ready: as.Ready(),
	}

	return json.Marshal(s)
}
