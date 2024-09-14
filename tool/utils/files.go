package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadFile(filename string, a ...any) ([]byte, error) {
	respFile, err := os.Open(fmt.Sprintf(filename, a...))
	if err != nil {
		return nil, err
	}
	defer respFile.Close()

	bytes, err := io.ReadAll(respFile)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func SaveFile(input []byte, dir, filename string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s", dir)
	}

	filename = fmt.Sprintf("%s/%s", dir, filename)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s", filename)
	}
	defer file.Close()

	if _, err = file.Write(input); err != nil {
		return fmt.Errorf("failed to save file content %s", filename)
	}

	return nil
}

func CopyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func ReadJSONFile[T any](path string, a ...any) (*T, error) {
	byteValue, err := ReadFile(path, a...)
	if err != nil {
		return nil, err
	}

	var dst T
	err = json.Unmarshal(byteValue, &dst)
	if err != nil {
		return nil, err
	}

	return &dst, nil
}

func SaveJSONFile(input any, dir, filename string) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling struct to JSON: %v", err)
	}

	return SaveFile(data, dir, filename)
}
