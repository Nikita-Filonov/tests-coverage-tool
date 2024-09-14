package coverage

import (
	"strings"
)

type ResultParameters struct {
	Covered    bool               `json:"covered"`
	Parameter  string             `json:"parameter"`
	Parameters []ResultParameters `json:"parameters,omitempty"`
	Deprecated bool               `json:"deprecated"`
}

type Result struct {
	Method   string             `json:"method"`
	Request  []ResultParameters `json:"request"`
	Response []ResultParameters `json:"response"`
}

func (r Result) GetLogicalService() string {
	return r.Method[:strings.LastIndex(r.Method, ".")]
}
