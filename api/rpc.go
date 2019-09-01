package api

import (
	"github.com/go-kit/kit/log"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
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
		return nil, xerrors.Errorf("failed to get current deployment: %w", err)
	}

	return &signalcdproto.CurrentDeploymentResponse{
		CurrentDeployment: &signalcdproto.Deployment{
			Number:  deployment.Number,
			Created: deployment.Created.Unix(),
		},
	}, nil
}

// DeploymentStatusSetter sets the phase for a specific Deployment by its number
type DeploymentStatusSetter interface {
	SetDeploymentStatus(context.Context, int64, signalcd.DeploymentPhase) error
}

// DeploymentStatus sets the phase for a specific Deployment when receiving a request
func (r *RPC) DeploymentStatus(ctx context.Context, req *signalcdproto.SetDeploymentStatusRequest) (*signalcdproto.SetDeploymentStatusResponse, error) {
	var phase signalcd.DeploymentPhase

	switch req.Phase {
	case signalcdproto.SetDeploymentStatusRequest_UNKNOWN:
		phase = signalcd.Unknown
	case signalcdproto.SetDeploymentStatusRequest_SUCCESS:
		phase = signalcd.Success
	case signalcdproto.SetDeploymentStatusRequest_FAILURE:
		phase = signalcd.Failure
	case signalcdproto.SetDeploymentStatusRequest_PROGRESS:
		phase = signalcd.Progress
	case signalcdproto.SetDeploymentStatusRequest_PENDING:
		phase = signalcd.Pending
	case signalcdproto.SetDeploymentStatusRequest_KILLED:
		phase = signalcd.Killed
	}

	err := r.DB.SetDeploymentStatus(ctx, req.Number, phase)
	if err != nil {
		return nil, xerrors.Errorf("failed to update status: %w", err)
	}

	return &signalcdproto.SetDeploymentStatusResponse{}, nil
}
