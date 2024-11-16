package models

type MethodRequestCoverage struct {
	Name                   string             `json:"name"`
	TotalCoverage          float64            `json:"totalCoverage,omitempty"`
	TotalParameters        int                `json:"totalParameters,omitempty"`
	ParametersCoverage     []ResultParameters `json:"parametersCoverage,omitempty"`
	TotalCoverageHistory   []CoverageHistory  `json:"totalCoverageHistory,omitempty"`
	TotalCoveredParameters int                `json:"totalCoveredParameters,omitempty"`
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
	Methods              []MethodCoverage  `json:"methods,omitempty"`
	TotalCoverage        float64           `json:"totalCoverage,omitempty"`
	LogicalService       string            `json:"logicalService"`
	TotalCoverageHistory []CoverageHistory `json:"totalCoverageHistory,omitempty"`
}

type ServiceCoverage struct {
	TotalCoverage        float64           `json:"totalCoverage,omitempty"`
	TotalCoverageHistory []CoverageHistory `json:"totalCoverageHistory,omitempty"`
}
