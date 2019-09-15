package signal

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
)

type Labels map[string]string

type Logger interface {
	Log(labels Labels)
}

type Signal struct {
	logger Logger
}

func New() *Signal {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC)
	logger = log.WithPrefix(logger, "caller", log.Caller(5))

	return &Signal{
		logger: &gokitLogger{
			logger: logger,
		},
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
