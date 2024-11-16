package models

import (
	"encoding/json"
	"time"

	"github.com/samber/lo"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
)

type CoverageState struct {
	Config                  config.Config                                  `json:"config"`
	CreatedAt               time.Time                                      `json:"createdAt"`
	ServiceCoverages        map[config.ServiceKey]ServiceCoverage          `json:"serviceCoverages"`
	LogicalServiceCoverages map[config.ServiceKey][]LogicalServiceCoverage `json:"logicalServiceCoverages"`
}

func NewCoverageState(conf config.Config) CoverageState {
	return CoverageState{
		Config:                  conf,
		CreatedAt:               time.Now(),
		ServiceCoverages:        make(map[config.ServiceKey]ServiceCoverage),
		LogicalServiceCoverages: make(map[config.ServiceKey][]LogicalServiceCoverage),
	}
}

func (s CoverageState) GetReportStateJSON() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s CoverageState) GetHistoryState() HistoryState {
	return lo.SliceToMap(s.Config.Services, s.buildServiceHistory)
}

func (s CoverageState) buildServiceHistory(service config.Service) (config.ServiceKey, History) {
	return service.Key, History{
		Service: s.buildServiceTotalCoverage(service.Key),
		LogicalServices: lo.SliceToMap(
			s.LogicalServiceCoverages[service.Key],
			s.buildLogicalServiceHistory,
		),
	}
}

func (s CoverageState) buildServiceTotalCoverage(serviceKey config.ServiceKey) ServiceHistory {
	return ServiceHistory{TotalCoverage: s.ServiceCoverages[serviceKey].TotalCoverageHistory}
}

func (s CoverageState) buildLogicalServiceHistory(coverage LogicalServiceCoverage) (string, LogicalServiceHistory) {
	return coverage.LogicalService, LogicalServiceHistory{
		Methods:       lo.SliceToMap(coverage.Methods, s.buildMethodHistory),
		TotalCoverage: coverage.TotalCoverageHistory,
	}
}

func (s CoverageState) buildMethodHistory(coverage MethodCoverage) (string, MethodHistory) {
	return coverage.Method, MethodHistory{
		RequestTotalCoverage:  coverage.RequestCoverage.TotalCoverageHistory,
		ResponseTotalCoverage: coverage.ResponseCoverage.TotalCoverageHistory,
	}
}
