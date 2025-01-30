package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMaps struct {
	globalKeyMap
	listKeyMap
	contentKeyMap
	reviewKeyMap
	reviewStackKeyMap
	promptKeyMap
	configSummaryKeyMap
	stateKeyMap
	contextKeyMap
	sourceListKeyMap
	messageKeyMap
}

func DefaultKeyMap() keyMaps {
	return keyMaps{
		globalKeyMap:        GetGlobalKeymap(),
		listKeyMap:          GetListKeymap(),
		contentKeyMap:       GetContentKeymap(),
		reviewKeyMap:        GetReviewKeymap(),
		reviewStackKeyMap:   GetReviewStackKeymap(),
		promptKeyMap:        GetPromptKeymap(),
		configSummaryKeyMap: GetConfigSummaryKeymap(),
		stateKeyMap:         GetStateKeymap(),
		contextKeyMap:       GetContextKeymap(),
		sourceListKeyMap:    GetSourceListKeymap(),
		messageKeyMap:       GetMessageKeymap(),
	}
}

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

func GetMessageKeymap() messageKeyMap {
	return MessageKeyMap
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
	FocusInstantPrompt        key.Binding
	FocusStatePanel           key.Binding
	FocusContextListPanel     key.Binding
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
		k.FocusStatePanel,
		k.FocusContextListPanel,
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
			k.FocusStatePanel,
			k.FocusContextListPanel,
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
	FocusContextListPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus context"),
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
	FocusConfigSummaryPanel key.Binding
	FocusStateSummaryPanel  key.Binding
}

func (k reviewStackKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusConfigSummaryPanel,
		k.FocusStateSummaryPanel,
	}
}

func (k reviewStackKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusConfigSummaryPanel,
			k.FocusStateSummaryPanel,
		},
	}
}

var ReviewStackKeyMap = reviewStackKeyMap{
	FocusConfigSummaryPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus config"),
	),
	FocusStateSummaryPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus state"),
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
	FocusReviewProgressPanel key.Binding
	FocusSourceListPanel     key.Binding
}

func (k configSummaryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusSourceListPanel,
		k.FocusReviewProgressPanel,
	}
}

func (k configSummaryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusSourceListPanel,
			k.FocusReviewProgressPanel,
		},
	}
}

var ConfigSummaryKeyMap = configSummaryKeyMap{
	FocusSourceListPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus source list"),
	),
	FocusReviewProgressPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus progress"),
	),
}

type stateKeyMap struct {
	FocusItemListPanel       key.Binding
	FocusReviewProgressPanel key.Binding
}

func (k stateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusItemListPanel,
		k.FocusReviewProgressPanel,
	}
}

func (k stateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusItemListPanel,
			k.FocusReviewProgressPanel,
		},
	}
}

var StateKeyMap = stateKeyMap{
	FocusItemListPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus list"),
	),
	FocusReviewProgressPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus config"),
	),
}

type contextKeyMap struct {
	FocusItemListPanel   key.Binding
	FocusSourceListPanel key.Binding
	RemoveContext        key.Binding
	CursorDown           key.Binding
	CursorUp             key.Binding
}

func (k contextKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.FocusItemListPanel,
		k.FocusSourceListPanel,
		k.RemoveContext,
		k.CursorDown,
		k.CursorUp,
	}
}
func (k contextKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.FocusItemListPanel,
			k.FocusSourceListPanel,
			k.RemoveContext,
			k.CursorDown,
			k.CursorUp,
		},
	}
}

var ContextKeyMap = contextKeyMap{
	FocusItemListPanel: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "focus stack"),
	),
	FocusSourceListPanel: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "focus source list"),
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

type messageKeyMap struct {
	Quit   key.Binding
	Return key.Binding
}

func (k messageKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Return,
	}
}

func (k messageKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Quit,
			k.Return,
		},
	}
}

var MessageKeyMap = messageKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	Return: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "return"),
	),
}

