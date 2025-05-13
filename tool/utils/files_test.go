package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type readFileTest struct {
	name     string
	want     []byte
	filename string
}

type saveFileTest struct {
	dir      string
	name     string
	input    []byte
	filename string
}

func TestReadFile(t *testing.T) {
	tests := []readFileTest{
		{
			name:     "File exist",
			want:     []byte("default\n"),
			filename: "../../testdata/default.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content, err := ReadFile(test.filename)

			assert.NoError(t, err)
			assert.Equal(t, test.want, content)
		})
	}
}

func TestReadFileNegative(t *testing.T) {
	tests := []readFileTest{
		{
			name:     "File does not exist",
			filename: "error.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content, err := ReadFile(test.filename)

			assert.Error(t, err)
			assert.Equal(t, test.want, content)
		})
	}
}

func TestSaveFile(t *testing.T) {
	tests := []saveFileTest{
		{
			dir:      ".",
			name:     "Input: default, dir: ., filename: default.txt",
			input:    []byte("default"),
			filename: "default.txt",
		},
		{
			dir:      ".",
			name:     "Input: {\"key\": \"value\"}, dir: ./json, filename: default.json",
			input:    []byte("default"),
			filename: "default.txt",
		},
		{
			dir:      ".",
			name:     "Input: <p>Content</p>, dir: ./html, filename: default.html",
			input:    []byte("default"),
			filename: "default.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NoError(t, SaveFile(test.input, test.dir, test.filename))

			content, err := ReadFile(test.filename)
			assert.NoError(t, err)
			assert.Equal(t, test.input, content)
		})
	}
}
