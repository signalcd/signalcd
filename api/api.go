package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/signalcd/signalcd/signalcd"
)

var fakeCurrentPipeline = struct {
	mu       sync.RWMutex
	Pipeline signalcd.Pipeline
}{}

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
		Checks: fakeChecks,
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
		Checks: fakeChecks,
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
		Checks: fakeChecks,
	},
}

func getPipeline(id string) (signalcd.Pipeline, error) {
	for _, p := range fakePipelines {
		if p.ID == id {
			return p, nil
		}
	}
	return signalcd.Pipeline{}, fmt.Errorf("pipeline not found")
}

const (
	PipelineCurrent       = "/pipeline"
	PipelineCurrentUpdate = "/pipeline/{id}"
	Pipelines             = "/pipelines"
	Pipeline              = "/pipelines/{id}"
	PipelinesStatus       = "/pipelines/status"
)

func New() *chi.Mux {
	fakeCurrentPipeline.Pipeline = fakePipelines[0]

	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)

	router.Get("/", index())

	router.Get(PipelineCurrent, pipelineCurrent())
	router.Patch(PipelineCurrentUpdate, updateCurrentPipeline())
	router.Get(Pipelines, pipelines())
	router.Get(Pipeline, pipeline())
	router.Get(PipelinesStatus, pipelineAgents())
	router.Post(PipelinesStatus, updatePipelineAgents())

	return router
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "index")
	}
}

func pipelines() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(fakePipelines)
		if err != nil {
			http.Error(w, "failed to marshal", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bytes)
	}
}

func pipeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		p, err := getPipeline(id)
		if err != nil {
			http.Error(w, "pipeline not found", http.StatusNotFound)
			return
		}

		payload, err := json.Marshal(p)
		if err != nil {
			http.Error(w, "failed to marshal", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(payload)
	}
}

func pipelineCurrent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fakeCurrentPipeline.mu.RLock()
		defer fakeCurrentPipeline.mu.RUnlock()

		bytes, err := json.Marshal(fakeCurrentPipeline.Pipeline)
		if err != nil {
			http.Error(w, "failed to marshal", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bytes)
	}
}

func updateCurrentPipeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		p, err := getPipeline(id)
		if err != nil {
			http.Error(w, "pipeline not found", http.StatusNotFound)
			return
		}

		fakeCurrentPipeline.mu.Lock()
		fakeCurrentPipeline.Pipeline = p
		fakeCurrentPipeline.mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	}
}

var agents = sync.Map{}

func pipelineAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var as []signalcd.Agent

		agents.Range(func(key, value interface{}) bool {
			as = append(as, value.(signalcd.Agent))
			return true
		})

		sort.Slice(as, func(i, j int) bool {
			return as[i].Name < as[j].Name
		})

		payload, err := json.Marshal(as)
		if err != nil {
			http.Error(w, "failed to marshal", http.StatusInternalServerError)
			return
		}

		w.Write(payload)
	}
}

func updatePipelineAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var agent signalcd.Agent
		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			http.Error(w, "failed to decode", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		agent.Heartbeat = time.Now()

		agents.Store(agent.Name, agent)
	}
}
