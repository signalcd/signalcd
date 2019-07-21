package main

import (
	"testing"

	"github.com/signalcd/signalcd/signalcd"
)

func TestCurrentDeployment(t *testing.T) {
	u := updater{}

	d := u.currentDeployment.get()
	if d != nil {
		t.Errorf("expected current deployment to be nil, got %v", d)
	}

	u.currentDeployment.set(signalcd.Deployment{Number: 42})

	d = u.currentDeployment.get()
	if d == nil {
		t.Fatalf("expected current deployment not to be nil")
	}
	if d.Number != 42 {
		t.Errorf("expected current deployment number to be 42, got %d", d.Number)
	}
}
