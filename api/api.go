package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_gokit "github.com/grpc-ecosystem/go-grpc-middleware/logging/kit"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

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
func NewV1(logger log.Logger, db SignalDB, events Events, addr string, certPath string, keyPath string) (*chi.Mux, error) {
	router := chi.NewRouter()

	var server *grpc.Server
	{
		var opts []grpc.ServerOption

		if certPath != "" && keyPath != "" {
			creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
			if err != nil {
				return nil, fmt.Errorf("failed to create credentials: %w", err)
			}

			opts = append(opts, grpc.Creds(creds))
		}

		opts = append(opts, grpc_middleware.WithUnaryServerChain(
			// TODO: we can still improve the key-value pairs logged to reduce noise
			grpc_gokit.UnaryServerInterceptor(logger,
				grpc_gokit.WithLevels(LoggerGRPC()),
			),
		))

		server = grpc.NewServer(opts...)
		signalcdproto.RegisterUIServiceServer(server, &UIServer{
			db:     db,
			logger: logger,
		})
	}

	var mux *runtime.ServeMux
	{
		var opts []grpc.DialOption

		if certPath != "" && keyPath != "" {
			creds, err := credentials.NewClientTLSFromFile(certPath, "localhost")
			if err != nil {
				return nil, fmt.Errorf("failed to load credentials: %w", err)
			}

			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithInsecure())
		}

		mux = runtime.NewServeMux()
		err := signalcdproto.RegisterUIServiceHandlerFromEndpoint(context.Background(), mux, addr, opts)
		if err != nil {
			return nil, err
		}
	}

	handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}), &http2.Server{})

	router.Mount("/", handler)
	router.Get("/api/v1/deployments/events", deploymentEventsHandler(logger, events))

	return router, nil
}

type UIServer struct {
	db     SignalDB
	logger log.Logger
}

// DeploymentLister lists all Deployments
type DeploymentLister interface {
	ListDeployments() ([]signalcd.Deployment, error)
}

func (s *UIServer) ListDeployment(context.Context, *signalcdproto.ListDeploymentRequest) (*signalcdproto.ListDeploymentResponse, error) {
	list, err := s.db.ListDeployments()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to list deployments: %w", err)
	}

	resp := &signalcdproto.ListDeploymentResponse{}

	for _, d := range list {
		dProto, err := signalcdproto.DeploymentProto(d)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert deployment to proto: %w", err)
		}
		resp.Deployments = append(resp.Deployments, dProto)
	}

	return resp, nil
}

// CurrentDeploymentGetter gets the current Deployment
type CurrentDeploymentGetter interface {
	GetCurrentDeployment() (signalcd.Deployment, error)
}

func (s *UIServer) GetCurrentDeployment(context.Context, *signalcdproto.GetCurrentDeploymentRequest) (*signalcdproto.GetCurrentDeploymentResponse, error) {
	panic("implement me")
}

// CurrentDeploymentSetter gets a Pipeline and then creates a new Deployments
type CurrentDeploymentSetter interface {
	PipelineGetter
	CreateDeployment(signalcd.Pipeline) (signalcd.Deployment, error)
}

func (s *UIServer) SetCurrentDeployment(ctx context.Context, req *signalcdproto.SetCurrentDeploymentRequest) (*signalcdproto.SetCurrentDeploymentResponse, error) {
	p, err := s.db.GetPipeline(req.GetId())
	if err != nil {

		return nil, status.Errorf(codes.NotFound, "failed to get pipeline: %w", err)
	}

	d, err := s.db.CreateDeployment(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create deployment: %w", err)
	}

	dProto, err := signalcdproto.DeploymentProto(d)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert deployment to proto: %w", err)
	}

	return &signalcdproto.SetCurrentDeploymentResponse{Deployment: dProto}, nil
}

// PipelinesLister returns a list of Pipelines
type PipelinesLister interface {
	ListPipelines() ([]signalcd.Pipeline, error)
}

func (s *UIServer) ListPipelines(context.Context, *signalcdproto.ListPipelinesRequest) (*signalcdproto.ListPipelinesResponse, error) {
	pipelines, err := s.db.ListPipelines()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to list pipelines: %w", err)
	}

	psProto := make([]*signalcdproto.Pipeline, len(pipelines))
	for i, p := range pipelines {
		pProto, err := signalcdproto.PipelineProto(p)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert pipeline to proto: %w", err)
		}
		psProto[i] = pProto
	}

	return &signalcdproto.ListPipelinesResponse{Pipelines: psProto}, nil
}

// PipelineCreator creates a new Pipeline
type PipelineCreator interface {
	CreatePipeline(signalcd.Pipeline) (signalcd.Pipeline, error)
}

func (s *UIServer) CreatePipeline(ctx context.Context, req *signalcdproto.CreatePipelineRequest) (*signalcdproto.CreatePipelineResponse, error) {
	p, err := signalcdproto.PipelineSignalCD(req.Pipeline)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed converting to internal pipeline: %w", err)
	}

	p, err = s.db.CreatePipeline(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed creating the pipeline: %w", err)
	}

	protoPipeline, err := signalcdproto.PipelineProto(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert to gRPC pipeline: %w", err)
	}

	return &signalcdproto.CreatePipelineResponse{Pipeline: protoPipeline}, nil

}

// PipelineGetter gets a new Pipeline
type PipelineGetter interface {
	GetPipeline(id string) (signalcd.Pipeline, error)
}

func (s *UIServer) GetPipeline(ctx context.Context, req *signalcdproto.GetPipelineRequest) (*signalcdproto.GetPipelineResponse, error) {
	p, err := s.db.GetPipeline(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get pipeline")
	}

	pProto, err := signalcdproto.PipelineProto(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert pipeline to proto: %w", err)
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
				model, err := signalcdproto.DeploymentProto(deployment)
				if err != nil {
					level.Warn(logger).Log("msg", "failed to convert deployment to proto", "err", err)
					return // TODO
				}
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

// LoggerGRPC overwrites to log normal requests as debug logs, defaults to grpc_gokit otherwise
func LoggerGRPC() func(code codes.Code, logger log.Logger) log.Logger {
	return func(code codes.Code, logger log.Logger) log.Logger {
		switch code {
		case codes.OK, codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
			return level.Debug(logger)
		default:
			return grpc_gokit.DefaultCodeToLevel(code, logger)
		}
	}
}
