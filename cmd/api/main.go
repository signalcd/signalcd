package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/signalcd/signalcd/database/boltdb"
	"github.com/signalcd/signalcd/signalcd"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

const (
	flagAddr         = "addr"
	flagAddrAgent    = "addr.agent"
	flagAddrInternal = "addr.internal"
	flagBoltPath     = "bolt.path"
	flagTLSCert      = "tls.cert"
	flagTLSKey       = "tls.key"
	flagUIAssets     = "ui.assets"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	app := cli.NewApp()
	app.Action = apiAction(logger)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagBoltPath,
			Value: "./development/data",
		},
		cli.StringFlag{
			Name:  flagAddr,
			Usage: "The address for the public HTTP/gRPC API server",
			Value: ":6660",
		},
		cli.StringFlag{
			Name:  flagAddrAgent,
			Usage: "The address for the agent gRPC server",
			Value: ":6661",
		},
		cli.StringFlag{
			Name:  flagAddrInternal,
			Usage: "The address for the internal HTTP server",
			Value: ":6662",
		},
		cli.StringFlag{
			Name:  flagTLSCert,
			Usage: "The path to the TLS certificate",
		},
		cli.StringFlag{
			Name:  flagTLSKey,
			Usage: "The path to the TLS key",
		},
		cli.StringFlag{
			Name:  flagUIAssets,
			Usage: "The path to the UI assets on disk",
			Value: "/assets",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running api", "err", err)
		os.Exit(1)
	}
}

func apiAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		registry := prometheus.NewRegistry()
		registry.MustRegister(
			prometheus.NewGoCollector(),
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		)

		events := signalcd.NewEvents()

		var db Database
		{
			bolt, dbClose, err := boltdb.New(c.String(flagBoltPath))
			if err != nil {
				return xerrors.Errorf("failed to create bolt db: %w", err)
			}
			defer dbClose()

			boltEvents := boltdb.NewEvents(bolt, events)

			db = boltEvents
		}

		var gr run.Group
		{
			apiV1, err := NewV1(logger, registry, db, events)
			if err != nil {
				return fmt.Errorf("failed to initialize api: %w", err)
			}

			r := chi.NewRouter()
			r.Use(Logger(logger))
			r.Use(HTTPMetrics(registry))

			{
				// Serving the HTTP/gRPC API
				r.Mount("/", apiV1)
			}
			{
				directory := c.String(flagUIAssets)
				if _, err := os.Stat(directory); os.IsNotExist(err) {
					return fmt.Errorf("assets directory does not exist: %s", directory)
				}

				// Serving the UI assets
				r.Get("/", file(directory, "index.html", "text/html"))
				r.Get("/bulma.min.css", file(directory, "bulma.min.css", "text/css"))
				r.Get("/main.dart.js", file(directory, "main.dart.js", "application/javascript"))
				r.NotFound(file(directory, "index.html", "text/html"))
			}

			s := http.Server{
				Addr:    c.String(flagAddr),
				Handler: r,
			}

			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running public HTTP/gRPC API server",
					"addr", s.Addr,
				)

				tlsCert, tlsKey := c.String(flagTLSCert), c.String(flagTLSKey)
				if tlsCert != "" && tlsKey != "" {
					return s.ListenAndServeTLS(tlsCert, tlsKey)
				} else {
					return s.ListenAndServe()
				}
			}, func(err error) {
				_ = s.Shutdown(context.TODO())
			})
		}
		{
			r := chi.NewRouter()
			r.Mount("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
			r.Mount("/debug", middleware.Profiler())

			s := http.Server{
				Addr:    c.String(flagAddrInternal),
				Handler: r,
			}
			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running internal HTTP API",
					"addr", s.Addr,
				)

				return s.ListenAndServe()
			}, func(err error) {
				_ = s.Shutdown(context.TODO())
			})
		}

		if err := gr.Run(); err != nil {
			return xerrors.Errorf("error running: %w", err)
		}

		return nil
	}
}

func file(directory, name, mime string) http.HandlerFunc {
	file, _ := ioutil.ReadFile(filepath.Join(directory, name))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime)
		_, _ = w.Write(file)
	}
}

// Logger returns a middleware to log HTTP requests
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
				"content", r.Header.Get("Content-Type"),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}

// HTTPMetrics returns a middleware to track HTTP requests with Prometheus metrics
func HTTPMetrics(registry *prometheus.Registry) func(next http.Handler) http.Handler {
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Tracks the latencies for HTTP requests.",
	}, nil)

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Tracks the number of HTTP requests.",
	}, []string{"code", "method"})

	registry.MustRegister(duration, counter)

	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerDuration(duration,
			promhttp.InstrumentHandlerCounter(counter, next),
		)
	}
}
