package ui

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/config"
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
	m.panels.itemListPanel.model.CursorDown()
	return m.onChangeListSelectedItem()
}

func (m *model) ListCursorUp() (tea.Model, tea.Cmd) {
	m.panels.itemListPanel.model.CursorUp()
	return m.onChangeListSelectedItem()
}

func (m *model) ReviewStack() (tea.Model, tea.Cmd) {
	item := m.panels.itemListPanel.model.SelectedItem().(listItem)
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
	item := m.panels.itemListPanel.model.SelectedItem().(listItem)
	index := findIndex(m.panels.contextListPanel.Items(), item.id)
	if index == -1 {
		return m.addContextStack(item.id)
	} else {
		return m.removeContextStack(item.id)
	}
}

func (m *model) EditContext() (tea.Model, tea.Cmd) {
	contextItem, ok := m.panels.contextListPanel.SelectedItem().(contextItem)
	if !ok {
		return *m, nil
	}
	m.panels.contextEditPanel.SetValue(contextItem.content)
	m.FocusContextEditPanel()
	m.panels.contextEditPanel.Focus()
	return *m, nil
}

func (m *model) SaveEditingContext() (tea.Model, tea.Cmd) {
	item, ok := m.panels.contextListPanel.SelectedItem().(contextItem)
	if !ok {
		return *m, nil
	}
	index := findIndex(m.panels.contextListPanel.Items(), item.id)
	text := m.panels.contextEditPanel.Value()
	newContextItem := contextItem{
		title:      item.title,
		param:      item.param,
		sourceName: item.sourceName,
		id:         item.id,
		content:    text,
		isEdited:   true,
	}
	m.panels.contextListPanel.SetItem(index, newContextItem)
	return *m, nil
}

func (m *model) ReloadItems() (tea.Model, tea.Cmd) {
	if m.panels.itemListPanel.model.FilterState() == list.Unfiltered {
		m.panels.itemListPanel.model.SetItems(getItems(m.conf, m.reviewList))
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
	m.prevFocusState = m.focusState
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

func (m *model) FocusContextEditPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(ContextEditPanelFocus)
}

func (m *model) FocusHelpPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(HelpPanelFocus)
}

func (m *model) BlurHelpPanel() (tea.Model, tea.Cmd) {
	return m.focusPanel(m.prevFocusState)
}

func (m *model) ExitMessagePanel() (tea.Model, tea.Cmd) {
	m.message = ""
	return m.focusPanel(ItemListPanelFocus)
}

func (m *model) InstantPromptHistoryPrev() (tea.Model, tea.Cmd) {
	if len(m.uiState.PromptHistory) > 0 && m.currentHistoryIndex > 0 {
		m.currentHistoryIndex--
		m.panels.promptPanel.SetValue(m.uiState.PromptHistory[m.currentHistoryIndex])
	}

	m.instantPrompt = m.panels.promptPanel.Value()

	return m, nil
}

func (m *model) InstantPromptHistoryNext() (tea.Model, tea.Cmd) {
	if len(m.uiState.PromptHistory) > 0 && m.currentHistoryIndex+1 < len(m.uiState.PromptHistory) {
		m.currentHistoryIndex++
		m.panels.promptPanel.SetValue(m.uiState.PromptHistory[m.currentHistoryIndex])
	} else {
		m.panels.promptPanel.SetValue("")
		m.currentHistoryIndex = len(m.uiState.PromptHistory)
	}

	m.instantPrompt = m.panels.promptPanel.Value()

	return m, nil
}

func (m *model) OpenCurrentReview() (tea.Model, tea.Cmd) {
	selectedItem, ok := m.panels.itemListPanel.model.SelectedItem().(listItem)
	if !ok {
		return m, nil
	}

	if m.getReviewIndex(selectedItem.id) == -1 {
		return m, nil
	}

	review := m.reviewList[m.getReviewIndex(selectedItem.id)].Review
	state.SaveTmpReview(m.conf.TmpReviewPath, review)

	cmdName, cmdArgs := config.ParseCommand(m.conf.Opener)

	cmdArgs = append(cmdArgs, m.conf.TmpReviewPath)
	c := exec.Command(cmdName, cmdArgs...)
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
	selectedItem := m.panels.sourceListPanel.SelectedItem().(config.Source)
	m.conf.ToggleSourceEnabled(selectedItem.Name)
	cmd := func() tea.Msg {
		return updateSourceListMsg{}
	}
	return m, cmd
}

func (m *model) DeleteReviewResult() (tea.Model, tea.Cmd) {
	selectedItem := m.panels.itemListPanel.model.SelectedItem()
	item, ok := selectedItem.(listItem)
	if !ok {
		return m, func() tea.Msg {
			return SendErrorMessage("Failed to delete review:", nil)
		}
	}
	m.deleteReview(item.id)
	index := m.panels.itemListPanel.model.Index()
	m.changeItemTitlePrefix(index, "☐ ")
	m.onChangeListSelectedItem()
	return m, nil
}

func (m *model) ToggleItemListViewStyle() (tea.Model, tea.Cmd) {
	m.panels.itemListPanel.showDescription = !m.panels.itemListPanel.showDescription

	if m.panels.itemListPanel.showDescription {
		m.panels.itemListPanel.model.SetDelegate(m.panels.itemListPanel.normalDelegate)
	} else {
		m.panels.itemListPanel.model.SetDelegate(m.panels.itemListPanel.narrowDelegate)
	}
	return m, nil
}

func (m *model) EditPromptInEditor() (tea.Model, tea.Cmd) {
	return m.OpenPromptInEditor()
}
