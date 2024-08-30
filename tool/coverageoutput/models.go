package coverageoutput

type ParameterCoverage struct {
	Parameter  string `json:"parameter"`
	Covered    bool   `json:"covered"`
	TotalCases int    `json:"totalCases"`
}

type MethodRequestCoverage struct {
	Name                   string              `json:"name"`
	TotalParameters        int                 `json:"totalParameters"`
	ParametersCoverage     []ParameterCoverage `json:"parametersCoverage"`
	TotalCoveredParameters int                 `json:"totalCoveredParameters"`
}

type MethodCoverage struct {
	Method           string                `json:"method"`
	Covered          bool                  `json:"covered"`
	TotalCases       int                   `json:"totalCases"`
	RequestCoverage  MethodRequestCoverage `json:"requestCoverage"`
	ResponseCoverage MethodRequestCoverage `json:"responseCoverage"`
}

type LogicalServiceCoverage struct {
	Methods        []MethodCoverage `json:"methods"`
	TotalCoverage  float64          `json:"totalCoverage"`
	LogicalService string           `json:"logicalService"`
}
