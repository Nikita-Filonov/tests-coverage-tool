package reflection

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func buildEnumExpectedResultParameters(descriptor protoreflect.FieldDescriptor) []models.ResultParameters {
	enumDescriptor := descriptor.Enum()

	var results []models.ResultParameters
	for index := 0; index < enumDescriptor.Values().Len(); index++ {
		value := enumDescriptor.Values().Get(index)

		result := models.ResultParameters{Parameter: string(value.Name())}
		if options, ok := value.Options().(*descriptorpb.EnumValueOptions); ok {
			result.Deprecated = options.GetDeprecated()
		}

		results = append(results, result)
	}

	return results
}

func buildExpectedResultParameters(descriptor protoreflect.MessageDescriptor, visited map[protoreflect.FullName]bool) []models.ResultParameters {
	if visited[descriptor.FullName()] {
		return nil
	}
	visited[descriptor.FullName()] = true

	var results []models.ResultParameters
	for index := 0; index < descriptor.Fields().Len(); index++ {
		field := descriptor.Fields().Get(index)

		result := models.ResultParameters{Parameter: string(field.Name())}
		if options, ok := field.Options().(*descriptorpb.FieldOptions); ok {
			result.Deprecated = options.GetDeprecated()
		}

		switch field.Kind() {
		case protoreflect.MessageKind:
			if field.IsMap() && field.MapValue().Kind() == protoreflect.MessageKind {
				result.Parameters = buildExpectedResultParameters(field.MapValue().Message(), visited)
			} else {
				result.Parameters = buildExpectedResultParameters(field.Message(), visited)
			}
		case protoreflect.EnumKind:
			result.Parameters = buildEnumExpectedResultParameters(field)
		}

		results = append(results, result)
	}

	delete(visited, descriptor.FullName())
	return results
}

func BuildExpectedResultParameters(descriptor protoreflect.MessageDescriptor) []models.ResultParameters {
	return buildExpectedResultParameters(descriptor, make(map[protoreflect.FullName]bool))
}
