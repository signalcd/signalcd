package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/metalmatze/cd/cd"
	"github.com/oklog/run"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const api = "http://localhost:6660"

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig")
	agentName := flag.String("name", "", "Name for this agent")
	flag.Parse()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to read Kubernetes config",
			"err", err,
		)
		os.Exit(2)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to create Kubernetes client",
			"err", err,
		)
		os.Exit(3)
	}

	var gr run.Group
	{
		c := &controller{client: client, logger: logger, agentName: *agentName}

		ctx, cancel := context.WithCancel(context.Background())

		gr.Add(func() error {
			return c.reconcileLoop(ctx)
		}, func(err error) {
			cancel()
		})
	}

	if err := gr.Run(); err != nil {
		level.Error(logger).Log(
			"msg", "error running",
			"err", err,
		)
		os.Exit(1)
	}
}

type controller struct {
	client *kubernetes.Clientset
	logger log.Logger

	agentName string
}

func (c *controller) reconcileLoop(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	level.Info(c.logger).Log("msg", "starting reconcile loop")

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			if err := c.reconcile(); err != nil {
				level.Warn(c.logger).Log(
					"msg", "failed to reconcile",
					"err", err,
				)
			}

		}
	}
}

func (c *controller) reconcile() error {
	deployment, err := c.client.AppsV1().Deployments("logging").Get("loki", metav1.GetOptions{})
	if err != nil {
		return err
	}

	w, err := c.workload()
	if err != nil {
		return err
	}

	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Image != w.Image {
			deployment.Spec.Template.Spec.Containers[i].Image = w.Image

			_, err := c.client.AppsV1().Deployments("logging").Update(deployment)
			if err != nil {
				return err
			}
		}
	}

	return c.workloadAgent(deployment.Status)
}

func (c *controller) workload() (cd.Workload, error) {
	var w cd.Workload

	resp, err := http.Get(api + "/workload")
	if err != nil {
		return w, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return w, err
	}

	return w, err
}

func (c *controller) workloadAgent(status appsv1.DeploymentStatus) error {
	payload, err := json.Marshal(cd.Agent{
		Name:   c.agentName,
		Status: status,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(api+"/workloads/agents", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("no 200 OK, but %s", resp.Status)
	}

	return nil
}
