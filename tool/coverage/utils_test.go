package coverage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mergeFilteredResultParametersTest struct {
	name       string
	want       []ResultParameters
	parameters [][]ResultParameters
}

type getTotalCoverageResultParametersTest struct {
	name       string
	want       int
	parameters []ResultParameters
}

func TestMergeFilteredResultParameters(t *testing.T) {
	tests := []mergeFilteredResultParametersTest{
		{
			name: "Parameters with multiple level depth",
			want: []ResultParameters{
				{
					Covered:   true,
					Parameter: "a",
					Parameters: []ResultParameters{
						{Covered: true, Parameter: "d", Parameters: []ResultParameters{}},
						{Covered: true, Parameter: "e", Parameters: []ResultParameters{}},
					},
				},
				{Covered: true, Parameter: "b", Parameters: []ResultParameters{}},
				{Covered: true, Parameter: "c", Parameters: []ResultParameters{}},
			},
			parameters: [][]ResultParameters{
				{
					{
						Covered:   true,
						Parameter: "a",
						Parameters: []ResultParameters{
							{Covered: true, Parameter: "d"},
							{Covered: false, Parameter: "e"},
						},
					},
					{Covered: false, Parameter: "b"},
					{Covered: true, Parameter: "c"},
				},
				{
					{
						Covered:   false,
						Parameter: "a",
						Parameters: []ResultParameters{
							{Covered: false, Parameter: "d"},
							{Covered: true, Parameter: "e"},
						},
					},
					{Covered: true, Parameter: "b"},
					{Covered: true, Parameter: "c"},
				},
			},
		},
		{
			name: "Empty parameters",
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MergeFilteredResultParameters(test.parameters)
			SortResultParameters(result)
			SortResultParameters(test.want)

			assert.Equal(t, test.want, result)
		})
	}
}

func TestGetTotalResultParameters(t *testing.T) {
	tests := []getTotalCoverageResultParametersTest{
		{
			name:       "Parameters with one level depth",
			want:       3,
			parameters: []ResultParameters{{}, {}, {}},
		},
		{
			name: "Parameters with multiple level depth",
			want: 7,
			parameters: []ResultParameters{
				{Parameters: []ResultParameters{}},
				{Parameters: []ResultParameters{}},
				{
					Parameters: []ResultParameters{
						{Parameters: []ResultParameters{}},
						{Parameters: []ResultParameters{{}, {}}},
					},
				},
			},
		},
		{
			name: "Empty parameters",
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, GetTotalResultParameters(test.parameters))
		})
	}
}

func TestGetTotalCoveredResultParameters(t *testing.T) {
	tests := []getTotalCoverageResultParametersTest{
		{
			name:       "Parameters with one level depth",
			want:       1,
			parameters: []ResultParameters{{Covered: true}, {}, {}},
		},
		{
			name: "Parameters with multiple level depth",
			want: 5,
			parameters: []ResultParameters{
				{Covered: true, Parameters: []ResultParameters{}},
				{Covered: true, Parameters: []ResultParameters{}},
				{
					Parameters: []ResultParameters{
						{Covered: true, Parameters: []ResultParameters{}},
						{Parameters: []ResultParameters{}},
						{Parameters: []ResultParameters{{Covered: true}, {Covered: true}}},
					},
				},
			},
		},
		{
			name: "Empty parameters",
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, GetTotalCoveredResultParameters(test.parameters))
		})
	}
}
