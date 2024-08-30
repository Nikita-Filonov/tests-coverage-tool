package coverageinupt

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"tests-coverage-tool/tool/utils"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const coverageResultsDir = "TESTS_COVERAGE_INPUT_RESULTS_DIR"

func getCoverageResultsDir() (string, error) {
	dir := os.Getenv(coverageResultsDir)
	if dir == "" {
		return "", fmt.Errorf("environment variable %s cannot be empty\n", coverageResultsDir)
	}

	return dir, nil
}

func getDescriptorFromProtoMessage(input interface{}) (protoreflect.MessageDescriptor, proto.Message) {
	if message, ok := input.(proto.Message); ok {
		reflectMessage := message.ProtoReflect()
		return reflectMessage.Descriptor(), message
	}

	return nil, nil
}

func getOnlyFilledProtoMessageFields(message proto.Message) []string {
	v := reflect.ValueOf(message)
	v = reflect.Indirect(v)
	t := v.Type()

	var filledFields []string
	for index := 0; index < v.NumField(); index++ {
		field := v.Field(index)
		fieldType := t.Field(index)

		if !fieldType.IsExported() {
			continue
		}

		if !field.IsZero() {
			filledFields = append(filledFields, utils.PascalCaseToSnakeCase(fieldType.Name))
		}
	}

	return filledFields
}

func buildCoverageResult(fullMethod string, req, reply interface{}) (InputCoverageResult, error) {
	methodParts := strings.Split(fullMethod, "/")
	method, service := methodParts[2], methodParts[1]

	requestDescriptor, requestMessage := getDescriptorFromProtoMessage(req)
	if requestDescriptor == nil || requestMessage == nil {
		return InputCoverageResult{}, fmt.Errorf("unable to build coverage result struct because of malformed request")
	}

	responseDescriptor, responseMessage := getDescriptorFromProtoMessage(reply)
	if responseDescriptor == nil || responseMessage == nil {
		return InputCoverageResult{}, fmt.Errorf("unable to build coverage result struct because of malformed response")
	}

	return InputCoverageResult{
		Method:             method,
		RequestName:        string(requestDescriptor.FullName()),
		ResponseName:       string(responseDescriptor.FullName()),
		LogicalService:     service,
		RequestParameters:  getOnlyFilledProtoMessageFields(requestMessage),
		ResponseParameters: getOnlyFilledProtoMessageFields(responseMessage),
	}, nil
}
