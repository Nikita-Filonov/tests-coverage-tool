package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func TestReportStateRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{JSONReportDir: dir, JSONReportFile: "report.json", HTMLReportDir: dir, HTMLReportFile: "report.html",
		Services: []config.Service{{Key: "api", Name: "Price $1 and ${value}"}}}
	state := models.NewCoverageState(cfg)
	state.ServiceCoverages["api"] = models.ServiceCoverage{TotalCoverage: 50}
	client := NewCoverageReportClient(cfg, state)
	require.NoError(t, client.SaveJSONReport())
	require.NoError(t, client.SaveHTMLReport())

	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	t.Setenv("TESTS_COVERAGE_JSON_REPORT_DIR", dir)
	t.Setenv("TESTS_COVERAGE_JSON_REPORT_FILE", "report.json")
	loaded, err := ReadCoverageReportState()
	require.NoError(t, err)
	assert.Equal(t, state.ServiceCoverages, loaded.ServiceCoverages)
	assert.Equal(t, cfg.Services, loaded.Config.Services)
	assert.True(t, state.CreatedAt.Equal(loaded.CreatedAt))

	html, err := os.ReadFile(filepath.Join(dir, "report.html"))
	require.NoError(t, err)
	match := regexp.MustCompile(`<script id="state" type="application/json">([\s\S]*?)</script>`).FindSubmatch(html)
	require.Len(t, match, 2)
	var embedded models.CoverageState
	require.NoError(t, json.Unmarshal(match[1], &embedded))
	assert.Equal(t, loaded.Config.Services, embedded.Config.Services)
	assert.Equal(t, loaded.ServiceCoverages, embedded.ServiceCoverages)
}

func TestReadCoverageReportStateErrors(t *testing.T) {
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	dir := t.TempDir()
	t.Setenv("TESTS_COVERAGE_JSON_REPORT_DIR", dir)
	t.Setenv("TESTS_COVERAGE_JSON_REPORT_FILE", "missing.json")
	_, err := ReadCoverageReportState()
	assert.ErrorIs(t, err, os.ErrNotExist)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "missing.json"), []byte("{"), 0o600))
	_, err = ReadCoverageReportState()
	assert.Error(t, err)

	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filepath.Join(dir, "missing.yaml"))
	_, err = ReadCoverageReportState()
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestReadCoverageReportStateDisabled(t *testing.T) {
	for _, content := range []string{"jsonReportDir: ''\n", "jsonReportFile: ''\n"} {
		t.Run(content, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "config.yaml")
			require.NoError(t, os.WriteFile(filename, []byte(content), 0o600))
			t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filename)
			state, err := ReadCoverageReportState()
			assert.NoError(t, err)
			assert.Nil(t, state)
		})
	}
}
