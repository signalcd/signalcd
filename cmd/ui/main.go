package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

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
			Name:  "addr",
			Usage: "The addr to run the UI HTTP server on",
			Value: "localhost:6670",
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

			router.Get("/", file("index.html", "text/html"))
			router.Get("/bulma.min.css", file("bulma.min.css", "text/css"))
			router.Get("/main.dart.js", file("main.dart.js", "application/javascript"))

			s := http.Server{
				Addr:    c.String("addr"),
				Handler: router,
			}

			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running UI HTTP server",
					"addr", s.Addr,
				)
				return s.ListenAndServeTLS("./development/signalcd.dev+6.pem", "./development/signalcd.dev+6-key.pem")
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

func file(name, mime string) http.HandlerFunc {
	file, _ := ioutil.ReadFile("./cmd/ui/assets/" + name)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime)
		_, _ = w.Write(file)
	}
}
