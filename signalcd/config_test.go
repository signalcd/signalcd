package signalcd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		expect Config
	}{
		{
			name:   "Empty",
			input:  ``,
			expect: Config{},
		},
		{
			name:   "NameOnly",
			input:  `name: foobar`,
			expect: Config{Name: "foobar"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := ParseConfig(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, output)
		})
	}
}
