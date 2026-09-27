package coverageinupt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBuildCoverageResult(t *testing.T) {
	result, err := buildCoverageResult("/grpc.health.v1.Health/Check",
		&grpc_health_v1.HealthCheckRequest{Service: "api"},
		&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING})
	require.NoError(t, err)
	assert.Equal(t, "grpc.health.v1.Health.Check", result.Method)
	require.Len(t, result.Request, 1)
	assert.True(t, result.Request[0].Covered)
	require.Len(t, result.Response, 1)
	assert.True(t, result.Response[0].Covered)
	assert.Equal(t, "status", result.Response[0].Parameter)
	require.Len(t, result.Response[0].Parameters, 4)
	assert.True(t, result.Response[0].Parameters[1].Covered)
}

func TestBuildCoverageResultMalformedMessages(t *testing.T) {
	for _, tc := range []struct {
		name       string
		req, reply any
	}{
		{name: "invalid request", req: "invalid", reply: &wrapperspb.StringValue{}},
		{name: "invalid response", req: &wrapperspb.StringValue{}, reply: nil},
		{name: "typed nil request", req: (*wrapperspb.StringValue)(nil), reply: &wrapperspb.StringValue{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := buildCoverageResult("/service/Method", tc.req, tc.reply)
			assert.Error(t, err)
		})
	}
}

func TestDoubleDefaultValue(t *testing.T) {
	for _, tc := range []struct {
		value   float64
		covered bool
	}{{0, false}, {1.5, true}} {
		parameters := buildActualResultParameters(wrapperspb.Double(tc.value))
		require.Len(t, parameters, 1)
		assert.Equal(t, tc.covered, parameters[0].Covered)
	}
}

func TestBuildActualMapAndRepeatedMessages(t *testing.T) {
	message, err := structpb.NewStruct(map[string]any{"name": "api", "count": 2.0})
	require.NoError(t, err)
	parameters := buildActualResultParameters(message)
	require.Len(t, parameters, 1)
	assert.True(t, parameters[0].Covered)
	require.Len(t, parameters[0].Parameters, 6)
	for _, child := range parameters[0].Parameters {
		if child.Parameter == "string_value" || child.Parameter == "number_value" {
			assert.True(t, child.Covered)
		}
	}

	list, err := structpb.NewList([]any{"api", 2.0})
	require.NoError(t, err)
	parameters = buildActualResultParameters(list)
	require.Len(t, parameters, 1)
	assert.True(t, parameters[0].Covered)
	require.Len(t, parameters[0].Parameters, 6)
	for _, child := range parameters[0].Parameters {
		if child.Parameter == "string_value" || child.Parameter == "number_value" {
			assert.True(t, child.Covered)
		}
	}
}
