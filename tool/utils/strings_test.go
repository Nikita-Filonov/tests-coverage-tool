package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type pascalCaseToSnakeCaseTest struct {
	name  string
	input string
	want  string
}

func TestPascalCaseToSnakeCase(t *testing.T) {
	tests := []pascalCaseToSnakeCaseTest{
		{
			name:  "Two words pascale case string",
			input: "PascaleCase",
			want:  "pascale_case",
		},
		{
			name:  "One word pascale case string",
			input: "Pascale",
			want:  "pascale",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := PascalCaseToSnakeCase(test.input)

			assert.Equal(t, test.want, result)
		})
	}
}
