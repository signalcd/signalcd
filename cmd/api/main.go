package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/signalcd/signalcd/api"
	"github.com/signalcd/signalcd/database/boltdb"
	"github.com/signalcd/signalcd/signalcd"
	signalcdproto "github.com/signalcd/signalcd/signalcd/proto"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	flagAddr         = "addr"
	flagAddrAgent    = "addr.agent"
	flagAddrInternal = "addr.internal"
	flagBoltPath     = "bolt.path"
	flagTLSCert      = "tls.cert"
	flagTLSKey       = "tls.key"
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
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log("msg", "failed running api", "err", err)
		os.Exit(1)
	}
}

func apiAction(logger log.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		events := signalcd.NewEvents()

		registry := prometheus.NewRegistry()
		registry.MustRegister(
			prometheus.NewGoCollector(),
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		)

		var db api.SignalDB
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
			apiV1, err := api.NewV1(
				log.WithPrefix(logger, "component", "api"),
				db,
				events,
				c.String(flagAddr),
				c.String(flagTLSCert),
				c.String(flagTLSKey),
			)
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
				// Serving the UI assets
				r.Get("/", file("index.html", "text/html"))
				r.Get("/bulma.min.css", file("bulma.min.css", "text/css"))
				r.Get("/main.dart.js", file("main.dart.js", "application/javascript"))
				r.NotFound(file("index.html", "text/html"))
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
			addr := c.String(flagAddrAgent)
			l, err := net.Listen("tcp", addr)
			if err != nil {
				return xerrors.Errorf("failed to listen on %s: %w", addr, err)
			}

			var server *grpc.Server
			{
				var opts []grpc.ServerOption

				tlsCert, tlsKey := c.String(flagTLSCert), c.String(flagTLSKey)
				if tlsCert != "" && tlsKey != "" {
					creds, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
					if err != nil {
						return fmt.Errorf("failed to create credentials: %w", err)
					}
					opts = append(opts, grpc.Creds(creds))

					level.Debug(logger).Log("msg", "serving requests with TLS", "cert", tlsCert, "key", tlsKey)
				}

				server = grpc.NewServer(opts...)
			}

			signalcdproto.RegisterAgentServiceServer(server,
				api.NewRPC(db, log.WithPrefix(logger, "component", "api")),
			)

			gr.Add(func() error {
				level.Info(logger).Log(
					"msg", "running agent gRPC API server",
					"addr", l.Addr().String(),
				)
				if err := server.Serve(l); err != nil {
					return xerrors.Errorf("failed to serve: %w", err)
				}
				return nil
			}, func(err error) {
				_ = l.Close()
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

func file(name, mime string) http.HandlerFunc {
	file, _ := ioutil.ReadFile("./cmd/api/assets/" + name)
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
