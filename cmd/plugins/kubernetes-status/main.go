package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := Main(ctx); err != nil {
		level.Error(logger).Log(
			"msg", "failed to run",
			"err", err,
		)
		os.Exit(1)
	}
}

func Main(ctx context.Context) error {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	labels := os.Getenv("PLUGIN_LABELS")

	watch, err := client.AppsV1().Deployments("default").Watch(metav1.ListOptions{
		LabelSelector: labels,
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
