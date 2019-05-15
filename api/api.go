package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-openapi/loads"
	restmiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/signalcd/signalcd/api/v1/models"
	"github.com/signalcd/signalcd/api/v1/restapi"
	"github.com/signalcd/signalcd/api/v1/restapi/operations"
	"github.com/signalcd/signalcd/api/v1/restapi/operations/deployments"
	"github.com/signalcd/signalcd/api/v1/restapi/operations/pipeline"
	"github.com/signalcd/signalcd/signalcd"
	"golang.org/x/xerrors"
)

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

var fakeDeployments = []signalcd.Deployment{
	{
		Number:  4,
		Created: time.Now().Add(-30 * time.Second),
		Status: signalcd.DeploymentStatus{
			Phase: signalcd.Success,
		},
		Pipeline: fakePipelines[2],
	},
	{
		Number:  3,
		Created: time.Now().Add(-3 * time.Minute),
		Status: signalcd.DeploymentStatus{
			Phase: signalcd.Success,
		},
		Pipeline: fakePipelines[0],
	},
	{
		Number:  2,
		Created: time.Now().Add(-8 * time.Minute),
		Status: signalcd.DeploymentStatus{
			Phase: signalcd.Failed,
		},
		Pipeline: fakePipelines[1],
	},
	{
		Number:  1,
		Created: time.Now().Add(-10 * time.Minute),
		Status: signalcd.DeploymentStatus{
			Phase: signalcd.Success,
		},
		Pipeline: fakePipelines[0],
	},
}

const (
	PipelineCurrent       = "/pipeline"
	PipelineCurrentUpdate = "/pipeline/{id}"
	Pipelines             = "/pipelines"
	Pipeline              = "/pipelines/{id}"
	PipelinesStatus       = "/pipelines/status"
)

func NewV1() (*chi.Mux, error) {
	router := chi.NewRouter()

	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, xerrors.Errorf("failed to load embedded swagger file: %w", err.Error())
	}

	api := operations.NewCdAPI(swaggerSpec)

	// Skip the  redoc middleware, only serving the OpenAPI specification and
	// the API itself via RoutesHandler. See:
	// https://github.com/go-swagger/go-swagger/issues/1779
	api.Middleware = func(b restmiddleware.Builder) http.Handler {
		return restmiddleware.Spec("", swaggerSpec.Raw(), api.Context().RoutesHandler(b))
	}

	api.DeploymentsDeploymentsHandler = getDeploymentsHandler()
	api.DeploymentsCurrentDeploymentHandler = getCurrentDeploymentHandler()
	api.PipelinePipelineHandler = getPipelineHandler()
	api.PipelinePipelinesHandler = getPipelinesHandler()

	//api.PipelinePipelineAgentsHandler
	//api.PipelineUpdatePipelineAgentsHandler

	router.Mount("/", api.Serve(nil))

	return router, nil
}

func getModelsPipeline(p signalcd.Pipeline) *models.Pipeline {
	id := strfmt.UUID(p.ID)
	mp := &models.Pipeline{
		ID:   &id,
		Name: p.Name,
	}

	for _, s := range p.Steps {
		mp.Steps = append(mp.Steps, &models.Step{
			Name:     &s.Name,
			Image:    &s.Image,
			Commands: s.Commands,
		})
	}

	for _, c := range p.Checks {
		var env []*models.CheckEnvironmentItems0
		for key, value := range c.Environment {
			env = append(env, &models.CheckEnvironmentItems0{
				Key:   key,
				Value: value,
			})
		}

		mp.Checks = append(mp.Checks, &models.Check{
			Name:        &c.Name,
			Image:       &c.Image,
			Duration:    c.Duration.Seconds(),
			Environment: env,
		})
	}

	return mp
}

func getPipelinesHandler() pipeline.PipelinesHandlerFunc {
	return func(params pipeline.PipelinesParams) restmiddleware.Responder {
		var payload []*models.Pipeline

		for _, fp := range fakePipelines {
			payload = append(payload, getModelsPipeline(fp))
		}

		return pipeline.NewPipelinesOK().WithPayload(payload)
	}
}

func getDeploymentStatusPhase(phase signalcd.DeploymentPhase) string {
	switch phase {
	case signalcd.Success:
		return models.DeploymentstatusPhaseSuccess
	case signalcd.Failed:
		return models.DeploymentstatusPhaseFailed
	case signalcd.Progress:
		return models.DeploymentstatusPhaseProgress
	default:
		return models.DeploymentstatusPhaseUnknown
	}
}

func getDeploymentsHandler() deployments.DeploymentsHandlerFunc {
	return func(params deployments.DeploymentsParams) restmiddleware.Responder {
		var payload []*models.Deployment

		for _, fd := range fakeDeployments {
			number := fd.Number
			d := &models.Deployment{
				Number:   &number,
				Created:  strfmt.DateTime(fd.Created),
				Pipeline: getModelsPipeline(fd.Pipeline),
				Status: &models.Deploymentstatus{
					Phase: getDeploymentStatusPhase(fd.Status.Phase),
				},
			}
			payload = append(payload, d)
		}

		return deployments.NewDeploymentsOK().WithPayload(payload)
	}
}

func getCurrentDeploymentHandler() deployments.CurrentDeploymentHandlerFunc {
	return func(params deployments.CurrentDeploymentParams) restmiddleware.Responder {
		return nil
	}
}

func getPipelineHandler() pipeline.PipelineHandlerFunc {
	return func(params pipeline.PipelineParams) restmiddleware.Responder {
		p, err := getPipeline(params.ID)
		if err != nil {
			return pipeline.NewPipelineInternalServerError()
		}
		return pipeline.NewPipelineOK().WithPayload(getModelsPipeline(p))
	}
}

var agents = sync.Map{}

func pipelineAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var as []signalcd.AgentServer

		agents.Range(func(key, value interface{}) bool {
			as = append(as, value.(signalcd.AgentServer))
			return true
		})

		sort.Slice(as, func(i, j int) bool {
			return as[i].Agent.Name < as[j].Agent.Name
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
		var agent signalcd.AgentServer
		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			http.Error(w, "failed to decode", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		agent.Heartbeat = time.Now()

		agents.Store(agent.Agent.Name, agent)
	}
}
