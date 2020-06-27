package boltdb

import (
	"github.com/signalcd/signalcd/signalcd"
)

// Events wraps BoltDB to publish updated structs
type Events struct {
	*BoltDB
	events *signalcd.Events
}

// NewEvents creates a new BoltDB wrapper that publishes events
func NewEvents(db *BoltDB, events *signalcd.Events) *Events {
	return &Events{
		BoltDB: db,
		events: events,
	}
}

// CreateDeployment wraps the underlying BoltDB func to publish successfully created Deployments
func (e *Events) CreateDeployment(pipeline signalcd.Pipeline) (signalcd.Deployment, error) {
	deployment, err := e.BoltDB.CreateDeployment(pipeline)
	if err != nil {
		return deployment, err
	}

	e.events.PublishDeployment(deployment)
	return deployment, nil
}

func (e *Events) UpdateDeploymentStatus(deploymentNumber int64, step int64, agent string, phase signalcd.Phase) (signalcd.Deployment, error) {
	deployment, err := e.BoltDB.UpdateDeploymentStatus(deploymentNumber, step, agent, phase)
	if err != nil {
		return deployment, err
	}

	e.events.PublishDeployment(deployment)
	return deployment, nil
}
