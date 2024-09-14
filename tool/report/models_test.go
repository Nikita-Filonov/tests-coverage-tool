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

func TestCoverageReportStateGetMapServiceNameToTotalCoverage(t *testing.T) {
	tests := []coverageReportStateTest[map[string]float64]{
		{
			name: "Multiple services",
			want: map[string]float64{"service1": 100.00, "service2": 66.67, "service3": 33.33},
			model: CoverageReportState{
				Config: config.Config{
					Services: []config.Service{
						{Name: "service1", Host: "http://localhost:1000"},
						{Name: "service2", Host: "http://localhost:2000"},
						{Name: "service3", Host: "http://localhost:3000"},
					},
				},
				ServiceCoverages: map[config.ServiceHost]coverageoutput.ServiceCoverage{
					"http://localhost:1000": {TotalCoverage: 100.00},
					"http://localhost:2000": {TotalCoverage: 66.67},
					"http://localhost:3000": {TotalCoverage: 33.33},
				},
			},
		},
		{
			name:  "Empty services",
			want:  map[string]float64{},
			model: CoverageReportState{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.model.GetMapServiceNameToTotalCoverage())
		})
	}
}
