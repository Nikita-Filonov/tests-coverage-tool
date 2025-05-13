package coverage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type getTotalCoverageResultParametersTest struct {
	name       string
	want       int
	parameters []models.ResultParameters
}

func TestGetTotalResultParameters(t *testing.T) {
	tests := []getTotalCoverageResultParametersTest{
		{
			name:       "Parameters with one level depth",
			want:       3,
			parameters: []models.ResultParameters{{}, {}, {}},
		},
		{
			name: "Parameters with multiple level depth",
			want: 7,
			parameters: []models.ResultParameters{
				{Parameters: []models.ResultParameters{}},
				{Parameters: []models.ResultParameters{}},
				{
					Parameters: []models.ResultParameters{
						{Parameters: []models.ResultParameters{}},
						{Parameters: []models.ResultParameters{{}, {}}},
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
			parameters: []models.ResultParameters{{Covered: true}, {}, {}},
		},
		{
			name: "Parameters with multiple level depth",
			want: 5,
			parameters: []models.ResultParameters{
				{Covered: true, Parameters: []models.ResultParameters{}},
				{Covered: true, Parameters: []models.ResultParameters{}},
				{
					Parameters: []models.ResultParameters{
						{Covered: true, Parameters: []models.ResultParameters{}},
						{Parameters: []models.ResultParameters{}},
						{Parameters: []models.ResultParameters{{Covered: true}, {Covered: true}}},
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
