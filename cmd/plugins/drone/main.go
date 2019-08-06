package main

import (
	stdlog "log"
	"net/url"
	"os"
	"time"

	"github.com/signalcd/signalcd/api/v1/client"
	"github.com/signalcd/signalcd/api/v1/client/pipeline"
	"github.com/signalcd/signalcd/api/v1/models"
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
	}

	if err := app.Run(os.Args); err != nil {
		stdlog.Fatal(err)
	}
}

func action(c *cli.Context) error {
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

	p := models.Pipeline{Name: "foobar"}

	params := &pipeline.CreateParams{Pipeline: &p}
	params = params.WithTimeout(15 * time.Second)

	ok, err := client.Pipeline.Create(params)
	if err != nil {
		return err
	}

	stdlog.Printf("Created pipeline: %+v\n", ok.Payload)

	return nil
}
