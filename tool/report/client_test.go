package report

import (
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
				config: config.Config{HTMLReportDir: ".", HTMLReportFile: "report.html"},
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
				config: config.Config{JSONReportDir: ".", JSONReportFile: "report.json"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.SaveJSONReport())
		})
	}
}
