package ui

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/state"
)

func (m *model) Quit() (tea.Model, tea.Cmd) {
	return *m, tea.Quit
}

func (m *model) FocusNextPanel() (tea.Model, tea.Cmd) {
	m.focusState = (m.focusState + 1) % 3
	return *m, nil
}

func (m *model) ZoomPanel() (tea.Model, tea.Cmd) {
	m.zoomState = (m.zoomState + 1) % 3
	return *m, nil
}

func (m *model) ListCursorDown() (tea.Model, tea.Cmd) {
	m.list.CursorDown()
	return m.onChangeListSelectedItem()
}

func (m *model) ListCursorUp() (tea.Model, tea.Cmd) {
	m.list.CursorUp()
	return m.onChangeListSelectedItem()
}

func (m *model) ReviewStack() (tea.Model, tea.Cmd) {
	item := m.list.SelectedItem().(listItem)
	var cmds []tea.Cmd
	cmds = append(cmds, func() tea.Msg {
		return reviewStackMsg{
			id:        item.id,
			operation: Add,
		}
	})
	cmds = append(cmds, m.reviewContent())
	return *m, tea.Batch(cmds...)
}

func (m *model) ToggleAiContext() (tea.Model, tea.Cmd) {
	item := m.list.SelectedItem().(listItem)
	index := findIndex(m.contextPanel.Items(), item.id)
	if index == -1 {
		return m.addContextStack(item.id)
	} else {
		return m.removeContextStack(item.id)
	}
}

func (m *model) ReloadItems() (tea.Model, tea.Cmd) {
	if m.list.FilterState() == list.Unfiltered {
		m.list.SetItems(getItems(m.conf, m.reviewList))
	}
	return *m, nil
}

func (m *model) ItemContentCursorDown() (tea.Model, tea.Cmd) {
	m.contentPanel.LineDown(1)
	return *m, nil
}

func (m *model) ItemContentCursorUp() (tea.Model, tea.Cmd) {
	m.contentPanel.LineUp(1)
	return *m, nil
}

func (m *model) ItemContentHalfViewDown() (tea.Model, tea.Cmd) {
	m.contentPanel.HalfViewDown()
	return *m, nil
}

func (m *model) ItemContentHalfViewUp() (tea.Model, tea.Cmd) {
	m.contentPanel.HalfViewUp()
	return *m, nil
}

func (m *model) ReviewContentCursorDown() (tea.Model, tea.Cmd) {
	m.reviewPanel.LineDown(1)
	return *m, nil
}

func (m *model) ReviewContentCursorUp() (tea.Model, tea.Cmd) {
	m.reviewPanel.LineUp(1)
	return *m, nil
}

func (m *model) ReviewContentHalfViewDown() (tea.Model, tea.Cmd) {
	m.reviewPanel.HalfViewDown()
	return *m, nil
}

func (m *model) ReviewContentHalfViewUp() (tea.Model, tea.Cmd) {
	m.reviewPanel.HalfViewUp()
	return *m, nil
}

func (m *model) ContextDetailCursorDown() (tea.Model, tea.Cmd) {
	m.contextDetailPanel.LineDown(1)
	return *m, nil
}

func (m *model) ContextDetailCursorUp() (tea.Model, tea.Cmd) {
	m.contextDetailPanel.LineUp(1)
	return *m, nil
}

func (m *model) FocusInstantPrompt() (tea.Model, tea.Cmd) {
	m.focusState = InstantPromptPanelFocus
	m.instantPromptPanel.Focus()
	return *m, nil
}

func (m *model) BlurInstantPrompt() (tea.Model, tea.Cmd) {
	m.focusState = ContentPanelFocus
	m.instantPromptPanel.Blur()
	return *m, nil
}

func (m *model) focusPanel(panel FocusState) (tea.Model, tea.Cmd) {
	m.focusState = panel
	cmd := func() tea.Msg {
		return updateFocusPanelMsg{
			target: panel,
		}
	}
	return *m, cmd
}

func (m *model) FocusListPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ListPanelFocus)
}

func (m *model) FocusContentPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ContentPanelFocus)
}

func (m *model) FocusReviewPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ReviewPanelFocus)
}

func (m *model) FocusReviewStackPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ReviewStackProgressPanelFocus)
}

func (m *model) FocusInstantPromptPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(InstantPromptPanelFocus)
}

func (m *model) FocusConfigSummaryPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ConfigSummaryPanelFocus)
}

func (m *model) FocusStatePanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(StatePanelFocus)
}

func (m *model) FocusContextPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ContextPanelFocus)
}

func (m *model) FocusSourceListPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(SourceListPanelFocus)
}

func (m *model) InstantPromptHistoryPrev() (tea.Model, tea.Cmd) {
	if m.currentHistoryIndex-1 >= 0 && len(m.uiState.PromptHistory) > 0 {
		m.instantPromptPanel.SetValue(m.uiState.PromptHistory[m.currentHistoryIndex-1])
		m.currentHistoryIndex--
	}

	return m, nil
}

func (m *model) InstantPromptHistoryNext() (tea.Model, tea.Cmd) {
	if m.currentHistoryIndex+1 < len(m.uiState.PromptHistory) && len(m.uiState.PromptHistory) > 0 {
		m.instantPromptPanel.SetValue(m.uiState.PromptHistory[m.currentHistoryIndex+1])
		m.currentHistoryIndex++
	} else {
		m.instantPromptPanel.SetValue("")
		m.currentHistoryIndex = len(m.uiState.PromptHistory)
	}

	return m, nil
}

func (m *model) OpenCurrentReview() (tea.Model, tea.Cmd) {
	selectedItem, ok := m.list.SelectedItem().(listItem)
	if !ok {
		return m, nil
	}

	if m.getReviewIndex(selectedItem.id) == -1 {
		return m, nil
	}

	review := m.reviewList[m.getReviewIndex(selectedItem.id)].Review
	state.SaveTmpReview(m.conf.TmpReviewPath, review)

	exec.Command(m.conf.Opener, m.conf.TmpReviewPath).Start()
	return m, nil
}

func (m *model) ToggleSourceEnabled() (tea.Model, tea.Cmd) {
	selectedItem := m.sourceListPanel.SelectedItem().(sourceItem)
	m.conf.ToggleSourceEnabled(selectedItem.name)
	cmd := func() tea.Msg {
		return updateSourceListMsg{}
	}
	return m, cmd
}
