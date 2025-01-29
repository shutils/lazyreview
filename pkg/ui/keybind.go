package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type globalKeyMap struct {
	Quit      key.Binding
	ZoomPanel key.Binding
}

func (k globalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.ZoomPanel}
}

func (k globalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.ZoomPanel},
	}
}

func GetGlobalKeymap() globalKeyMap {
	return GlobalKeyMap
}

func GetListKeymap() listKeyMap {
	return ListKeyMap
}

func GetContentKeymap() contentKeyMap {
	return ContentKeyMap
}

func GetReviewKeymap() reviewKeyMap {
	return ReviewKeyMap
}

func GetReviewStackKeymap() reviewStackKeyMap {
	return ReviewStackKeyMap
}

func GetPromptKeymap() promptKeyMap {
	return PromptKeyMap
}

func GetConfigSummaryKeymap() configSummaryKeyMap {
	return ConfigSummaryKeyMap
}

func GetStateKeymap() stateKeyMap {
	return StateKeyMap
}

func GetContextKeymap() contextKeyMap {
	return ContextKeyMap
}

func GetSourceListKeymap() sourceListKeyMap {
	return SourceListKeyMap
}

var GlobalKeyMap = globalKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	ZoomPanel: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "zoom"),
	),
}

type listKeyMap struct {
	ListCursorDown            key.Binding
	ListCursorUp              key.Binding
	StartFilter               key.Binding
	ReviewStack               key.Binding
	ReloadItems               key.Binding
	FocusContentPanel         key.Binding
	FocusReviewStackPanel     key.Binding
	FocusInstantPrompt        key.Binding
	FocusStatePanel           key.Binding
	ReviewContentCursorDown   key.Binding
	ReviewContentCursorUp     key.Binding
	ReviewContentHalfViewDown key.Binding
	ReviewContentHalfViewUp   key.Binding
	OpenReview                key.Binding
	ToggleAiContext           key.Binding
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ListCursorDown,
		k.ListCursorUp,
		k.StartFilter,
		k.ReviewStack,
		k.ReloadItems,
		k.FocusContentPanel,
		k.FocusReviewStackPanel,
		k.FocusStatePanel,
		k.FocusInstantPrompt,
		k.OpenReview,
		k.ToggleAiContext,
		// k.ReviewContentCursorDown,
		// k.ReviewContentCursorUp,
		// k.ReviewContentHalfViewDown,
		// k.ReviewContentHalfViewUp,
	}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.ListCursorDown,
			k.ListCursorUp,
			k.StartFilter,
			k.ReviewStack,
			k.ReloadItems,
			k.FocusContentPanel,
			k.FocusReviewStackPanel,
			k.FocusStatePanel,
			k.OpenReview,
			// k.ReviewContentCursorDown,
			// k.ReviewContentCursorUp,
			// k.ReviewContentHalfViewDown,
			// k.ReviewContentHalfViewUp,
		},
	}
}

var ListKeyMap = listKeyMap{
	ListCursorDown: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	ListCursorUp: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	FocusStatePanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus state"),
	),
	FocusReviewStackPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus stack"),
	),
	StartFilter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	ReviewStack: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "review"),
	),
	ReloadItems: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "reload"),
	),
	FocusContentPanel: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "focus content"),
	),
	ReviewContentCursorDown: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "down review"),
	),
	ReviewContentCursorUp: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "up review"),
	),
	ReviewContentHalfViewDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "half down review"),
	),
	ReviewContentHalfViewUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "half up review"),
	),
	OpenReview: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open review"),
	),
	FocusInstantPrompt: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "focus prompt"),
	),
	ToggleAiContext: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add context"),
	),
}

type contentKeyMap struct {
	ItemContentCursorDown   key.Binding
	ItemContentCursorUp     key.Binding
	ItemContentHalfViewDown key.Binding
	ItemContentHalfViewUp   key.Binding
	ReviewStack             key.Binding
	FocusInstantPrompt      key.Binding
	FocusReviewPanel        key.Binding
	FocusListPanel          key.Binding
}

func (k contentKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ItemContentCursorDown,
		k.ItemContentCursorUp,
		k.ItemContentHalfViewDown,
		k.ItemContentHalfViewUp,
		k.ReviewStack,
		k.FocusInstantPrompt,
		k.FocusReviewPanel,
		k.FocusListPanel,
	}
}

