package main

import (
	stdlog "log"
	"net/url"
	"os"
	"time"

	"github.com/signalcd/signalcd/api/v1/client"
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
			Name:   "signalcd.file",
			Usage:  "The path to the SignalCD file to use",
			EnvVar: "PLUGIN_SIGNALCD_FILE",
			Value:  ".signalcd.yaml",
		},
	}

	if err := app.Run(os.Args); err != nil {
		stdlog.Fatal(err)
	}
}

func action(c *cli.Context) error {
	path := c.String("signalcd.file")
	file, err := os.Open(path)
	if err != nil {
		return xerrors.Errorf("failed to read SignalCD file from: %s", path)
	}

	config, err := signalcd.ParseConfig(file)
	if err != nil {
		return xerrors.Errorf("failed to parse SignalCD config: %w", err)
	}

	apiURLFlag := c.String("api.url")
	if apiURLFlag == "" {
		return xerrors.New("no API URL provided")
	}

	apiURL, err := url.Parse(apiURLFlag)
	if err != nil {
		return xerrors.Errorf("failed to parse API URL: %w", err)
	}

	client := client.NewHTTPClientWithConfig(
		nil,
		client.DefaultTransportConfig().
			WithSchemes([]string{apiURL.Scheme}).
			WithHost(apiURL.Host).
			WithBasePath(apiURL.Path),
	)

	params := &pipeline.CreateParams{Pipeline: configToPipeline(config)}
	params = params.WithTimeout(15 * time.Second)

	ok, err := client.Pipeline.Create(params)
	if err != nil {
		return xerrors.Errorf("failed creating pipeline with the API: %w", err)
	}

	stdlog.Printf("Created pipeline: %+v\n", ok.Payload)

	return nil
}

func configToPipeline(config signalcd.Config) *models.Pipeline {
	p := models.Pipeline{Name: config.Name}

	for _, s := range config.Steps {
		p.Steps = append(p.Steps, &models.Step{
			Name:     &s.Name,
			Image:    &s.Image,
			Commands: s.Commands,
		})
	}
	for _, c := range config.Checks {
		p.Checks = append(p.Checks, &models.Check{
			Name:        &c.Name,
			Image:       &c.Image,
			Duration:    c.Duration.Seconds(),
			Environment: nil,
		})
	}

	return &p
}
