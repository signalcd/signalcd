package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	apiclient "github.com/signalcd/signalcd/api/client/go"
	"github.com/signalcd/signalcd/signalcd"
	"github.com/urfave/cli"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	apiURL      = "localhost:6663"
	flagTLSCert = "tls.cert"
	flagTLSKey  = "tls.key"
)

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
			Usage: "Full URL to API, like http://localhost:6660",
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
		cli.StringFlag{
			Name:  flagTLSCert,
			Usage: "The path to the certificate to use when making requests",
		},
		cli.StringFlag{
			Name:  flagTLSKey,
			Usage: "The path to the key to use then making requests",
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

		if c.String("api.url") == "" {
			return errors.New("no api.url to API gRPC endpoint given, use --api.url flag")
		}

		apiURL, err := url.Parse(c.String("api.url"))
		if err != nil {
			return fmt.Errorf("failed to parse API URL: %w", err)
		}

		clientCfg := apiclient.NewConfiguration()
		clientCfg.Scheme = apiURL.Scheme
		clientCfg.Host = apiURL.Host
		clientCfg.BasePath = path.Join(apiURL.Path, "/api/v1")

		client := apiclient.NewAPIClient(clientCfg)

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

		events := make(chan signalcd.Deployment, 3)

		var gr run.Group
		{
			l := listener{client: client}

			gr.Add(func() error {
				return l.listen(ctx, events)
			}, func(err error) {
				close(events)
				cancel()
			})
		}
		{
			u := &runner{
				client: client,
				klient: klient,
				logger: logger,

				agentName:      c.String("name"),
				namespace:      namespace,
				serviceAccount: c.String("serviceaccount"),
			}

			gr.Add(func() error {
				return u.run(ctx, events)
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

type listener struct {
	client *apiclient.APIClient
}

func (l *listener) listen(ctx context.Context, deployments chan<- signalcd.Deployment) error {
	var gr run.Group
	{
		t := time.NewTicker(time.Minute)

		gr.Add(func() error {
			deployment, _, err := l.client.DeploymentApi.GetCurrentDeployment(context.Background())
			if err != nil {
				return err
			}
			deployments <- signalDeployment(deployment)

			for {
				select {
				case <-t.C:
					ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
					defer cancel()
					deployment, _, err := l.client.DeploymentApi.GetCurrentDeployment(ctx)
					if err != nil {
						return err
					}
					deployments <- signalDeployment(deployment)
				case <-ctx.Done():
					return nil
				}
			}
		}, func(err error) {
			t.Stop()
		})
	}
	{
		u := url.URL{
			Scheme: l.client.GetConfig().Scheme,
			Host:   l.client.GetConfig().Host,
			Path:   path.Join(l.client.GetConfig().BasePath, "deployments/events"),
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		gr.Add(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				return err
			}

			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			scanner := bufio.NewScanner(resp.Body)
			scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
				if atEOF || len(data) == 0 {
					return 0, nil, nil
				}

				index := bytes.Index(data, []byte("\n\n"))
				if index > 0 {
					return index + 2, data[0:index], nil
				}

				if atEOF {
					return len(data), data, nil
				}

				return 0, nil, nil
			})

			for scanner.Scan() {
				var d apiclient.Deployment

				text := bytes.TrimPrefix(scanner.Bytes(), []byte("data:"))
				text = bytes.TrimSpace(text)

				if err := json.Unmarshal(text, &d); err != nil {
					fmt.Println(scanner.Text())
					return err
				}
				deployments <- signalDeployment(d)
			}
			if scanner.Err() != nil {
				return err
			}

			return nil
		}, func(err error) {
			cancel()
		})
	}

	return gr.Run()
}

func signalDeployment(d apiclient.Deployment) signalcd.Deployment {
	var steps []signalcd.Step
	for _, step := range d.Pipeline.Steps {
		steps = append(steps, signalcd.Step{
			Name:             step.Name,
			Image:            step.Image,
			ImagePullSecrets: step.ImagePullSecrets,
			Commands:         step.Commands,
		})
	}

	var status map[string]*signalcd.Status
	if d.Status != nil {
		status = make(map[string]*signalcd.Status, len(d.Status))
		for agent, s := range d.Status {
			var steps []signalcd.StepStatus
			for _, step := range s.Steps {
				var phase signalcd.Phase
				switch step.Phase {
				case "unknown":
					phase = signalcd.Unknown
				case "success":
					phase = signalcd.Success
				case "failure":
					phase = signalcd.Failure
				case "progress":
					phase = signalcd.Progress
				case "pending":
					phase = signalcd.Pending
				case "killed":
					phase = signalcd.Killed
				}

				var stopped *time.Time
				if !step.Stopped.IsZero() {
					stopped = &step.Stopped
				}

				steps = append(steps, signalcd.StepStatus{
					Phase:    phase,
					ExitCode: 0,
					Started:  step.Started,
					Stopped:  stopped,
				})
			}

			status[agent] = &signalcd.Status{
				Steps: steps,
			}
		}
	}

	return signalcd.Deployment{
		Number:  d.Number,
		Created: d.Created,
		Pipeline: signalcd.Pipeline{
			ID:      d.Pipeline.Id,
			Name:    d.Pipeline.Name,
			Steps:   steps,
			Checks:  nil,
			Created: d.Pipeline.Created,
		},
		Status: status,
	}
}

type runner struct {
	client *apiclient.APIClient
	klient *kubernetes.Clientset
	logger log.Logger

	agentName      string
	namespace      string
	serviceAccount string

	currentDeployment currentDeployment
}

func (r *runner) run(ctx context.Context, events <-chan signalcd.Deployment) error {
	level.Info(r.logger).Log("msg", "runner starting, waiting for events")

	for {
		select {
		case <-ctx.Done():
			return nil
		case deployment := <-events:
			level.Debug(r.logger).Log("msg", "received deployment event", "number", deployment.Number)
			if err := r.poll(ctx, deployment); err != nil {
				level.Warn(r.logger).Log("msg", "failed to run pipeline", "err", err)
			}
		}
	}
}

func (r *runner) poll(ctx context.Context, deployment signalcd.Deployment) error {
	if r.currentDeployment.get() == nil {
		loaded, err := r.loadDeployment()
		if err != nil && !apierrors.IsNotFound(err) {
			return r.sendStatus(ctx, deployment.Number, fmt.Errorf("failed to load pipeline: %v", err))
		}
		level.Debug(r.logger).Log("msg", "loaded pipeline from ConfigMap")

		r.currentDeployment.set(deployment)

		// if running deployment id in cluster equals to wanted deployment
		// we don't need to run the pipeline
		if loaded.Number == deployment.Number {
			level.Debug(r.logger).Log("msg", "ConfigMap has same deployment ID", "number", deployment.Number)
			return nil
		}

		level.Info(r.logger).Log("msg", "unknown pipeline, deploying...", "deployment", deployment.Number)

		if err := r.runPipeline(ctx, deployment.Number, deployment.Pipeline); err != nil {
			return r.sendStatus(ctx, deployment.Number, fmt.Errorf("failed to run pipeline: %w", err))
		}

		err = r.saveDeployment(deployment)
		if err != nil {
			return r.sendStatus(ctx, deployment.Number, fmt.Errorf("failed to save pipeline: %w", err))
		}

		level.Info(r.logger).Log("msg", "finished updating deployment", "number", deployment.Number)

		return r.sendStatus(ctx, deployment.Number, nil)
	}

	if r.currentDeployment.get().Number != deployment.Number {
		r.currentDeployment.set(deployment)
		level.Info(r.logger).Log("msg", "update deployment", "number", deployment.Number)

		if err := r.runPipeline(ctx, deployment.Number, deployment.Pipeline); err != nil {
			return r.sendStatus(ctx, deployment.Number, fmt.Errorf("failed to run deployment: %w", err))
		}

		if err := r.saveDeployment(deployment); err != nil {
			return r.sendStatus(ctx, deployment.Number, fmt.Errorf("failed to save deployment: %w", err))
		}

		level.Info(r.logger).Log("msg", "finished updating deployment", "number", deployment.Number)

		return r.sendStatus(ctx, deployment.Number, nil)
	}

	return nil
}

func (r *runner) sendStatus(ctx context.Context, number int64, err error) error {
	return nil
	// TODO: Introduce in OpenAPI again
	//	if err != nil {
	//		_, err2 := r.client.SetDeploymentStatus(ctx, &signalcdproto.SetDeploymentStatusRequest{
	//			Number: number,
	//			Status: &signalcdproto.DeploymentStatus{
	//				Phase: signalcdproto.DeploymentStatus_FAILURE,
	//			},
	//		})
	//		if err2 != nil {
	//			return fmt.Errorf("failed to update deployment status: %w - original error: %w", err2, err)
	//		}
	//	} else {
	//		_, err2 := r.client.SetDeploymentStatus(ctx, &signalcdproto.SetDeploymentStatusRequest{
	//			Number: number,
	//			Status: &signalcdproto.DeploymentStatus{
	//				Phase: signalcdproto.DeploymentStatus_SUCCESS,
	//			},
	//		})
	//		if err2 != nil {
	//			return fmt.Errorf("failed to update deployment status: %w", err2)
	//		}
	//	}
	//
	//	return err
	//}
}

func (r *runner) runPipeline(ctx context.Context, deploymentNumber int64, p signalcd.Pipeline) error {
	level.Info(r.logger).Log("msg", "running steps")
	if err := r.runSteps(ctx, deploymentNumber, p); err != nil {
		return fmt.Errorf("failed to run steps: %w", err)
	}

	level.Info(r.logger).Log("msg", "cleaning checks")
	if err := r.cleanChecks(p); err != nil {
		return fmt.Errorf("failed to clean old checks: %w", err)
	}

	level.Info(r.logger).Log("msg", "running checks")
	if err := r.runChecks(ctx, deploymentNumber, p); err != nil {
		return fmt.Errorf("failed to run checks: %w", err)
	}

	return nil
}

func (r *runner) runSteps(ctx context.Context, deploymentNumber int64, p signalcd.Pipeline) error {
	for i, s := range p.Steps {
		level.Debug(r.logger).Log(
			"msg", "running step",
			"pipeline", p.Name,
			"step", s.Name,
		)

		if err := r.runStep(ctx, deploymentNumber, int64(i), p, s); err != nil {
			_, _, err := r.client.DeploymentApi.UpdateDeploymentStatus(context.Background(), deploymentNumber, apiclient.DeploymentStatusUpdate{
				Agent: r.agentName,
				Step:  int64(i),
				Phase: string(signalcd.Failure),
			})
			if err != nil {
				// TODO: Log and increase error metrics for failed step
			}

			return fmt.Errorf("failed to run pipeline %s step %s: %w", p.Name, s.Name, err)
		}

		_, _, err := r.client.DeploymentApi.UpdateDeploymentStatus(context.Background(), deploymentNumber, apiclient.DeploymentStatusUpdate{
			Agent: r.agentName,
			Step:  int64(i),
			Phase: string(signalcd.Success),
		})
		if err != nil {
			// TODO: Log and increase error metrics for failed step
		}
	}

	return nil
}

func (r *runner) runStep(ctx context.Context, deploymentNumber int64, stepNumber int64, pipeline signalcd.Pipeline, step signalcd.Step) error {
	_, _, err := r.client.DeploymentApi.UpdateDeploymentStatus(context.Background(), deploymentNumber, apiclient.DeploymentStatusUpdate{
		Agent: r.agentName,
		Step:  stepNumber,
		Phase: string(signalcd.Progress),
	})
	if err != nil {
		// TODO: Log and increase error metrics for failed step
	}

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
			Namespace: r.namespace,
			Labels: map[string]string{
				"signalcd": "step",
				"pipeline": pipeline.Name,
				"step":     step.Name,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: r.serviceAccount,
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

	podLogger := log.With(r.logger, "namespace", r.namespace, "pod", p.Name)

	// Clean up previous runs if the pods still exists
	err = r.klient.CoreV1().Pods(r.namespace).Delete(p.Name, nil)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete previous pod: %w", err)
	}

	_, err = r.klient.CoreV1().Pods(r.namespace).Create(&p)
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}

	defer func(p *corev1.Pod) {
		logs, err := r.podLogs(p.Name)
		if err != nil {
			level.Warn(podLogger).Log("msg", "failed to get pod logs", "err", err)
		}

		level.Debug(podLogger).Log("msg", "step logs", "logs", string(logs))

		// TODO: Introduce in OpenAPI again
		//_, err = r.client.StepLogs(ctx, &signalcdproto.StepLogsRequest{
		//	Number: deploymentNumber,
		//	Step:   stepNumber,
		//	Logs:   logs,
		//})
		//if err != nil {
		//	level.Warn(podLogger).Log("msg", "failed to ship logs", "err", err)
		//}

		_ = r.klient.CoreV1().Pods(r.namespace).Delete(p.Name, nil)
		level.Debug(podLogger).Log("msg", "deleted pod")
	}(&p)

	level.Debug(podLogger).Log("msg", "created pod")

	timeout := int64(time.Minute.Seconds())
	watch, err := r.klient.CoreV1().Pods(r.namespace).Watch(metav1.ListOptions{
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

func (r *runner) podLogs(name string) ([]byte, error) {
	limit := int64(1048576)
	reader, err := r.klient.CoreV1().Pods(r.namespace).GetLogs(name, &corev1.PodLogOptions{LimitBytes: &limit}).Stream()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(reader)
}

func (r *runner) cleanChecks(pipeline signalcd.Pipeline) error {
	err := r.klient.CoreV1().Pods(r.namespace).DeleteCollection(nil, metav1.ListOptions{
		LabelSelector: labelsSelector(checkLabels),
	})
	if err != nil {
		return fmt.Errorf("failed to delete pods: %w", err)
	}

	return nil
}

func (r *runner) runChecks(ctx context.Context, deploymentNumber int64, p signalcd.Pipeline) error {
	for i, c := range p.Checks {
		if err := r.runCheck(ctx, deploymentNumber, int64(i), p, c); err != nil {
			return fmt.Errorf("failed to run pipeline %s check %s: %w", p.Name, c.Name, err)
		}
	}

	return nil
}

var checkLabels = map[string]string{
	"cd": "check",
}

func (r *runner) runCheck(ctx context.Context, deploymentNumber, checkNumber int64, pipeline signalcd.Pipeline, check signalcd.Check) error {
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
			Namespace: r.namespace,
			Labels:    checkLabels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: r.serviceAccount,
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

	_, err := r.klient.CoreV1().Pods(r.namespace).Create(&p)
	if err != nil {
		return fmt.Errorf("failed to create check pod: %w", err)
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

func (r *runner) loadDeployment() (signalcd.Deployment, error) {
	cm, err := r.klient.CoreV1().ConfigMaps(r.namespace).Get(configMapName, metav1.GetOptions{})
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

func (r *runner) saveDeployment(d signalcd.Deployment) error {
	b, err := json.Marshal(&d)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment configmap: %w", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: r.namespace,
		},
		Data: map[string]string{
			configMapFilename: string(b),
		},
	}

	_, err = r.klient.CoreV1().ConfigMaps(r.namespace).Get(configMapName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = r.klient.CoreV1().ConfigMaps(r.namespace).Create(cm)
		if err != nil {
			return fmt.Errorf("failed to create ConfigMap: %v", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get ConfigMap: %v", err)
	}

	_, err = r.klient.CoreV1().ConfigMaps(r.namespace).Update(cm)
	if err != nil {
		return fmt.Errorf("failed to update ConfigMap: %v", err)
	}
	return nil
}
