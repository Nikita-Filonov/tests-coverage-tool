package coverageinupt

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func getProtoMessage(input interface{}) proto.Message {
	if message, ok := input.(proto.Message); ok {
		return message
	}

	return nil
}

func buildCoverageResult(method string, req, reply interface{}) (models.Result, error) {
	requestMessage := getProtoMessage(req)
	if requestMessage == nil {
		return models.Result{}, fmt.Errorf("unable to build coverage result struct because of malformed request")
	}

	responseMessage := getProtoMessage(reply)
	if responseMessage == nil {
		return models.Result{}, fmt.Errorf("unable to build coverage result struct because of malformed response")
	}

	return models.Result{
		Method:   strings.ReplaceAll(strings.TrimPrefix(method, "/"), "/", "."),
		Request:  buildActualResultParameters(requestMessage),
		Response: buildActualResultParameters(responseMessage),
	}, nil
}
