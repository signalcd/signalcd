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

	steps := func(steps1 []signalcd.Step) []*signalcdproto.Step {
		var steps2 []*signalcdproto.Step
		for _, s := range steps1 {
			steps2 = append(steps2, &signalcdproto.Step{
				Name:     s.Name,
				Image:    s.Image,
				Commands: s.Commands,
			})
		}
		return steps2
	}

	checks := func(checks1 []signalcd.Check) []*signalcdproto.Check {
		var checks2 []*signalcdproto.Check
		for _, c := range checks1 {
			checks2 = append(checks2, &signalcdproto.Check{
				Name:     c.Name,
				Image:    c.Image,
				Duration: int64(c.Duration.Seconds()),
			})
		}
		return checks2
	}

	return &signalcdproto.CurrentDeploymentResponse{
		CurrentDeployment: &signalcdproto.Deployment{
			Number:  deployment.Number,
			Created: deployment.Created.Unix(),

			Pipeline: &signalcdproto.Pipeline{
				Id:     deployment.Pipeline.ID,
				Name:   deployment.Pipeline.Name,
				Steps:  steps(deployment.Pipeline.Steps),
				Checks: checks(deployment.Pipeline.Checks),
			},
		},
	}, nil
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
