package coverageinupt

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

func TestInterceptorPreservesRPCError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
	t.Setenv("TESTS_COVERAGE_RESULTS_DIR", dir)
	rpcErr := errors.New("RPC failed")
	called := false
	invoker := func(_ context.Context, method string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		called = true
		assert.Equal(t, "/service/Method", method)
		return rpcErr
	}
	err := CoverageInterceptor()(context.Background(), "/service/Method", wrapperspb.String("request"), wrapperspb.String("response"), nil, invoker)
	assert.ErrorIs(t, err, rpcErr)
	assert.True(t, called)
	files, err := os.ReadDir(filepath.Join(dir, "coverage-results"))
	require.NoError(t, err)
	require.Len(t, files, 1)
	result, err := utils.ReadJSONFile[models.Result](filepath.Join(dir, "coverage-results", files[0].Name()))
	require.NoError(t, err)
	assert.Equal(t, "service.Method", result.Method)
}

func TestInterceptorIgnoresCollectionErrors(t *testing.T) {
	for _, name := range []string{"malformed message", "invalid config", "unwritable results"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("TESTS_COVERAGE_CONFIG_FILE", "")
			t.Setenv("TESTS_COVERAGE_RESULTS_DIR", dir)
			var request any = wrapperspb.String("request")
			switch name {
			case "malformed message":
				request = "invalid"
			case "invalid config":
				t.Setenv("TESTS_COVERAGE_CONFIG_FILE", filepath.Join(dir, "missing.yaml"))
			case "unwritable results":
				require.NoError(t, os.WriteFile(filepath.Join(dir, "coverage-results"), []byte("file"), 0o600))
			}
			invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error { return nil }
			err := CoverageInterceptor()(context.Background(), "/service/Method", request, wrapperspb.String("response"), nil, invoker)
			assert.NoError(t, err)
		})
	}
}
