package coverage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type enrichWithUncoveredResultParametersTest[T any] struct {
	name  string
	want  T
	input T
}

var wantWithUncoveredChildrenParameters = models.ResultParameters{
	Covered:   true,
	Parameter: "a",
	Parameters: []models.ResultParameters{
		{
			Covered:   true,
			Parameter: "b",
			Parameters: []models.ResultParameters{
				{Covered: true, Parameter: "e"},
				{Covered: false, Parameter: "d"},
			},
			HasUncoveredParameters: true,
		},
		{Parameter: "c", Covered: false},
	},
	HasUncoveredParameters: true,
}
var inputWithUncoveredChildrenParameters = models.ResultParameters{
	Covered:   true,
	Parameter: "a",
	Parameters: []models.ResultParameters{
		{
			Covered:   true,
			Parameter: "b",
			Parameters: []models.ResultParameters{
				{Covered: true, Parameter: "e"},
				{Covered: false, Parameter: "d"},
			},
		},
		{Parameter: "c", Covered: false},
	},
}

var wantWithoutUncoveredChildrenParameters = models.ResultParameters{
	Covered:   true,
	Parameter: "a",
	Parameters: []models.ResultParameters{
		{Covered: true, Parameter: "b"},
		{Covered: true, Parameter: "c"},
	},
}
var inputWithoutUncoveredChildrenParameters = models.ResultParameters{
	Covered:   true,
	Parameter: "a",
	Parameters: []models.ResultParameters{
		{Covered: true, Parameter: "b"},
		{Covered: true, Parameter: "c"},
	},
}

func TestEnrichWithUncoveredResultParameters(t *testing.T) {
	tests := []enrichWithUncoveredResultParametersTest[models.ResultParameters]{
		{
			name:  "With uncovered children parameters",
			want:  wantWithUncoveredChildrenParameters,
			input: inputWithUncoveredChildrenParameters,
		},
		{
			name:  "Without uncovered children parameters",
			want:  wantWithoutUncoveredChildrenParameters,
			input: inputWithoutUncoveredChildrenParameters,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enrichWithUncoveredResultParameters(&test.input)
			assert.Equal(t, test.want, test.input)
		})
	}
}

func TestEnrichSliceWithUncoveredResultParameters(t *testing.T) {
	tests := []enrichWithUncoveredResultParametersTest[[]models.ResultParameters]{
		{
			name:  "With uncovered children parameters",
			want:  []models.ResultParameters{wantWithUncoveredChildrenParameters},
			input: []models.ResultParameters{inputWithUncoveredChildrenParameters},
		},
		{
			name:  "Without uncovered children parameters",
			want:  []models.ResultParameters{wantWithoutUncoveredChildrenParameters},
			input: []models.ResultParameters{inputWithoutUncoveredChildrenParameters},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			EnrichSliceWithUncoveredResultParameters(test.input)
			assert.Equal(t, test.want, test.input)
		})
	}
}
