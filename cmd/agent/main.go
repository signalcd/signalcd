package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const apiURL = "localhost:6661"

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Action = agentAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "The name for this specific agent instance",
		},
		cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "Path to the kubeconfig",
		},
		cli.StringFlag{
			Name:  "api.url",
			Usage: "Full URL to API, like http://localhost:6661",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "The namespace to deploy to",
		},
		cli.StringFlag{
			Name:  "serviceaccount",
			Usage: "The name of the ServiceAccount to use",
			Value: "signalcd",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running agent", "err", err)
		os.Exit(1)
	}
}

func agentAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		namespace := c.String("namespace")
		if namespace == "" {
			return errors.New("no namespace given, use --namespace flag")
		}

		conn, err := grpc.Dial(c.String("api.url"), grpc.WithInsecure())
		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to dial gRPC target",
				"err", err,
			)
			os.Exit(2)
		}
		defer conn.Close()

		client := signalcdproto.NewAgentServiceClient(conn)

		konfig, err := clientcmd.BuildConfigFromFlags("", c.String("kubeconfig"))
		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to read Kubernetes config",
				"err", err,
			)
			os.Exit(2)
		}

		klient, err := kubernetes.NewForConfig(konfig)
		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to create Kubernetes client",
				"err", err,
			)
			os.Exit(3)
		}

		ctx, cancel := context.WithCancel(context.Background())

		var gr run.Group
		{
			u := &updater{
				client: client,
				klient: klient,
				logger: logger,

				agentName:      c.String("name"),
				namespace:      namespace,
				serviceAccount: c.String("serviceaccount"),
			}

			gr.Add(func() error {
				return u.pollLoop(ctx)
			}, func(err error) {
				cancel()
			})
		}
		{
			// Start gRPC Server here
		}
		{
			sig := make(chan os.Signal)

			gr.Add(func() error {
				signal.Notify(sig, os.Interrupt)
				<-sig
				cancel()
				return nil
			}, func(err error) {
				close(sig)
			})
		}

		if err := gr.Run(); err != nil {
			level.Error(logger).Log(
				"msg", "error running",
				"err", err,
			)
			os.Exit(1)
		}

		return nil
	}
}

type currentDeployment struct {
	mu         sync.RWMutex
	deployment *signalcd.Deployment
}

func (cd *currentDeployment) get() *signalcd.Deployment {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.deployment
}

func (cd *currentDeployment) set(d signalcd.Deployment) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.deployment = &d
}

type updater struct {
	client signalcdproto.AgentServiceClient
	klient *kubernetes.Clientset
	logger log.Logger

	agentName      string
	namespace      string
	serviceAccount string

	currentDeployment currentDeployment
}

func (u *updater) pollLoop(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	level.Info(u.logger).Log("msg", "starting poll loop")

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := u.poll(ctx)
			if err != nil {
				level.Warn(u.logger).Log(
					"msg", "failed to poll",
					"err", err,
				)
			}
		}
	}
}

func deploymentStatusPhase(phase signalcdproto.DeploymentStatus_Phase) signalcd.DeploymentPhase {
	switch phase {
	case signalcdproto.DeploymentStatus_SUCCESS:
		return signalcd.Success
	case signalcdproto.DeploymentStatus_FAILURE:
		return signalcd.Failure
	case signalcdproto.DeploymentStatus_PROGRESS:
		return signalcd.Progress
	case signalcdproto.DeploymentStatus_PENDING:
		return signalcd.Pending
	case signalcdproto.DeploymentStatus_KILLED:
		return signalcd.Killed
	default:
		return signalcd.Unknown
	}
}

func deploymentFromRPC(deployment *signalcdproto.Deployment) signalcd.Deployment {
	return signalcd.Deployment{
		Number:   deployment.GetNumber(),
		Created:  time.Unix(deployment.GetCreated(), 0),
		Pipeline: pipelineFromRPC(deployment.GetPipeline()),
		Status: signalcd.DeploymentStatus{
			Phase: deploymentStatusPhase(deployment.GetStatus().GetPhase()),
		},
	}
}

