package history

import (
	"log"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

type OutputHistoryClient struct {
	state  models.CoverageState
	config config.Config
}

func NewOutputHistoryClient(conf config.Config, state models.CoverageState) OutputHistoryClient {
	return OutputHistoryClient{config: conf, state: state}
}

func (c OutputHistoryClient) SaveHistory() error {
	log.Println("Starting to make history file")

	if c.config.HistoryDir == "" {
		logger.EnvVariableEmptySkipping(config.HistoryDir.String())
		return nil
	}

	if c.config.HistoryFile == "" {
		logger.EnvVariableEmptySkipping(config.HistoryFile.String())
		return nil
	}

	err := utils.SaveJSONFile(c.state.GetHistoryState(), c.config.HistoryDir, c.config.HistoryFile)
	if err != nil {
		log.Println("Error making history file")
		return err
	}

	log.Println("Successfully made history file")
	return nil
}
