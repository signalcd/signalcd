package main

import (
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net/url"
	"os"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/signalcd/signalcd/api/v1/client"
	"github.com/signalcd/signalcd/api/v1/client/deployments"
	"github.com/signalcd/signalcd/api/v1/client/pipeline"
	"github.com/signalcd/signalcd/api/v1/models"
	"github.com/signalcd/signalcd/signalcd"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

func main() {
	app := cli.NewApp()
	app.Name = "SignalCD Drone plugin"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api.url",
			Usage:  "The URL to talk to the SignalCD API at",
			EnvVar: "PLUGIN_API_URL",
		},
		cli.StringFlag{
			Name:   "signalcd.file,f",
			Usage:  "The path to the SignalCD file to use",
			EnvVar: "PLUGIN_SIGNALCD_FILE",
			Value:  ".signalcd.yaml",
		},
		cli.StringFlag{
			Name:   "basicauth.username",
			Usage:  "The username to authenticate with",
			EnvVar: "PLUGIN_BASICAUTH_USERNAME",
		},
		cli.StringFlag{
			Name:   "basicauth.password",
			Usage:  "The user's password to authenticate with",
			EnvVar: "PLUGIN_BASICAUTH_PASSWORD",
		},
	}

	if err := app.Run(os.Args); err != nil {
		stdlog.Fatal(err)
	}
}

func action(c *cli.Context) error {
	path := c.String("signalcd.file")
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read SignalCD file from: %s", path)
	}

	config, err := signalcd.ParseConfig(string(fileContent))
	if err != nil {
		return fmt.Errorf("failed to parse SignalCD config: %w", err)
	}

	apiURLFlag := c.String("api.url")
	if apiURLFlag == "" {
		return xerrors.New("no API URL provided")
	}

	apiURL, err := url.Parse(apiURLFlag)
	if err != nil {
		return fmt.Errorf("failed to parse API URL: %w", err)
	}

	client := client.NewHTTPClientWithConfig(
		nil,
		client.DefaultTransportConfig().
			WithSchemes([]string{apiURL.Scheme}).
			WithHost(apiURL.Host).
			WithBasePath(apiURL.Path),
	)

	auth := httptransport.PassThroughAuth

	username := c.String("basicauth.username")
	password := c.String("basicauth.password")
	if username != "" && password != "" {
		auth = httptransport.BasicAuth(username, password)
	}

	pipelineParams := &pipeline.CreateParams{Pipeline: configToPipeline(config)}
	pipelineParams = pipelineParams.WithTimeout(15 * time.Second)
	pipeline, err := client.Pipeline.Create(pipelineParams, auth)
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	deploymentParams := &deployments.SetCurrentDeploymentParams{
		Pipeline: pipeline.Payload.ID.String(),
	}
	deploymentParams.WithTimeout(15 * time.Second)
	_, err = client.Deployments.SetCurrentDeployment(deploymentParams, auth)
	if err != nil {
		return fmt.Errorf("failed to set current deployment: %w", err)
	}

	return nil
}

func configToPipeline(config signalcd.Config) *models.Pipeline {
	p := models.Pipeline{Name: config.Name}

	for _, s := range config.Steps {
		p.Steps = append(p.Steps, &models.Step{
			Name:             &s.Name,
			Image:            &s.Image,
			ImagePullSecrets: s.ImagePullSecrets,
			Commands:         s.Commands,
		})
	}
	for _, c := range config.Checks {
		p.Checks = append(p.Checks, &models.Check{
			Name:             &c.Name,
			Image:            &c.Image,
			ImagePullSecrets: c.ImagePullSecrets,
			Duration:         c.Duration.Seconds(),
			Environment:      nil,
		})
	}

	return &p
}
