package service

import (
	"testing"

	"gotest.tools/assert"
)

func TestSanitizeClusterName(t *testing.T) {
	cases := map[string]string{
		"foo/bar":                "foo_bar",
		"löl":                    "l_l",
		"nothing_to-see_here-42": "nothing_to-see_here-42",
	}

	for input, expectedOutput := range cases {
		assert.Equal(t, expectedOutput, sanitizeClusterName(input))
	}
}
