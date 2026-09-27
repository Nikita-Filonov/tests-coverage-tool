package coverageinupt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestScalarDefaults(t *testing.T) {
	for _, tc := range []struct {
		name        string
		zero, value proto.Message
	}{
		{"bool", wrapperspb.Bool(false), wrapperspb.Bool(true)},
		{"int32", wrapperspb.Int32(0), wrapperspb.Int32(-1)},
		{"int64", wrapperspb.Int64(0), wrapperspb.Int64(1)},
		{"uint32", wrapperspb.UInt32(0), wrapperspb.UInt32(1)},
		{"uint64", wrapperspb.UInt64(0), wrapperspb.UInt64(1)},
		{"float", wrapperspb.Float(0), wrapperspb.Float(1.5)},
		{"double", wrapperspb.Double(0), wrapperspb.Double(1.5)},
		{"string", wrapperspb.String(""), wrapperspb.String("api")},
		{"bytes", wrapperspb.Bytes(nil), wrapperspb.Bytes([]byte("api"))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parameters := buildActualResultParameters(tc.zero)
			require.Len(t, parameters, 1)
			assert.False(t, parameters[0].Covered)
			parameters = buildActualResultParameters(tc.value)
			require.Len(t, parameters, 1)
			assert.True(t, parameters[0].Covered)
		})
	}
	assert.Nil(t, buildActualResultParameters(nil))
	assert.Nil(t, buildActualResultParameters((*wrapperspb.StringValue)(nil)))
}

func coverageDescriptor(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()
	field := func(name string, number int32, kind descriptorpb.FieldDescriptorProto_Type, typeName string) *descriptorpb.FieldDescriptorProto {
		result := &descriptorpb.FieldDescriptorProto{Name: proto.String(name), Number: proto.Int32(number), Type: kind.Enum()}
		if typeName != "" {
			result.TypeName = proto.String(typeName)
		}
		return result
	}
	labels := field("labels", 1, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".coverage.Root.LabelsEntry")
	labels.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
	enums := field("statuses", 2, descriptorpb.FieldDescriptorProto_TYPE_ENUM, ".coverage.Status")
	enums.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
	empties := field("empties", 5, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".coverage.Empty")
	empties.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
	name := field("name", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, "")
	name.Options = &descriptorpb.FieldOptions{Deprecated: proto.Bool(true)}
	file, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name: proto.String("coverage.proto"), Package: proto.String("coverage"), Syntax: proto.String("proto3"),
		EnumType: []*descriptorpb.EnumDescriptorProto{
			{Name: proto.String("Status"), Value: []*descriptorpb.EnumValueDescriptorProto{
				{Name: proto.String("UNKNOWN"), Number: proto.Int32(0)},
				{Name: proto.String("ACTIVE"), Number: proto.Int32(1), Options: &descriptorpb.EnumValueOptions{Deprecated: proto.Bool(true)}},
			}},
			{Name: proto.String("Only"), Value: []*descriptorpb.EnumValueDescriptorProto{{Name: proto.String("ONLY"), Number: proto.Int32(0)}}},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: proto.String("Empty")},
			{Name: proto.String("Child"), Field: []*descriptorpb.FieldDescriptorProto{name}},
			{Name: proto.String("Root"), Field: []*descriptorpb.FieldDescriptorProto{labels, enums,
				field("child", 3, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".coverage.Child"),
				field("only", 4, descriptorpb.FieldDescriptorProto_TYPE_ENUM, ".coverage.Only"), empties},
				NestedType: []*descriptorpb.DescriptorProto{{Name: proto.String("LabelsEntry"), Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(true)},
					Field: []*descriptorpb.FieldDescriptorProto{field("key", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""), field("value", 2, descriptorpb.FieldDescriptorProto_TYPE_STRING, "")}}}},
		},
	}, nil)
	require.NoError(t, err)
	return file.Messages().ByName("Root")
}

func TestNestedMessageDefaults(t *testing.T) {
	descriptor := coverageDescriptor(t)
	childField := descriptor.Fields().ByName("child")
	root := dynamicpb.NewMessage(descriptor)
	result := buildFieldResult(childField, root.Get(childField))
	assert.False(t, result.Covered)
	assert.Empty(t, result.Parameters)
	child := dynamicpb.NewMessage(childField.Message())
	root.Set(childField, protoreflect.ValueOfMessage(child))
	result = buildFieldResult(childField, root.Get(childField))
	assert.False(t, result.Covered)
	child.Set(child.Descriptor().Fields().ByName("name"), protoreflect.ValueOfString("api"))
	result = buildFieldResult(childField, root.Get(childField))
	assert.True(t, result.Covered)
	require.Len(t, result.Parameters, 1)
	assert.True(t, result.Parameters[0].Covered)
	assert.True(t, result.Parameters[0].Deprecated)
}

func TestPrimitiveMapCoverage(t *testing.T) {
	root := dynamicpb.NewMessage(coverageDescriptor(t))
	field := root.Descriptor().Fields().ByName("labels")
	assert.False(t, buildFieldResult(field, root.Get(field)).Covered)
	root.Mutable(field).Map().Set(protoreflect.ValueOfString("key").MapKey(), protoreflect.ValueOfString("value"))
	result := buildFieldResult(field, root.Get(field))
	assert.True(t, result.Covered)
	assert.Empty(t, result.Parameters)
}

func TestRepeatedEnumCoverageMergesValues(t *testing.T) {
	root := dynamicpb.NewMessage(coverageDescriptor(t))
	field := root.Descriptor().Fields().ByName("statuses")
	assert.False(t, buildFieldResult(field, root.Get(field)).Covered)
	values := root.Mutable(field).List()
	values.Append(protoreflect.ValueOfEnum(0))
	values.Append(protoreflect.ValueOfEnum(1))
	values.Append(protoreflect.ValueOfEnum(1))
	result := buildFieldResult(field, root.Get(field))
	assert.True(t, result.Covered)
	require.Len(t, result.Parameters, 2)
	for _, value := range result.Parameters {
		assert.True(t, value.Covered)
		assert.Equal(t, value.Parameter == "ACTIVE", value.Deprecated)
	}
	field = root.Descriptor().Fields().ByName("only")
	result = buildFieldResult(field, root.Get(field))
	assert.True(t, result.Covered)
	require.Len(t, result.Parameters, 1)
	assert.True(t, result.Parameters[0].Covered)
}

func TestRepeatedEmptyMessagesAreCovered(t *testing.T) {
	root := dynamicpb.NewMessage(coverageDescriptor(t))
	field := root.Descriptor().Fields().ByName("empties")
	assert.False(t, buildFieldResult(field, root.Get(field)).Covered)
	values := root.Mutable(field).List()
	values.Append(protoreflect.ValueOfMessage(dynamicpb.NewMessage(field.Message())))
	result := buildFieldResult(field, root.Get(field))
	assert.True(t, result.Covered)
	assert.Empty(t, result.Parameters)
}
