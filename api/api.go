package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
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

// SignalDB is the union of all necessary interfaces for the API
type SignalDB interface {
	DeploymentLister
	DeploymentCreator
	DeploymentStatusSetter
	CurrentDeploymentGetter
	PipelinesLister
	PipelineCreator
}

// NewV1 creates a new v1 API
func NewV1(db SignalDB, logger log.Logger) (*chi.Mux, error) {
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

	api.DeploymentsDeploymentsHandler = getDeploymentsHandler(db)
	api.DeploymentsCurrentDeploymentHandler = getCurrentDeploymentHandler(db)
	api.DeploymentsSetCurrentDeploymentHandler = setCurrentDeploymentHandler(db, logger)
	api.PipelinePipelineHandler = getPipelineHandler(db)
	api.PipelinePipelinesHandler = getPipelinesHandler(db)
	api.PipelineCreateHandler = createPipelineHandler(db)

	router.Mount("/", api.Serve(nil))

	return router, nil
}

func getModelsPipeline(p signalcd.Pipeline) *models.Pipeline {
	mp := &models.Pipeline{
		ID:   strfmt.UUID(p.ID),
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

// PipelinesLister returns a list of Pipelines
type PipelinesLister interface {
	ListPipelines() ([]signalcd.Pipeline, error)
}

func getPipelinesHandler(lister PipelinesLister) pipeline.PipelinesHandlerFunc {
	return func(params pipeline.PipelinesParams) restmiddleware.Responder {
		var payload []*models.Pipeline

		pipelines, err := lister.ListPipelines()
		if err != nil {
			return pipeline.NewPipelinesInternalServerError()
		}

		for _, p := range pipelines {
			payload = append(payload, getModelsPipeline(p))
		}

		return pipeline.NewPipelinesOK().WithPayload(payload)
	}
}

func getDeploymentStatusPhase(phase signalcd.DeploymentPhase) string {
	switch phase {
	case signalcd.Success:
		return models.DeploymentstatusPhaseSuccess
	case signalcd.Failure:
		return models.DeploymentstatusPhaseFailure
	case signalcd.Progress:
		return models.DeploymentstatusPhaseProgress
	default:
		return models.DeploymentstatusPhaseUnknown
	}
}

// DeploymentLister lists all Deployments
type DeploymentLister interface {
	ListDeployments() ([]signalcd.Deployment, error)
}

func getDeploymentsHandler(lister DeploymentLister) deployments.DeploymentsHandlerFunc {
	return func(params deployments.DeploymentsParams) restmiddleware.Responder {
		var payload []*models.Deployment

		list, err := lister.ListDeployments()
		if err != nil {
			return deployments.NewDeploymentsInternalServerError()
		}

		for _, d := range list {
			payload = append(payload, getModelsDeployment(d))
		}

		return deployments.NewDeploymentsOK().WithPayload(payload)
	}
}

func getModelsDeployment(fd signalcd.Deployment) *models.Deployment {
	return &models.Deployment{
		Number:   &fd.Number,
		Created:  strfmt.DateTime(fd.Created),
		Pipeline: getModelsPipeline(fd.Pipeline),
		Status: &models.Deploymentstatus{
			Phase: getDeploymentStatusPhase(fd.Status.Phase),
		},
	}
}

// CurrentDeploymentGetter gets the current Deployment
type CurrentDeploymentGetter interface {
	GetCurrentDeployment() (signalcd.Deployment, error)
}

func getCurrentDeploymentHandler(getter CurrentDeploymentGetter) deployments.CurrentDeploymentHandlerFunc {
	return func(params deployments.CurrentDeploymentParams) restmiddleware.Responder {
		d, err := getter.GetCurrentDeployment()
		if err != nil {
			return deployments.NewSetCurrentDeploymentInternalServerError()
		}

		return deployments.NewCurrentDeploymentOK().WithPayload(getModelsDeployment(d))
	}
}

// DeploymentCreator gets a Pipeline and then creates a new Deployments
type DeploymentCreator interface {
	PipelineGetter
	CreateDeployment(signalcd.Pipeline) (signalcd.Deployment, error)
}

func setCurrentDeploymentHandler(creator DeploymentCreator, logger log.Logger) deployments.SetCurrentDeploymentHandlerFunc {
	return func(params deployments.SetCurrentDeploymentParams) restmiddleware.Responder {
		p, err := creator.GetPipeline(params.Pipeline)
		if err != nil {
			logger.Log("msg", "failed to get pipeline", "id", params.Pipeline, "err", err)
			return deployments.NewSetCurrentDeploymentInternalServerError()
		}

		d, err := creator.CreateDeployment(p)
		if err != nil {
			logger.Log("msg", "failed to create deployment", "err", err)
			return deployments.NewSetCurrentDeploymentInternalServerError()
		}

		return deployments.NewSetCurrentDeploymentOK().WithPayload(getModelsDeployment(d))
	}
}

// PipelineGetter gets a new Pipeline
type PipelineGetter interface {
	GetPipeline(id string) (signalcd.Pipeline, error)
}

func getPipelineHandler(getter PipelineGetter) pipeline.PipelineHandlerFunc {
	return func(params pipeline.PipelineParams) restmiddleware.Responder {
		p, err := getter.GetPipeline(params.ID)
		if err != nil {
			return pipeline.NewPipelineInternalServerError()
		}
		return pipeline.NewPipelineOK().WithPayload(getModelsPipeline(p))
	}
}

// PipelineCreator creates a new Pipeline
type PipelineCreator interface {
	CreatePipeline(signalcd.Pipeline) (signalcd.Pipeline, error)
}

func createPipelineHandler(creator PipelineCreator) pipeline.CreateHandlerFunc {
	return func(params pipeline.CreateParams) restmiddleware.Responder {
		p, err := creator.CreatePipeline(fromModelPipeline(params.Pipeline))
		if err != nil {
			return pipeline.NewCreateInternalServerError()
		}

		return pipeline.NewCreateOK().WithPayload(getModelsPipeline(p))
	}
}

func fromModelPipeline(m *models.Pipeline) signalcd.Pipeline {
	p := signalcd.Pipeline{
		ID:   m.ID.String(),
		Name: m.Name,
	}

	for _, c := range m.Checks {
		check := signalcd.Check{
			Name:     *c.Name,
			Image:    *c.Image,
			Duration: time.Duration(c.Duration) * time.Second,
		}

		for _, env := range c.Environment {
			check.Environment[env.Key] = env.Value
		}

		p.Checks = append(p.Checks, check)
	}

	for _, s := range m.Steps {
		step := signalcd.Step{
			Name:     *s.Name,
			Image:    *s.Image,
			Commands: s.Commands,
		}

		p.Steps = append(p.Steps, step)
	}

	return p
}
