package api

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"

	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
)

// RPC implement the gRPC server connecting it to a SignalDB
type RPC struct {
	DB     SignalDB
	logger log.Logger
}

// NewRPC creates a new gRPC Server
func NewRPC(db SignalDB, logger log.Logger) *RPC {
	return &RPC{
		DB:     db,
		logger: logger,
	}
}

// CurrentDeployment returns the current Deployment
func (r *RPC) CurrentDeployment(ctx context.Context, req *signalcdproto.CurrentDeploymentRequest) (*signalcdproto.CurrentDeploymentResponse, error) {
	deployment, err := r.DB.GetCurrentDeployment()
	if err != nil {
		return nil, fmt.Errorf("failed to get current deployment: %w", err)
	}

	dProto, err := signalcdproto.DeploymentProto(deployment)
	if err != nil {
		return nil, fmt.Errorf("failed to convert deployment to proto: %w", err)
	}

	return &signalcdproto.CurrentDeploymentResponse{CurrentDeployment: dProto}, nil
}

// DeploymentStatusSetter sets the phase for a specific Deployment by its number
type DeploymentStatusSetter interface {
	SetDeploymentStatus(context.Context, int64, signalcd.DeploymentPhase) (signalcd.Deployment, error)
}

// SetDeploymentStatus sets the phase for a specific Deployment when receiving a request
func (r *RPC) SetDeploymentStatus(ctx context.Context, req *signalcdproto.SetDeploymentStatusRequest) (*signalcdproto.SetDeploymentStatusResponse, error) {
	var phase signalcd.DeploymentPhase

	switch req.Status.Phase {
	case signalcdproto.DeploymentStatus_UNKNOWN:
		phase = signalcd.Unknown
	case signalcdproto.DeploymentStatus_SUCCESS:
		phase = signalcd.Success
	case signalcdproto.DeploymentStatus_FAILURE:
		phase = signalcd.Failure
	case signalcdproto.DeploymentStatus_PROGRESS:
		phase = signalcd.Progress
	case signalcdproto.DeploymentStatus_PENDING:
		phase = signalcd.Pending
	case signalcdproto.DeploymentStatus_KILLED:
		phase = signalcd.Killed
	}

	// TODO: Use returned Deployment in SetDeploymentStatusResponse
	_, err := r.DB.SetDeploymentStatus(ctx, req.Number, phase)
	if err != nil {
		return nil, xerrors.Errorf("failed to update status: %w", err)
	}

	return &signalcdproto.SetDeploymentStatusResponse{}, nil
}

type StepStatusSetter interface {
	SetStepStatus(deployment int64, step int64, status signalcd.Status) error
}

func (r *RPC) StepStatus(ctx context.Context, req *signalcdproto.StepStatusRequest) (*signalcdproto.StepStatusResponse, error) {
	status := signalcd.Status{
		ExitCode: req.GetStatus().GetExitCode(),
	}

	if req.GetStatus().GetStarted() != nil {
		started, err := ptypes.Timestamp(req.GetStatus().GetStarted())
		if err != nil {
			return nil, err
		}
		status.Started = started
	}
	if req.GetStatus().GetStopped() != nil {
		stopped, err := ptypes.Timestamp(req.GetStatus().GetStopped())
		if err != nil {
			return nil, err
		}
		status.Stopped = stopped
	}

	err := r.DB.SetStepStatus(req.GetDeployment(), req.GetStep(), status)
	return &signalcdproto.StepStatusResponse{}, err
}

// StepLogsSaver saves the logs for a Deployment step by its number
type StepLogsSaver interface {
	SaveStepLogs(ctx context.Context, deployment int64, step int64, logs []byte) error
}

//StepLogs saves the logs for a specific deployment and step coming from an agent
func (r *RPC) StepLogs(ctx context.Context, req *signalcdproto.StepLogsRequest) (*signalcdproto.StepLogsResponse, error) {
	err := r.DB.SaveStepLogs(ctx, req.GetDeployment(), req.GetStep(), req.GetLogs())
	if err != nil {
		return nil, fmt.Errorf("failed to save logs: %w", err)
	}

	return &signalcdproto.StepLogsResponse{}, nil
}
