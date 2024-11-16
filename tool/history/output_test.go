package history

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
)

type outputHistoryClientTest[T any] struct {
	name   string
	want   T
	client OutputHistoryClient
}

func TestOutputHistoryClientSaveHistory(t *testing.T) {
	tests := []outputHistoryClientTest[error]{
		{
			name:   "Empty history dir variable",
			want:   nil,
			client: OutputHistoryClient{config: config.Config{}},
		},
		{
			name:   "Empty history file variable",
			want:   nil,
			client: OutputHistoryClient{config: config.Config{HistoryDir: "."}},
		},
		{
			name: "Empty state",
			want: nil,
			client: OutputHistoryClient{
				config: config.Config{HistoryDir: ".", HistoryFile: "history.json"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.SaveHistory())
		})
	}
}
