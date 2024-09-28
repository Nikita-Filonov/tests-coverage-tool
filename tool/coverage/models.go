package coverage

import (
	"strings"
)

type ResultParameters struct {
	Covered                bool               `json:"covered,omitempty"`
	Parameter              string             `json:"parameter"`
	Parameters             []ResultParameters `json:"parameters,omitempty"`
	Deprecated             bool               `json:"deprecated,omitempty"`
	HasUncoveredParameters bool               `json:"hasUncoveredParameters,omitempty"`
}

type Result struct {
	Method   string             `json:"method"`
	Request  []ResultParameters `json:"request"`
	Response []ResultParameters `json:"response"`
}

func (r Result) GetLogicalService() string {
	return r.Method[:strings.LastIndex(r.Method, ".")]
}
