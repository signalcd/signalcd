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
	"github.com/oklog/run"
	"github.com/urfave/cli"
)

var (
	config = struct {
		Duration  time.Duration
		URL       string
		Threshold uint
	}{}

	flags = []cli.Flag{
		cli.DurationFlag{
			Name:        "duration",
			Usage:       "How long these probes should be run in total",
			EnvVar:      "PLUGIN_DURATION",
			Value:       time.Minute,
			Destination: &config.Duration,
		},
		cli.StringFlag{
			Name:        "url",
			EnvVar:      "PLUGIN_URL",
			Destination: &config.URL,
		},
		cli.UintFlag{
			Name:        "threshold",
			Usage:       "The maximum number of probes allowed to fail",
			EnvVar:      "PLUGIN_THRESHOLD",
			Value:       2,
			Destination: &config.Threshold,
		},
	}
)

type update struct {
	Success bool
}

func main() {
	app := cli.NewApp()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)

	app.Name = "kubernetes status"
	app.Action = action(logger)
	app.Flags = flags

	if err := app.Run(os.Args); err != nil {
		level.Error(logger).Log(
			"msg", "failed to run",
			"err", err,
		)
		os.Exit(1)
	}
}

func action(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {

		updates := make(chan update)

		var gr run.Group

		gr.Add(func() error {
			return probe(logger, updates)
		}, func(err error) {
			close(updates)
		})

		gr.Add(func() error {
			timer := time.NewTimer(5 * time.Second)

			for {
				select {
				case <-timer.C:
					fmt.Println("send status: checking...")
				case update := <-updates:
					if update.Success {
						fmt.Println("send status: success")
					} else {
						fmt.Println("send status: failure")
					}
					return nil
				}
			}
		}, func(err error) {
		})

		return gr.Run()
	}
}

func probe(logger log.Logger, updates chan<- update) error {
	u, err := url.Parse(config.URL)
	if err != nil {
		return fmt.Errorf("provided URL is not valid: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	var failed uint

	ticker := time.NewTicker(5 * time.Second)

	level.Info(logger).Log("msg", "starting to probe")
	for {
		select {
		case <-ctx.Done():
			level.Info(logger).Log("msg", "done", "status", "success")
			level.Debug(logger).Log("msg", "failed probes", "failed", failed, "threshold", config.Threshold)
			updates <- update{Success: true}
			return nil
		case <-ticker.C:
			level.Debug(logger).Log("msg", "probing", "addr", u.String())
			if err := request(u); err != nil {
				failed++
				if failed >= config.Threshold {
					updates <- update{Success: false}
					return fmt.Errorf("failed to many times: %d", failed)
				}
			}
			level.Debug(logger).Log("msg", "current number of failed probes", "failed", failed, "threshold", config.Threshold)
		}
	}
}

func request(u *url.URL) error {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("response did not return 200 OK but: %s", resp.Status)
	}

	return nil
}

func status() error {
	return nil
}
