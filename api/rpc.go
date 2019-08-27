package api

import (
	"github.com/go-kit/kit/log"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
)

type RPC struct {
	DB     SignalDB
	logger log.Logger
}

func NewRPC(db SignalDB, logger log.Logger) *RPC {
	return &RPC{
		DB:     db,
		logger: logger,
	}
}

func (rpc *RPC) CurrentDeployment(ctx context.Context, req *signalcdproto.CurrentDeploymentRequest) (*signalcdproto.CurrentDeploymentResponse, error) {
	deployment, err := rpc.DB.GetCurrentDeployment()
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
