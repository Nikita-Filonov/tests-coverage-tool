package coverageinupt

import (
	"context"
	"fmt"
	"log"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func CoverageInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		invokerErr := invoker(ctx, method, req, reply, cc, opts...)

		result, err := buildCoverageResult(method, req, reply)
		if err != nil {
			log.Printf("Error building coverage result: %v", err)
			return invokerErr
		}

		toolConfig, err := config.NewConfig()
		if err != nil {
			log.Printf("Error building config: %v", err)
			return invokerErr
		}

		filename := fmt.Sprintf("%s.json", uuid.New().String())
		resultsDir := toolConfig.GetResultsDir()
		if err = utils.SaveJSONFile(result, resultsDir, filename); err != nil {
			log.Printf("Error saving coverage result: %v", err)
		}

		return invokerErr
	}
}
