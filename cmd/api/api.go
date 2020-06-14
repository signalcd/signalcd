package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	openapi "github.com/signalcd/signalcd/api/server/go/go"
	"github.com/signalcd/signalcd/signalcd"
)

type Database interface {
	DeploymentLister
	CurrentDeploymentGetter
	PipelineLister
	PipelineGetter
}

func NewV1(registry *prometheus.Registry, database Database, events *signalcd.Events) (http.Handler, error) {
	instrument := instrument(registry)

	routes := []openapi.Router{
		openapi.NewDeploymentApiController(&Deployments{lister: database, currentGetter: database}),
		openapi.NewPipelineApiController(&Pipelines{lister: database, getter: database}),
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, api := range routes {
		for _, route := range api.Routes() {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(instrument(
					route.HandlerFunc,
					route.Name,
				))
		}
	}

	return router, nil
}

func instrument(r *prometheus.Registry) func(next http.Handler, name string) http.Handler {
	requests := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "openapi_http_request_duration_seconds",
		Help: "Latency for HTTP calls to OpenAPI handlers",
	}, []string{"code", "method", "name"})
	r.MustRegister(requests)

	return func(next http.Handler, name string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			requests.WithLabelValues(
				fmt.Sprintf("%d", ww.Status()),
				r.Method,
				name,
			).Observe(time.Since(start).Seconds())
		})
	}
}

type Deployments struct {
	lister        DeploymentLister
	currentGetter CurrentDeploymentGetter
}

type DeploymentLister interface {
	ListDeployments() ([]signalcd.Deployment, error)
}

func (d *Deployments) ListDeployments() (interface{}, error) {
	deployments, err := d.lister.ListDeployments()
	if err != nil {
		return nil, err
	}

	var deploys []openapi.Deployment
	for _, d := range deployments {
		deploys = append(deploys, deploymentOpenAPI(d))
	}

	return deploys, nil
}

type CurrentDeploymentGetter interface {
	GetCurrentDeployment() (signalcd.Deployment, error)
}

func (d *Deployments) GetCurrentDeployment() (interface{}, error) {
	deployment, err := d.currentGetter.GetCurrentDeployment()
	if err != nil {
		return nil, err
	}
	return deploymentOpenAPI(deployment), nil
}

func deploymentOpenAPI(d signalcd.Deployment) openapi.Deployment {
	return openapi.Deployment{
		Number:   d.Number,
		Created:  d.Created,
		Started:  d.Started,
		Finished: d.Finished,
	}
}

func (d *Deployments) SetCurrentDeployment(params openapi.InlineObject) (interface{}, error) {
	fmt.Println(params.PipelineID)
	return nil, nil
}

type Pipelines struct {
	lister PipelineLister
	getter PipelineGetter
}

type PipelineLister interface {
	ListPipelines() ([]signalcd.Pipeline, error)
}

func (p *Pipelines) ListPipelines() (interface{}, error) {
	pipelines, err := p.lister.ListPipelines()
	if err != nil {
		return nil, err
	}

	var pipes []openapi.Pipeline
	for _, p := range pipelines {
		pipes = append(pipes, pipelineOpenAPI(p))
	}

	return pipes, nil
}

type PipelineGetter interface {
	GetPipeline(id string) (signalcd.Pipeline, error)
}

func (p *Pipelines) GetPipeline(id string) (interface{}, error) {
	pipeline, err := p.getter.GetPipeline(id)
	return pipelineOpenAPI(pipeline), err
}

func pipelineOpenAPI(p signalcd.Pipeline) openapi.Pipeline {
	return openapi.Pipeline{
		Id:      p.ID,
		Name:    p.Name,
		Created: p.Created,
	}
}
