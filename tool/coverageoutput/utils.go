package coverageoutput

import (
	"github.com/samber/lo"
)

func getCoveragePercent(original, actual []string) float64 {
	left, _ := lo.Difference(original, actual)

	totalCovered := len(original) - len(left)

	coveragePercent := (float64(totalCovered) / float64(len(original))) * 100
	return coveragePercent
}
