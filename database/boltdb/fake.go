package boltdb

import (
	"time"

	"github.com/signalcd/signalcd/signalcd"
)

var fakePipelines = []signalcd.Pipeline{
	{
		ID:   "eee4047d-3826-4bf0-a7f1-b0b339521a52",
		Name: "cheese0",
		Steps: []signalcd.Step{
			{
				Name:     "cheese0",
				Image:    "quay.io/signalcd/examples:cheese0",
				Commands: []string{"kubectl apply -f /data"},
			},
		},
		Checks:  fakeChecks,
		Created: time.Now().Add(-10 * time.Minute),
	},
	{
		ID:   "6151e283-99b6-4611-bbc4-8aa4d3ddf8fd",
		Name: "cheese1",
		Steps: []signalcd.Step{
			{
				Name:     "cheese1",
				Image:    "quay.io/signalcd/examples:cheese1",
				Commands: []string{"kubectl apply -f /data"},
			},
		},
		Checks:  fakeChecks,
		Created: time.Now().Add(-8 * time.Minute),
	},
	{
		ID:   "a7cae189-400e-4d8c-a982-f0e9a5b4901f",
		Name: "cheese2",
		Steps: []signalcd.Step{
			{
				Name:     "cheese2",
				Image:    "quay.io/signalcd/examples:cheese2",
				Commands: []string{"kubectl apply -f /data"},
			},
		},
		Checks:  fakeChecks,
		Created: time.Now().Add(-45 * time.Second),
	},
}

var fakeChecks = []signalcd.Check{
	{
		Name:     "kubernetes-status",
		Image:    "quay.io/signalcd/kubernetes-status",
		Duration: time.Minute,
		Environment: map[string]string{
			"PLUGIN_LABELS": "app=cheese",
		},
	},
}
