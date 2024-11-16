package coverage

import (
	"sort"

	"github.com/samber/lo"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

func SortResultParameters(params []models.ResultParameters) {
	sort.Slice(params, func(i, j int) bool {
		return params[i].Parameter < params[j].Parameter
	})

	for index := range params {
		if params[index].Parameters != nil {
			SortResultParameters(params[index].Parameters)
		}
	}
}

func GetTotalResultParameters(results []models.ResultParameters) int {
	if len(results) == 0 {
		return 0
	}

	count := len(results)
	for _, param := range results {
		count += GetTotalResultParameters(param.Parameters)
	}

	return count
}

func GetTotalCoveredResultParameters(results []models.ResultParameters) int {
	if len(results) == 0 {
		return 0
	}

	count := lo.SumBy(results, func(item models.ResultParameters) int {
		return lo.Ternary(item.Covered, 1, 0)
	})
	for _, param := range results {
		count += GetTotalCoveredResultParameters(param.Parameters)
	}

	return count
}
