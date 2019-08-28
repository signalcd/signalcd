package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	config = struct {
		API        string
		KubeConfig string
		Duration   time.Duration
		Labels     string
	}{}

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Path to the kubeconfig which should be sued",
			EnvVar:      "PLUGIN_KUBECONFIG",
			Value:       "",
			Destination: &config.KubeConfig,
		},
		cli.StringFlag{
			Name:        "api",
			EnvVar:      "PLUGIN_API",
			Destination: &config.API,
		},
		cli.DurationFlag{
			Name:        "duration",
			EnvVar:      "PLUGIN_DURATION",
			Value:       time.Minute,
			Destination: &config.Duration,
		},
		cli.StringFlag{
			Name:        "labels",
			EnvVar:      "PLUGIN_LABELS",
			Destination: &config.Labels,
		},
	}
)

func main() {
	app := cli.NewApp()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	app.Name = "kubernetes status"
	app.Action = run
	app.Flags = flags

	if err := app.Run(os.Args); err != nil {
		level.Error(logger).Log(
			"msg", "failed to run",
			"err", err,
		)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	var klient kubernetes.Interface
	{
		konfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
		if err != nil {
			return xerrors.Errorf("failed to create kubernetes config: %w", err)
		}

		klient, err = kubernetes.NewForConfig(konfig)
		if err != nil {
			return xerrors.Errorf("failed to create kubernetes client: %w", err)
		}
	}
	{
		conn, err := grpc.Dial(config.API)
		if err != nil {
			return xerrors.Errorf("failed to create gRPC connection to API: %w", err)
		}
		defer conn.Close()
	}

	//message, err := proto.Marshal(&signalcdproto.CheckMessage{
	//	Message: "OK",
	//})
	//if err != nil {
	//	// TODO
	//}

	//message := proto.CheckMessage{Message: "OK"}
	//proto.

	watch, err := klient.AppsV1().Deployments("default").Watch(metav1.ListOptions{
		LabelSelector: config.Labels,
		Watch:         true,
	})
	if err != nil {
		return err
	}

	ticker := time.NewTicker(5 * time.Second)

	var status v1.DeploymentStatus

	for {
		select {
		case <-ticker.C:
			printStatus(status)
		case event := <-watch.ResultChan():
			if event.Object == nil {
				continue
			}
			d := event.Object.(*v1.Deployment)
			status = d.Status
			printStatus(status)
		case <-ctx.Done():
			ticker.Stop()
			watch.Stop()
			return nil
		}
	}
}

func printStatus(status v1.DeploymentStatus) {
	fmt.Printf("%d of %d replicas are ready\n", status.ReadyReplicas, status.Replicas)
}
