package coverageinupt

import (
	"context"
	"log"

	"tests-coverage-tool/tool/utils"

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

		outputDir, err := getCoverageResultsDir()
		if err != nil {
			log.Printf("Error getting coverage results dir: %v", err)
			return invokerErr
		}

		if err = utils.SaveJSONFile(result, outputDir, uuid.New().String()); err != nil {
			log.Printf("Error saving coverage result: %v", err)
		}

		return invokerErr
	}
}
