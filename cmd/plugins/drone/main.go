package main

import (
	"context"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net/url"
	"os"

	"github.com/ghodss/yaml"
	apiclient "github.com/signalcd/signalcd/api/client/go"
	"github.com/signalcd/signalcd/signalcd"
	"github.com/urfave/cli"
)

const (
	flagFile         = "signalcd.file"
	flagAPIURL       = "api.url"
	flagAuthPassword = "basicauth.password"
	flagAuthUsername = "basicauth.username"
	flagTLSCert      = "tls.cert"
)

func main() {

	fileFlag := cli.StringFlag{
		Name:   flagFile + ",f",
		Usage:  "The path to the SignalCD file to use",
		EnvVar: "PLUGIN_SIGNALCD_FILE",
		Value:  ".signalcd.yaml",
	}

	app := cli.NewApp()
	app.Name = "SignalCD Drone plugin"
	app.Action = action
	app.Flags = []cli.Flag{
		fileFlag,
		cli.StringFlag{
			Name:   flagAPIURL,
			Usage:  "The URL to talk to the SignalCD API at",
			EnvVar: "PLUGIN_API_URL",
		},
		cli.StringFlag{
			Name:   flagAuthUsername,
			Usage:  "The username to authenticate with",
			EnvVar: "PLUGIN_BASICAUTH_USERNAME",
		},
		cli.StringFlag{
			Name:   flagAuthPassword,
			Usage:  "The user's password to authenticate with",
			EnvVar: "PLUGIN_BASICAUTH_PASSWORD",
		},
		cli.StringFlag{
			Name:  flagTLSCert,
			Usage: "The path to the certificate to use when making requests",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "eval",
			Usage: "Evaluate the given signalcd configuration file",
			Flags: []cli.Flag{
				fileFlag,
			},
			Action: evalAction,
		},
	}

	if err := app.Run(os.Args); err != nil {
		stdlog.Fatal(err)
	}
}

func action(c *cli.Context) error {
	path := c.String(flagFile)
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read SignalCD file from: %s", path)
	}

	config, err := signalcd.ParseConfig(string(fileContent))
	if err != nil {
		return fmt.Errorf("failed to parse SignalCD config: %w", err)
	}

	apiURL, err := url.Parse(c.String(flagAPIURL))
	if err != nil {
		return fmt.Errorf("failed to parse API URL: %w", err)
	}

	clientCfg := apiclient.NewConfiguration()
	clientCfg.Scheme = apiURL.Scheme
	clientCfg.Host = apiURL.Host
	clientCfg.BasePath = apiURL.Path

	client := apiclient.NewAPIClient(clientCfg)

	//certPath := c.String(flagTLSCert)
	//if certPath != "" {
	//	caCert, err := ioutil.ReadFile(certPath)
	//	if err != nil {
	//		return fmt.Errorf("failed to read TLS cert from %s: %w", path, err)
	//	}
	//
	//	caCertPool := x509.NewCertPool()
	//	caCertPool.AppendCertsFromPEM(caCert)
	//
	//	client.Transport = &http.Transport{
	//		TLSClientConfig: &tls.Config{
	//			RootCAs: caCertPool,
	//		},
	//	}
	//}

	var pipelineSteps []apiclient.PipelineSteps
	for _, step := range config.Steps {
		pipelineSteps = append(pipelineSteps, apiclient.PipelineSteps{
			Name:     step.Name,
			Image:    step.Image,
			Commands: step.Commands,
		})
	}

	pipeline, _, err := client.PipelineApi.CreatePipeline(context.Background(), apiclient.Pipeline{
		Name:  config.Name,
		Steps: pipelineSteps,
	})
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	deployment, _, err := client.DeploymentApi.SetCurrentDeployment(context.Background(), apiclient.SetCurrentDeployment{
		PipelineID: pipeline.Id,
	})

	fmt.Println("Created pipeline:", pipeline.Id)
	fmt.Println("Set current deployment:", deployment.Number)

	return nil
}

func evalAction(c *cli.Context) error {
	path := c.String(flagFile)
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read SignalCD file from: %s", path)
	}

	config, err := signalcd.ParseConfig(string(fileContent))
	if err != nil {
		return fmt.Errorf("failed to parse SignalCD config: %w", err)
	}

	// Ignoring error, as this YAML is only for debug printing
	configYAML, _ := yaml.Marshal(config)

	fmt.Println(string(configYAML))

	return nil
}
