package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	state "github.com/shutils/lazyreview/pkg/state"
)

// JSONレビュー情報
type reviewInfo struct {
	Param  string `json:"param"`
	Review string `json:"review"`
	State  string `json:"state"`
}

type ReviewState int

type reviewStateMsg struct {
	state ReviewState
}

type reviewMsg struct {
	param   string
	content string
}

type reviewStackMsg struct {
	param     string
	operation ReviewStackOperation
}

type ReviewStackOperation int

const (
	Add ReviewStackOperation = iota
	Remove
)

const (
	NoAction ReviewState = iota
	Reviewing
)

func (m *model) saveReviews() {
	var reviews []reviewInfo
	for _, review := range m.reviewList {
		reviews = append(reviews, reviewInfo{Param: review.Param, Review: review.Review, State: "finish"})
	}
	jsonData, _ := json.MarshalIndent(reviews, "", "  ")
	_ = os.WriteFile(m.conf.Output, jsonData, 0644)
}

func (m *model) loadReviews() {
	data, err := os.ReadFile(m.outputFile)
	if os.IsNotExist(err) {
		// ファイルが存在しない場合、新しいマップを初期化
		m.reviewList = []reviewInfo{}
		return
	}
	if err != nil {
		log.Fatalf("Failed to read reviews: %v", err)
	}
	if err := json.Unmarshal(data, &m.reviewList); err != nil {
		log.Fatalf("Failed to unmarshal reviews: %v", err)
	}
}

func (m *model) getReviewIndex(param string) int {
	for i, review := range m.reviewList {
		if review.Param == param {
			return i
		}
	}
	return -1
}

func (m *model) reviewContent() tea.Cmd {
	return func() tea.Msg {
		selectedItem, ok := m.list.SelectedItem().(listItem)

		var (
			review string
			err    error
		)

		if ok {
			content := ""
			if m.conf.Previewer != "" {
				content = customPreviewer(m.conf.Previewer, selectedItem.param)
			} else {
				content = defaultPreviewer(selectedItem.param)
			}

			if m.instantPrompt == "" {
				review, err = m.client.Getreviewfromchatgpt(content, m.conf)
			} else {
				review, err = m.client.GetReviewFromChatGPTWithPrompt(content, m.conf, m.instantPrompt)
				m.uiState.PromptHistory = append(m.uiState.PromptHistory, m.instantPrompt)
				state.SaveState(m.stateFile, m.uiState)
			}
			if err != nil {
				review = fmt.Sprintf("Failed to get review: %v", err)
			}
		}
		return reviewMsg{
			param:   selectedItem.param,
			content: review,
		}
	}
}
