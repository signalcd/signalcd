package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/metalmatze/cd/api"
	"github.com/oklog/run"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	var gr run.Group
	{
		router := api.New()

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
