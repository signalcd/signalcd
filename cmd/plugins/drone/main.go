package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/ptypes"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
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

	apiURLFlag := c.String(flagAPIURL)
	if apiURLFlag == "" {
		return xerrors.New("no API URL provided")
	}

	err = createPipeline(apiURLFlag, config)
	if err != nil {
		return err
	}

	//setCurrentD

	//var client signalcdproto.UIServiceClient
	//{
	//	var opts []grpc.DialOption
	//
	//	tlsCert := c.String(flagTLSCert)
	//	if tlsCert != "" {
	//		creds, err := credentials.NewClientTLSFromFile(tlsCert, "")
	//		if err != nil {
	//			return fmt.Errorf("failed to load credentials: %w", err)
	//		}
	//
	//		opts = append(opts, grpc.WithTransportCredentials(creds))
	//
	//		stdlog.Println("Making requests with TLS")
	//	} else {
	//		stdlog.Println("Making requests unencrypted")
	//		opts = append(opts, grpc.WithInsecure())
	//	}
	//
	//	dialCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//	defer cancel()
	//	conn, err := grpc.DialContext(dialCtx, apiURLFlag, opts...)
	//	if err != nil {
	//		return fmt.Errorf("failed to connect to the api: %w", err)
	//	}
	//	defer conn.Close()
	//
	//	client = signalcdproto.NewUIServiceClient(conn)
	//}
	//
	//var pipelineResp *signalcdproto.CreatePipelineResponse
	//{
	//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer cancel()
	//
	//	pipelineResp, err = client.CreatePipeline(ctx, &signalcdproto.CreatePipelineRequest{
	//		Pipeline: configToPipeline(config),
	//	})
	//	if err != nil {
	//		return fmt.Errorf("failed to create pipeline: %w", err)
	//	}
	//}
	//
	//var deploymentResp *signalcdproto.SetCurrentDeploymentResponse
	//{
	//	deploymentResp, err = client.SetCurrentDeployment(context.Background(), &signalcdproto.SetCurrentDeploymentRequest{
	//		Id: pipelineResp.GetPipeline().GetId(),
	//	})
	//	if err != nil {
	//		return fmt.Errorf("failed to set pipeline as current deployment: %w", err)
	//	}
	//}

	//fmt.Printf("Crated and applied pipeline %s as deployment %d\n", pipelineResp.GetPipeline().GetId(), deploymentResp.Deployment.Number)

	return nil
}

func createPipeline(api string, config signalcd.Config) error {
	payload := &signalcdproto.CreatePipelineRequest{
		Pipeline: configToPipeline(config),
	}

	payloadBytes, err := json.Marshal(payload.GetPipeline())
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, api+"/api/v1/pipelines", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

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
