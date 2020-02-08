package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/ptypes"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fileFlag := cli.StringFlag{
		Name:   "signalcd.file,f",
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
			Name:   "api.url",
			Usage:  "The URL to talk to the SignalCD API at",
			EnvVar: "PLUGIN_API_URL",
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

	var client signalcdproto.UIServiceClient
	{
		pair, err := tls.LoadX509KeyPair("./development/signalcd.dev+6.pem", "./development/signalcd.dev+6-key.pem")
		if err != nil {
			return err
		}

		cert, err := ioutil.ReadFile("./development/signalcd.dev+6.pem")
		if err != nil {
			return err
		}

		pool := x509.NewCertPool()

		ok := pool.AppendCertsFromPEM(cert)
		if !ok {
			return fmt.Errorf("failed to appened certificate")
		}

		creds := credentials.NewTLS(&tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{pair},
		})
		opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

		dialCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(dialCtx, apiURLFlag, opts...)
		if err != nil {
			return fmt.Errorf("failed to connect to the api: %w", err)
		}
		defer conn.Close()

		client = signalcdproto.NewUIServiceClient(conn)
	}

	var pipelineResp *signalcdproto.CreatePipelineResponse
	{
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pipelineResp, err = client.CreatePipeline(ctx, &signalcdproto.CreatePipelineRequest{
			Pipeline: configToPipeline(config),
		})
		if err != nil {
			return fmt.Errorf("failed to create pipeline: %w", err)
		}
	}

	var deploymentResp *signalcdproto.SetCurrentDeploymentResponse
	{
		deploymentResp, err = client.SetCurrentDeployment(context.Background(), &signalcdproto.SetCurrentDeploymentRequest{
			Id: pipelineResp.GetPipeline().GetId(),
		})
		if err != nil {
			return fmt.Errorf("failed to set pipeline as current deployment: %w", err)
		}
	}

	fmt.Printf("Crated and applied pipeline %s as deployment %d\n", pipelineResp.GetPipeline().GetId(), deploymentResp.Deployment.Number)

	return nil
}

func configToPipeline(config signalcd.Config) *signalcdproto.Pipeline {
	p := &signalcdproto.Pipeline{}

	p.Name = config.Name

	for _, s := range config.Steps {
		p.Steps = append(p.Steps, &signalcdproto.Step{
			Name:             s.Name,
			Image:            s.Image,
			ImagePullSecrets: s.ImagePullSecrets,
			Commands:         s.Commands,
		})
	}

	for _, c := range config.Checks {
		p.Checks = append(p.Checks, &signalcdproto.Check{
			Name:             c.Name,
			Image:            c.Image,
			ImagePullSecrets: c.ImagePullSecrets,
			Duration:         ptypes.DurationProto(c.Duration),
		})
	}

	return p
}

func evalAction(c *cli.Context) error {
	path := c.String("signalcd.file")
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
