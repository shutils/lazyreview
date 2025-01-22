package state

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type State struct {
	PromptHistory []string `json:"promptHistory"`
}

func LoadState(stateFilePath string) State {
	s := State{}
	data, err := os.ReadFile(stateFilePath)
	if os.IsNotExist(err) {
		return s
	}
	if err != nil {
		log.Fatalf("Failed to read state file: %v", err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		log.Fatalf("Failed to unmarshal state: %v", err)
	}
	return s
}

func SaveState(filePath string, state State) {
	if filePath == "" {
		log.Fatalf("File path is empty")
	}

	jsonData, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal state: %v", err)
	}

	// Ensure the directory exists
	if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	if err = os.WriteFile(filePath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to save state: %v", err)
	}
}

func SaveTmpReview(filePath string, review string) {
	if filePath == "" {
		log.Fatalf("File path is empty")
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(review), 0644); err != nil {
		log.Fatalf("Failed to save review: %v", err)
	}
}
