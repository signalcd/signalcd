package signalcdproto

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/signalcd/signalcd/signalcd"
)

func DeploymentSignalCD(deployment *Deployment) (signalcd.Deployment, error) {
	pipeline, err := PipelineSignalCD(deployment.GetPipeline())
	if err != nil {
		return signalcd.Deployment{}, err
	}

	created, err := ptypes.Timestamp(deployment.GetCreated())
	if err != nil {
		return signalcd.Deployment{}, err
	}
	started, err := ptypes.Timestamp(deployment.GetStarted())
	if err != nil {
		return signalcd.Deployment{}, err
	}
	finished, err := ptypes.Timestamp(deployment.GetFinished())
	if err != nil {
		return signalcd.Deployment{}, err
	}

	var phase signalcd.DeploymentPhase
	switch deployment.GetStatus().GetPhase() {
	case DeploymentStatus_UNKNOWN:
		phase = signalcd.Unknown
	case DeploymentStatus_SUCCESS:
		phase = signalcd.Success
	case DeploymentStatus_FAILURE:
		phase = signalcd.Failure
	case DeploymentStatus_PROGRESS:
		phase = signalcd.Progress
	case DeploymentStatus_PENDING:
		phase = signalcd.Pending
	case DeploymentStatus_KILLED:
		phase = signalcd.Killed
	}

	return signalcd.Deployment{
		Number:   deployment.GetNumber(),
		Created:  created,
		Started:  started,
		Finished: finished,
		Pipeline: pipeline,
		Status: signalcd.DeploymentStatus{
			Phase: phase,
		},
	}, nil
}

func DeploymentProto(d signalcd.Deployment) (*Deployment, error) {
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
	p, err := PipelineProto(d.Pipeline)
	if err != nil {
		return nil, err
	}

	var phase DeploymentStatus_Phase
	switch d.Status.Phase {
	case signalcd.Unknown:
		phase = DeploymentStatus_UNKNOWN
	case signalcd.Success:
		phase = DeploymentStatus_SUCCESS
	case signalcd.Failure:
		phase = DeploymentStatus_FAILURE
	case signalcd.Progress:
		phase = DeploymentStatus_PROGRESS
	case signalcd.Pending:
		phase = DeploymentStatus_PENDING
	case signalcd.Killed:
		phase = DeploymentStatus_KILLED
	}

	return &Deployment{
		Number:   d.Number,
		Created:  created,
		Started:  started,
		Finished: finished,
		Status: &DeploymentStatus{
			Phase: phase,
		},
		Pipeline: p,
	}, nil
}

func PipelineSignalCD(p *Pipeline) (signalcd.Pipeline, error) {
	steps := make([]signalcd.Step, len(p.GetSteps()))
	for i, s := range p.GetSteps() {
		steps[i] = signalcd.Step{
			Name:             s.GetName(),
			Image:            s.GetImage(),
			ImagePullSecrets: s.GetImagePullSecrets(),
			Commands:         s.GetCommands(),
		}
	}

	checks := make([]signalcd.Check, len(p.GetChecks()))
	for i, c := range p.GetChecks() {
		duration, err := ptypes.Duration(c.GetDuration())
		if err != nil {
			return signalcd.Pipeline{}, err
		}

		checks[i] = signalcd.Check{
			Name:             c.GetName(),
			Image:            c.GetImage(),
			ImagePullSecrets: c.GetImagePullSecrets(),
			Duration:         duration,
		}
	}

	pipeline := signalcd.Pipeline{
		ID:     p.GetId(),
		Name:   p.GetName(),
		Steps:  steps,
		Checks: checks,
	}

	if p.GetCreated() != nil {
		created, err := ptypes.Timestamp(p.GetCreated())
		if err != nil {
			return signalcd.Pipeline{}, err
		}
		pipeline.Created = created
	}

	return pipeline, nil
}

func PipelineProto(p signalcd.Pipeline) (*Pipeline, error) {
	created, err := ptypes.TimestampProto(p.Created)
	if err != nil {
		return nil, err
	}

	steps := make([]*Step, len(p.Steps))
	for i, s := range p.Steps {
		steps[i] = &Step{
			Name:             s.Name,
			Image:            s.Image,
			ImagePullSecrets: s.ImagePullSecrets,
			Commands:         s.Commands,
		}
	}

	checks := make([]*Check, len(p.Checks))
	for i, c := range p.Checks {
		checks[i] = &Check{
			Name:             c.Name,
			Image:            c.Image,
			ImagePullSecrets: c.ImagePullSecrets,
			Duration:         ptypes.DurationProto(c.Duration),
		}
	}

	return &Pipeline{
		Id:      p.ID,
		Name:    p.Name,
		Created: created,
		Steps:   steps,
		Checks:  checks,
	}, nil
}
