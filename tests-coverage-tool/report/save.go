package report

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageinupt"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageoutput"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/history"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/reflection"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/report"
)

func SaveReport() {
	if err := saveReport(context.Background()); err != nil {
		logger.FatalSavingReport("coverage", err)
	}
}

func saveReport(ctx context.Context) error {
	toolConfig, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	inputCoverageClient, err := coverageinupt.NewInputCoverageClient(toolConfig.GetResultsDir())
	if err != nil {
		return fmt.Errorf("failed to build input coverage client: %w", err)
	}

	inputHistoryClientFactory, err := history.NewInputHistoryClientFactory(toolConfig)
	if err != nil {
		return fmt.Errorf("failed to build input history client factory: %w", err)
	}

	coverageState := models.NewCoverageState(toolConfig)
	for _, service := range toolConfig.Services {
		reflectionClient, err := reflection.NewGRPCReflectionClient(ctx, service)
		if err != nil {
			return fmt.Errorf("failed to build grpc reflection client: %w", err)
		}
		defer func() { _ = reflectionClient.Close() }()

		inputHistoryClient := inputHistoryClientFactory.NewClient(service.Key)

		outputCoverageClient, err := coverageoutput.NewOutputCoverageClient(
			reflectionClient, inputHistoryClient, inputCoverageClient,
		)
		if err != nil {
			return fmt.Errorf("failed to build output coverage client: %w", err)
		}

		serviceCoverage, err := outputCoverageClient.GetServiceCoverage()
		if err != nil {
			return fmt.Errorf("failed to get service coverage: %w", err)
		}

		logicalServiceCoverage, err := outputCoverageClient.GetLogicalServiceCoverages()
		if err != nil {
			return fmt.Errorf("failed to get logical service coverages: %w", err)
		}

		coverageState.ServiceCoverages[service.Key] = serviceCoverage
		coverageState.LogicalServiceCoverages[service.Key] = logicalServiceCoverage
	}

	outputHistoryClient := history.NewOutputHistoryClient(toolConfig, coverageState)
	if err = outputHistoryClient.SaveHistory(); err != nil {
		return fmt.Errorf("failed to save history: %w", err)
	}

	coverageReportClient := report.NewCoverageReportClient(toolConfig, coverageState)

	if err = coverageReportClient.SaveHTMLReport(); err != nil {
		return fmt.Errorf("failed to save HTML report: %w", err)
	}

	if err = coverageReportClient.SaveJSONReport(); err != nil {
		return fmt.Errorf("failed to save JSON report: %w", err)
	}
	return nil
}

func NewSaveReportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "save-report",
		Short: "Saves a report",
		RunE:  func(cmd *cobra.Command, _ []string) error { return saveReport(cmd.Context()) },
	}
}