func (k contentKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.ItemContentCursorDown,
			k.ItemContentCursorUp,
			k.ItemContentHalfViewDown,
			k.ItemContentHalfViewUp,
			k.ReviewStack,
			k.FocusInstantPrompt,
			k.FocusReviewPanel,
			k.FocusListPanel,
		},
	}
}

var ContentKeyMap = contentKeyMap{
	ItemContentCursorDown: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	ItemContentCursorUp: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	ItemContentHalfViewDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "half down"),
	),
	ItemContentHalfViewUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "half up"),
	),
	ReviewStack: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "review"),
	),
	FocusInstantPrompt: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "focus prompt"),
	),
	FocusReviewPanel: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "focus review"),
	),
	FocusListPanel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "focus list"),
	),
}

type reviewKeyMap struct {
	ReviewContentCursorDown   key.Binding
	ReviewContentCursorUp     key.Binding
	ReviewContentHalfViewDown key.Binding
	ReviewContentHalfViewUp   key.Binding
	ReviewStack               key.Binding
	FocusInstantPrompt        key.Binding
	FocusContentPanel         key.Binding
	FocusListPanel            key.Binding
}

func (k reviewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ReviewContentCursorDown,
		k.ReviewContentCursorUp,
		k.ReviewContentHalfViewDown,
		k.ReviewContentHalfViewUp,
		k.ReviewStack,
		k.FocusInstantPrompt,
		k.FocusContentPanel,
		k.FocusListPanel,
	}
}

func (k reviewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.ReviewContentCursorDown,
			k.ReviewContentCursorUp,
			k.ReviewContentHalfViewDown,
			k.ReviewContentHalfViewUp,
			k.ReviewStack,
			k.FocusInstantPrompt,
			k.FocusContentPanel,
			k.FocusListPanel,
		},
	}
}

var ReviewKeyMap = reviewKeyMap{
	ReviewContentCursorDown: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	ReviewContentCursorUp: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	ReviewContentHalfViewDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "half down"),
	),
	ReviewContentHalfViewUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "half up"),
	),
	ReviewStack: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "review"),
	),
	FocusInstantPrompt: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "focus prompt"),
	),
	FocusContentPanel: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "focus content"),
	),
	FocusListPanel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "focus list"),
	),
}

type reviewStackKeyMap struct {
	FocusListPanel    key.Binding
	FocusContextPanel key.Binding
}

func (k reviewStackKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusListPanel,
		k.FocusContextPanel,
	}
}

func (k reviewStackKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusListPanel,
			k.FocusContextPanel,
		},
	}
}

var ReviewStackKeyMap = reviewStackKeyMap{
	FocusListPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus list"),
	),
	FocusContextPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus context"),
	),
}

type promptKeyMap struct {
	Blur                     key.Binding
	InstantPromptHistoryPrev key.Binding
	InstantPromptHistoryNext key.Binding
	ReviewStack              key.Binding
}

func (k promptKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Blur, k.InstantPromptHistoryPrev, k.InstantPromptHistoryNext, k.ReviewStack}
}

func (k promptKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Blur, k.InstantPromptHistoryPrev, k.InstantPromptHistoryNext, k.ReviewStack},
	}
}

var PromptKeyMap = promptKeyMap{
	Blur: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "blur"),
	),
	InstantPromptHistoryPrev: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "history prev"),
	),
	InstantPromptHistoryNext: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "history next"),
	),
	ReviewStack: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "review"),
	),
}

type configSummaryKeyMap struct {
	FocusStateSummaryPanel key.Binding
	FocusSourceList        key.Binding
}

func (k configSummaryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusStateSummaryPanel,
		k.FocusSourceList,
	}
}

func (k configSummaryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusStateSummaryPanel,
			k.FocusSourceList,
		},
	}
}

var ConfigSummaryKeyMap = configSummaryKeyMap{
	FocusSourceList: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus source list"),
	),
	FocusStateSummaryPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus state"),
	),
}

type stateKeyMap struct {
	FocusListPanel          key.Binding
	FocusConfigSummaryPanel key.Binding
}

func (k stateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusListPanel,
		k.FocusConfigSummaryPanel,
	}
}

func (k stateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusListPanel,
			k.FocusConfigSummaryPanel,
		},
	}
}

var StateKeyMap = stateKeyMap{
	FocusListPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus list"),
	),
	FocusConfigSummaryPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus config"),
	),
}

type contextKeyMap struct {
	FocusSourceListPanel  key.Binding
	FocusReviewStackPanel key.Binding
	RemoveContext         key.Binding
	CursorDown            key.Binding
	CursorUp              key.Binding
}

