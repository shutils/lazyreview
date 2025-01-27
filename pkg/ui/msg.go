package ui

type AiContextMethod int

const (
	AddContext AiContextMethod = iota
	RemoveContext
)

type aiContextMsg struct {
	method AiContextMethod
	id     string
}
