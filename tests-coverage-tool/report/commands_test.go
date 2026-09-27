package report

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	grpcreflection "google.golang.org/grpc/reflection"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

func reportConfig(t *testing.T) config.Config {
	t.Helper()
	dir := t.TempDir()
	cfg := config.Config{ResultsDir: dir, HistoryDir: dir, HistoryFile: "history.json", HistoryRetentionLimit: 2,
		HTMLReportDir: dir, HTMLReportFile: "report.html", JSONReportDir: dir, JSONReportFile: "report.json"}
	require.NoError(t, os.MkdirAll(cfg.GetResultsDir(), 0o755))
	return cfg
}

func useConfig(t *testing.T, cfg config.Config) {
	t.Helper()
	data, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	filename := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(filename, data, 0o600))
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filename)
}

func TestSaveReportCommand(t *testing.T) {
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	grpcreflection.Register(server)
	httpServer := httptest.NewUnstartedServer(server)
	httpServer.EnableHTTP2 = true
	httpServer.StartTLS()
	t.Cleanup(httpServer.Close)
	t.Cleanup(server.Stop)
	cfg := reportConfig(t)
	cfg.Services = []config.Service{{Key: "api", Host: strings.TrimPrefix(httpServer.URL, "https://")}}
	useConfig(t, cfg)
	require.NoError(t, utils.SaveJSONFile(models.Result{
		Method: "grpc.health.v1.Health.Check", Request: []models.ResultParameters{{Parameter: "service", Covered: true}},
		Response: []models.ResultParameters{{Parameter: "status", Covered: true}},
	}, cfg.GetResultsDir(), "result.json"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, NewSaveReportCommand().ExecuteContext(ctx))
	state, err := utils.ReadJSONFile[models.CoverageState](cfg.GetJSONReportFile())
	require.NoError(t, err)
	assert.Greater(t, state.ServiceCoverages["api"].TotalCoverage, 0.0)
	require.Len(t, state.LogicalServiceCoverages["api"], 1)
	assert.Equal(t, 1, state.LogicalServiceCoverages["api"][0].TotalCoveredMethods)
	assert.FileExists(t, filepath.Join(cfg.HTMLReportDir, cfg.HTMLReportFile))
	assert.FileExists(t, cfg.GetHistoryFile())

	SaveReport()
	state, err = utils.ReadJSONFile[models.CoverageState](cfg.GetJSONReportFile())
	require.NoError(t, err)
	assert.Len(t, state.ServiceCoverages["api"].TotalCoverageHistory, 2)
}

func TestSaveReportCommandErrors(t *testing.T) {
	for _, tc := range []struct {
		name      string
		configure func(*testing.T, *config.Config)
		message   string
	}{
		{name: "missing results", configure: func(t *testing.T, cfg *config.Config) { cfg.ResultsDir = t.TempDir() }, message: "input coverage client"},
		{name: "corrupt history", configure: func(t *testing.T, cfg *config.Config) {
			require.NoError(t, os.WriteFile(cfg.GetHistoryFile(), []byte("{"), 0o600))
		}, message: "input history client factory"},
		{name: "history write failure", configure: func(_ *testing.T, cfg *config.Config) { cfg.HistoryFile = "missing/history.json" }, message: "save history"},
		{name: "HTML write failure", configure: func(t *testing.T, cfg *config.Config) {
			filename := filepath.Join(t.TempDir(), "file")
			require.NoError(t, os.WriteFile(filename, nil, 0o600))
			cfg.HTMLReportDir = filename
		}, message: "save HTML"},
		{name: "JSON write failure", configure: func(t *testing.T, cfg *config.Config) {
			filename := filepath.Join(t.TempDir(), "file")
			require.NoError(t, os.WriteFile(filename, nil, 0o600))
			cfg.JSONReportDir = filename
		}, message: "save JSON"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := reportConfig(t)
			tc.configure(t, &cfg)
			useConfig(t, cfg)
			err := NewSaveReportCommand().Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.message)
		})
	}
	t.Run("invalid config", func(t *testing.T) {
		t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.yaml"))
		err := saveReport(context.Background())
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("cancelled reflection", func(t *testing.T) {
		cfg := reportConfig(t)
		cfg.Services = []config.Service{{Key: "api", Host: "localhost:1"}}
		useConfig(t, cfg)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := NewSaveReportCommand().ExecuteContext(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "service coverage")
	})
}

func TestCopyReportCommand(t *testing.T) {
	t.Chdir(t.TempDir())
	err := NewCopyReportCommand().Execute()
	assert.ErrorIs(t, err, os.ErrNotExist)
	require.NoError(t, os.MkdirAll("submodules/tests-coverage-report/build", 0o755))
	require.NoError(t, os.MkdirAll("tool/report/templates", 0o755))
	input := []byte("<html>updated report</html>")
	require.NoError(t, os.WriteFile("submodules/tests-coverage-report/build/index.html", input, 0o600))
	require.NoError(t, NewCopyReportCommand().Execute())
	CopyReport()
	actual, err := os.ReadFile("tool/report/templates/index.html")
	require.NoError(t, err)
	assert.Equal(t, input, actual)
}
