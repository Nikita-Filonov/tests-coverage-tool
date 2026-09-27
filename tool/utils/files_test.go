package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			dir:      t.TempDir(),
			name:     "Input: default, dir: ., filename: default.txt",
			input:    []byte("default"),
			filename: "default.txt",
		},
		{
			dir:      t.TempDir(),
			name:     "Input: {\"key\": \"value\"}, dir: ./json, filename: default.json",
			input:    []byte("default"),
			filename: "default.txt",
		},
		{
			dir:      t.TempDir(),
			name:     "Input: <p>Content</p>, dir: ./html, filename: default.html",
			input:    []byte("default"),
			filename: "default.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NoError(t, SaveFile(test.input, test.dir, test.filename))

			content, err := ReadFile(filepath.Join(test.dir, test.filename))
			assert.NoError(t, err)
			assert.Equal(t, test.input, content)
		})
	}
}

func TestJSONFileRoundTrip(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "reports")
	input := map[string]string{"name": "api", "value": "<script>$1</script>"}
	require.NoError(t, SaveJSONFile(input, dir, "state.json"))
	output, err := ReadJSONFile[map[string]string](filepath.Join(dir, "state.json"))
	require.NoError(t, err)
	assert.Equal(t, input, *output)
}

func TestJSONFileErrors(t *testing.T) {
	dir := t.TempDir()
	assert.Error(t, SaveJSONFile(make(chan int), dir, "state.json"))
	filename := filepath.Join(dir, "invalid.json")
	require.NoError(t, os.WriteFile(filename, []byte("{"), 0o600))
	_, err := ReadJSONFile[map[string]string](filename)
	assert.Error(t, err)
	_, err = ReadJSONFile[map[string]string](filepath.Join(dir, "missing.json"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	source, destination := filepath.Join(dir, "source.txt"), filepath.Join(dir, "destination.txt")
	require.NoError(t, os.WriteFile(source, []byte("content"), 0o600))
	require.NoError(t, CopyFile(source, destination))
	content, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, []byte("content"), content)
	assert.ErrorIs(t, CopyFile(filepath.Join(dir, "missing.txt"), destination), os.ErrNotExist)
	assert.Error(t, CopyFile(source, dir))
}

func TestSaveFileErrors(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "file.txt")
	require.NoError(t, os.WriteFile(filename, []byte("content"), 0o600))
	assert.Error(t, SaveFile([]byte("new"), filename, "output.txt"))
	assert.Error(t, SaveFile([]byte("new"), dir, "missing/output.txt"))
}

func TestReadAndCopyDirectoryErrors(t *testing.T) {
	dir := t.TempDir()
	_, err := ReadFile(dir)
	assert.Error(t, err)
	assert.Error(t, CopyFile(dir, filepath.Join(t.TempDir(), "destination")))
}
