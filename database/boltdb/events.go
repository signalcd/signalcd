package boltdb

import (
	"context"

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

	if err == nil {
		e.events.PublishDeployment(deployment)
	}

	return deployment, err
}

// SetDeploymentStatus wraps t he underlying BoltDB func to publish successfully updated Deployments
func (e *Events) SetDeploymentStatus(ctx context.Context, number int64, phase signalcd.DeploymentPhase) (signalcd.Deployment, error) {
	deployment, err := e.BoltDB.SetDeploymentStatus(ctx, number, phase)

	if err == nil {
		e.events.PublishDeployment(deployment)
	}

	return deployment, err
}
