package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/signalcd/signalcd/api/v1/models"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
)

const addr = "localhost:6660"

// SignalDB is the union of all necessary interfaces for the API
type SignalDB interface {
	DeploymentLister
	DeploymentStatusSetter
	CurrentDeploymentGetter
	CurrentDeploymentSetter
	PipelinesLister
	PipelineCreator
	StepLogsSaver
}

// Events to Deployments that should be sent via SSE (Server Sent Events)
type Events interface {
	SubscribeDeployments(channel chan signalcd.Deployment) signalcd.Subscription
	UnsubscribeDeployments(s signalcd.Subscription)
}

// NewV1 creates a new v1 API
func NewV1(logger log.Logger, db SignalDB, events Events) (*chi.Mux, error) {
	router := chi.NewRouter()

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	cert, err := ioutil.ReadFile("./development/signalcd.dev+6.pem")
	if err != nil {
		return nil, err
	}

	ok := pool.AppendCertsFromPEM(cert)
	if !ok {
		return nil, fmt.Errorf("failed to appened certificate")
	}

	var server *grpc.Server
	{
		opts := []grpc.ServerOption{
			grpc.Creds(credentials.NewClientTLSFromCert(pool, addr)),
		}

		server = grpc.NewServer(opts...)
		signalcdproto.RegisterUIServiceServer(server, &UIServer{
			db:     db,
			logger: logger,
		})
	}

	var mux *runtime.ServeMux
	{
		keyPair, err := tls.LoadX509KeyPair("./development/signalcd.dev+6.pem", "./development/signalcd.dev+6-key.pem")
		if err != nil {
			return nil, err
		}

		creds := credentials.NewTLS(&tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{keyPair},
		})
		opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

		mux = runtime.NewServeMux()
		err = signalcdproto.RegisterUIServiceHandlerFromEndpoint(context.Background(), mux, addr, opts)
		if err != nil {
			return nil, err
		}
	}

	router.Mount("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}))

	router.Get("/api/v1/deployments/events", deploymentEventsHandler(logger, events))

	return router, nil
}

type UIServer struct {
	db     SignalDB
	logger log.Logger
}

func (s *UIServer) ListDeployment(context.Context, *signalcdproto.ListDeploymentRequest) (*signalcdproto.ListDeploymentResponse, error) {
	list, err := s.db.ListDeployments()
	if err != nil {
		return nil, err
	}

	resp := &signalcdproto.ListDeploymentResponse{}

	for _, d := range list {
		dProto, err := deployment(d)
		if err != nil {
			return nil, err
		}
		resp.Deployments = append(resp.Deployments, dProto)
	}

	return resp, nil
}

func (s *UIServer) GetCurrentDeployment(context.Context, *signalcdproto.GetCurrentDeploymentRequest) (*signalcdproto.GetCurrentDeploymentResponse, error) {
	panic("implement me")
}

