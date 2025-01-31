package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Save text file
func SaveTmpFile(filePath string, text string) {
	if filePath == "" {
		log.Fatalf("File path is empty")
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(text), 0644); err != nil {
		log.Fatalf("Failed to save review: %v", err)
	}
}

func LoadTmpFile(filePath string) string {
	if filePath == "" {
		log.Fatalf("File path is empty")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err.Error()
	}
	return strings.TrimRight(string(data), "\n")
}
