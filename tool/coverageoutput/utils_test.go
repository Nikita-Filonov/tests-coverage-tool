package coverageoutput

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type getCoveragePercentTest struct {
	name     string
	actual   []string
	original []string
	result   float64
}

type getRequestCoveragePercentTest struct {
	name         string
	total        int
	totalCovered int
	covered      bool
	result       float64
}

func TestGetCoveragePercent(t *testing.T) {
	tests := []getCoveragePercentTest{
		{
			name:     "100% coverage",
			actual:   []string{"a", "b", "c"},
			original: []string{"a", "b", "c"},
			result:   100.00,
		},
		{
			name:     "67.67% coverage",
			actual:   []string{"a", "b"},
			original: []string{"a", "b", "c"},
			result:   66.67,
		},
		{
			name:     "0% coverage",
			actual:   []string{},
			original: []string{"a", "b", "c"},
			result:   00.00,
		},
		{
			name:     "More that 100% coverage",
			actual:   []string{"a", "b", "c", "d"},
			original: []string{"a", "b", "c"},
			result:   100.00,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.result, getCoveragePercent(test.original, test.actual))
		})
	}
}

func TestGetRequestCoveragePercent(t *testing.T) {
	tests := []getRequestCoveragePercentTest{
		{
			name:         "100% coverage",
			total:        3,
			totalCovered: 3,
			covered:      true,
			result:       100.00,
		},
		{
			name:         "66.67% coverage",
			total:        3,
			totalCovered: 2,
			covered:      true,
			result:       66.67,
		},
		{
			name:         "0% coverage",
			total:        3,
			totalCovered: 0,
			covered:      true,
			result:       00.00,
		},
		{
			name:         "More than 100% coverage",
			total:        3,
			totalCovered: 10,
			covered:      true,
			result:       100.00,
		},
		{
			name:         "Total is zero",
			total:        0,
			totalCovered: 0,
			covered:      true,
			result:       100.00,
		},
		{
			name:         "Not covered",
			total:        0,
			totalCovered: 0,
			covered:      false,
			result:       00.00,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.result, getRequestCoveragePercent(test.covered, test.total, test.totalCovered))
		})
	}
}
