package signalcd

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		expect Config
		error  error
	}{
		{
			name:   "Empty",
			input:  ``,
			expect: Config{},
			error:  io.EOF,
		},
		{
			name:   "NameOnly",
			input:  `name: foobar`,
			expect: Config{Name: "foobar"},
			error:  nil,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := ParseConfig(bytes.NewBufferString(tc.input))
			assert.Equal(t, tc.expect, output)
			if tc.error == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tc.error))
			}
		})
	}
}
