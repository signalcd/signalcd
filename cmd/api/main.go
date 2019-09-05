package main

import (
	"context"
	"fmt"
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
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Action = apiAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bolt.path",
			Value: "./development/data",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running api", "err", err)
		os.Exit(1)
	}
}

func apiAction(logger log.Logger) cli.ActionFunc {
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
			apiV1, err := api.NewV1(log.WithPrefix(logger, "component", "api"), db, events)
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
				level.Info(logger).Log(
					"msg", "running HTTP API",
					"addr", s.Addr,
				)
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
				level.Info(logger).Log(
					"msg", "running gRPC API",
					"addr", l.Addr().String(),
				)
				if err := s.Serve(l); err != nil {
					return xerrors.Errorf("failed to serve: %w", err)
				}
				return nil
			}, func(err error) {
				_ = l.Close()
			})
		}
		{
			for i := 0; i < 3; i++ {
				deployments := make(chan signalcd.Deployment)
				subscription := events.SubscribeDeployments(deployments)

				gr.Add(func(i int) func() error {
					return func() error {
						for d := range deployments {
							fmt.Printf("deployment event %d: %+v\n", i, d)
						}
						return nil
					}
				}(i), func(err error) {
					events.UnsubscribeDeployments(subscription)
				},
				)
			}
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
