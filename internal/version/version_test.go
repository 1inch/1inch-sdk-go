package version

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "Version is an exact release version",
			pattern: `^v\d+\.\d+\.\d+$`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Regexp(t, regexp.MustCompile(tc.pattern), Version)
		})
	}
}
