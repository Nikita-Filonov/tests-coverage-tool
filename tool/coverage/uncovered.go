package coverage

import (
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func enrichWithUncoveredResultParameters(parameter *models.ResultParameters) bool {
	hasUncoveredChild := false
	for index := range parameter.Parameters {
		if childUncovered := enrichWithUncoveredResultParameters(&parameter.Parameters[index]); childUncovered {
			hasUncoveredChild = true
		}
	}

	if parameter.Covered && hasUncoveredChild {
		parameter.HasUncoveredParameters = true
	}

	return hasUncoveredChild || !parameter.Covered
}

func EnrichSliceWithUncoveredResultParameters(parameters []models.ResultParameters) {
	for index := range parameters {
		enrichWithUncoveredResultParameters(&parameters[index])
	}
}
