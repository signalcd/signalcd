package main

import (
	stdlog "log"
	"os"

	"github.com/go-openapi/strfmt"

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
	client := client.NewHTTPClientWithConfig(
		nil,
		client.DefaultTransportConfig().
			WithSchemes([]string{"http"}).
			WithHost(c.String("api.url")),
	)

	id := strfmt.UUID("f05ae07b-7f6d-4ffe-82ba-30dc4c5e1e31")

	ok, err := client.Pipeline.Create(&pipeline.CreateParams{
		Pipeline: &models.Pipeline{
			ID:   &id,
			Name: "foobar",
		},
	})
	if err != nil {
		return err
	}

	stdlog.Println("Create pipeline:", ok.Payload.ID)

	return nil
}
