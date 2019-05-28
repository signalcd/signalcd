package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/signalcd/signalcd/cmd/agent/client"
	"github.com/signalcd/signalcd/cmd/agent/models"
	"github.com/signalcd/signalcd/signalcd"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const apiURL = "localhost:6660"
const namespace = "default"

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
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running agent", "err", err)
		os.Exit(1)
	}
}

func agentAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		transportCfg := client.DefaultTransportConfig().
			WithSchemes([]string{"http"}).
			WithHost(apiURL)

		client := client.NewHTTPClientWithConfig(nil, transportCfg)

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
				client:    client,
				klient:    klient,
				logger:    logger,
				agentName: c.String("name"),
			}

			gr.Add(func() error {
				return u.pollLoop(ctx)
			}, func(err error) {
				cancel()
			})
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

func (cd currentDeployment) get() *signalcd.Deployment {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.deployment
}

func (cd currentDeployment) set(deployment *signalcd.Deployment) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.deployment = deployment
}

type updater struct {
	client    *client.SignalCDSwaggerSpec
	klient    *kubernetes.Clientset
	logger    log.Logger
	agentName string

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
				// TODO
				//if err := u.pipelineStatus(signalcd.Failed); err != nil {
				//	level.Warn(u.logger).Log(
				//		"msg", "failed to update deployment status",
				//		"err", err,
				//	)
				//}
			} else {
				// TODO
				//if err := u.pipelineStatus(signalcd.Success); err != nil {
				//	level.Warn(u.logger).Log(
				//		"msg", "failed to update deployment status",
				//		"err", err,
				//	)
				//}
			}
		}
	}
}

func deploymentStatusPhase(phase string) signalcd.DeploymentPhase {
	switch phase {
	case models.DeploymentstatusPhaseSuccess:
		return signalcd.Success
	case models.DeploymentstatusPhaseFailure:
		return signalcd.Failed
	case models.DeploymentstatusPhaseProgress:
		return signalcd.Progress
	default:
		return signalcd.Unknown
	}
}

func deploymentFromAPI(deployment *models.Deployment) signalcd.Deployment {
	return signalcd.Deployment{
		Number:  *deployment.Number,
		Created: time.Time(deployment.Created),
		Status: signalcd.DeploymentStatus{
			Phase: deploymentStatusPhase(deployment.Status.Phase),
		},
		Pipeline: pipelineFromAPI(deployment.Pipeline),
	}
}

func pipelineFromAPI(pipeline *models.Pipeline) signalcd.Pipeline {
	p := signalcd.Pipeline{
		ID:   pipeline.ID.String(),
		Name: pipeline.Name,
	}

	for _, step := range pipeline.Steps {
		p.Steps = append(p.Steps, signalcd.Step{
			Name:     *step.Name,
			Image:    *step.Image,
			Commands: step.Commands,
		})
	}

	for _, check := range pipeline.Checks {
		env := map[string]string{}
		for _, item := range check.Environment {
			env[item.Key] = item.Value
		}

		p.Checks = append(p.Checks, signalcd.Check{
			Name:        *check.Name,
			Image:       *check.Image,
			Duration:    time.Duration(check.Duration) * time.Second,
			Environment: env,
		})
	}

	return p
}

func (u *updater) poll(ctx context.Context) error {
	deploymentOK, err := u.client.Deployments.CurrentDeployment(nil)
	if err != nil {
		return xerrors.Errorf("failed to get current deployment: %w", err)
	}
	deployment := deploymentFromAPI(deploymentOK.Payload)

	if u.currentDeployment.get() == nil {
		loaded, err := u.loadDeployment()
		if err != nil && !apierrors.IsNotFound(err) {
			return xerrors.Errorf("failed to load pipeline: %v", err)
		}
		level.Debug(u.logger).Log("msg", "loaded pipeline from ConfigMap")

		u.currentDeployment.set(&deployment)

		// if running deployment id in cluster equals to wanted deployment
		// we don't need to run the pipeline
		if loaded.Number == deployment.Number {
			level.Debug(u.logger).Log("msg", "ConfigMap has same deployment ID", "number", deployment.Number)
			return nil
		}

		level.Info(u.logger).Log("msg", "unknown pipeline, deploying...", "deployment", deployment.Number)

		// TODO
		//if err := u.pipelineStatus(signalcd.Progress); err != nil {
		//	level.Warn(u.logger).Log(
		//		"msg", "failed to update pipeline status",
		//		"err", err,
		//	)
		//}

		if err := u.runPipeline(ctx, deployment.Pipeline); err != nil {
			return xerrors.Errorf("failed to run pipeline: %w", err)
		}

		err = u.saveDeployment(deployment)
		if err != nil {
			return xerrors.Errorf("failed to save pipeline: %w", err)
		}

		level.Info(u.logger).Log("msg", "finished updating deployment", "number", deployment.Number)

		return nil
	}

	if u.currentDeployment.get().Number != deployment.Number {
		u.currentDeployment.set(&deployment)
		level.Info(u.logger).Log("msg", "update deployment", "number", deployment.Number)

		// TODO
		//if err := u.pipelineStatus(signalcd.Progress); err != nil {
		//	level.Warn(u.logger).Log(
		//		"msg", "failed to update pipeline status",
		//		"err", err,
		//	)
		//}

		if err := u.runPipeline(ctx, deployment.Pipeline); err != nil {
			return xerrors.Errorf("failed to run deployment: %w", err)
		}

		err := u.saveDeployment(deployment)
		if err != nil {
			return xerrors.Errorf("failed to save deployment: %w", err)
		}

		level.Info(u.logger).Log("msg", "finished updating deployment", "number", deployment.Number)

		return nil
	}

	return nil
}

