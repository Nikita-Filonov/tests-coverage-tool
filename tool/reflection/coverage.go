package reflection

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func buildEnumExpectedResultParameters(descriptor protoreflect.FieldDescriptor) []models.ResultParameters {
	enumDescriptor := descriptor.Enum()
	var enumResults []models.ResultParameters
	for index := 0; index < enumDescriptor.Values().Len(); index++ {
		enumValue := enumDescriptor.Values().Get(index)
		enumOptions := enumValue.Options().(*descriptorpb.EnumValueOptions)

		enumResults = append(enumResults, models.ResultParameters{
			Parameter:  string(enumValue.Name()),
			Deprecated: enumOptions.GetDeprecated(),
		})
	}

	return enumResults
}

func BuildExpectedResultParameters(descriptor protoreflect.MessageDescriptor) []models.ResultParameters {
	var results []models.ResultParameters

	for index := 0; index < descriptor.Fields().Len(); index++ {
		field := descriptor.Fields().Get(index)
		result := models.ResultParameters{
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
