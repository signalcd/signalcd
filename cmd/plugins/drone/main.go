package main

import (
	stdlog "log"
	"net/url"
	"os"

	"golang.org/x/xerrors"

	"github.com/signalcd/signalcd/api/v1/client/pipeline"
	"github.com/signalcd/signalcd/api/v1/models"

	"github.com/signalcd/signalcd/api/v1/client"
	"github.com/urfave/cli"
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

	ok, err := client.Pipeline.Create(&pipeline.CreateParams{
		Pipeline: &models.Pipeline{
			Name: "foobar",
		},
	})
	if err != nil {
		return err
	}

	stdlog.Println("Create pipeline:", ok.Payload.ID)

	return nil
}
