package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigDefaults(t *testing.T) {
	for _, key := range []string{
		"TESTS_COVERAGE_CONFIG_FILE", "TESTS_COVERAGE_RESULTS_DIR", "TESTS_COVERAGE_HISTORY_DIR",
		"TESTS_COVERAGE_HISTORY_FILE", "TESTS_COVERAGE_HTML_REPORT_DIR", "TESTS_COVERAGE_JSON_REPORT_DIR",
		"TESTS_COVERAGE_HTML_REPORT_FILE", "TESTS_COVERAGE_JSON_REPORT_FILE", "TESTS_COVERAGE_HISTORY_RETENTION_LIMIT",
	} {
		t.Setenv(key, "")
	}
	cfg, err := NewConfig()
	require.NoError(t, err)
	assert.Equal(t, ".", cfg.ResultsDir)
	assert.Equal(t, "coverage-history.json", cfg.HistoryFile)
	assert.Equal(t, "index.html", cfg.HTMLReportFile)
	assert.Equal(t, "coverage-report.json", cfg.JSONReportFile)
	assert.Equal(t, 30, cfg.HistoryRetentionLimit)
}

func TestNewConfigYAMLOverridesEnvironment(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(filename, []byte(`services:
  - key: service
    host: localhost:50051
resultsDir: yaml-results
historyDir: ""
historyFile: ""
historyRetentionLimit: 10
`), 0o600))
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filename)
	t.Setenv("TESTS_COVERAGE_RESULTS_DIR", "env-results")
	t.Setenv("TESTS_COVERAGE_JSON_REPORT_FILE", "env-report.json")
	cfg, err := NewConfig()
	require.NoError(t, err)
	assert.Equal(t, "yaml-results", cfg.ResultsDir)
	assert.Equal(t, "env-report.json", cfg.JSONReportFile)
	assert.Empty(t, cfg.HistoryDir)
	assert.Empty(t, cfg.HistoryFile)
	assert.Equal(t, 10, cfg.HistoryRetentionLimit)
	require.Len(t, cfg.Services, 1)
	assert.Equal(t, ServiceKey("service"), cfg.Services[0].Key)
}

func TestNewConfigErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		yaml string
		want error
	}{
		{name: "malformed YAML", yaml: "services: ["},
		{name: "empty service key", yaml: "services:\n  - key: ''\n", want: ErrServiceKeysShouldNotContainEmptyValues},
		{name: "duplicate service key", yaml: "services:\n  - key: api\n  - key: api\n", want: ErrDuplicateServiceKeysFoundInConfiguration},
	} {
		t.Run(tc.name, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "config.yaml")
			require.NoError(t, os.WriteFile(filename, []byte(tc.yaml), 0o600))
			t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filename)
			_, err := NewConfig()
			require.Error(t, err)
			if tc.want != nil {
				assert.ErrorIs(t, err, tc.want)
			}
		})
	}
	t.Run("missing file", func(t *testing.T) {
		t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.yaml"))
		_, err := NewConfig()
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("invalid environment", func(t *testing.T) {
		t.Setenv("TESTS_COVERAGE_HISTORY_RETENTION_LIMIT", "invalid")
		_, err := NewConfig()
		assert.Error(t, err)
	})
}
