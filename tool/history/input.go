package history

import (
	"fmt"
	"sort"
	"time"

	"github.com/samber/lo"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type InputHistoryClient struct {
	config            config.Config
	history           models.History
	createdAt         time.Time
	isHistoryDisabled bool
}

type InputHistoryClientFactory struct {
	state             models.HistoryState
	config            config.Config
	createdAt         time.Time
	isHistoryDisabled bool
}

func NewInputHistoryClientFactory(conf config.Config) (*InputHistoryClientFactory, error) {
	state, err := ReadHistoryState(conf)
	if err != nil {
		return nil, fmt.Errorf("error reading history state: %w", err)
	}

	return &InputHistoryClientFactory{
		state:             state,
		config:            conf,
		createdAt:         time.Now(),
		isHistoryDisabled: lo.IsNil(state),
	}, nil
}

func (f *InputHistoryClientFactory) NewClient(key config.ServiceKey) *InputHistoryClient {
	return &InputHistoryClient{
		config:            f.config,
		history:           f.state[key],
		createdAt:         f.createdAt,
		isHistoryDisabled: f.isHistoryDisabled,
	}
}

func (c *InputHistoryClient) buildCoverageHistory(totalCoverage float64) models.CoverageHistory {
	if totalCoverage > 100 {
		totalCoverage = 100
	}

	return models.CoverageHistory{TotalCoverage: totalCoverage, CreatedAt: c.createdAt}
}

func (c *InputHistoryClient) appendCoverageHistory(history []models.CoverageHistory, totalCoverage float64) []models.CoverageHistory {
	if c.isHistoryDisabled {
		return []models.CoverageHistory{}
	}

	if totalCoverage == 0 {
		return history
	}

	result := append(history, c.buildCoverageHistory(totalCoverage))
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })

	if len(result) > c.config.HistoryRetentionLimit {
		result = result[len(result)-c.config.HistoryRetentionLimit:]
	}

	return result
}

func (c *InputHistoryClient) BuildServiceHistoryTotalCoverage(totalCoverage float64) []models.CoverageHistory {
	return c.appendCoverageHistory(c.history.Service.TotalCoverage, totalCoverage)
}

func (c *InputHistoryClient) BuildLogicalServiceHistoryTotalCoverage(logicalService string, totalCoverage float64) []models.CoverageHistory {
	return c.appendCoverageHistory(c.history.LogicalServices[logicalService].TotalCoverage, totalCoverage)
}

func (c *InputHistoryClient) BuildMethodHistoryRequestTotalCoverage(logicalService, method string, totalCoverage float64) []models.CoverageHistory {
	return c.appendCoverageHistory(
		c.history.LogicalServices[logicalService].Methods[method].RequestTotalCoverage,
		totalCoverage,
	)
}

func (c *InputHistoryClient) BuildMethodHistoryResponseTotalCoverage(logicalService, method string, totalCoverage float64) []models.CoverageHistory {
	return c.appendCoverageHistory(
		c.history.LogicalServices[logicalService].Methods[method].ResponseTotalCoverage,
		totalCoverage,
	)
}