func (k contextKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusSourceListPanel,
		k.FocusReviewStackPanel,
		k.RemoveContext,
		k.CursorDown,
		k.CursorUp,
	}
}
func (k contextKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusSourceListPanel,
			k.FocusReviewStackPanel,
			k.RemoveContext,
			k.CursorDown,
			k.CursorUp,
		},
	}
}

var ContextKeyMap = contextKeyMap{
	FocusSourceListPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus source list"),
	),
	FocusReviewStackPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus stack"),
	),
	RemoveContext: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "remove context"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("shift+j", "down"),
	),
	CursorUp: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("shift+k", "down"),
	),
}

type sourceListKeyMap struct {
	FocusContextPanel   key.Binding
	FocusConfigPanel    key.Binding
	ToggleSourceEnabled key.Binding
}

func (k sourceListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusContextPanel,
		k.FocusConfigPanel,
		k.ToggleSourceEnabled,
	}
}

func (k sourceListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusContextPanel,
			k.FocusConfigPanel,
			k.ToggleSourceEnabled,
		},
	}
}

var SourceListKeyMap = sourceListKeyMap{
	FocusContextPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus context"),
	),
	FocusConfigPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus config"),
	),
	ToggleSourceEnabled: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle source enabled"),
	),
}

func MakeBottomLine(globalHelp string, panelHelp string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left, globalHelp, " | ", panelHelp)
}

func (m *model) handleGlobalKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.globalKeyMap.Quit):
			return m.Quit
		case key.Matches(msg, m.globalKeyMap.ZoomPanel):
			return m.ZoomPanel
		}
	}
	return nil
}

func (m *model) handleListKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	if m.list.FilterState() == list.Filtering {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.listKeyMap.ListCursorDown):
			return m.ListCursorDown
		case key.Matches(msg, m.listKeyMap.ListCursorUp):
			return m.ListCursorUp
		case key.Matches(msg, m.listKeyMap.ReviewStack):
			return m.ReviewStack
		case key.Matches(msg, m.listKeyMap.ReloadItems):
			return m.ReloadItems
		case key.Matches(msg, m.listKeyMap.FocusContentPanel):
			return m.FocusContentPanel
		case key.Matches(msg, m.listKeyMap.FocusReviewStackPanel):
			return m.FocusReviewStackPanel
		case key.Matches(msg, m.listKeyMap.FocusStatePanel):
			return m.FocusStatePanel
		case key.Matches(msg, m.listKeyMap.FocusInstantPrompt):
			return m.FocusInstantPrompt
		case key.Matches(msg, m.listKeyMap.ReviewContentCursorDown):
			return m.ReviewContentCursorDown
		case key.Matches(msg, m.listKeyMap.ReviewContentCursorUp):
			return m.ReviewContentCursorUp
		case key.Matches(msg, m.listKeyMap.ReviewContentHalfViewDown):
			return m.ReviewContentHalfViewDown
		case key.Matches(msg, m.listKeyMap.ReviewContentHalfViewUp):
			return m.ReviewContentHalfViewUp
		case key.Matches(msg, m.listKeyMap.OpenReview):
			return m.OpenCurrentReview
		case key.Matches(msg, m.listKeyMap.ToggleAiContext):
			return m.ToggleAiContext
		}
	}
	return nil
}

func (m *model) handleContentKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.contentKeyMap.ItemContentHalfViewDown):
			return m.ItemContentHalfViewDown
		case key.Matches(msg, m.contentKeyMap.ItemContentHalfViewUp):
			return m.ItemContentHalfViewUp
		case key.Matches(msg, m.contentKeyMap.ItemContentCursorDown):
			return m.ItemContentCursorDown
		case key.Matches(msg, m.contentKeyMap.ItemContentCursorUp):
			return m.ItemContentCursorUp
		case key.Matches(msg, m.contentKeyMap.ReviewStack):
			return m.ReviewStack
		case key.Matches(msg, m.contentKeyMap.FocusInstantPrompt):
			return m.FocusInstantPrompt
		case key.Matches(msg, m.contentKeyMap.FocusReviewPanel):
			return m.FocusReviewPanel
		case key.Matches(msg, m.contentKeyMap.FocusListPanel):
			return m.FocusListPanel
		}
	}
	return nil
}