func pipelineFromRPC(pipeline *signalcdproto.Pipeline) signalcd.Pipeline {
	p := signalcd.Pipeline{
		ID:   pipeline.GetId(),
		Name: pipeline.GetName(),
	}

	for _, step := range pipeline.GetSteps() {
		p.Steps = append(p.Steps, signalcd.Step{
			Name:             step.GetName(),
			Image:            step.GetImage(),
			ImagePullSecrets: step.GetImagePullSecrets(),
			Commands:         step.GetCommands(),
		})
	}

	for _, check := range pipeline.GetChecks() {
		//env := map[string]string{}
		//for _, item := range check.Environment {
		//	env[item.Key] = item.Value
		//}

		p.Checks = append(p.Checks, signalcd.Check{
			Name:             check.GetName(),
			Image:            check.GetImage(),
			ImagePullSecrets: check.GetImagePullSecrets(),
			Duration:         time.Duration(check.Duration) * time.Second,
			Environment:      map[string]string{},
		})
	}

	return p
}

func (u *updater) poll(ctx context.Context) error {
	resp, err := u.client.CurrentDeployment(ctx, &signalcdproto.CurrentDeploymentRequest{})
	if err != nil {
		return xerrors.Errorf("failed to get current deployment: %w", err)
	}

	deployment := deploymentFromRPC(resp.GetCurrentDeployment())

	if u.currentDeployment.get() == nil {
		loaded, err := u.loadDeployment()
		if err != nil && !apierrors.IsNotFound(err) {
			return u.sendStatus(ctx, deployment.Number, xerrors.Errorf("failed to load pipeline: %v", err))
		}
		level.Debug(u.logger).Log("msg", "loaded pipeline from ConfigMap")

		u.currentDeployment.set(deployment)

		// if running deployment id in cluster equals to wanted deployment
		// we don't need to run the pipeline
		if loaded.Number == deployment.Number {
			level.Debug(u.logger).Log("msg", "ConfigMap has same deployment ID", "number", deployment.Number)
			return nil
		}

		level.Info(u.logger).Log("msg", "unknown pipeline, deploying...", "deployment", deployment.Number)

		_, err = u.client.SetDeploymentStatus(ctx, &signalcdproto.SetDeploymentStatusRequest{
			Number: deployment.Number,
			Status: &signalcdproto.DeploymentStatus{
				Phase: signalcdproto.DeploymentStatus_PROGRESS,
			},
		})
		if err != nil {
			return xerrors.Errorf("failed to update deployment status: %w", err)
		}

		if err := u.runPipeline(ctx, deployment.Pipeline); err != nil {
			return u.sendStatus(ctx, deployment.Number, xerrors.Errorf("failed to run pipeline: %w", err))
		}

		err = u.saveDeployment(deployment)
		if err != nil {
			return u.sendStatus(ctx, deployment.Number, xerrors.Errorf("failed to save pipeline: %w", err))
		}

		level.Info(u.logger).Log("msg", "finished updating deployment", "number", deployment.Number)

		return u.sendStatus(ctx, deployment.Number, nil)
	}

	if u.currentDeployment.get().Number != deployment.Number {
		u.currentDeployment.set(deployment)
		level.Info(u.logger).Log("msg", "update deployment", "number", deployment.Number)

		_, err := u.client.SetDeploymentStatus(ctx, &signalcdproto.SetDeploymentStatusRequest{
			Number: deployment.Number,
			Status: &signalcdproto.DeploymentStatus{
				Phase: signalcdproto.DeploymentStatus_PROGRESS,
			},
		})
		if err != nil {
			return xerrors.Errorf("failed to update deployment status: %w", err)
		}

		if err := u.runPipeline(ctx, deployment.Pipeline); err != nil {
			return u.sendStatus(ctx, deployment.Number, xerrors.Errorf("failed to run deployment: %w", err))
		}

		if err := u.saveDeployment(deployment); err != nil {
			return u.sendStatus(ctx, deployment.Number, xerrors.Errorf("failed to save deployment: %w", err))
		}

		level.Info(u.logger).Log("msg", "finished updating deployment", "number", deployment.Number)

		return u.sendStatus(ctx, deployment.Number, nil)
	}

	return nil
}

