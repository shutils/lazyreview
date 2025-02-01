package ui

import (
	"math"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/width"
)

type ZoomState int
type FocusState int

type panels struct {
	itemListPanel       itemListPanel
	itemPreviewPanel    viewport.Model
	itemReviewPanel     viewport.Model
	stateSummaryPanel   viewport.Model
	stateDetailPanel    viewport.Model
	configSummaryPanel  viewport.Model
	configDetailPanel   viewport.Model
	reviewProgressPanel progress.Model
	reviewStackPanel    viewport.Model
	sourceListPanel     list.Model
	sourceDetailPanel   viewport.Model
	contextListPanel    list.Model
	contextDetailPanel  viewport.Model
	promptPanel         textarea.Model
	spinner             spinner.Model
	messagePanel        viewport.Model
}

func NewPanels() panels {
	p := panels{
		itemListPanel:       NewItemListPanel(),
		itemPreviewPanel:    viewport.New(0, 0),
		itemReviewPanel:     viewport.New(0, 0),
		stateSummaryPanel:   viewport.New(0, 0),
		stateDetailPanel:    viewport.New(0, 0),
		configSummaryPanel:  viewport.New(0, 0),
		configDetailPanel:   viewport.New(0, 0),
		reviewProgressPanel: progress.New(),
		reviewStackPanel:    viewport.New(0, 0),
		sourceListPanel:     list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		sourceDetailPanel:   viewport.New(0, 0),
		contextListPanel:    list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		contextDetailPanel:  viewport.New(0, 0),
		promptPanel:         textarea.New(),
		spinner:             spinner.New(),
		messagePanel:        viewport.New(0, 0),
	}

	p.setInitSetting()
	return p
}

func (p *panels) setInitSetting() {
	p.sourceListPanel = setListInitSetting(p.sourceListPanel)
	p.contextListPanel = setListInitSetting(p.contextListPanel)

	p.itemListPanel.model.SetShowHelp(false)
	p.itemListPanel.model.SetShowTitle(false)
	p.itemListPanel.model.KeyMap.Quit.Unbind()
}

func setListInitSetting(l list.Model) list.Model {
	l.SetDelegate(list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	})
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.KeyMap.Quit.Unbind()
	l.KeyMap.Filter.Unbind()
	return l
}

const (
	Normal ZoomState = iota
	Middle
	Max
)

const (
	ItemListPanelFocus FocusState = iota
	ContentPanelFocus
	ReviewPanelFocus
	ReviewStackProgressPanelFocus
	InstantPromptPanelFocus
	ConfigSummaryPanelFocus
	StatePanelFocus
	ContextPanelFocus
	SourceListPanelFocus
	MessagePanelFocus
	Other
)

const (
	stateSummaryPanelHeight   = 1
	configSummaryPanelHeight  = 1
	reviewProgressPanelHeight = 1
	contextListPanelMaxHeight = 5
	instantPromptPanelHeight  = 5
	sourceListPanelMaxHeight  = 5
	footerHeight              = 1

	listPaginationHeight = 2

	borderHeight = 1
	borderWidth  = 1
)

var (
	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
)

func runeWidth(r rune) int {
	prop := width.LookupRune(r)
	switch prop.Kind() {
	case width.EastAsianFullwidth, width.EastAsianWide:
		return 2
	default:
		return 1
	}
}

func replacePrefix(text string, prefix string) string {
	runes := []rune(text)
	prefixRunes := []rune(prefix)

	if len(runes) < len(prefixRunes) {
		return prefix
	}

	updatedTitle := string(prefixRunes) + string(runes[len(prefixRunes):])
	return updatedTitle
}

func visibleRunes(s string) []rune {
	ansiEscapePattern := `\x1b\[[0-9;]*m`
	re := regexp.MustCompile(ansiEscapePattern)

	cleaned := re.ReplaceAllString(s, "")

	return []rune(cleaned)
}

func stringWidth(runes []rune) int {
	l := 0
	for _, r := range runes {
		l = l + runeWidth(r)
	}
	return l
}

