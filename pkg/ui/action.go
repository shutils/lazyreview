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
	m.panels.itemListPanel.CursorDown()
	return m.onChangeListSelectedItem()
}

func (m *model) ListCursorUp() (tea.Model, tea.Cmd) {
	m.panels.itemListPanel.CursorUp()
	return m.onChangeListSelectedItem()
}

func (m *model) ReviewStack() (tea.Model, tea.Cmd) {
	item := m.panels.itemListPanel.SelectedItem().(listItem)
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
	item := m.panels.itemListPanel.SelectedItem().(listItem)
	index := findIndex(m.panels.contextListPanel.Items(), item.id)
	if index == -1 {
		return m.addContextStack(item.id)
	} else {
		return m.removeContextStack(item.id)
	}
}

func (m *model) ReloadItems() (tea.Model, tea.Cmd) {
	if m.panels.itemListPanel.FilterState() == list.Unfiltered {
		m.panels.itemListPanel.SetItems(getItems(m.conf, m.reviewList))
	}
	return *m, nil
}

func (m *model) ItemContentCursorDown() (tea.Model, tea.Cmd) {
	m.panels.itemPreviewPanel.LineDown(1)
	return *m, nil
}

func (m *model) ItemContentCursorUp() (tea.Model, tea.Cmd) {
	m.panels.itemPreviewPanel.LineUp(1)
	return *m, nil
}

func (m *model) ItemContentHalfViewDown() (tea.Model, tea.Cmd) {
	m.panels.itemPreviewPanel.HalfViewDown()
	return *m, nil
}

func (m *model) ItemContentHalfViewUp() (tea.Model, tea.Cmd) {
	m.panels.itemPreviewPanel.HalfViewUp()
	return *m, nil
}

func (m *model) ReviewContentCursorDown() (tea.Model, tea.Cmd) {
	m.panels.itemReviewPanel.LineDown(1)
	return *m, nil
}

func (m *model) ReviewContentCursorUp() (tea.Model, tea.Cmd) {
	m.panels.itemReviewPanel.LineUp(1)
	return *m, nil
}

func (m *model) ReviewContentHalfViewDown() (tea.Model, tea.Cmd) {
	m.panels.itemReviewPanel.HalfViewDown()
	return *m, nil
}

func (m *model) ReviewContentHalfViewUp() (tea.Model, tea.Cmd) {
	m.panels.itemReviewPanel.HalfViewUp()
	return *m, nil
}

func (m *model) ContextDetailCursorDown() (tea.Model, tea.Cmd) {
	m.panels.contextDetailPanel.LineDown(1)
	return *m, nil
}

func (m *model) ContextDetailCursorUp() (tea.Model, tea.Cmd) {
	m.panels.contextDetailPanel.LineUp(1)
	return *m, nil
}

func (m *model) FocusInstantPrompt() (tea.Model, tea.Cmd) {
	m.focusState = InstantPromptPanelFocus
	m.panels.promptPanel.Focus()
	return *m, nil
}

func (m *model) BlurInstantPrompt() (tea.Model, tea.Cmd) {
	m.focusState = ContentPanelFocus
	m.panels.promptPanel.Blur()
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

func (m *model) FocusItemListPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ItemListPanelFocus)
}

func (m *model) FocusContentPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ContentPanelFocus)
}

func (m *model) FocusReviewPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ReviewPanelFocus)
}

func (m *model) FocusReviewProgressPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ReviewStackProgressPanelFocus)
}

func (m *model) FocusInstantPromptPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(InstantPromptPanelFocus)
}

func (m *model) FocusConfigSummaryPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ConfigSummaryPanelFocus)
}

func (m *model) FocusStateSummaryPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(StatePanelFocus)
}

func (m *model) FocusContextPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ContextPanelFocus)
}

func (m *model) FocusSourceListPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(SourceListPanelFocus)
}

func (m *model) ExitMessagePanel() (tea.Model, tea.Cmd) {
	m.message = ""
	return m.focusPanel(ItemListPanelFocus)
}

func (m *model) InstantPromptHistoryPrev() (tea.Model, tea.Cmd) {
	if m.currentHistoryIndex-1 >= 0 && len(m.uiState.PromptHistory) > 0 {
		m.panels.promptPanel.SetValue(m.uiState.PromptHistory[m.currentHistoryIndex-1])
		m.currentHistoryIndex--
	}

	return m, nil
}

func (m *model) InstantPromptHistoryNext() (tea.Model, tea.Cmd) {
	if m.currentHistoryIndex+1 < len(m.uiState.PromptHistory) && len(m.uiState.PromptHistory) > 0 {
		m.panels.promptPanel.SetValue(m.uiState.PromptHistory[m.currentHistoryIndex+1])
		m.currentHistoryIndex++
	} else {
		m.panels.promptPanel.SetValue("")
		m.currentHistoryIndex = len(m.uiState.PromptHistory)
	}

	return m, nil
}

func (m *model) OpenCurrentReview() (tea.Model, tea.Cmd) {
	selectedItem, ok := m.panels.itemListPanel.SelectedItem().(listItem)
	if !ok {
		return m, nil
	}

	if m.getReviewIndex(selectedItem.id) == -1 {
		return m, nil
	}

	review := m.reviewList[m.getReviewIndex(selectedItem.id)].Review
	state.SaveTmpReview(m.conf.TmpReviewPath, review)

	c := exec.Command(m.conf.Opener, m.conf.TmpReviewPath)
	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		var message = ""
		if err != nil {
			message = err.Error()
		}
		return showMessageMsg{
			message: message,
		}
	})
}

func (m *model) ToggleSourceEnabled() (tea.Model, tea.Cmd) {
	selectedItem := m.panels.sourceListPanel.SelectedItem().(sourceItem)
	m.conf.ToggleSourceEnabled(selectedItem.name)
	cmd := func() tea.Msg {
		return updateSourceListMsg{}
	}
	return m, cmd
}

func (m *model) DeleteReviewResult() (tea.Model, tea.Cmd) {
	selectedItem := m.panels.itemListPanel.SelectedItem()
	item, ok := selectedItem.(listItem)
	if !ok {
		return m, func() tea.Msg {
			return SendErrorMessage("Failed to delete review:", nil)
		}
	}
	m.deleteReview(item.id)
	index := m.panels.itemListPanel.Index()
	m.changeItemTitlePrefix(index, "☐ ")
	m.onChangeListSelectedItem()
	return m, nil
}

func (m *model) EditPromptInEditor() (tea.Model, tea.Cmd) {
	return m.OpenPromptInEditor()
}
