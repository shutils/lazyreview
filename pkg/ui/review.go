package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openai/openai-go"
	"github.com/shutils/lazyreview/pkg/state"
)

// JSONレビュー情報
type reviewInfo struct {
	ID     string `json:"id"`
	Param  string `json:"param"`
	Review string `json:"review"`
	State  string `json:"state"`
}

type ReviewState int

type reviewStateMsg struct {
	state ReviewState
}

type reviewMsg struct {
	id      string
	content string
}

type reviewStackMsg struct {
	id        string
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

func (m *model) saveReviews() tea.Cmd {
	var reviews []reviewInfo
	for _, review := range m.reviewList {
		reviews = append(reviews, reviewInfo{
			ID:     review.ID,
			Param:  review.Param,
			Review: review.Review,
			State:  "finish",
		})
	}
	jsonData, err := json.MarshalIndent(reviews, "", "  ")
	if err != nil {
		return func() tea.Msg {
			return SendErrorMessage("Failed to save marshal reviews", err)
		}
	}
	err = os.WriteFile(m.conf.Output, jsonData, 0644)
	if err != nil {
		return func() tea.Msg {
			return SendErrorMessage("Failed to save json", err)
		}
	}
	return nil
}

func (m *model) loadReviews() (*model, tea.Cmd) {
	data, err := os.ReadFile(m.outputFile)
	if os.IsNotExist(err) {
		m.reviewList = []reviewInfo{}
		return m, nil
	}
	if err != nil {
		return m, func() tea.Msg {
			return SendErrorMessage("Failed to read reviews", err)
		}
	}
	if err := json.Unmarshal(data, &m.reviewList); err != nil {
		return m, func() tea.Msg {
			return SendErrorMessage("Failed to unmarshal reviews", err)
		}
	}
	return m, nil
}

func (m *model) getReviewIndex(id string) int {
	for i, review := range m.reviewList {
		if review.ID == id {
			return i
		}
	}
	return -1
}

func (m *model) reviewContent() tea.Cmd {
	return func() tea.Msg {
		selectedItem, ok := m.panels.itemListPanel.SelectedItem().(listItem)

		var (
			chat   *openai.ChatCompletion
			review string
			err    error
		)
		if ok {
			context := m.getContextString()
			// Generate content by including contextItems
			content := previewContent(selectedItem, m.conf.Sources)
			content = context + content
			if m.instantPrompt == "" {
				chat, err = m.client.Getreviewfromchatgpt(content, m.conf)
			} else {
				chat, err = m.client.GetReviewFromChatGPTWithPrompt(content, m.conf, m.instantPrompt)
			}
			if err != nil {
				review = fmt.Sprintf("Failed to get review: %v", err)
			} else {
				review = chat.Choices[0].Message.Content
			}

			if m.instantPrompt != "" {
				m.uiState.PromptHistory = append(m.uiState.PromptHistory, m.instantPrompt)
			}
			promptToken := m.uiState.Usage.PromptTokens + chat.Usage.PromptTokens
			completionTokens := m.uiState.Usage.CompletionTokens + chat.Usage.CompletionTokens
			m.uiState.Usage = state.Usage{
				PromptTokens:     promptToken,
				CompletionTokens: completionTokens,
			}
			m.UpdateState()
			state.SaveState(m.stateFile, m.uiState)
		}
		return reviewMsg{
			id:      selectedItem.id,
			content: review,
		}

	}
}

func (m *model) getContextString() string {
	items := m.panels.contextListPanel.Items()
	if len(items) == 0 {
		return ""
	}
	var contextItems []string
	for _, item := range items {
		item, ok := item.(listItem)
		if ok {
			contextItems = append(contextItems, item.param+"\n"+previewContent(item, m.conf.Sources))
		}
	}
	return strings.Join(contextItems, "\n\n")
}

func (m *model) deleteReview(reviewID string) tea.Cmd {
	index := m.getReviewIndex(reviewID)
	if index == -1 {
		return func() tea.Msg {
			return SendErrorMessage(fmt.Sprintf("Review with ID %s not found", reviewID), nil)
		}
	}

	m.reviewList = append(m.reviewList[:index], m.reviewList[index+1:]...)

	m.saveReviews()
	return nil
}
