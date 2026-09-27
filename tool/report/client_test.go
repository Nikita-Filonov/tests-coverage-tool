package report

import (
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
)

type coverageReportClientTest[T any] struct {
	name   string
	want   T
	client CoverageReportClient
}

func TestCoverageReportClientSaveHTMLReport(t *testing.T) {
	tests := []coverageReportClientTest[error]{
		{
			name:   "Empty HTML report dir variable",
			want:   nil,
			client: CoverageReportClient{config: config.Config{}},
		},
		{
			name:   "Empty HTML report file variable",
			want:   nil,
			client: CoverageReportClient{config: config.Config{HTMLReportDir: "."}},
		},
		{
			name: "Empty state",
			want: nil,
			client: CoverageReportClient{
				config: config.Config{HTMLReportDir: t.TempDir(), HTMLReportFile: "report.html"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.SaveHTMLReport())
		})
	}
}

func TestCoverageReportClientSaveJSONReport(t *testing.T) {
	tests := []coverageReportClientTest[error]{
		{
			name:   "Empty JSON report dir variable",
			want:   nil,
			client: CoverageReportClient{config: config.Config{}},
		},
		{
			name:   "Empty JSON report file variable",
			want:   nil,
			client: CoverageReportClient{config: config.Config{JSONReportDir: "."}},
		},
		{
			name: "Empty state",
			want: nil,
			client: CoverageReportClient{
				config: config.Config{JSONReportDir: t.TempDir(), JSONReportFile: "report.json"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.SaveJSONReport())
		})
	}
}

func TestReportSerializationErrors(t *testing.T) {
	state := models.NewCoverageState(config.Config{})
	state.ServiceCoverages["api"] = models.ServiceCoverage{TotalCoverage: math.NaN()}
	dir := t.TempDir()
	client := NewCoverageReportClient(config.Config{
		HTMLReportDir: dir, HTMLReportFile: "report.html", JSONReportDir: dir, JSONReportFile: "report.json",
	}, state)
	assert.Error(t, client.SaveHTMLReport())
	assert.Error(t, client.SaveJSONReport())
	assert.NoFileExists(t, filepath.Join(dir, "report.html"))
	assert.NoFileExists(t, filepath.Join(dir, "report.json"))
}

func TestReportWriteErrors(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	require.NoError(t, os.WriteFile(file, nil, 0o600))
	client := NewCoverageReportClient(config.Config{
		HTMLReportDir: file, HTMLReportFile: "report.html", JSONReportDir: file, JSONReportFile: "report.json",
	}, models.NewCoverageState(config.Config{}))
	assert.Error(t, client.SaveHTMLReport())
	assert.Error(t, client.SaveJSONReport())
}
