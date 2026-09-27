package reflection

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestExpectedEnumParameters(t *testing.T) {
	parameters := BuildExpectedResultParameters((&grpc_health_v1.HealthCheckResponse{}).ProtoReflect().Descriptor())
	require.Len(t, parameters, 1)
	assert.Equal(t, "status", parameters[0].Parameter)
	require.Len(t, parameters[0].Parameters, 4)
	assert.Equal(t, "SERVING", parameters[0].Parameters[1].Parameter)
	assert.False(t, parameters[0].Parameters[1].Covered)
}

func TestExpectedMapAndRecursiveParameters(t *testing.T) {
	parameters := BuildExpectedResultParameters((&structpb.Struct{}).ProtoReflect().Descriptor())
	require.Len(t, parameters, 1)
	assert.Equal(t, "fields", parameters[0].Parameter)
	require.NotEmpty(t, parameters[0].Parameters)
	// Struct -> Value -> Struct and Value -> ListValue -> Value must terminate.
	assert.Less(t, len(parameters[0].Parameters), 10)
}

func TestExpectedDeprecatedAndSiblingParameters(t *testing.T) {
	file, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name: proto.String("coverage.proto"), Syntax: proto.String("proto3"), Package: proto.String("coverage"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: proto.String("Child"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: proto.String("name"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					Options: &descriptorpb.FieldOptions{Deprecated: proto.Bool(true)}},
			}},
			{Name: proto.String("Parent"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: proto.String("first"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".coverage.Child")},
				{Name: proto.String("second"), Number: proto.Int32(2), Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".coverage.Child")},
			}},
		},
	}, nil)
	require.NoError(t, err)
	parameters := BuildExpectedResultParameters(file.Messages().ByName("Parent"))
	require.Len(t, parameters, 2)
	for _, field := range parameters {
		require.Len(t, field.Parameters, 1)
		assert.Equal(t, "name", field.Parameters[0].Parameter)
		assert.True(t, field.Parameters[0].Deprecated)
	}
}
