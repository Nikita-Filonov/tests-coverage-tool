package coverageoutput

import (
	"math"

	"github.com/samber/lo"
)

func getCoveragePercent(original, actual []string) float64 {
	left, _ := lo.Difference(original, actual)

	totalCovered := len(original) - len(left)

	percent := (float64(totalCovered) / float64(len(original))) * 100
	return getRoundedTotalCoverage(percent)
}

func getRequestCoveragePercent(covered bool, total, totalCovered int) float64 {
	if !covered {
		return 0
	}

	if covered && total == 0 {
		return 100
	}

	percent := (float64(totalCovered) / float64(total)) * 100
	return getRoundedTotalCoverage(lo.Ternary(percent > 100, 100, percent))
}

func getRoundedTotalCoverage(percent float64) float64 {
	return math.Round(percent*100) / 100
}
