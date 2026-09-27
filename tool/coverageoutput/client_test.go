package coverageoutput

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	grpcreflection "google.golang.org/grpc/reflection"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageinupt"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/history"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/reflection"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/report"
)

func TestCoveragePipeline(t *testing.T) {
	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("api", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	grpcreflection.Register(server)
	httpServer := httptest.NewUnstartedServer(server)
	httpServer.EnableHTTP2 = true
	httpServer.StartTLS()
	t.Cleanup(httpServer.Close)
	t.Cleanup(server.Stop)

	dir := t.TempDir()
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	t.Setenv("TESTS_COVERAGE_RESULTS_DIR", dir)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.NewClient(strings.TrimPrefix(httpServer.URL, "https://"),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
		grpc.WithChainUnaryInterceptor(coverageinupt.CoverageInterceptor()))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })
	response, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "api"})
	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)

	service := config.Service{Key: "api", Host: strings.TrimPrefix(httpServer.URL, "https://")}
	cfg := config.Config{Services: []config.Service{service}, ResultsDir: dir,
		HistoryDir: dir, HistoryFile: "history.json", HistoryRetentionLimit: 2,
		JSONReportDir: dir, JSONReportFile: "report.json", HTMLReportDir: dir, HTMLReportFile: "report.html"}
	input, err := coverageinupt.NewInputCoverageClient(cfg.GetResultsDir())
	require.NoError(t, err)
	factory, err := history.NewInputHistoryClientFactory(cfg)
	require.NoError(t, err)
	reflectionClient, err := reflection.NewGRPCReflectionClient(ctx, service)
	require.NoError(t, err)
	defer func() { require.NoError(t, reflectionClient.Close()) }()
	services, err := reflectionClient.GetServices()
	require.NoError(t, err)
	assert.Equal(t, []string{"grpc.health.v1.Health"}, services)
	_, err = reflectionClient.GetServiceMethods("missing.Service")
	assert.Error(t, err)
	client, err := NewOutputCoverageClient(reflectionClient, factory.NewClient(service.Key), input)
	require.NoError(t, err)
	state := models.NewCoverageState(cfg)
	state.ServiceCoverages[service.Key], err = client.GetServiceCoverage()
	require.NoError(t, err)
	state.LogicalServiceCoverages[service.Key], err = client.GetLogicalServiceCoverages()
	require.NoError(t, err)

	logical := state.LogicalServiceCoverages[service.Key]
	require.Len(t, logical, 1)
	assert.Equal(t, 1, logical[0].TotalCoveredMethods)
	assert.GreaterOrEqual(t, logical[0].TotalMethods, 2)
	assert.Equal(t, getRoundedTotalCoverage(100/float64(logical[0].TotalMethods)), logical[0].TotalCoverage)
	assert.Equal(t, logical[0].TotalCoverage, state.ServiceCoverages[service.Key].TotalCoverage)
	foundCheck := false
	for _, method := range logical[0].Methods {
		if method.Method == "Check" {
			foundCheck = true
			assert.True(t, method.Covered)
			assert.Equal(t, 1, method.TotalCases)
			assert.Equal(t, 100.0, method.RequestCoverage.TotalCoverage)
			assert.Len(t, method.RequestCoverage.TotalCoverageHistory, 1)
		} else {
			assert.False(t, method.Covered)
			assert.Zero(t, method.TotalCases)
		}
	}
	assert.True(t, foundCheck)

	output := report.NewCoverageReportClient(cfg, state)
	require.NoError(t, output.SaveJSONReport())
	require.NoError(t, output.SaveHTMLReport())
	historyOutput := history.NewOutputHistoryClient(cfg, state)
	require.NoError(t, historyOutput.SaveHistory())
	loaded, err := history.ReadHistoryState(cfg)
	require.NoError(t, err)
	expectedJSON, err := json.Marshal(state.GetHistoryState())
	require.NoError(t, err)
	loadedJSON, err := json.Marshal(loaded)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedJSON), string(loadedJSON))
}

func TestNewOutputCoverageClientRejectsNilInputs(t *testing.T) {
	for _, tc := range []struct {
		name       string
		reflection *reflection.GRPCReflectionClient
		history    *history.InputHistoryClient
		input      *coverageinupt.InputCoverageClient
	}{
		{name: "reflection"},
		{name: "history", reflection: &reflection.GRPCReflectionClient{}},
		{name: "coverage", reflection: &reflection.GRPCReflectionClient{}, history: &history.InputHistoryClient{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewOutputCoverageClient(tc.reflection, tc.history, tc.input)
			assert.Nil(t, client)
			assert.Error(t, err)
		})
	}
}

func TestCoverageFailsWhenServiceDescriptorCannotBeResolved(t *testing.T) {
	server := grpc.NewServer()
	server.RegisterService(&grpc.ServiceDesc{ServiceName: "unknown.Service", HandlerType: (*interface{})(nil), Metadata: "missing.proto"}, nil)
	grpcreflection.Register(server)
	httpServer := httptest.NewUnstartedServer(server)
	httpServer.EnableHTTP2 = true
	httpServer.StartTLS()
	t.Cleanup(httpServer.Close)
	t.Cleanup(server.Stop)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	reflectionClient, err := reflection.NewGRPCReflectionClient(ctx, config.Service{Host: strings.TrimPrefix(httpServer.URL, "https://")})
	require.NoError(t, err)
	defer func() { require.NoError(t, reflectionClient.Close()) }()
	client, err := NewOutputCoverageClient(reflectionClient, &history.InputHistoryClient{}, &coverageinupt.InputCoverageClient{})
	require.NoError(t, err)
	_, err = client.getLogicalServiceCoverage("missing.Service")
	assert.Error(t, err)
	_, err = client.GetServiceCoverage()
	assert.Error(t, err)
	coverages, err := client.GetLogicalServiceCoverages()
	assert.Error(t, err)
	assert.Nil(t, coverages)
}

func TestCoveragePropagatesCancelledReflection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	reflectionClient, err := reflection.NewGRPCReflectionClient(ctx, config.Service{Host: "localhost:1"})
	require.NoError(t, err)
	defer func() { require.NoError(t, reflectionClient.Close()) }()
	client, err := NewOutputCoverageClient(reflectionClient, &history.InputHistoryClient{}, &coverageinupt.InputCoverageClient{})
	require.NoError(t, err)
	_, err = client.GetServiceCoverage()
	assert.Error(t, err)
	_, err = client.GetLogicalServiceCoverages()
	assert.Error(t, err)
}
