package coverageoutput

import (
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
)

type MethodRequestCoverage struct {
	Name                   string                      `json:"name"`
	TotalCoverage          float64                     `json:"totalCoverage"`
	TotalParameters        int                         `json:"totalParameters"`
	ParametersCoverage     []coverage.ResultParameters `json:"parametersCoverage,omitempty"`
	TotalCoveredParameters int                         `json:"totalCoveredParameters"`
}

type MethodCoverage struct {
	Method           string                `json:"method"`
	Covered          bool                  `json:"covered"`
	TotalCases       int                   `json:"totalCases"`
	Deprecated       bool                  `json:"deprecated"`
	RequestCoverage  MethodRequestCoverage `json:"requestCoverage"`
	ResponseCoverage MethodRequestCoverage `json:"responseCoverage"`
}

type LogicalServiceCoverage struct {
	Methods        []MethodCoverage `json:"methods"`
	TotalCoverage  float64          `json:"totalCoverage"`
	LogicalService string           `json:"logicalService"`
}

type ServiceCoverage struct {
	TotalCoverage float64 `json:"totalCoverage"`
}