func MakeBottomLine(globalHelp string, panelHelp string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left, globalHelp, " | ", panelHelp)
}

func (m *model) handleGlobalKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.globalKeyMap.Quit):
			return m.Quit
		case key.Matches(msg, m.keyMaps.globalKeyMap.ZoomPanel):
			return m.ZoomPanel
		}
	}
	return nil
}

func (m *model) handleItemListKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	if m.panels.itemListPanel.FilterState() == list.Filtering {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.listKeyMap.ListCursorDown):
			return m.ListCursorDown
		case key.Matches(msg, m.keyMaps.listKeyMap.ListCursorUp):
			return m.ListCursorUp
		case key.Matches(msg, m.keyMaps.listKeyMap.ReviewStack):
			return m.ReviewStack
		case key.Matches(msg, m.keyMaps.listKeyMap.ReloadItems):
			return m.ReloadItems
		case key.Matches(msg, m.keyMaps.listKeyMap.FocusContentPanel):
			return m.FocusContentPanel
		case key.Matches(msg, m.keyMaps.listKeyMap.FocusContextListPanel):
			return m.FocusContextPanel
		case key.Matches(msg, m.keyMaps.listKeyMap.FocusStatePanel):
			return m.FocusStateSummaryPanel
		case key.Matches(msg, m.keyMaps.listKeyMap.FocusInstantPrompt):
			return m.FocusInstantPrompt
		case key.Matches(msg, m.keyMaps.listKeyMap.ReviewContentCursorDown):
			return m.ReviewContentCursorDown
		case key.Matches(msg, m.keyMaps.listKeyMap.ReviewContentCursorUp):
			return m.ReviewContentCursorUp
		case key.Matches(msg, m.keyMaps.listKeyMap.ReviewContentHalfViewDown):
			return m.ReviewContentHalfViewDown
		case key.Matches(msg, m.keyMaps.listKeyMap.ReviewContentHalfViewUp):
			return m.ReviewContentHalfViewUp
		case key.Matches(msg, m.keyMaps.listKeyMap.OpenReview):
			return m.OpenCurrentReview
		case key.Matches(msg, m.keyMaps.listKeyMap.ToggleAiContext):
			return m.ToggleAiContext
		}
	}
	return nil
}

func (m *model) handleContentKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.contentKeyMap.ItemContentHalfViewDown):
			return m.ItemContentHalfViewDown
		case key.Matches(msg, m.keyMaps.contentKeyMap.ItemContentHalfViewUp):
			return m.ItemContentHalfViewUp
		case key.Matches(msg, m.keyMaps.contentKeyMap.ItemContentCursorDown):
			return m.ItemContentCursorDown
		case key.Matches(msg, m.keyMaps.contentKeyMap.ItemContentCursorUp):
			return m.ItemContentCursorUp
		case key.Matches(msg, m.keyMaps.contentKeyMap.ReviewStack):
			return m.ReviewStack
		case key.Matches(msg, m.keyMaps.contentKeyMap.FocusInstantPrompt):
			return m.FocusInstantPrompt
		case key.Matches(msg, m.keyMaps.contentKeyMap.FocusReviewPanel):
			return m.FocusReviewPanel
		case key.Matches(msg, m.keyMaps.contentKeyMap.FocusListPanel):
			return m.FocusItemListPanel
		}
	}
	return nil
}

func (m *model) handleReviewKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.reviewKeyMap.ReviewContentCursorDown):
			return m.ReviewContentCursorDown
		case key.Matches(msg, m.keyMaps.reviewKeyMap.ReviewContentCursorUp):
			return m.ReviewContentCursorUp
		case key.Matches(msg, m.keyMaps.reviewKeyMap.ReviewContentHalfViewDown):
			return m.ReviewContentHalfViewDown
		case key.Matches(msg, m.keyMaps.reviewKeyMap.ReviewContentHalfViewUp):
			return m.ReviewContentHalfViewUp
		case key.Matches(msg, m.keyMaps.reviewKeyMap.ReviewStack):
			return m.ReviewStack
		case key.Matches(msg, m.keyMaps.reviewKeyMap.FocusInstantPrompt):
			return m.FocusInstantPrompt
		case key.Matches(msg, m.keyMaps.reviewKeyMap.FocusContentPanel):
			return m.FocusContentPanel
		case key.Matches(msg, m.keyMaps.reviewKeyMap.FocusListPanel):
			return m.FocusItemListPanel
		}
	}
	return nil
}