func (m *model) handleReviewKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.reviewKeyMap.ReviewContentCursorDown):
			return m.ReviewContentCursorDown
		case key.Matches(msg, m.reviewKeyMap.ReviewContentCursorUp):
			return m.ReviewContentCursorUp
		case key.Matches(msg, m.reviewKeyMap.ReviewContentHalfViewDown):
			return m.ReviewContentHalfViewDown
		case key.Matches(msg, m.reviewKeyMap.ReviewContentHalfViewUp):
			return m.ReviewContentHalfViewUp
		case key.Matches(msg, m.reviewKeyMap.ReviewStack):
			return m.ReviewStack
		case key.Matches(msg, m.reviewKeyMap.FocusInstantPrompt):
			return m.FocusInstantPrompt
		case key.Matches(msg, m.reviewKeyMap.FocusContentPanel):
			return m.FocusContentPanel
		case key.Matches(msg, m.reviewKeyMap.FocusListPanel):
			return m.FocusListPanel
		}
	}
	return nil
}

func (m *model) handleReviewStackKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.reviewStackKeyMap.FocusListPanel):
			return m.FocusListPanel
		case key.Matches(msg, m.reviewStackKeyMap.FocusContextPanel):
			return m.FocusContextPanel
		}
	}
	return nil
}

func (m *model) handlePromptKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.promptKeyMap.Blur):
			return m.BlurInstantPrompt
		case key.Matches(msg, m.promptKeyMap.InstantPromptHistoryPrev):
			return m.InstantPromptHistoryPrev
		case key.Matches(msg, m.promptKeyMap.InstantPromptHistoryNext):
			return m.InstantPromptHistoryNext
		case key.Matches(msg, m.promptKeyMap.ReviewStack):
			return m.ReviewStack
		}
	}
	return nil
}

func (m *model) handleConfigSummaryKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.configSummaryKeyMap.FocusStateSummaryPanel):
			return m.FocusStatePanel
		case key.Matches(msg, m.configSummaryKeyMap.FocusSourceList):
			return m.FocusSourceListPanel
		}
	}
	return nil
}

func (m *model) handleSourceListKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.sourceListKeyMap.FocusConfigPanel):
			return m.FocusConfigSummaryPanel
		case key.Matches(msg, m.sourceListKeyMap.FocusContextPanel):
			return m.FocusContextPanel
		case key.Matches(msg, m.sourceListKeyMap.ToggleSourceEnabled):
			return m.ToggleSourceEnabled
		}
	}
	return nil
}

func (m *model) handleStateKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.stateKeyMap.FocusConfigSummaryPanel):
			return m.FocusConfigSummaryPanel
		case key.Matches(msg, m.stateKeyMap.FocusListPanel):
			return m.FocusListPanel
		}
	}
	return nil
}

func (m *model) handleContextKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.contextKeyMap.FocusReviewStackPanel):
			return m.FocusReviewStackPanel
		case key.Matches(msg, m.contextKeyMap.FocusSourceListPanel):
			return m.FocusSourceListPanel
		case key.Matches(msg, m.contextKeyMap.RemoveContext):
			currItem := m.contextPanel.SelectedItem()
			currListItem, ok := currItem.(listItem)
			if !ok {
				return func() (tea.Model, tea.Cmd) {
					return m, nil
				}
			}
			return func() (tea.Model, tea.Cmd) {
				return m.removeContextStack(currListItem.id)
			}
		case key.Matches(msg, m.contextKeyMap.CursorDown):
			return m.ContextDetailCursorDown
		case key.Matches(msg, m.contextKeyMap.CursorUp):
			return m.ContextDetailCursorUp
		}
	}
	return nil
}

func (m *model) handleKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	if action := m.handleGlobalKey(msg); action != nil {
		return func() (tea.Model, tea.Cmd) {
			return action()
		}
	}

	switch m.focusState {
	case ListPanelFocus:
		if action := m.handleListKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case ContentPanelFocus:
		if action := m.handleContentKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case ReviewPanelFocus:
		if action := m.handleReviewKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case ReviewStackPanelFocus:
		if action := m.handleReviewStackKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case InstantPromptPanelFocus:
		if action := m.handlePromptKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case ConfigSummaryPanelFocus:
		if action := m.handleConfigSummaryKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case StatePanelFocus:
		if action := m.handleStateKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case ContextPanelFocus:
		if action := m.handleContextKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	case SourceListPanelFocus:
		if action := m.handleSourceListKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	}
	return nil
}
