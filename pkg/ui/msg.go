package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type AiContextMethod int

const (
	AddContext AiContextMethod = iota
	RemoveContext
)

type updateFocusPanelMsg struct {
	target FocusState
}

type closedEditorMsg struct {
	err error
}

type setPromptMsg struct {
	text string
}

type showMessageMsg struct {
	message string
}

func SendErrorMessage(desc string, err error) tea.Msg {
	var sb strings.Builder
	sb.WriteString(desc)
	sb.WriteString(": \n\n")
	sb.WriteString(err.Error())
	message := sb.String()
	return showMessageMsg{
		message: message,
	}
}
