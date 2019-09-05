package boltdb

import (
	"context"

	"github.com/signalcd/signalcd/signalcd"
)

type Events struct {
	*BoltDB
	events *signalcd.Events
}

func NewEvents(db *BoltDB, events *signalcd.Events) *Events {
	return &Events{
		BoltDB: db,
		events: events,
	}
}

func (e *Events) CreateDeployment(pipeline signalcd.Pipeline) (signalcd.Deployment, error) {
	deployment, err := e.BoltDB.CreateDeployment(pipeline)

	if err == nil {
		e.events.PublishDeployment(deployment)
	}

	return deployment, err
}

func (e *Events) SetDeploymentStatus(ctx context.Context, number int64, phase signalcd.DeploymentPhase) (signalcd.Deployment, error) {
	deployment, err := e.BoltDB.SetDeploymentStatus(ctx, number, phase)

	if err == nil {
		e.events.PublishDeployment(deployment)
	}

	return deployment, err
}
