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

func TestConfigValidate(t *testing.T) {
	tests := []configTest[error]{
		{
			name:   "Empty config",
			want:   nil,
			config: Config{},
		},
		{
			name:   "Valid services",
			want:   nil,
			config: Config{Services: []Service{{Key: "1"}, {Key: "2"}}},
		},
		{
			name:   "Services with empty keys",
			want:   ServiceKeysShouldNotContainEmptyValuesError,
			config: Config{Services: []Service{{}, {}}},
		},
		{
			name:   "Services with duplicated keys",
			want:   DuplicateServiceKeysFoundInConfigurationError,
			config: Config{Services: []Service{{Key: "1"}, {Key: "1"}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.config.validate())
		})
	}
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

func TestConfigGetHistoryFile(t *testing.T) {
	tests := []configTest[string]{
		{
			name:   "History dir set, history file set",
			want:   "./history.json",
			config: Config{HistoryDir: ".", HistoryFile: "history.json"},
		},
		{
			name:   "History dir unset, history file unset",
			want:   "/",
			config: Config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.config.GetHistoryFile())
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
