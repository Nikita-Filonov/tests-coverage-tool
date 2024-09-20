package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageoutput"
)

type coverageReportStateTest[T any] struct {
	name  string
	want  T
	model CoverageReportState
}

func TestCoverageReportStateGetState(t *testing.T) {
	tests := []coverageReportStateTest[map[StateKey]any]{
		{
			name: "New coverage report state",
			want: map[StateKey]any{
				configKey:                  config.Config{},
				createdAtKey:               time.Now(),
				serviceCoveragesKey:        map[config.ServiceHost]coverageoutput.ServiceCoverage{},
				logicalServiceCoveragesKey: map[config.ServiceHost][]coverageoutput.LogicalServiceCoverage{},
			},
			model: NewCoverageReportState(config.Config{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.want[createdAtKey] = test.model.CreatedAt

			assert.Equal(t, test.want, test.model.getState())
		})
	}
}
