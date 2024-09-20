package report

import (
	"encoding/json"
	"time"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageoutput"
)

type CoverageReportState struct {
	Config                  config.Config                                                  `json:"config"`
	CreatedAt               time.Time                                                      `json:"createdAt"`
	ServiceCoverages        map[config.ServiceHost]coverageoutput.ServiceCoverage          `json:"serviceCoverages"`
	LogicalServiceCoverages map[config.ServiceHost][]coverageoutput.LogicalServiceCoverage `json:"logicalServiceCoverages"`
}

func NewCoverageReportState(conf config.Config) CoverageReportState {
	return CoverageReportState{
		Config:                  conf,
		CreatedAt:               time.Now(),
		ServiceCoverages:        make(map[config.ServiceHost]coverageoutput.ServiceCoverage),
		LogicalServiceCoverages: make(map[config.ServiceHost][]coverageoutput.LogicalServiceCoverage),
	}
}

func (s CoverageReportState) getState() map[StateKey]any {
	return map[StateKey]any{
		configKey:                  s.Config,
		createdAtKey:               s.CreatedAt,
		serviceCoveragesKey:        s.ServiceCoverages,
		logicalServiceCoveragesKey: s.LogicalServiceCoverages,
	}
}

func (s CoverageReportState) getStateJSON() (string, error) {
	data, err := json.Marshal(s.getState())
	if err != nil {
		return "", err
	}

	return string(data), nil
}
