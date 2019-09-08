package signalcd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventsDeploymentsSubscription(t *testing.T) {
	events := NewEvents()
	assert.Len(t, events.deploymentSubscribers, 0)
	subs1 := events.SubscribeDeployments(make(chan Deployment))
	assert.Len(t, events.deploymentSubscribers, 1)
	subs2 := events.SubscribeDeployments(make(chan Deployment))
	assert.Len(t, events.deploymentSubscribers, 2)
	events.UnsubscribeDeployments(subs1)
	assert.Len(t, events.deploymentSubscribers, 1)
	events.UnsubscribeDeployments(subs2)
	assert.Len(t, events.deploymentSubscribers, 0)
}

func TestEventsDeploymentOneSubscriber(t *testing.T) {
	events := NewEvents()
	channel := make(chan Deployment, 1)
	_ = events.SubscribeDeployments(channel)

	events.PublishDeployment(Deployment{Number: 123})

	d := <-channel
	assert.Equal(t, int64(123), d.Number)
}

func TestEventsDeploymentMultipleSubscribers(t *testing.T) {
	events := NewEvents()
	channel1 := make(chan Deployment, 1)
	_ = events.SubscribeDeployments(channel1)
	channel2 := make(chan Deployment, 1)
	_ = events.SubscribeDeployments(channel2)

	events.PublishDeployment(Deployment{Number: 123})

	d1 := <-channel1
	assert.Equal(t, int64(123), d1.Number)

	d2 := <-channel2
	assert.Equal(t, int64(123), d2.Number)

	events.PublishDeployment(Deployment{Number: 666})

	d1 = <-channel1
	assert.Equal(t, int64(666), d1.Number)

	d2 = <-channel2
	assert.Equal(t, int64(666), d2.Number)
}
