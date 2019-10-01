package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/urfave/cli"
)

var (
	config = struct {
		URL string
	}{}

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url",
			EnvVar:      "PLUGIN_URL",
			Destination: &config.URL,
		},
	}
)

func main() {
	app := cli.NewApp()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	app.Name = "kubernetes status"
	app.Action = run
	app.Flags = flags

	if err := app.Run(os.Args); err != nil {
		level.Error(logger).Log(
			"msg", "failed to run",
			"err", err,
		)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	u, err := url.Parse(config.URL)
	if err != nil {
		return fmt.Errorf("provided URL is not valid: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	threshold := 5
	counter := 0

	ticker := time.NewTicker(5 * time.Second)

	println("starting to probe")
	for {
		select {
		case <-ctx.Done():
			println("done")
			return nil
		case <-ticker.C:
			println("probing...")
			if err := probe(u); err != nil {
				counter++
				if counter == threshold {
					return fmt.Errorf("failed to many times: %d", counter)
				}
			}
		}
	}
}

func probe(u *url.URL) error {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("response did not return 200 OK but: %s", resp.Status)
	}

	return nil
}
