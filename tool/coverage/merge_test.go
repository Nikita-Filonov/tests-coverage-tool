package coverage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type mergeResultParametersTest struct {
	name    string
	want    []models.ResultParameters
	master  []models.ResultParameters
	replica []models.ResultParameters
}

type mergeFilteredResultParametersTest struct {
	name       string
	want       []models.ResultParameters
	parameters [][]models.ResultParameters
}

var replicaResultParameters = []models.ResultParameters{
	{Covered: true, Parameter: "a", Deprecated: true},
	{Covered: true, Parameter: "b", Deprecated: false},
	{
		Covered:    false,
		Parameter:  "c",
		Deprecated: true,
		Parameters: []models.ResultParameters{
			{Covered: false, Parameter: "c.a", Deprecated: true},
			{Covered: false, Parameter: "c.b", Deprecated: false},
			{Covered: true, Parameter: "c.c", Deprecated: false},
		},
	},
}

var masterResultParameters = []models.ResultParameters{
	{Covered: false, Parameter: "a", Deprecated: true},
	{Covered: false, Parameter: "b", Deprecated: false},
	{
		Covered:    true,
		Parameter:  "c",
		Deprecated: true,
		Parameters: []models.ResultParameters{
			{Covered: true, Parameter: "c.a", Deprecated: true},
			{Covered: true, Parameter: "c.b", Deprecated: false},
			{Covered: false, Parameter: "c.c", Deprecated: false},
		},
	},
}

func TestMergeResultParameters(t *testing.T) {
	tests := []mergeResultParametersTest{
		{
			name:    "Empty master and replica",
			want:    []models.ResultParameters{},
			master:  []models.ResultParameters{},
			replica: []models.ResultParameters{},
		},
		{
			name:    "Empty master, set replica",
			want:    replicaResultParameters,
			master:  []models.ResultParameters{},
			replica: replicaResultParameters,
		},
		{
			name:    "Set master, empty replica",
			want:    masterResultParameters,
			master:  masterResultParameters,
			replica: []models.ResultParameters{},
		},
		{
			name: "Set master, set replica",
			want: []models.ResultParameters{
				{Covered: true, Parameter: "a", Deprecated: true, Parameters: []models.ResultParameters{}},
				{Covered: true, Parameter: "b", Deprecated: false, Parameters: []models.ResultParameters{}},
				{
					Covered:    true,
					Parameter:  "c",
					Deprecated: true,
					Parameters: []models.ResultParameters{
						{Covered: true, Parameter: "c.a", Deprecated: true, Parameters: []models.ResultParameters{}},
						{Covered: true, Parameter: "c.b", Deprecated: false, Parameters: []models.ResultParameters{}},
						{Covered: true, Parameter: "c.c", Deprecated: false, Parameters: []models.ResultParameters{}},
					},
				},
			},
			master:  masterResultParameters,
			replica: replicaResultParameters,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MergeResultParameters(test.master, test.replica)
			SortResultParameters(result)
			SortResultParameters(test.want)

			assert.Equal(t, test.want, result)
		})
	}
}

func TestMergeFilteredResultParameters(t *testing.T) {
	tests := []mergeFilteredResultParametersTest{
		{
			name: "Parameters with multiple level depth",
			want: []models.ResultParameters{
				{
					Covered:   true,
					Parameter: "a",
					Parameters: []models.ResultParameters{
						{Covered: true, Parameter: "d", Parameters: []models.ResultParameters{}},
						{Covered: true, Parameter: "e", Parameters: []models.ResultParameters{}},
					},
				},
				{Covered: true, Parameter: "b", Parameters: []models.ResultParameters{}},
				{Covered: true, Parameter: "c", Parameters: []models.ResultParameters{}},
			},
			parameters: [][]models.ResultParameters{
				{
					{
						Covered:   true,
						Parameter: "a",
						Parameters: []models.ResultParameters{
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
						Parameters: []models.ResultParameters{
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
