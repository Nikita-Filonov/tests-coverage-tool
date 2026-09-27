package config

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintConfigCommand(t *testing.T) {
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	t.Setenv("TESTS_COVERAGE_RESULTS_DIR", "custom-results")
	var output bytes.Buffer
	original := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(original) })
	require.NoError(t, NewPrintConfigCommand().Execute())
	assert.Contains(t, output.String(), "ResultsDir:custom-results")
	PrintConfig()
	assert.Contains(t, output.String(), "Tests coverage config")
}

func TestPrintConfigCommandInvalidConfig(t *testing.T) {
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.yaml"))
	assert.ErrorIs(t, NewPrintConfigCommand().Execute(), os.ErrNotExist)
}