// POC
func InsertTitleWithOffset(rendered, title string) string {
	borderStart := "╭"
	borderEnd := "╮"
	borderBar := "─"

	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return rendered
	}

	firstLine := lines[0]
	ansiEscaseStart := firstLine[:strings.Index(firstLine, borderStart)+len(borderStart)]
	ansiEscaseEnd := firstLine[strings.Index(firstLine, borderEnd)+len(borderEnd):]

	titleWidth := stringWidth([]rune(title))
	runes := visibleRunes(firstLine)

	if titleWidth < len(runes) {
		remainingContent := string(runes[titleWidth+2:])
		lines[0] = ansiEscaseStart + borderBar + title + remainingContent + ansiEscaseEnd
	}
	return strings.Join(lines, "\n")
}

func (m *model) handleWindowSize(msg tea.WindowSizeMsg) (model, tea.Cmd) {
	m.winSize.height = msg.Height
	m.winSize.width = msg.Width
	m.setPanelSize()
	return *m, nil
}

func (m *model) makeView() string {
	m.setPanelSize()
	if m.message != "" {
		m.focusState = MessagePanelFocus
		m.panels.messagePanel.SetContent(m.message)
		return m.panels.messagePanel.View()
	}

	state := "/"
	if m.reviewState == Reviewing {
		state = m.panels.spinner.View()
	}

	helpModel := help.New()
	globalHelp := helpModel.View(m.keyMaps.globalKeyMap)
	helpString := m.getHelpString(helpModel, globalHelp)

	listPanel := m.buildPanel(m.panels.itemListPanel.model.View(), m.getPanelStyle(ItemListPanelFocus), m.panels.itemListPanel.model.Width(), m.panels.itemListPanel.model.Height(), "List")
	contentPanel := m.buildPanel(m.panels.itemPreviewPanel.View(), m.getPanelStyle(ContentPanelFocus), m.panels.itemPreviewPanel.Width, m.panels.itemPreviewPanel.Height, "Content")
	reviewPanel := m.buildPanel(m.panels.itemReviewPanel.View(), m.getPanelStyle(ReviewPanelFocus), m.panels.itemReviewPanel.Width, m.panels.itemReviewPanel.Height, "Review")
	reviewStackPanel := m.buildPanel(m.panels.reviewStackPanel.View(), m.getPanelStyle(Other), m.panels.reviewStackPanel.Width, m.panels.reviewStackPanel.Height, "Review stack")
	configPanel := m.buildPanel(m.panels.configSummaryPanel.View(), m.getPanelStyle(ConfigSummaryPanelFocus), m.panels.configSummaryPanel.Width, m.panels.configSummaryPanel.Height, "Config")
	configContentPanel := m.buildPanel(m.panels.configDetailPanel.View(), m.getPanelStyle(Other), m.panels.configDetailPanel.Width, m.panels.configDetailPanel.Height, "Config content")
	statePanel := m.buildPanel(m.panels.stateSummaryPanel.View(), m.getPanelStyle(StatePanelFocus), m.panels.stateSummaryPanel.Width, m.panels.stateSummaryPanel.Height, "State")
	stateDetailPanel := m.buildPanel(m.panels.stateDetailPanel.View(), m.getPanelStyle(Other), m.panels.stateDetailPanel.Width, m.panels.stateDetailPanel.Height, "State detail")
	instantPromptPanel := m.buildPanel(m.panels.promptPanel.View(), m.getPanelStyle(InstantPromptPanelFocus), m.panelSize.secondlyPanelWidth, instantPromptPanelHeight, "Instant prompt")
	contextPanel := m.buildPanel(m.panels.contextListPanel.View(), m.getPanelStyle(ContextPanelFocus), m.panels.contextListPanel.Width(), m.panels.contextListPanel.Height(), "Context")
	sourceListPanel := m.buildPanel(m.panels.sourceListPanel.View(), m.getPanelStyle(SourceListPanelFocus), m.panels.sourceListPanel.Width(), m.panels.sourceListPanel.Height(), "Source list")
	sourceDetailPanel := m.buildPanel(m.panels.sourceDetailPanel.View(), m.getPanelStyle(Other), m.panels.sourceDetailPanel.Width, m.panels.sourceDetailPanel.Height, "Source detail")
	contextDetailPanel := m.buildPanel(m.panels.contextDetailPanel.View(), m.getPanelStyle(Other), m.panels.contextDetailPanel.Width, m.panels.contextDetailPanel.Height, "Context detail")
	reviewProgressPanel := m.buildPanel(m.panels.reviewProgressPanel.View(), m.getPanelStyle(ReviewStackProgressPanelFocus), m.panels.reviewProgressPanel.Width, 1, "Review progress")

	primaryPanels := m.buildPrimaryPanels(statePanel, listPanel, reviewProgressPanel, contextPanel, sourceListPanel, configPanel)
	reviewCombiPanels := m.buildReviewCombiPanels(contentPanel, reviewPanel, instantPromptPanel)
	contentPanels := m.buildContentPanels(contentPanel, instantPromptPanel)
	reviewPanels := m.buildReviewPanels(reviewPanel, instantPromptPanel)

	bottomLine := lipgloss.JoinHorizontal(lipgloss.Left, state, " ", helpString)
	return m.buildWindowBasedOnFocus(
		primaryPanels,
		configContentPanel,
		stateDetailPanel,
		reviewCombiPanels,
		contentPanels,
		reviewPanels,
		sourceDetailPanel,
		contextDetailPanel,
		reviewStackPanel,
		bottomLine,
	)
}

