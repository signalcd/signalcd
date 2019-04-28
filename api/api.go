package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	v1 "k8s.io/api/apps/v1"

	"github.com/go-chi/chi"
	"github.com/metalmatze/cd/cd"
)

var fakeCurrentWorkload = struct {
	mu       sync.RWMutex
	Workload cd.Workload
}{}

// Example, but real images as workload
var fakeWorkloads = []cd.Workload{
	{Image: "grafana/loki:master-e506f16"},
	{Image: "grafana/loki:master-e2b2561"},
	{Image: "grafana/loki:master-9440dc9"},
	{Image: "grafana/loki:master-80fdece"},
	{Image: "grafana/loki:master-199746a"},
}

func New() *chi.Mux {
	fakeCurrentWorkload.Workload = fakeWorkloads[0]

	router := chi.NewRouter()

	router.Get("/", index())
	router.Get("/workload", workload())
	router.Get("/workloads", workloads())
	router.Patch("/workloads/{index}", setWorkload())

	router.Get("/workloads/agents", workloadAgents())
	router.Post("/workloads/agents", updateWorkloadAgents())

	return router
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "index")
	}
}

func workloads() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(fakeWorkloads)
		if err != nil {
			http.Error(w, "failed to marshal", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bytes)
	}
}

func workload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fakeCurrentWorkload.mu.RLock()
		defer fakeCurrentWorkload.mu.RUnlock()

		bytes, err := json.Marshal(fakeCurrentWorkload.Workload)
		if err != nil {
			http.Error(w, "failed to marshal", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(bytes)
	}
}

func setWorkload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index := chi.URLParam(r, "index")
		i, err := strconv.Atoi(index)
		if err != nil {
			return
		}

		fakeCurrentWorkload.mu.Lock()
		defer fakeCurrentWorkload.mu.Unlock()
		// TODO: This will panic with unbound index
		fakeCurrentWorkload.Workload = fakeWorkloads[i]
	}
}

var agents = sync.Map{}

func workloadAgents() http.HandlerFunc {
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

func updateWorkloadAgents() http.HandlerFunc {
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
