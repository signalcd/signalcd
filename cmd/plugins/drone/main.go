package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/ptypes"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
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

	apiURLFlag := c.String(flagAPIURL)
	if apiURLFlag == "" {
		return fmt.Errorf("no API URL provided")
	}

	client := &http.Client{
		Timeout: time.Minute,
	}

	certPath := c.String(flagTLSCert)
	if certPath != "" {
		caCert, err := ioutil.ReadFile(certPath)
		if err != nil {
			return fmt.Errorf("failed to read TLS cert from %s: %w", path, err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}
	}

	username := c.String(flagAuthUsername)
	password := c.String(flagAuthPassword)

	pipelineID, err := createPipeline(client, apiURLFlag, username, password, config)
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	deploymentNumber, err := setCurrentDeployment(client, apiURLFlag, username, password, pipelineID)
	if err != nil {
		return fmt.Errorf("failed to set current deployment pipeline: %w", err)
	}

	fmt.Printf("Crated and applied pipeline %s as deployment %s\n", pipelineID, deploymentNumber)

	return nil
}

func createPipeline(client *http.Client, api string, username string, password string, config signalcd.Config) (string, error) {
	payload := &signalcdproto.CreatePipelineRequest{
		Pipeline: configToPipeline(config),
	}

	payloadBytes, err := json.Marshal(payload.GetPipeline())
	if err != nil {
		return "", err
	}

	url := api + "/api/v1/pipelines"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("unexpected error: %s", resp.Status)
	}

	// TODO: Most likely it's better to generate swagger based Go client, but seems overkill for 2 call right now...
	var respPayload struct {
		Pipeline struct {
			ID string `json:"id"`
		} `json:"pipeline"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
		return "", err
	}

	return respPayload.Pipeline.ID, nil
}

func setCurrentDeployment(client *http.Client, api string, username string, password string, pipelineID string) (string, error) {
	payload := &signalcdproto.SetCurrentDeploymentRequest{Id: pipelineID}

	payloadBytes, err := json.Marshal(payload.Id)
	if err != nil {
		return "", err
	}

	url := api + "/api/v1/deployments/current"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("unexpected error: %s", resp.Status)
	}

	// TODO: Most likely it's better to generate swagger based Go client, but seems overkill for 2 call right now...
	var respPayload struct {
		Deployment struct {
			Number string `json:"number"`
		} `json:"deployment"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
		return "", err
	}

	return respPayload.Deployment.Number, nil
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
