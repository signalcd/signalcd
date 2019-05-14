package signalcd

import "time"

type PipelineStatus string

const (
	Unknown  PipelineStatus = "unknown"
	Success  PipelineStatus = "success"
	Failed   PipelineStatus = "failed"
	Progress PipelineStatus = "progress"
)

type Agent struct {
	Name      string         `json:"name"`
	Status    PipelineStatus `json:"status"`
	Heartbeat time.Time      `json:"heartbeat"`

	Pipeline Pipeline `json:"pipeline"`
}
