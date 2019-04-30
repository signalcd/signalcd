package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/metalmatze/cd/cd"
	"github.com/oklog/run"
	appsv1 "k8s.io/api/apps/v1"
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

	ctx, cancel := context.WithCancel(context.Background())

	var gr run.Group
	{
		u := &updater{client: client, logger: logger, agentName: *agentName}

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
}

type updater struct {
	client    *kubernetes.Clientset
	logger    log.Logger
	agentName string

	currentPipeline cd.Pipeline
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
			if err := u.poll(); err != nil {
				level.Warn(u.logger).Log(
					"msg", "failed to poll",
					"err", err,
				)
			}
		}
	}
}

func (u *updater) poll() error {
	p, err := u.pipeline()
	if err != nil {
		return err
	}

	if u.currentPipeline.ID == "" {
		level.Info(u.logger).Log("msg", "unknown pipeline", "pipeline", p.ID)
		u.currentPipeline = p
	}
	if u.currentPipeline.ID != p.ID {
		level.Info(u.logger).Log("msg", "updated pipeline", "pipeline", p.ID)
		u.currentPipeline = p
	}

	return nil
}

func (u *updater) pipeline() (cd.Pipeline, error) {
	var w cd.Pipeline

	resp, err := http.Get(api + "/pipeline")
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

func (u *updater) pipelineAgents(status appsv1.DeploymentStatus) error {
	payload, err := json.Marshal(cd.Agent{
		Name:   u.agentName,
		Status: status,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(api+"/pipelines/agents", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("no 200 OK, but %s", resp.Status)
	}

	return nil
}
