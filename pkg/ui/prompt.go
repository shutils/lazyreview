package ui

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) OpenPromptInEditor() (tea.Model, tea.Cmd) {
	text := m.panels.promptPanel.Value()

	SaveTmpFile(m.conf.TmpPromptPath, text)

	c := exec.Command(m.conf.Opener, m.conf.TmpPromptPath)
	return *m, tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return closedEditorMsg{
				err: err,
			}
		}
		editedText := LoadTmpFile(m.conf.TmpPromptPath)
		return setPromptMsg{
			text: editedText,
		}
	})
}
