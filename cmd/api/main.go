package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/signalcd/signalcd/api"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	var gr run.Group
	{
		router, err := api.NewV1()
		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to initialize api",
				"err", err,
			)
			os.Exit(1)
		}

		s := http.Server{
			Addr:    ":6660",
			Handler: router,
		}

		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "running api",
				"addr", s.Addr,
			)

			return s.ListenAndServe()
		}, func(err error) {
			_ = s.Shutdown(context.TODO())
		})
	}

	if err := gr.Run(); err != nil {
		level.Error(logger).Log(
			"msg", "error running",
			"err", err,
		)
	}
}
