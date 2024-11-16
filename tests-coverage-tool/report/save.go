package report

import (
	"context"

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
	ctx := context.Background()

	toolConfig, err := config.NewConfig()
	if err != nil {
		logger.FatalBuildingNewClient("config", err)
	}

	inputCoverageClient, err := coverageinupt.NewInputCoverageClient(toolConfig.GetResultsDir())
	if err != nil {
		logger.FatalBuildingNewClient("input coverage client", err)
	}

	inputHistoryClientFactory, err := history.NewInputHistoryClientFactory(toolConfig)
	if err != nil {
		logger.FatalBuildingNewClient("input history client factory", err)
	}

	coverageState := models.NewCoverageState(toolConfig)
	for _, service := range toolConfig.Services {
		reflectionClient, err := reflection.NewGRPCReflectionClient(ctx, service)
		if err != nil {
			logger.FatalBuildingNewClient("grpc reflection client", err)
		}

		inputHistoryClient := inputHistoryClientFactory.NewClient(service.Key)

		outputCoverageClient, err := coverageoutput.NewOutputCoverageClient(
			reflectionClient, inputHistoryClient, inputCoverageClient,
		)
		if err != nil {
			logger.FatalBuildingNewClient("output coverage client", err)
		}

		serviceCoverage, err := outputCoverageClient.GetServiceCoverage()
		if err != nil {
			logger.FatalGettingEntity("service coverage", err)
		}

		logicalServiceCoverage, err := outputCoverageClient.GetLogicalServiceCoverages()
		if err != nil {
			logger.FatalGettingEntity("logical service coverages", err)
		}

		coverageState.ServiceCoverages[service.Key] = serviceCoverage
		coverageState.LogicalServiceCoverages[service.Key] = logicalServiceCoverage
	}

	outputHistoryClient := history.NewOutputHistoryClient(toolConfig, coverageState)
	if err = outputHistoryClient.SaveHistory(); err != nil {
		logger.FatalBuildingNewClient("output history client", err)
	}

	coverageReportClient := report.NewCoverageReportClient(toolConfig, coverageState)

	if err = coverageReportClient.SaveHTMLReport(); err != nil {
		logger.FatalSavingReport("HTML", err)
	}

	if err = coverageReportClient.SaveJSONReport(); err != nil {
		logger.FatalSavingReport("JSON", err)
	}
}

func NewSaveReportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "save-report",
		Short: "Saves a report",
		Run:   func(_ *cobra.Command, _ []string) { SaveReport() },
	}
}
