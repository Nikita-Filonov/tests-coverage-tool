package report

import (
	"log"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

func ReadCoverageReportState() (*CoverageReportState, error) {
	toolConfig, err := config.NewConfig()
	if err != nil {
		log.Printf("Error building config: %v", err)
		return nil, err
	}

	if toolConfig.JSONReportDir == "" {
		logger.EnvVariableEmptySkipping(config.JSONReportDir.String())
		return nil, nil
	}

	if toolConfig.JSONReportFile == "" {
		logger.EnvVariableEmptySkipping(config.JSONReportFile.String())
		return nil, nil
	}

	return utils.ReadJSONFile[CoverageReportState](toolConfig.GetJSONReportFile())
}
