package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	var output bytes.Buffer
	cmd := newRootCommand()
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--help"})
	require.NoError(t, cmd.Execute())
	for _, name := range []string{"copy-report", "save-report", "print-config"} {
		assert.Contains(t, output.String(), name)
	}
	cmd = newRootCommand()
	cmd.SetArgs([]string{"unknown-command"})
	require.Error(t, cmd.Execute())
}

func TestMainPrintConfig(t *testing.T) {
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	t.Setenv("TESTS_COVERAGE_RESULTS_DIR", "cli-results")
	args := os.Args
	os.Args = []string{"tests-coverage-tool", "print-config"}
	t.Cleanup(func() { os.Args = args })
	var output bytes.Buffer
	writer := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(writer) })
	main()
	assert.Contains(t, output.String(), "ResultsDir:cli-results")
}
