package coverage

import (
	"github.com/samber/lo"
)

func MergeResultParameters(master, replica []ResultParameters) []ResultParameters {
	mergedMap := lo.SliceToMap(master, func(item ResultParameters) (string, ResultParameters) { return item.Parameter, item })

	for _, param := range replica {
		if existingParam, exists := mergedMap[param.Parameter]; exists {
			mergedMap[param.Parameter] = ResultParameters{
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

func MergeFilteredResultParameters(results [][]ResultParameters) []ResultParameters {
	if len(results) == 0 {
		return nil
	}

	mergedResult := results[0]
	for _, result := range results[1:] {
		mergedResult = MergeResultParameters(mergedResult, result)
	}
	return mergedResult
}