func (u *updater) sendStatus(ctx context.Context, number int64, err error) error {
	if err != nil {
		_, err2 := u.client.SetDeploymentStatus(ctx, &signalcdproto.SetDeploymentStatusRequest{
			Number: number,
			Status: &signalcdproto.DeploymentStatus{
				Phase: signalcdproto.DeploymentStatus_FAILURE,
			},
		})
		if err2 != nil {
			return xerrors.Errorf("failed to update deployment status: %w - original error: %w", err2, err)
		}
	} else {
		_, err2 := u.client.SetDeploymentStatus(ctx, &signalcdproto.SetDeploymentStatusRequest{
			Number: number,
			Status: &signalcdproto.DeploymentStatus{
				Phase: signalcdproto.DeploymentStatus_SUCCESS,
			},
		})
		if err2 != nil {
			return xerrors.Errorf("failed to update deployment status: %w", err2)
		}
	}

	return err
}

func (u *updater) runPipeline(ctx context.Context, p signalcd.Pipeline) error {
	level.Info(u.logger).Log("msg", "running steps")
	if err := u.runSteps(ctx, p); err != nil {
		return fmt.Errorf("failed to run steps: %w", err)
	}

	level.Info(u.logger).Log("msg", "cleaning checks")
	if err := u.cleanChecks(p); err != nil {
		return fmt.Errorf("failed to clean old checks: %w", err)
	}

	level.Info(u.logger).Log("msg", "running checks")
	if err := u.runChecks(p); err != nil {
		return fmt.Errorf("failed to run checks: %w", err)
	}

	return nil
}

func (u *updater) runSteps(ctx context.Context, p signalcd.Pipeline) error {
	for _, s := range p.Steps {
		level.Debug(u.logger).Log(
			"msg", "running step",
			"pipeline", p.Name,
			"step", s.Name,
		)
		if err := u.runStep(ctx, p, s); err != nil {
			return fmt.Errorf("failed to run pipeline %s step %s: %w", p.Name, s.Name, err)
		}
	}

	return nil
}

func (u *updater) runStep(ctx context.Context, pipeline signalcd.Pipeline, step signalcd.Step) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	args := []string{"-c"}
	for _, c := range step.Commands {
		args = append(args, c)
	}

	podName := strings.ToLower(pipeline.Name + "-" + step.Name)

	var imagePullSecrets []corev1.LocalObjectReference
	for _, secret := range step.ImagePullSecrets {
		imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{Name: secret})
	}

	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: u.namespace,
			Labels: map[string]string{
				"signalcd": "step",
				"pipeline": pipeline.Name,
				"step":     step.Name,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: u.serviceAccount,
			Containers: []corev1.Container{{
				Name:            step.Name,
				Image:           step.Image,
				ImagePullPolicy: corev1.PullAlways,
				Command:         []string{"sh"},
				Args:            args,
			}},
			ImagePullSecrets: imagePullSecrets,
			RestartPolicy:    corev1.RestartPolicyNever,
		},
	}

	podLogger := log.With(u.logger, "namespace", u.namespace, "pod", p.Name)

	_, err := u.klient.CoreV1().Pods(u.namespace).Create(&p)
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}
	defer func(p *corev1.Pod) {
		_ = u.klient.CoreV1().Pods(u.namespace).Delete(p.Name, nil)
		level.Debug(podLogger).Log("msg", "deleted pod")
	}(&p)

	level.Debug(podLogger).Log("msg", "created pod")

	timeout := int64(time.Minute.Seconds())
	watch, err := u.klient.CoreV1().Pods(u.namespace).Watch(metav1.ListOptions{
		LabelSelector:  labelsSelector(p.GetLabels()),
		Watch:          true,
		TimeoutSeconds: &timeout,
	})
	if err != nil {
		return fmt.Errorf("failed to watch pods %s: %w", labelsSelector(p.GetLabels()), err)
	}

	for event := range watch.ResultChan() {
		pod := event.Object.(*corev1.Pod)

		if pod.Status.Phase == corev1.PodSucceeded {
			return nil
		}
		if pod.Status.Phase == corev1.PodFailed {
			return fmt.Errorf("step failed")
		}
	}

	return nil
}

