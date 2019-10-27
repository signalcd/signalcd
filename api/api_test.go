package api

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/signalcd/signalcd/signalcd"
	"github.com/stretchr/testify/assert"
)

func Test_getModelsPipeline(t *testing.T) {
	var testcases = []struct {
		Name     string
		Pipeline signalcd.Pipeline
		JSON     string
	}{
		{
			Name:     "Empty",
			Pipeline: signalcd.Pipeline{},
			JSON:     `{"steps":[],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
		{
			Name:     "ID",
			Pipeline: signalcd.Pipeline{ID: "5f152ab1-605a-4f0f-abab-e864d26d07c8"},
			JSON:     `{"id":"5f152ab1-605a-4f0f-abab-e864d26d07c8","steps":[],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
		{
			Name:     "Name",
			Pipeline: signalcd.Pipeline{Name: "Name"},
			JSON:     `{"name":"Name","steps":[],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
		{
			Name:     "Created",
			Pipeline: signalcd.Pipeline{Created: time.Date(2019, 10, 27, 18, 41, 21, 0, time.UTC)},
			JSON:     `{"steps":[],"checks":[],"created":"2019-10-27T18:41:21.000Z"}`,
		},
		{
			Name: "StepEmpty",
			Pipeline: signalcd.Pipeline{Steps: []signalcd.Step{{
				Name:             "Step Name",
				Image:            "someImage",
				ImagePullSecrets: nil,
				Commands:         nil,
				Status:           nil,
			}}},
			JSON: `{"steps":[{"name":"Step Name","image":"someImage","commands":null,"imagePullSecrets":[]}],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
		{
			Name: "StepImagePullSecrets",
			Pipeline: signalcd.Pipeline{Steps: []signalcd.Step{{
				Name:             "Step Name",
				Image:            "someImage",
				ImagePullSecrets: []string{"somepullsecret"},
				Commands:         nil,
				Status:           nil,
			}}},
			JSON: `{"steps":[{"imagePullSecrets":["somepullsecret"],"name":"Step Name","image":"someImage","commands":null}],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
		{
			Name: "StepCommands",
			Pipeline: signalcd.Pipeline{Steps: []signalcd.Step{{
				Name:             "Step Name",
				Image:            "someImage",
				ImagePullSecrets: nil,
				Commands:         []string{"ls"},
				Status:           nil,
			}}},
			JSON: `{"steps":[{"commands":["ls"],"name":"Step Name","image":"someImage","imagePullSecrets":[]}],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
		{
			Name: "StepStatus",
			Pipeline: signalcd.Pipeline{Steps: []signalcd.Step{{
				Name:             "Step Name",
				Image:            "someImage",
				ImagePullSecrets: nil,
				Commands:         nil,
				Status: &signalcd.Status{
					Logs: []byte("some\nlogs"),
				},
			}}},
			JSON: `{"steps":[{"status":{"logs":"some\nlogs"},"name":"Step Name","image":"someImage","commands":null,"imagePullSecrets":[]}],"checks":[],"created":"0001-01-01T00:00:00.000Z"}`,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			m := getModelsPipeline(tc.Pipeline)
			j, err := json.Marshal(m)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.JSON, string(j))
		})
	}
}