func (m *model) getHelpString(helpModel help.Model, globalHelp string) string {
	switch m.focusState {
	case ItemListPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.listKeyMap))
	case ContentPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.contentKeyMap))
	case ReviewPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.reviewKeyMap))
	case ReviewStackProgressPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.reviewStackKeyMap))
	case InstantPromptPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.promptKeyMap))
	case ConfigSummaryPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.configSummaryKeyMap))
	case StatePanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.stateKeyMap))
	case ContextPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.contextKeyMap))
	case SourceListPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.keyMaps.sourceListKeyMap))
	default:
		return ""
	}
}

func (m *model) getPanelStyle(focus FocusState) lipgloss.Style {
	style := baseStyle
	if m.focusState == focus {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return style
}

func (m *model) buildWindowBasedOnFocus(primaryPanels,
	configContentPanel,
	stateDetailPanel,
	reviewCombiPanels,
	contentPanels,
	reviewPanels,
	sourceDetailPanel,
	contextDetailPanel,
	reviewStackPanel,
	bottomLine string,
) string {
	switch m.focusState {
	case ConfigSummaryPanelFocus:
		return m.buildWindow(primaryPanels, configContentPanel, bottomLine)
	case StatePanelFocus:
		return m.buildWindow(primaryPanels, stateDetailPanel, bottomLine)
	case ItemListPanelFocus:
		return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
	case ContentPanelFocus:
		return m.buildWindowBasedOnZoom(primaryPanels, reviewCombiPanels, contentPanels, bottomLine)
	case ReviewPanelFocus, InstantPromptPanelFocus:
		return m.buildWindowBasedOnZoom(primaryPanels, reviewCombiPanels, reviewPanels, bottomLine)
	case SourceListPanelFocus:
		return m.buildWindow(primaryPanels, sourceDetailPanel, bottomLine)
	case ContextPanelFocus:
		return m.buildWindow(primaryPanels, contextDetailPanel, bottomLine)
	case ReviewStackProgressPanelFocus:
		return m.buildWindow(primaryPanels, reviewStackPanel, bottomLine)
	default:
		return ""
	}
}

func (m *model) buildWindowBasedOnZoom(primaryPanels, reviewCombiPanels, zoomPanels, bottomLine string) string {
	switch m.zoomState {
	case Normal:
		return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
	case Middle:
		return m.buildWindowZoom(reviewCombiPanels, bottomLine)
	case Max:
		return m.buildWindowZoom(zoomPanels, bottomLine)
	default:
		return ""
	}
}

func (m *model) calcAreaSize() (int, int) {
	var primaryAreaWidth, secondlyAreaWidth int
	switch m.zoomState {
	case Normal:
		if isFocusPrimary(m.focusState) {
			primaryAreaWidth = m.winSize.width / 3
		} else {
			primaryAreaWidth = m.winSize.width / 5
		}
	case Middle:
		if isFocusPrimary(m.focusState) {
			primaryAreaWidth = m.winSize.width / 2
		} else {
			primaryAreaWidth = 0
		}
	case Max:
		if isFocusPrimary(m.focusState) {
			primaryAreaWidth = m.winSize.width
		} else {
			primaryAreaWidth = 0
		}
	}
	secondlyAreaWidth = m.winSize.width - primaryAreaWidth
	return primaryAreaWidth, secondlyAreaWidth
}

func (m *model) setPanelSize() (tea.Model, tea.Cmd) {
	m.setPrimaryPanelSizes()
	m.setSecondaryPanelSizes()

	m.panels.messagePanel.Width = m.winSize.width
	m.panels.messagePanel.Height = m.winSize.height
	return m, nil
}

func (m *model) setPrimaryPanelSizes() {
	const (
		stateSummaryPanelOuterHeight  = stateSummaryPanelHeight + borderHeight*2
		reviewProgresPanelOuterHeight = reviewProgressPanelHeight + borderHeight*2
		configSummaryPanelOuterHeight = configSummaryPanelHeight + borderHeight*2
	)
	contextListPanelHeight := m.getListPanelHeight(m.panels.contextListPanel.Items(), contextListPanelMaxHeight)
	contextListPanelOuterHeight := contextListPanelHeight + listPaginationHeight
	if len(m.panels.contextListPanel.Items()) > contextListPanelMaxHeight {
		m.panels.contextListPanel.SetShowPagination(true)
	} else {
		m.panels.contextListPanel.SetShowPagination(false)
	}

	sourceListPanelHeight := m.getListPanelHeight(m.panels.sourceListPanel.Items(), sourceListPanelMaxHeight)
	sourceListPanelOuterHeight := sourceListPanelHeight + listPaginationHeight
	if len(m.panels.sourceListPanel.Items()) > sourceListPanelMaxHeight {
		m.panels.sourceListPanel.SetShowPagination(true)
	} else {
		m.panels.sourceListPanel.SetShowPagination(false)
	}

	listPanelHeight := m.winSize.height - stateSummaryPanelOuterHeight - reviewProgresPanelOuterHeight - contextListPanelOuterHeight - configSummaryPanelOuterHeight - footerHeight - listPaginationHeight - sourceListPanelOuterHeight
	primaryAreaWidth, _ := m.calcAreaSize()

	m.panels.itemListPanel.model.SetSize(primaryAreaWidth-borderWidth*2, listPanelHeight)
	m.panels.configSummaryPanel.Width = primaryAreaWidth - borderWidth*2
	m.panels.configSummaryPanel.Height = configSummaryPanelHeight
	m.panels.reviewProgressPanel.Width = primaryAreaWidth - borderWidth*2
	m.panels.stateSummaryPanel.Width = primaryAreaWidth - borderWidth*2
	m.panels.stateSummaryPanel.Height = stateSummaryPanelHeight
	m.panels.contextListPanel.SetSize(primaryAreaWidth-borderWidth*2, contextListPanelHeight)
	m.panels.sourceListPanel.SetSize(primaryAreaWidth-borderWidth*2, sourceListPanelHeight)
}

func (m *model) setSecondaryPanelSizes() {
	const (
		instantPromptPanelOuterHeight = instantPromptPanelHeight + borderHeight*2
	)

	_, secondlyAreaWidth := m.calcAreaSize()

	if m.zoomState == Max {
		if isFocusItemPreviewPanel(m.focusState) {
			m.panelSize.itemPreviewPanelWidth = secondlyAreaWidth
		} else {
			m.panelSize.itemPreviewPanelWidth = 0
		}
	} else {
		m.panelSize.itemPreviewPanelWidth = secondlyAreaWidth / 2
	}
	m.panelSize.itemReviewPanelWidth = secondlyAreaWidth - m.panelSize.itemPreviewPanelWidth

	m.panels.itemPreviewPanel.Width = m.panelSize.itemPreviewPanelWidth - borderWidth*2
	m.panels.itemPreviewPanel.Height = m.winSize.height - instantPromptPanelOuterHeight - borderHeight*2 - footerHeight

	m.panels.configDetailPanel.Width = secondlyAreaWidth - borderWidth*2
	m.panels.configDetailPanel.Height = m.winSize.height - borderHeight*2 - footerHeight

	m.panels.itemReviewPanel.Width = m.panelSize.itemReviewPanelWidth - borderWidth*2
	m.panels.itemReviewPanel.Height = m.winSize.height - instantPromptPanelOuterHeight - borderHeight*2 - footerHeight

	m.panels.promptPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.promptPanel.SetHeight(instantPromptPanelHeight)

	m.panels.reviewStackPanel.Width = secondlyAreaWidth - borderWidth*2
	m.panels.reviewStackPanel.Height = m.winSize.height - borderHeight*2 - footerHeight

	m.panels.stateDetailPanel.Width = secondlyAreaWidth - borderWidth*2
	m.panels.stateDetailPanel.Height = m.winSize.height - borderHeight*2 - footerHeight

	m.panels.sourceDetailPanel.Width = secondlyAreaWidth - borderWidth*2
	m.panels.sourceDetailPanel.Height = m.winSize.height - borderHeight*2 - footerHeight

	m.panels.contextDetailPanel.Width = secondlyAreaWidth - borderWidth*2
	m.panels.contextDetailPanel.Height = m.winSize.height - borderHeight*2 - footerHeight
}

func (m *model) buildPrimaryPanels(statePanel, listPanel, reviewStackPanel, contextPanel, sourceListPanel, configPanel string) string {
	return lipgloss.JoinVertical(lipgloss.Top, statePanel, listPanel, contextPanel, sourceListPanel, configPanel, reviewStackPanel)
}

func (m *model) buildReviewCombiPanels(contentPanel, reviewPanel, instantPromptPanel string) string {
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left, contentPanel, reviewPanel), instantPromptPanel)
}

