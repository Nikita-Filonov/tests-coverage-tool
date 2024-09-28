package coverageoutput

import (
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
)

type MethodRequestCoverage struct {
	Name                   string                      `json:"name"`
	TotalCoverage          float64                     `json:"totalCoverage,omitempty"`
	TotalParameters        int                         `json:"totalParameters,omitempty"`
	ParametersCoverage     []coverage.ResultParameters `json:"parametersCoverage,omitempty"`
	TotalCoveredParameters int                         `json:"totalCoveredParameters,omitempty"`
}

type MethodCoverage struct {
	Method           string                `json:"method"`
	Covered          bool                  `json:"covered,omitempty"`
	TotalCases       int                   `json:"totalCases,omitempty"`
	Deprecated       bool                  `json:"deprecated,omitempty"`
	RequestCoverage  MethodRequestCoverage `json:"requestCoverage,omitempty"`
	ResponseCoverage MethodRequestCoverage `json:"responseCoverage,omitempty"`
}

type LogicalServiceCoverage struct {
	Methods        []MethodCoverage `json:"methods,omitempty"`
	TotalCoverage  float64          `json:"totalCoverage,omitempty"`
	LogicalService string           `json:"logicalService"`
}

type ServiceCoverage struct {
	TotalCoverage float64 `json:"totalCoverage,omitempty"`
}
