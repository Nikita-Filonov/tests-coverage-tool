package history

import (
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type readHistoryStateTest struct {
	err    error
	name   string
	state  models.HistoryState
	config config.Config
}

func TestReadHistoryState(t *testing.T) {
	tests := []readHistoryStateTest{
		{
			err:    nil,
			name:   "Empty history dir variable",
			state:  nil,
			config: config.Config{},
		},
		{
			err:    nil,
			name:   "Empty history file variable",
			state:  nil,
			config: config.Config{HistoryDir: "."},
		},
		{
			err:    nil,
			name:   "History file does not exists",
			state:  models.HistoryState{},
			config: config.Config{HistoryDir: t.TempDir(), HistoryFile: "history.json"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state, err := ReadHistoryState(test.config)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.state, state)
		})
	}
}

func TestHistoryReadAndWriteErrors(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{Services: []config.Service{{Key: "api"}}, HistoryDir: dir, HistoryFile: "history.json"}
	require.NoError(t, os.WriteFile(cfg.GetHistoryFile(), []byte("{"), 0o600))
	_, err := ReadHistoryState(cfg)
	assert.Error(t, err)
	_, err = NewInputHistoryClientFactory(cfg)
	assert.Error(t, err)
	state := models.NewCoverageState(cfg)
	state.ServiceCoverages["api"] = models.ServiceCoverage{TotalCoverageHistory: []models.CoverageHistory{{TotalCoverage: math.NaN()}}}
	client := NewOutputHistoryClient(cfg, state)
	assert.Error(t, client.SaveHistory())
}

func TestHistoryDisabledFactory(t *testing.T) {
	factory, err := NewInputHistoryClientFactory(config.Config{})
	require.NoError(t, err)
	client := factory.NewClient("api")
	assert.Empty(t, client.BuildServiceHistoryTotalCoverage(50))
}
