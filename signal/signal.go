package signal

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Labels map[string]string

type Logger interface {
	Log(labels Labels)
}

type Signal struct {
	logger   Logger
	registry *prometheus.Registry
}

func New() *Signal {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.Caller(5))

	return &Signal{
		logger:   &gokitLogger{logger: logger},
		registry: prometheus.NewRegistry(),
	}
}

func (s *Signal) Log(labels Labels) {
	s.logger.Log(labels)
}

type gokitLogger struct {
	logger log.Logger
}

func (g *gokitLogger) Log(labels Labels) {
	var keyvals []interface{}
	for key, value := range labels {
		keyvals = append(keyvals, key)
		keyvals = append(keyvals, value)
	}

	_ = g.logger.Log(keyvals...)
}

func (s *Signal) WithContext(ctx context.Context) *Signal {
	// TODO: create a new signal that adds context aware things
	return s
}

func (s *Signal) MetricsRegister(cs ...prometheus.Collector) {
	s.registry.MustRegister(cs...)
}

func (s *Signal) Serve(addr string) error {
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))

	s.Log(Labels{"msg": "running internal signal server", "addr": addr})

	return http.ListenAndServe(addr, m)
}
