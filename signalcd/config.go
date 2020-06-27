package signalcd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/drone/envsubst"
	"github.com/ghodss/yaml"
)

// Config for a SignalCD Pipeline
type Config struct {
	Kind   string        `json:"kind"`
	Name   string        `json:"name"`
	Steps  []ConfigStep  `json:"steps"`
	Checks []ConfigCheck `json:"checks"`
}

// ConfigStep for a SignalCD Pipeline Step
type ConfigStep struct {
	Name             string   `json:"name"`
	Image            string   `json:"image"`
	ImagePullSecrets []string `json:"imagePullSecrets"`
	Commands         []string `json:"commands"`
}

// ConfigCheck for a SignalCD Pipeline Check
type ConfigCheck struct {
	Name             string        `json:"name"`
	Image            string        `json:"image"`
	ImagePullSecrets []string      `json:"imagePullSecrets"`
	Duration         time.Duration `json:"duration"`
}

// ParseConfig decodes a io.Reader into a SignalCD Config
func ParseConfig(s string) (Config, error) {
	return parseConfigEnv(s, os.Getenv)
}

func parseConfigEnv(s string, env func(string) string) (Config, error) {
	var c Config

	if s == "" {
		return c, io.EOF
	}

	s, err := envsubst.Eval(s, env)
	if err != nil {
		return c, fmt.Errorf("failed to substitute environment variable: %w", err)
	}

	err = yaml.Unmarshal([]byte(s), &c)
	if err != nil {
		return c, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return c, nil
}
