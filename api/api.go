package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
	"github.com/metalmatze/cd/cd"
	v1 "k8s.io/api/apps/v1"
)

var fakeCurrentPipeline = struct {
	mu       sync.RWMutex
	Pipeline cd.Pipeline
}{}

var fakePipelines = []cd.Pipeline{
	{
		ID: "eee4047d-3826-4bf0-a7f1-b0b339521a52",
		Config: cd.Config{
			Steps: []cd.Step{
				{
					Name:     "echo",
					Image:    "alpine:3.7",
					Commands: []string{"echo 'hi'", "uname -a"},
				},
			},
		},
	},
	{
		ID: "6151e283-99b6-4611-bbc4-8aa4d3ddf8fd",
		Config: cd.Config{
			Steps: []cd.Step{
				{
					Name:     "echo",
					Image:    "alpine:3.6",
					Commands: []string{"echo 'hi'", "uname -a"},
				},
			},
		},
	},
}

func New() *chi.Mux {
	fakeCurrentPipeline.Pipeline = fakePipelines[0]

	router := chi.NewRouter()

	router.Get("/", index())
	router.Get("/pipeline", pipeline())
	router.Get("/pipelines", pipelines())
	router.Patch("/pipelines/{index}", setPipeline())

	router.Get("/pipelines/agents", pipelineAgents())
	router.Post("/pipelines/agents", updatePipelineAgents())

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

func setPipeline() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index := chi.URLParam(r, "index")
		i, err := strconv.Atoi(index)
		if err != nil {
			return
		}

		fakeCurrentPipeline.mu.Lock()
		defer fakeCurrentPipeline.mu.Unlock()
		// TODO: This will panic with unbound index
		fakeCurrentPipeline.Pipeline = fakePipelines[i]
	}
}

var agents = sync.Map{}

func pipelineAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var as []cd.Agent

		agents.Range(func(key, value interface{}) bool {
			as = append(as, cd.Agent{
				Name:   key.(string),
				Status: value.(v1.DeploymentStatus),
			})

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
		var agent cd.Agent
		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			http.Error(w, "failed to decode", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		agents.Store(agent.Name, agent.Status)
	}
}
