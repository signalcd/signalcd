package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/signalcd/signalcd/api"
	"github.com/signalcd/signalcd/database/boltdb"
	"github.com/signalcd/signalcd/signal"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func main() {
	signals := signal.New()

	app := cli.NewApp()
	app.Action = apiAction(signals)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bolt.path",
			Value: "./development/data",
		},
	}

	if err := app.Run(os.Args); err != nil {
		signals.Log(signal.Labels{"msg": "failed running api", "err": err.Error()})
		os.Exit(1)
	}
}

func apiAction(signals *signal.Signal) cli.ActionFunc {
	logger := log.NewNopLogger()

	return func(c *cli.Context) error {
		events := signalcd.NewEvents()

		var db api.SignalDB
		{
			bolt, dbClose, err := boltdb.New(c.String("bolt.path"))
			if err != nil {
				return xerrors.Errorf("failed to create bolt db: %w", err)
			}
			defer dbClose()

			boltEvents := boltdb.NewEvents(bolt, events)

			db = boltEvents
		}

		var gr run.Group
		{
			apiV1, err := api.NewV1(signals, db, events)
			if err != nil {
				return xerrors.Errorf("failed to initialize api: %w", err)
			}

			r := chi.NewRouter()
			r.Use(Logger(logger))
			r.Mount("/", apiV1)

			s := http.Server{
				Addr:    ":6660",
				Handler: r,
			}

			gr.Add(func() error {
				signals.Log(signal.Labels{"msg": "running HTTP API", "addr": s.Addr})
				return s.ListenAndServe()
			}, func(err error) {
				_ = s.Shutdown(context.TODO())
			})
		}
		{
			const addr = ":6661"
			l, err := net.Listen("tcp", addr)
			if err != nil {
				return xerrors.Errorf("failed to listen on %s: %w", addr, err)
			}

			s := grpc.NewServer()

			signalcdproto.RegisterAgentServiceServer(s,
				api.NewRPC(db, log.WithPrefix(logger, "component", "api")),
			)

			gr.Add(func() error {
				signals.Log(signal.Labels{"msg": "running gRPC API", "addr": l.Addr().String()})
				if err := s.Serve(l); err != nil {
					return xerrors.Errorf("failed to serve: %w", err)
				}
				return nil
			}, func(err error) {
				_ = l.Close()
			})
		}

		if err := gr.Run(); err != nil {
			return xerrors.Errorf("error running: %w", err)
		}

		return nil
	}
}

func Logger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			level.Debug(logger).Log(
				"proto", r.Proto,
				"method", r.Method,
				"status", ww.Status(),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}
