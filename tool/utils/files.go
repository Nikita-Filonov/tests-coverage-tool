package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func ReadJSONFile[T any](path string, a ...any) (*T, error) {
	path = fmt.Sprintf(path, a...)

	respFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer respFile.Close()

	byteValue, err := io.ReadAll(respFile)
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

func SaveFile(input []byte, dir, filename string) error {
	log.Printf("Starting to save file %s into directory %s", filename, dir)

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

	log.Printf("Successfully saved file %s into directory %s", filename, dir)

	return nil
}
