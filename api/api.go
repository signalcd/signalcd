package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"sync"

	"github.com/go-chi/chi/middleware"

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
		ID:       "eee4047d-3826-4bf0-a7f1-b0b339521a52",
		Artifact: cd.Artifact{URL: "cheese0.tar.gz"},
		Steps: []cd.Step{
			{
				Name:     "cheese0",
				Image:    "kubeciio/kubectl",
				Commands: []string{"kubectl version"},
			},
		},
	},
	{
		ID:       "6151e283-99b6-4611-bbc4-8aa4d3ddf8fd",
		Artifact: cd.Artifact{URL: "cheese1.tar.gz"},
		Steps: []cd.Step{
			{
				Name:     "cheese1",
				Image:    "kubeciio/kubectl",
				Commands: []string{"kubectl version"},
			},
		},
	},
	{
		ID:       "a7cae189-400e-4d8c-a982-f0e9a5b4901f",
		Artifact: cd.Artifact{URL: "cheese2.tar.gz"},
		Steps: []cd.Step{
			{
				Name:     "cheese2",
				Image:    "kubeciio/kubectl",
				Commands: []string{"kubectl version"},
			},
		},
	},
}

func getPipeline(id string) (cd.Pipeline, error) {
	for _, p := range fakePipelines {
		if p.ID == id {
			return p, nil
		}
	}
	return cd.Pipeline{}, fmt.Errorf("pipeline not found")
}

const (
	PipelineCurrent       = "/pipeline"
	PipelineCurrentUpdate = "/pipeline/{id}"
	Pipelines             = "/pipelines"
	Pipeline              = "/pipelines/{id}"
	PipelineArtifact      = "/pipelines/{id}/artifacts/{name}"
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
	router.Get(PipelineArtifact, pipelineArtifact())
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

func pipelineArtifact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		artifact := chi.URLParam(r, "name")

		_, err := getPipeline(id)
		if err != nil {
			http.Error(w, "pipeline not found", http.StatusNotFound)
			return
		}

		path := filepath.Join("examples", artifact)
		fmt.Println(path)
		http.ServeFile(w, r, path)
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
