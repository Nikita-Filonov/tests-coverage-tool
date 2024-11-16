package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
)

type newCoverageStateTest[T any] struct {
	name string
	want T
	conf config.Config
}

type coverageStateTest[T any] struct {
	name  string
	want  T
	state CoverageState
}

var createdAt = time.Now()

var defaultCoverageHistory = []CoverageHistory{
	{CreatedAt: createdAt, TotalCoverage: 99.99},
	{CreatedAt: createdAt, TotalCoverage: 77.77},
}

var defaultLogicalServiceCoverage = LogicalServiceCoverage{
	Methods: []MethodCoverage{
		{
			Method:           "method-1",
			RequestCoverage:  MethodRequestCoverage{TotalCoverageHistory: defaultCoverageHistory},
			ResponseCoverage: MethodRequestCoverage{TotalCoverageHistory: defaultCoverageHistory},
		},
	},
	LogicalService:       "logical-service-1",
	TotalCoverageHistory: defaultCoverageHistory,
}

func TestNewCoverageState(t *testing.T) {
	tests := []newCoverageStateTest[CoverageState]{
		{
			name: "New coverage report state",
			want: CoverageState{
				Config:                  config.Config{},
				CreatedAt:               time.Now(),
				ServiceCoverages:        map[config.ServiceKey]ServiceCoverage{},
				LogicalServiceCoverages: map[config.ServiceKey][]LogicalServiceCoverage{},
			},
			conf: config.Config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NewCoverageState(test.conf)
			test.want.CreatedAt = result.CreatedAt

			assert.Equal(t, test.want, result)
		})
	}
}

func TestCoverageStateGetHistoryState(t *testing.T) {
	tests := []coverageStateTest[HistoryState]{
		{
			name:  "Empty state",
			want:  HistoryState{},
			state: CoverageState{},
		},
		{
			name: "Existing state",
			want: HistoryState{
				"service-1": History{
					Service: ServiceHistory{TotalCoverage: defaultCoverageHistory},
					LogicalServices: map[string]LogicalServiceHistory{
						"logical-service-1": {
							Methods: map[string]MethodHistory{
								"method-1": {
									RequestTotalCoverage:  defaultCoverageHistory,
									ResponseTotalCoverage: defaultCoverageHistory,
								},
							},
							TotalCoverage: defaultCoverageHistory,
						},
					},
				},
				"service-2": History{
					Service: ServiceHistory{TotalCoverage: defaultCoverageHistory},
					LogicalServices: map[string]LogicalServiceHistory{
						"logical-service-1": {
							Methods: map[string]MethodHistory{
								"method-1": {
									RequestTotalCoverage:  defaultCoverageHistory,
									ResponseTotalCoverage: defaultCoverageHistory,
								},
							},
							TotalCoverage: defaultCoverageHistory,
						},
					},
				},
			},
			state: CoverageState{
				Config: config.Config{
					Services: []config.Service{{Key: "service-1"}, {Key: "service-2"}},
				},
				ServiceCoverages: map[config.ServiceKey]ServiceCoverage{
					"service-1": {TotalCoverageHistory: defaultCoverageHistory},
					"service-2": {TotalCoverageHistory: defaultCoverageHistory},
				},
				LogicalServiceCoverages: map[config.ServiceKey][]LogicalServiceCoverage{
					"service-1": {defaultLogicalServiceCoverage},
					"service-2": {defaultLogicalServiceCoverage},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.state.GetHistoryState())
		})
	}
}
