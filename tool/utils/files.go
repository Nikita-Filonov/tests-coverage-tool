package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", file, err)
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", file, err)
	}

	return data, nil
}

func SaveFile(input []byte, dir, filename string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	filename = fmt.Sprintf("%s/%s", dir, filename)

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() { _ = f.Close() }()

	if _, err = f.Write(input); err != nil {
		return fmt.Errorf("failed to save file content %s: %w", filename, err)
	}

	return nil
}

func CopyFile(source, destination string) error {
	sf, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", source, err)
	}
	defer func() { _ = sf.Close() }()

	df, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destination, err)
	}
	defer func() { _ = df.Close() }()

	_, err = io.Copy(df, sf)
	if err != nil {
		return fmt.Errorf("failed to copy from %s to %s: %w", source, destination, err)
	}

	err = df.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file %s: %w", destination, err)
	}

	return nil
}

func ReadJSONFile[T any](file string) (*T, error) {
	byteValue, err := ReadFile(file)
	if err != nil {
		return nil, err
	}

	var dst T
	err = json.Unmarshal(byteValue, &dst)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", file, err)
	}

	return &dst, nil
}

func SaveJSONFile(input any, dir, filename string) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %w", err)
	}

	return SaveFile(data, dir, filename)
}
