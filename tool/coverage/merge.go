package coverage

import (
	"github.com/samber/lo"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func MergeResultParameters(master, replica []models.ResultParameters) []models.ResultParameters {
	mergedMap := lo.SliceToMap(master, func(item models.ResultParameters) (string, models.ResultParameters) { return item.Parameter, item })

	for _, param := range replica {
		if existingParam, exists := mergedMap[param.Parameter]; exists {
			mergedMap[param.Parameter] = models.ResultParameters{
				Covered:    existingParam.Covered || param.Covered,
				Parameter:  param.Parameter,
				Parameters: MergeResultParameters(existingParam.Parameters, param.Parameters),
				Deprecated: param.Deprecated,
			}
		} else {
			mergedMap[param.Parameter] = param
		}
	}

	return lo.Values(mergedMap)
}

func MergeFilteredResultParameters(results [][]models.ResultParameters) []models.ResultParameters {
	if len(results) == 0 {
		return nil
	}

	mergedResult := results[0]
	for _, result := range results[1:] {
		mergedResult = MergeResultParameters(mergedResult, result)
	}
	return mergedResult
}
