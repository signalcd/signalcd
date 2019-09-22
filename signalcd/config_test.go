package signalcd

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		env    func(string) string
		expect Config
		error  error
	}{
		{
			name:   "Empty",
			input:  ``,
			env:    func(s string) string { return "" },
			expect: Config{},
			error:  io.EOF,
		},
		{
			name:   "NameOnly",
			input:  `name: foobar`,
			env:    func(s string) string { return "" },
			expect: Config{Name: "foobar"},
			error:  nil,
		},
		{
			name: "StepTagEnvvar",
			input: `
name: foobar
steps:
- name: foobar
  image: foobar:${SOME_VAR}
		`,
			env: func(s string) string {
				if s == "SOME_VAR" {
					return "baz"
				}
				return ""
			},
			expect: Config{
				Name: "foobar",
				Steps: []ConfigStep{{
					Name:  "foobar",
					Image: "foobar:baz",
				}},
			},
			error: nil,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := parseConfigEnv(strings.TrimSpace(tc.input), tc.env)
			assert.Equal(t, tc.expect, output)
			if tc.error == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tc.error))
			}
		})
	}
}