func (m *model) handleReviewStackKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.reviewStackKeyMap.FocusConfigSummaryPanel):
			return m.FocusConfigSummaryPanel
		case key.Matches(msg, m.keyMaps.reviewStackKeyMap.FocusStateSummaryPanel):
			return m.FocusStateSummaryPanel
		}
	}
	return nil
}

func (m *model) handlePromptKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.promptKeyMap.Blur):
			return m.BlurInstantPrompt
		case key.Matches(msg, m.keyMaps.promptKeyMap.InstantPromptHistoryPrev):
			return m.InstantPromptHistoryPrev
		case key.Matches(msg, m.keyMaps.promptKeyMap.InstantPromptHistoryNext):
			return m.InstantPromptHistoryNext
		case key.Matches(msg, m.keyMaps.promptKeyMap.ReviewStack):
			return m.ReviewStack
		}
	}
	return nil
}

func (m *model) handleConfigSummaryKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.configSummaryKeyMap.FocusReviewProgressPanel):
			return m.FocusReviewProgressPanel
		case key.Matches(msg, m.keyMaps.configSummaryKeyMap.FocusSourceListPanel):
			return m.FocusSourceListPanel
		}
	}
	return nil
}

func (m *model) handleSourceListKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.sourceListKeyMap.FocusConfigPanel):
			return m.FocusConfigSummaryPanel
		case key.Matches(msg, m.keyMaps.sourceListKeyMap.FocusContextPanel):
			return m.FocusContextPanel
		case key.Matches(msg, m.keyMaps.sourceListKeyMap.ToggleSourceEnabled):
			return m.ToggleSourceEnabled
		}
	}
	return nil
}

func (m *model) handleStateKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.stateKeyMap.FocusReviewProgressPanel):
			return m.FocusReviewProgressPanel
		case key.Matches(msg, m.keyMaps.stateKeyMap.FocusItemListPanel):
			return m.FocusItemListPanel
		}
	}
	return nil
}

func (m *model) handleContextKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.contextKeyMap.FocusItemListPanel):
			return m.FocusItemListPanel
		case key.Matches(msg, m.keyMaps.contextKeyMap.FocusSourceListPanel):
			return m.FocusSourceListPanel
		case key.Matches(msg, m.keyMaps.contextKeyMap.RemoveContext):
			currItem := m.panels.contextListPanel.SelectedItem()
			currListItem, ok := currItem.(listItem)
			if !ok {
				return func() (tea.Model, tea.Cmd) {
					return m, nil
				}
			}
			return func() (tea.Model, tea.Cmd) {
				return m.removeContextStack(currListItem.id)
			}
		case key.Matches(msg, m.keyMaps.contextKeyMap.CursorDown):
			return m.ContextDetailCursorDown
		case key.Matches(msg, m.keyMaps.contextKeyMap.CursorUp):
			return m.ContextDetailCursorUp
		}
	}
	return nil
}

func (m *model) handleMessageKey(msg tea.Msg) func() (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.messageKeyMap.Quit):
			return m.Quit
		case key.Matches(msg, m.keyMaps.messageKeyMap.Return):
			return m.ExitMessagePanel
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
	case ItemListPanelFocus:
		if action := m.handleItemListKey(msg); action != nil {
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
	case ReviewStackProgressPanelFocus:
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
	case MessagePanelFocus:
		if action := m.handleMessageKey(msg); action != nil {
			return func() (tea.Model, tea.Cmd) {
				return action()
			}
		}
	}
	return nil
}
