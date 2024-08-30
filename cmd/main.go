package main

import (
	"context"
	"log"

	"tests-coverage-tool/tool/config"
	"tests-coverage-tool/tool/coverageinupt"
	"tests-coverage-tool/tool/coverageoutput"
	"tests-coverage-tool/tool/reflection"
	"tests-coverage-tool/tool/report"
)

func main() {
	ctx := context.Background()

	toolConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error building config: %v", err)
	}

	reflectionClient, err := reflection.NewGRPCReflectionClient(ctx, toolConfig)
	if err != nil {
		log.Fatalf("Error building grpc reflection client: %v", err)
	}

	inputCoverageClient, err := coverageinupt.NewInputCoverageClient(toolConfig.InputResultsDir)
	if err != nil {
		log.Fatalf("Error building input coverage client: %v", err)
	}

	outputCoverageClient, err := coverageoutput.NewOutputCoverageClient(reflectionClient, inputCoverageClient)
	if err != nil {
		log.Fatalf("Error building output coverage client: %v", err)
	}

	coverages, err := outputCoverageClient.SaveResults(toolConfig.OutputResultsDir)
	if err != nil {
		log.Fatalf("Error saving output coverage results: %v", err)
	}

	if err = toolConfig.SaveResults(); err != nil {
		log.Fatalf("Error saving config: %v", err)
	}

	if err = report.SaveHTMLReport(toolConfig, coverages); err != nil {
		log.Fatalf("Error saving html report: %v", err)
	}
}
