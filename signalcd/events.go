package signalcd

import (
	"sync"
)

// Events fan out updates to Deployments to all subscribers
type Events struct {
	deploymentSubscribers map[int64]chan Deployment
	deploymentLock        sync.RWMutex
}

// NewEvents creates a new Event bus for SignalCD
func NewEvents() *Events {
	return &Events{
		deploymentSubscribers: make(map[int64]chan Deployment),
	}
}

// Subscription identifies a specific registered channel
type Subscription struct {
	id int64
}

// SubscribeDeployments adds a channel to subscribers and  returns a Subscription
func (e *Events) SubscribeDeployments(channel chan Deployment) Subscription {
	e.deploymentLock.Lock()
	defer e.deploymentLock.Unlock()

	id := int64(len(e.deploymentSubscribers)) + 1
	e.deploymentSubscribers[id] = channel
	return Subscription{id: id}
}

// UnsubscribeDeployments removes a Subscription/channel from the subscribers list
func (e *Events) UnsubscribeDeployments(s Subscription) {
	e.deploymentLock.Lock()
	defer e.deploymentLock.Unlock()

	delete(e.deploymentSubscribers, s.id)
}

// PublishDeployment fans out an updated Deployment to all active subscribers
func (e *Events) PublishDeployment(d Deployment) {
	e.deploymentLock.RLock()
	defer e.deploymentLock.RUnlock()

	for _, s := range e.deploymentSubscribers {
		s <- d
	}
}
