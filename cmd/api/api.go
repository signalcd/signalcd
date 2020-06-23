package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	openapi "github.com/signalcd/signalcd/api/server/go/go"
	"github.com/signalcd/signalcd/signalcd"
)

type Database interface {
	DeploymentLister
	CurrentDeploymentGetter
	CurrentDeploymentSetter
	PipelineLister
	PipelineGetter
	PipelineCreater
}

func NewV1(logger log.Logger, registry *prometheus.Registry, database Database, events Events) (http.Handler, error) {
	instrument := instrument(registry)

	routes := []openapi.Router{
		openapi.NewDeploymentApiController(&Deployments{
			lister:        database,
			currentGetter: database,
			currentSetter: database,
		}),
		openapi.NewPipelineApiController(&Pipelines{
			lister:  database,
			getter:  database,
			creator: database,
		}),
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, api := range routes {
		for _, route := range api.Routes() {
			router.Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(instrument(
					route.HandlerFunc,
					route.Name,
				))
		}
	}

	router.Methods(http.MethodGet).
		Path("/api/v1/deployments/events").
		HandlerFunc(deploymentEventsHandler(logger, registry, events))

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
	currentSetter CurrentDeploymentSetter
}

type DeploymentLister interface {
	ListDeployments() ([]signalcd.Deployment, error)
}

func (d *Deployments) ListDeployments() (interface{}, error) {
	deployments, err := d.lister.ListDeployments()
	if err != nil {
		return nil, err
	}

	deploys := []openapi.Deployment{}
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

// CurrentDeploymentSetter gets a Pipeline and then creates a new Deployments
type CurrentDeploymentSetter interface {
	PipelineGetter
	CreateDeployment(signalcd.Pipeline) (signalcd.Deployment, error)
}

func (d *Deployments) SetCurrentDeployment(params openapi.SetCurrentDeployment) (interface{}, error) {
	pipeline, err := d.currentSetter.GetPipeline(params.PipelineID)
	if err != nil {
		return nil, err
	}

	deployment, err := d.currentSetter.CreateDeployment(pipeline)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func deploymentOpenAPI(d signalcd.Deployment) openapi.Deployment {
	return openapi.Deployment{
		Number:   d.Number,
		Created:  d.Created,
		Pipeline: pipelineOpenAPI(d.Pipeline),
	}
}

type Pipelines struct {
	lister  PipelineLister
	getter  PipelineGetter
	creator PipelineCreater
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

type PipelineCreater interface {
	CreatePipeline(signalcd.Pipeline) (signalcd.Pipeline, error)
}

func (p *Pipelines) CreatePipeline(newPipeline openapi.Pipeline) (interface{}, error) {
	// TODO: Validate pipeline input
	pipeline, err := p.creator.CreatePipeline(pipelineSignalCD(newPipeline))
	if err != nil {
		return nil, err
	}

	return pipelineOpenAPI(pipeline), nil
}

func pipelineOpenAPI(p signalcd.Pipeline) openapi.Pipeline {
	var steps []openapi.PipelineSteps
	for _, s := range p.Steps {
		steps = append(steps, openapi.PipelineSteps{
			Name:     s.Name,
			Image:    s.Image,
			Commands: s.Commands,
		})
	}

	return openapi.Pipeline{
		Id:      p.ID,
		Name:    p.Name,
		Created: p.Created,
		Steps:   steps,
	}
}

func pipelineSignalCD(p openapi.Pipeline) signalcd.Pipeline {
	var steps []signalcd.Step
	for _, s := range p.Steps {
		steps = append(steps, signalcd.Step{
			Name:     s.Name,
			Image:    s.Image,
			Commands: s.Commands,
		})
	}

	return signalcd.Pipeline{
		ID:      p.Id,
		Name:    p.Name,
		Created: p.Created,
		Steps:   steps,
	}
}
