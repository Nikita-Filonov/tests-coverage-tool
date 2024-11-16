package history

import (
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
			config: config.Config{HistoryDir: ".", HistoryFile: "history.json"},
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
