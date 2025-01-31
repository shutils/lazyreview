package state

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shutils/lazyreview/pkg/config"
)

type Usage struct {
	PromptTokens, CompletionTokens int64
}

type State struct {
	PromptHistory []string `json:"promptHistory"`
	Usage         Usage    `json:"usage"`
}

func (s *State) ShowUsage(cost config.ModelCost) string {
	if cost.Input == 0 || cost.Output == 0 {
		return ""
	}
	inputCost := float64(s.Usage.PromptTokens) * (cost.Input) / 1000_000
	outputCost := float64(s.Usage.CompletionTokens) * (cost.Output) / 1000_000
	return "Cost: $" + strconv.FormatFloat(inputCost+outputCost, 'f', -1, 64)
}

func (s *State) ShowUsedToken() string {
	title := "Used tokens:"
	inputStr := "  Input: " + strconv.Itoa(int(s.Usage.PromptTokens))
	outputStr := "  Output: " + strconv.Itoa(int(s.Usage.CompletionTokens))
	return title + "\n" + inputStr + "\n" + outputStr
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
