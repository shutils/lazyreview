package ui

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/config"
)

func (m *model) OpenPromptInEditor() (tea.Model, tea.Cmd) {
	text := m.panels.promptPanel.Value()

	SaveTmpFile(m.conf.TmpPromptPath, text)

	cmdName, cmdArgs := config.ParseCommand(m.conf.Opener)

	cmdArgs = append(cmdArgs, m.conf.TmpPromptPath)
	c := exec.Command(cmdName, cmdArgs...)
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

func (m *model) OpenContextInEditor() (tea.Model, tea.Cmd) {
	text := m.panels.contextEditPanel.Value()

	SaveTmpFile(m.conf.TmpContextPath, text)

	cmdName, cmdArgs := config.ParseCommand(m.conf.Opener)

	cmdArgs = append(cmdArgs, m.conf.TmpContextPath)
	c := exec.Command(cmdName, cmdArgs...)
	return *m, tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return closedEditorMsg{
				err: err,
			}
		}
		editedText := LoadTmpFile(m.conf.TmpContextPath)
		return setEditedContextMsg{
			text: editedText,
		}
	})
}