func (m *model) buildContentPanels(contentPanel, instantPromptPanel string) string {
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left, contentPanel), instantPromptPanel)
}

func (m *model) buildReviewPanels(reviewPanel, instantPromptPanel string) string {
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left, reviewPanel), instantPromptPanel)
}

func (m *model) buildPanel(p string, s lipgloss.Style, w, h int, title string) string {
	return InsertTitleWithOffset(s.Width(w).Height(h).Render(p), title)
}

func (m *model) buildWindow(primary, secondly, bottom string) string {
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, primary, secondly), bottom)
}

func (m *model) buildWindowZoom(panels, bottom string) string {
	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, panels), bottom)
}

func (m *model) getListPanelHeight(items []list.Item, max int) int {
	return int(math.Max(math.Min(float64(len(items)), float64(max)), 1))
}

func getRendered(text string, style string, width int) string {
	if style != "" {
		r, _ := glamour.NewTermRenderer(
			glamour.WithStylePath(style),
			glamour.WithWordWrap(width),
		)
		rendered, err := r.Render(text)
		if err != nil {
			rendered = "Error rendering review content"
		}
		return rendered
	} else {
		return text
	}
}

func isFocusPrimary(state FocusState) bool {
	if state == ItemListPanelFocus || state == ConfigSummaryPanelFocus || state == StatePanelFocus || state == ContextPanelFocus || state == ReviewStackProgressPanelFocus || state == SourceListPanelFocus {
		return true
	}
	return false
}

func isFocusItemPreviewPanel(state FocusState) bool {
	return state == ContentPanelFocus
}

type itemListPanel struct {
	model           list.Model
	showDescription bool
	normalDelegate  list.DefaultDelegate
	narrowDelegate  list.DefaultDelegate
}

func NewItemListPanel() itemListPanel {
	l := itemListPanel{}
	l.model = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.showDescription = true
	l.normalDelegate = l.NewDefaultNormalDelegate()
	l.narrowDelegate = l.NewDefaultNarrowDelegate()
	return l
}

func (l itemListPanel) NewDefaultNormalDelegate() list.DefaultDelegate {
	return list.NewDefaultDelegate()
}

func (l itemListPanel) NewDefaultNarrowDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	return delegate
}
