package coverageinupt

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func isDefaultValue(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
	if field.IsList() {
		return value.List().Len() == 0
	}

	if field.IsMap() {
		return value.Map().Len() == 0
	}

	switch field.Kind() {
	case protoreflect.BoolKind:
		return !value.Bool()
	case protoreflect.EnumKind:
		return value.Enum() == field.Enum().Values().Get(0).Number()
	case protoreflect.FloatKind:
		return value.Float() == 0.0
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int() == 0
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint() == 0
	case protoreflect.StringKind:
		return value.String() == ""
	case protoreflect.BytesKind:
		return len(value.Bytes()) == 0
	case protoreflect.MessageKind:
		return !value.Message().IsValid() || isMessageDefault(value.Message())
	default:
		return false
	}
}

func isMessageDefault(message protoreflect.Message) bool {
	fields := message.Descriptor().Fields()
	for index := 0; index < fields.Len(); index++ {
		field := fields.Get(index)
		fieldValue := message.Get(field)

		if !isDefaultValue(field, fieldValue) {
			return false
		}
	}
	return true
}

func buildMapActualResultParameters(field protoreflect.FieldDescriptor, value protoreflect.Value) models.ResultParameters {
	var subResults []models.ResultParameters
	value.Map().Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		subMessage := v.Message()
		subResults = append(subResults, buildActualResultParameters(subMessage.Interface())...)
		return true
	})

	return models.ResultParameters{
		Covered:    !isDefaultValue(field, value),
		Parameter:  string(field.Name()),
		Parameters: subResults,
	}
}

func buildArrayActualResultParameters(field protoreflect.FieldDescriptor, value protoreflect.Value) models.ResultParameters {
	mergedSubResults := map[string]models.ResultParameters{}
	for index := 0; index < value.List().Len(); index++ {
		subResults := buildActualResultParameters(value.List().Get(index).Message().Interface())

		for _, subResult := range subResults {
			if existingResult, exists := mergedSubResults[subResult.Parameter]; exists {
				mergedSubResults[subResult.Parameter] = models.ResultParameters{
					Covered:    existingResult.Covered || subResult.Covered,
					Parameter:  subResult.Parameter,
					Parameters: coverage.MergeResultParameters(existingResult.Parameters, subResult.Parameters),
					Deprecated: subResult.Deprecated,
				}
			} else {
				mergedSubResults[subResult.Parameter] = subResult
			}
		}
	}

	finalSubResults := lo.Values(mergedSubResults)

	return models.ResultParameters{
		Covered:    len(finalSubResults) > 0,
		Parameter:  string(field.Name()),
		Parameters: finalSubResults,
	}
}

func getEnumParameters(value protoreflect.EnumNumber, descriptor protoreflect.EnumDescriptor) []models.ResultParameters {
	var enumResults []models.ResultParameters

	for index := 0; index < descriptor.Values().Len(); index++ {
		enumValue := descriptor.Values().Get(index)
		enumCovered := enumValue.Number() == value
		enumOptions := enumValue.Options().(*descriptorpb.EnumValueOptions)

		enumResults = append(enumResults, models.ResultParameters{
			Covered:    enumCovered,
			Parameter:  string(enumValue.Name()),
			Deprecated: enumOptions.GetDeprecated(),
		})
	}

	return enumResults
}

func buildEnumActualResultParameters(field protoreflect.FieldDescriptor, value protoreflect.Value) models.ResultParameters {
	enumDescriptor := field.Enum()

	var enumResults []models.ResultParameters
	if field.IsList() {
		for index := 0; index < value.List().Len(); index++ {
			enumResults = append(enumResults, getEnumParameters(value.List().Get(index).Enum(), enumDescriptor)...)
		}
	} else {
		enumResults = append(enumResults, getEnumParameters(value.Enum(), enumDescriptor)...)
	}

	covered := !isDefaultValue(field, value)
	if len(enumResults) == 1 {
		covered = enumResults[0].Covered
	}

	return models.ResultParameters{
		Covered:    covered,
		Parameter:  string(field.Name()),
		Parameters: enumResults,
	}
}

func buildFieldResult(field protoreflect.FieldDescriptor, value protoreflect.Value) models.ResultParameters {
	switch {
	case field.Kind() == protoreflect.MessageKind && field.IsList():
		return buildArrayActualResultParameters(field, value)
	case field.Kind() == protoreflect.MessageKind && field.IsMap():
		return buildMapActualResultParameters(field, value)
	case field.Kind() == protoreflect.EnumKind:
		return buildEnumActualResultParameters(field, value)
	case field.Kind() == protoreflect.MessageKind:
		return models.ResultParameters{
			Covered:    !isDefaultValue(field, value),
			Parameter:  string(field.Name()),
			Parameters: buildActualResultParameters(value.Message().Interface()),
		}
	default:
		return models.ResultParameters{Covered: !isDefaultValue(field, value), Parameter: string(field.Name())}
	}
}

func buildActualResultParameters(message proto.Message) []models.ResultParameters {
	var results []models.ResultParameters
	fields := message.ProtoReflect().Descriptor().Fields()

	for index := 0; index < fields.Len(); index++ {
		field := fields.Get(index)
		fieldValue := message.ProtoReflect().Get(field)
		fieldOptions := field.Options().(*descriptorpb.FieldOptions)

		result := buildFieldResult(field, fieldValue)
		result.Deprecated = fieldOptions.GetDeprecated()
		results = append(results, result)
	}

	return results
}
