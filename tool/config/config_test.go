package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type configTest[T any] struct {
	name   string
	want   T
	config Config
}

func TestConfigGetResultsDir(t *testing.T) {
	tests := []configTest[string]{
		{
			name:   "Results dir set",
			want:   "./coverage-results",
			config: Config{ResultsDir: "."},
		},
		{
			name:   "Results dir unset",
			want:   "/coverage-results",
			config: Config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.config.GetResultsDir())
		})
	}
}

func TestConfigGetJSONReportFile(t *testing.T) {
	tests := []configTest[string]{
		{
			name:   "JSON report dir set, JSON report file set",
			want:   "./report.json",
			config: Config{JSONReportDir: ".", JSONReportFile: "report.json"},
		},
		{
			name:   "JSON report dir unset, JSON report file unset",
			want:   "/",
			config: Config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.config.GetJSONReportFile())
		})
	}
}