func (s *UIServer) SetCurrentDeployment(ctx context.Context, req *signalcdproto.SetCurrentDeploymentRequest) (*signalcdproto.SetCurrentDeploymentResponse, error) {
	p, err := s.db.GetPipeline(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	d, err := s.db.CreateDeployment(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	dProto, err := deployment(d)
	if err != nil {
		return nil, err
	}

	return &signalcdproto.SetCurrentDeploymentResponse{Deployment: dProto}, nil
}

func (s *UIServer) ListPipelines(context.Context, *signalcdproto.ListPipelinesRequest) (*signalcdproto.ListPipelinesResponse, error) {
	//pipelines, _ := s.db.ListPipelines()
	return &signalcdproto.ListPipelinesResponse{}, nil
}

func (s *UIServer) CreatePipeline(ctx context.Context, req *signalcdproto.CreatePipelineRequest) (*signalcdproto.CreatePipelineResponse, error) {
	p, err := pipeline(req.Pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed converting to internal pipeline: %w", err)
	}

	p, err = s.db.CreatePipeline(p)
	if err != nil {
		return nil, fmt.Errorf("failed creating the pipeline: %w", err)
	}

	protoPipeline, err := pipelineProto(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to gRPC pipeline: %w", err)
	}

	return &signalcdproto.CreatePipelineResponse{Pipeline: protoPipeline}, nil

}

func deployment(d signalcd.Deployment) (*signalcdproto.Deployment, error) {
	created, err := ptypes.TimestampProto(d.Created)
	if err != nil {
		return nil, err
	}
	started, err := ptypes.TimestampProto(d.Started)
	if err != nil {
		return nil, err
	}
	finished, err := ptypes.TimestampProto(d.Finished)
	if err != nil {
		return nil, err
	}
	p, err := pipelineProto(d.Pipeline)
	if err != nil {
		return nil, err
	}

	return &signalcdproto.Deployment{
		Number:   d.Number,
		Created:  created,
		Started:  started,
		Finished: finished,
		Status: &signalcdproto.DeploymentStatus{
			Phase: signalcdproto.DeploymentStatus_UNKNOWN,
		},
		Pipeline: p,
	}, nil
}

func pipeline(p *signalcdproto.Pipeline) (signalcd.Pipeline, error) {
	created, err := ptypes.Timestamp(p.GetCreated())
	if err != nil {
		return signalcd.Pipeline{}, err
	}

	steps := make([]signalcd.Step, len(p.GetSteps()))
	for _, s := range p.GetSteps() {
		steps = append(steps, signalcd.Step{
			Name:             s.GetName(),
			Image:            s.GetImage(),
			ImagePullSecrets: s.GetImagePullSecrets(),
			Commands:         s.GetCommands(),
		})
	}

	checks := make([]signalcd.Check, len(p.GetChecks()))
	for _, c := range p.GetChecks() {
		duration, err := ptypes.Duration(c.GetDuration())
		if err != nil {
			return signalcd.Pipeline{}, err
		}

		check := signalcd.Check{
			Name:             c.GetName(),
			Image:            c.GetImage(),
			ImagePullSecrets: c.GetImagePullSecrets(),
			Duration:         duration,
		}

		checks = append(checks, check)
	}

	return signalcd.Pipeline{
		ID:      p.GetId(),
		Name:    p.GetName(),
		Created: created,
		Steps:   steps,
		Checks:  checks,
	}, nil
}

func pipelineProto(p signalcd.Pipeline) (*signalcdproto.Pipeline, error) {
	created, err := ptypes.TimestampProto(p.Created)
	if err != nil {
		return nil, err
	}

	steps := make([]*signalcdproto.Step, len(p.Steps))
	for _, s := range p.Steps {
		steps = append(steps, &signalcdproto.Step{
			Name:             s.Name,
			Image:            s.Image,
			ImagePullSecrets: s.ImagePullSecrets,
			Commands:         s.Commands,
		})
	}

	checks := make([]*signalcdproto.Check, len(p.Checks))
	for _, c := range p.Checks {
		checks = append(checks, &signalcdproto.Check{
			Name:             c.Name,
			Image:            c.Image,
			ImagePullSecrets: c.ImagePullSecrets,
			Duration:         ptypes.DurationProto(c.Duration),
		})
	}

	return &signalcdproto.Pipeline{
		Id:      p.ID,
		Name:    p.Name,
		Created: created,
		Steps:   steps,
		Checks:  checks,
	}, nil
}

func (s *UIServer) GetPipeline(ctx context.Context, req *signalcdproto.GetPipelineRequest) (*signalcdproto.GetPipelineResponse, error) {
	p, err := s.db.GetPipeline(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline")
	}

	pProto, err := pipelineProto(p)
	if err != nil {
		return nil, err
	}

	return &signalcdproto.GetPipelineResponse{Pipeline: pProto}, nil
}

func deploymentEventsHandler(logger log.Logger, events Events) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		level.Debug(logger).Log("msg", "streaming deployment http connection just opened")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		deploymentEvents := make(chan signalcd.Deployment, 8)
		subscription := events.SubscribeDeployments(deploymentEvents)

		defer func() {
			events.UnsubscribeDeployments(subscription)
			level.Debug(logger).Log("msg", "streaming deployment http connection just closed")
		}()

		for {
			select {
			case <-ctx.Done():
				close(deploymentEvents)
				return
			case deployment := <-deploymentEvents:
				model := getModelsDeployment(deployment)
				j, err := json.Marshal(model)
				if err != nil {
					return // TODO
				}

				_, err = fmt.Fprintf(w, "data: %s\n\n", j)
				if err != nil {
					return // TODO
				}

				flusher.Flush()
			}
		}
	}
}

func getModelsPipeline(p signalcd.Pipeline) *models.Pipeline {
	mp := &models.Pipeline{
		ID:      strfmt.UUID(p.ID),
		Name:    p.Name,
		Created: strfmt.DateTime(p.Created),
		Steps:   []*models.Step{},
		Checks:  []*models.Check{},
	}

	for _, s := range p.Steps {
		imagePullSecrets := []string{}
		if len(s.ImagePullSecrets) > 0 {
			imagePullSecrets = s.ImagePullSecrets
		}

		ms := &models.Step{
			Name:             &s.Name,
			Image:            &s.Image,
			ImagePullSecrets: imagePullSecrets,
			Commands:         s.Commands,
		}

		if s.Status != nil {
			ms.Status = &models.StepStatus{
				Logs: string(s.Status.Logs),
			}
		}

		mp.Steps = append(mp.Steps, ms)
	}

	for _, c := range p.Checks {
		var env []*models.CheckEnvironmentItems0
		for key, value := range c.Environment {
			env = append(env, &models.CheckEnvironmentItems0{
				Key:   key,
				Value: value,
			})
		}

		imagePullSecrets := []string{}
		if len(c.ImagePullSecrets) > 0 {
			imagePullSecrets = c.ImagePullSecrets
		}

		mp.Checks = append(mp.Checks, &models.Check{
			Name:             &c.Name,
			Image:            &c.Image,
			ImagePullSecrets: imagePullSecrets,
			Duration:         c.Duration.Seconds(),
			Environment:      env,
		})
	}

	return mp
}

// PipelinesLister returns a list of Pipelines
type PipelinesLister interface {
	ListPipelines() ([]signalcd.Pipeline, error)
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

func getModelsDeployment(fd signalcd.Deployment) *models.Deployment {
	return &models.Deployment{
		Number:   &fd.Number,
		Created:  strfmt.DateTime(fd.Created),
		Started:  strfmt.DateTime(fd.Started),
		Finished: strfmt.DateTime(fd.Finished),
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

// CurrentDeploymentSetter gets a Pipeline and then creates a new Deployments
type CurrentDeploymentSetter interface {
	PipelineGetter
	CreateDeployment(signalcd.Pipeline) (signalcd.Deployment, error)
}

// PipelineGetter gets a new Pipeline
type PipelineGetter interface {
	GetPipeline(id string) (signalcd.Pipeline, error)
}

// PipelineCreator creates a new Pipeline
type PipelineCreator interface {
	CreatePipeline(signalcd.Pipeline) (signalcd.Pipeline, error)
}
