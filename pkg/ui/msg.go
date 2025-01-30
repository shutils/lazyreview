package ui

type AiContextMethod int

const (
	AddContext AiContextMethod = iota
	RemoveContext
)

type updateFocusPanelMsg struct {
	target FocusState
}
