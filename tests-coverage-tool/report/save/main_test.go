package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyMain(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	for _, key := range []string{"TESTS_COVERAGE_RESULTS_DIR", "TESTS_COVERAGE_HISTORY_DIR", "TESTS_COVERAGE_HTML_REPORT_DIR", "TESTS_COVERAGE_JSON_REPORT_DIR"} {
		t.Setenv(key, dir)
	}
	require.NoError(t, os.Mkdir(filepath.Join(dir, "coverage-results"), 0o755))
	main()
	data, err := os.ReadFile(filepath.Join(dir, "coverage-report.json"))
	require.NoError(t, err)
	assert.True(t, json.Valid(data))
	assert.FileExists(t, filepath.Join(dir, "index.html"))
	assert.FileExists(t, filepath.Join(dir, "coverage-history.json"))
}
