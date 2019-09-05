package signalcd

import (
	"sync"
)

type Events struct {
	deploymentSubscribers map[int64]chan Deployment
	deploymentLock        sync.RWMutex

	pipelineSubscribers map[int64]chan Pipeline
	pipelineLock        sync.RWMutex
}

func NewEvents() *Events {
	return &Events{
		deploymentSubscribers: make(map[int64]chan Deployment),
		pipelineSubscribers:   make(map[int64]chan Pipeline),
	}
}

type Subscription struct {
	id int64
}

func (e *Events) SubscribeDeployments(channel chan Deployment) Subscription {
	e.deploymentLock.Lock()
	defer e.deploymentLock.Unlock()

	id := int64(len(e.deploymentSubscribers)) + 1
	e.deploymentSubscribers[id] = channel
	return Subscription{id: id}
}

func (e *Events) UnsubscribeDeployments(s Subscription) {
	e.deploymentLock.Lock()
	defer e.deploymentLock.Unlock()

	delete(e.deploymentSubscribers, s.id)
}

func (e *Events) PublishDeployment(d Deployment) {
	e.deploymentLock.RLock()
	defer e.deploymentLock.RUnlock()

	for _, s := range e.deploymentSubscribers {
		s <- d
	}
}
