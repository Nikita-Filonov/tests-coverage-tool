package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyMain(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll("submodules/tests-coverage-report/build", 0o755))
	require.NoError(t, os.MkdirAll("tool/report/templates", 0o755))
	input := []byte("<html>report</html>")
	require.NoError(t, os.WriteFile("submodules/tests-coverage-report/build/index.html", input, 0o600))
	main()
	actual, err := os.ReadFile("tool/report/templates/index.html")
	require.NoError(t, err)
	assert.Equal(t, input, actual)
}
