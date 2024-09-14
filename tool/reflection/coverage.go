package reflection

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
)

func buildEnumExpectedResultParameters(descriptor protoreflect.FieldDescriptor) []coverage.ResultParameters {
	enumDescriptor := descriptor.Enum()
	var enumResults []coverage.ResultParameters
	for index := 0; index < enumDescriptor.Values().Len(); index++ {
		enumValue := enumDescriptor.Values().Get(index)
		enumOptions := enumValue.Options().(*descriptorpb.EnumValueOptions)

		enumResults = append(enumResults, coverage.ResultParameters{
			Parameter:  string(enumValue.Name()),
			Deprecated: enumOptions.GetDeprecated(),
		})
	}

	return enumResults
}

func BuildExpectedResultParameters(descriptor protoreflect.MessageDescriptor) []coverage.ResultParameters {
	var results []coverage.ResultParameters

	for index := 0; index < descriptor.Fields().Len(); index++ {
		field := descriptor.Fields().Get(index)
		result := coverage.ResultParameters{
			Parameter:  string(field.Name()),
			Deprecated: field.Options().(*descriptorpb.FieldOptions).GetDeprecated(),
		}

		switch field.Kind() {
		case protoreflect.MessageKind:
			if field.IsMap() {
				result.Parameters = BuildExpectedResultParameters(field.MapValue().Message())
			} else {
				result.Parameters = BuildExpectedResultParameters(field.Message())
			}
		case protoreflect.EnumKind:
			result.Parameters = buildEnumExpectedResultParameters(field)
		}

		results = append(results, result)
	}

	return results
}
