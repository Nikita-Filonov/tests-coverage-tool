package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type resultTest[T any] struct {
	name  string
	want  T
	model Result
}

func TestResultGetLogicalService(t *testing.T) {
	tests := []resultTest[string]{
		{
			name:  "Service.Method",
			want:  "Service",
			model: Result{Method: "Service.Method"},
		},
		{
			name:  "Company.Team.Version.Service.Method",
			want:  "Company.Team.Version.Service",
			model: Result{Method: "Company.Team.Version.Service.Method"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.model.GetLogicalService())
		})
	}
}
