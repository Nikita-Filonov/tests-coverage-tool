package history

import (
	"os"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

func ReadHistoryState(conf config.Config) (models.HistoryState, error) {
	if conf.HistoryDir == "" {
		logger.EnvVariableEmptySkipping(config.HistoryDir.String())
		return nil, nil
	}

	if conf.HistoryFile == "" {
		logger.EnvVariableEmptySkipping(config.HistoryFile.String())
		return nil, nil
	}

	state, err := utils.ReadJSONFile[models.HistoryState](conf.GetHistoryFile())
	if err != nil {
		if os.IsNotExist(err) {
			return models.HistoryState{}, nil
		}

		return nil, err
	}

	return *state, nil
}
