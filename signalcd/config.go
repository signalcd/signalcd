package signalcd

import (
	"io"
	"time"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Kind   string        `json:"kind"`
	Name   string        `json:"name"`
	Steps  []ConfigStep  `json:"steps"`
	Checks []ConfigCheck `json:"checks"`
}

type ConfigStep struct {
	Name     string   `json:"name"`
	Image    string   `json:"image"`
	Commands []string `json:"commands"`
}

type ConfigCheck struct {
	Name     string        `json:"name"`
	Image    string        `json:"image"`
	Duration time.Duration `json:"duration"`
}

func ParseConfig(r io.Reader) (Config, error) {
	var c Config

	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return c, xerrors.Errorf("failed to unmarshal config: %w", err)
	}

	return c, nil
}