func (u *updater) cleanChecks(pipeline signalcd.Pipeline) error {
	err := u.klient.CoreV1().Pods(u.namespace).DeleteCollection(nil, metav1.ListOptions{
		LabelSelector: labelsSelector(checkLabels),
	})
	if err != nil {
		return xerrors.Errorf("failed to delete pods: %w", err)
	}

	return nil
}

func (u *updater) runChecks(p signalcd.Pipeline) error {
	for _, c := range p.Checks {
		if err := u.runCheck(p, c); err != nil {
			return fmt.Errorf("failed to run pipeline %s check %s: %w", p.Name, c.Name, err)
		}
	}

	return nil
}

var checkLabels = map[string]string{
	"cd": "check",
}

func (u *updater) runCheck(pipeline signalcd.Pipeline, check signalcd.Check) error {
	// Add PLUGIN_API for checks to find the API
	check.Environment["PLUGIN_API"] = apiURL

	var env []corev1.EnvVar
	for name, value := range check.Environment {
		env = append(env, corev1.EnvVar{Name: name, Value: value})
	}

	var imagePullSecrets []corev1.LocalObjectReference
	for _, secret := range check.ImagePullSecrets {
		imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{Name: secret})
	}

	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(pipeline.Name + "-" + check.Name),
			Namespace: u.namespace,
			Labels:    checkLabels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: u.serviceAccount,
			Containers: []corev1.Container{{
				Name:            strings.ToLower(check.Name),
				Image:           check.Image,
				ImagePullPolicy: corev1.PullAlways,
				Env:             env,
			}},
			ImagePullSecrets: imagePullSecrets,
			RestartPolicy:    corev1.RestartPolicyNever,
		},
	}

	_, err := u.klient.CoreV1().Pods(u.namespace).Create(&p)
	if err != nil {
		return xerrors.Errorf("failed to create check pod: %w", err)
	}

	return nil
}

func labelsSelector(ls map[string]string) string {
	var selectors []string
	for key, value := range ls {
		selectors = append(selectors, key+"="+value)
	}
	return strings.Join(selectors, ",")
}

const configMapName = "signalcd"
const configMapFilename = "deployment.json"

func (u *updater) loadDeployment() (signalcd.Deployment, error) {
	cm, err := u.klient.CoreV1().ConfigMaps(u.namespace).Get(configMapName, metav1.GetOptions{})
	if err != nil {
		return signalcd.Deployment{}, err
	}

	b := cm.Data[configMapFilename]

	var d signalcd.Deployment
	if err := json.Unmarshal([]byte(b), &d); err != nil {
		return d, err
	}

	return d, nil
}

func (u *updater) saveDeployment(d signalcd.Deployment) error {
	b, err := json.Marshal(&d)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment configmap: %w", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: u.namespace,
		},
		Data: map[string]string{
			configMapFilename: string(b),
		},
	}

	_, err = u.klient.CoreV1().ConfigMaps(u.namespace).Get(configMapName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = u.klient.CoreV1().ConfigMaps(u.namespace).Create(cm)
		if err != nil {
			return xerrors.Errorf("failed to create ConfigMap: %v", err)
		}
		return nil
	}
	if err != nil {
		return xerrors.Errorf("failed to get ConfigMap: %v", err)
	}

	_, err = u.klient.CoreV1().ConfigMaps(u.namespace).Update(cm)
	if err != nil {
		return xerrors.Errorf("failed to update ConfigMap: %v", err)
	}
	return nil
}