func (u *updater) runPipeline(ctx context.Context, p signalcd.Pipeline) error {
	println("running steps")
	if err := u.runSteps(ctx, p); err != nil {
		return err
	}

	println("cleaning checks")
	if err := u.cleanChecks(p); err != nil {
		return xerrors.Errorf("failed to clean old checks: %w", err)
	}

	println("running checks")
	if err := u.runChecks(p); err != nil {
		return err
	}

	return nil
}

//func (u *updater) pipelineStatus(status signalcd.DeploymentPhase) error {
//	payload, err := json.Marshal(signalcd.Agent{
//		Name:     u.agentName,
//		Pipeline: u.currentDeployment,
//	})
//	if err != nil {
//		return err
//	}
//
//	resp, err := http.Post(apiURL+api.PipelinesStatus, "application/json", bytes.NewBuffer(payload))
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("no 200 OK, but %s", resp.Status)
//	}
//
//	return nil
//}

func (u *updater) runSteps(ctx context.Context, p signalcd.Pipeline) error {
	for _, s := range p.Steps {
		if err := u.runStep(ctx, p, s); err != nil {
			return err
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

	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "cd",
			Containers: []corev1.Container{{
				Name:            step.Name,
				Image:           step.Image,
				ImagePullPolicy: corev1.PullAlways,
				Command:         []string{"sh"},
				Args:            args,
			}},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	_, err := u.klient.CoreV1().Pods(namespace).Create(&p)
	if err != nil {
		return err
	}
	defer func(p *corev1.Pod) {
		_ = u.klient.CoreV1().Pods(namespace).Delete(p.Name, nil)
	}(&p)

	timeout := int64(time.Minute.Seconds())
	watch, err := u.klient.CoreV1().Pods(namespace).Watch(metav1.ListOptions{
		LabelSelector:  labelsSelector(p.GetLabels()),
		Watch:          true,
		TimeoutSeconds: &timeout,
	})
	if err != nil {
		return err
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
	err := u.klient.CoreV1().Pods(namespace).DeleteCollection(nil, metav1.ListOptions{
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
			return err
		}
	}

	return nil
}

var checkLabels = map[string]string{
	"cd": "check",
}

func (u *updater) runCheck(pipeline signalcd.Pipeline, check signalcd.Check) error {
	// Add PLUGIN_API for plugins to find the API
	check.Environment["PLUGIN_API"] = apiURL

	var env []corev1.EnvVar
	for name, value := range check.Environment {
		env = append(env, corev1.EnvVar{Name: name, Value: value})
	}

	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(pipeline.Name + "-" + check.Name),
			Namespace: namespace,
			Labels:    checkLabels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "cd",
			Containers: []corev1.Container{{
				Name:            strings.ToLower(check.Name),
				Image:           check.Image,
				ImagePullPolicy: corev1.PullAlways,
				Env:             env,
			}},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	_, err := u.klient.CoreV1().Pods(namespace).Create(&p)
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
	cm, err := u.klient.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
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
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			configMapFilename: string(b),
		},
	}

	_, err = u.klient.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = u.klient.CoreV1().ConfigMaps(namespace).Create(cm)
		if err != nil {
			return xerrors.Errorf("failed to create ConfigMap: %v", err)
		}
		return nil
	}
	if err != nil {
		return xerrors.Errorf("failed to get ConfigMap: %v", err)
	}

	_, err = u.klient.CoreV1().ConfigMaps(namespace).Update(cm)
	if err != nil {
		return xerrors.Errorf("failed to update ConfigMap: %v", err)
	}
	return nil
}
