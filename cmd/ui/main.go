package main

import (
	"context"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Action = uiAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "assets.path",
			Usage: "Path on the filesystem to the folder containing the assets",
			Value: "/assets",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running ui", "err", err)
		os.Exit(1)
	}
}

func uiAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		var gr run.Group
		{
			router := chi.NewRouter()

			folder := c.String("assets.path")
			router.Get("/", file(folder, "index.html", "text/html"))
			router.Get("/bulma.min.css", file(folder, "bulma.min.css", "text/css"))
			router.Get("/main.dart.js", file(folder, "main.dart.js", "application/javascript"))

			s := http.Server{
				Addr:    ":6662",
				Handler: router,
			}

			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running ui",
					"addr", s.Addr,
				)
				return s.ListenAndServe()
			}, func(err error) {
				_ = s.Shutdown(context.Background())
			})
		}

		if err := gr.Run(); err != nil {
			return xerrors.Errorf("error running: %w", err)
		}

		return nil
	}
}

func file(folder, name, mime string) http.HandlerFunc {
	fp := filepath.Join(folder, name)
	file, err := ioutil.ReadFile(fp)
	if err != nil {
		stdlog.Fatalf("failed to read assets: %s", fp)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime)
		_, _ = w.Write(file)
	}
}
