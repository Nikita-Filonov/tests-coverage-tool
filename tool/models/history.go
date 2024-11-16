package models

import (
	"time"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
)

type CoverageHistory struct {
	CreatedAt     time.Time `json:"createdAt"`
	TotalCoverage float64   `json:"totalCoverage,omitempty"`
}

type MethodHistory struct {
	RequestTotalCoverage  []CoverageHistory `json:"requestCoverage,omitempty"`
	ResponseTotalCoverage []CoverageHistory `json:"responseCoverage,omitempty"`
}

type LogicalServiceHistory struct {
	Methods       map[string]MethodHistory `json:"methods,omitempty"`
	TotalCoverage []CoverageHistory        `json:"totalCoverage,omitempty"`
}

type ServiceHistory struct {
	TotalCoverage []CoverageHistory `json:"totalCoverage,omitempty"`
}

type History struct {
	Service         ServiceHistory                   `json:"service,omitempty"`
	LogicalServices map[string]LogicalServiceHistory `json:"logicalServices,omitempty"`
}

type HistoryState map[config.ServiceKey]History
