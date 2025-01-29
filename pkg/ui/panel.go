package ui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/width"
)

type ZoomState int
type FocusState int

const (
	Normal ZoomState = iota
	Middle
	Max
)

const (
	ListPanelFocus FocusState = iota
	ContentPanelFocus
	ReviewPanelFocus
	ReviewStackPanelFocus
	InstantPromptPanelFocus
	ConfigSummaryPanelFocus
	StatePanelFocus
	ContextPanelFocus
	SourceListPanelFocus
	Other
)

const (
	stateSummaryPanelHeight  = 1
	configSummaryPanelHeight = 1
	reviewStackPanelHeight   = 5
	contextPanelHeight       = 5
	instantPromptPanelHeight = 5
	sourceListPanelHeight    = 6
	footerHeight             = 1

	borderHeight = 1
	borderWidth  = 1
)

var (
	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
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
	borderStart := "┌"
	borderEnd := "┐"
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

func (m *model) handleWindowSize(msg tea.WindowSizeMsg) {
	m.winSize.height = msg.Height
	m.winSize.width = msg.Width
}

func (m *model) makeView() string {
	m.setPanelSize()

	state := "/"
	if m.reviewState == Reviewing {
		state = m.spinner.View()
	}

	helpModel := help.New()
	globalHelp := helpModel.View(m.globalKeyMap)
	helpString := m.getHelpString(helpModel, globalHelp)

	listPanel := m.buildPanel(m.list.View(), m.getPanelStyle(ListPanelFocus), m.list.Width(), m.list.Height(), "List")
	contentPanel := m.buildPanel(m.contentPanel.View(), m.getPanelStyle(ContentPanelFocus), m.contentPanel.Width, m.contentPanel.Height, "Content")
	reviewPanel := m.buildPanel(m.reviewPanel.View(), m.getPanelStyle(ReviewPanelFocus), m.reviewPanel.Width, m.reviewPanel.Height, "Review")
	reviewStackPanel := m.buildPanel(m.reviewStackPanel.View(), m.getPanelStyle(ReviewStackPanelFocus), m.reviewStackPanel.Width, m.reviewStackPanel.Height, "Review stack")
	configPanel := m.buildPanel(m.configSummaryPanel.View(), m.getPanelStyle(ConfigSummaryPanelFocus), m.configSummaryPanel.Width, m.configSummaryPanel.Height, "Config")
	configContentPanel := m.buildPanel(m.configContentPanel.View(), m.getPanelStyle(Other), m.configContentPanel.Width, m.configContentPanel.Height, "Config content")
	statePanel := m.buildPanel(m.statePanel.View(), m.getPanelStyle(StatePanelFocus), m.statePanel.Width, m.statePanel.Height, "State")
	stateDetailPanel := m.buildPanel(m.stateDetailPanel.View(), m.getPanelStyle(Other), m.stateDetailPanel.Width, m.stateDetailPanel.Height, "State detail")
	instantPromptPanel := m.buildPanel(m.instantPromptPanel.View(), m.getPanelStyle(InstantPromptPanelFocus), m.panelSize.secondlyPanelWidth, instantPromptPanelHeight, "Instant prompt")
	contextPanel := m.buildPanel(m.contextPanel.View(), m.getPanelStyle(ContextPanelFocus), m.contextPanel.Width, m.contextPanel.Height, "Context")
	sourceListPanel := m.buildPanel(m.sourceListPanel.View(), m.getPanelStyle(SourceListPanelFocus), m.sourceListPanel.Width(), m.sourceListPanel.Height(), "Source list")
	sourceDetailPanel := m.buildPanel(m.sourceDetailPanel.View(), m.getPanelStyle(Other), m.sourceDetailPanel.Width, m.sourceDetailPanel.Height, "Source detail")

	primaryPanels := m.buildPrimaryPanels(statePanel, listPanel, reviewStackPanel, contextPanel, sourceListPanel, configPanel)
	reviewCombiPanels := m.buildReviewCombiPanels(contentPanel, reviewPanel, instantPromptPanel)
	contentPanels := m.buildContentPanels(contentPanel, instantPromptPanel)
	reviewPanels := m.buildReviewPanels(reviewPanel, instantPromptPanel)

	bottomLine := lipgloss.JoinHorizontal(lipgloss.Left, state, " ", helpString)
	return m.buildWindowBasedOnFocus(primaryPanels, configContentPanel, stateDetailPanel, reviewCombiPanels, contentPanels, reviewPanels, sourceDetailPanel, bottomLine)
}

func (m *model) getHelpString(helpModel help.Model, globalHelp string) string {
	switch m.focusState {
	case ListPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.listKeyMap))
	case ContentPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.contentKeyMap))
	case ReviewPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.reviewKeyMap))
	case ReviewStackPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.reviewStackKeyMap))
	case InstantPromptPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.promptKeyMap))
	case ConfigSummaryPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.configSummaryKeyMap))
	case StatePanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.stateKeyMap))
	case ContextPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.contextKeyMap))
	case SourceListPanelFocus:
		return MakeBottomLine(globalHelp, helpModel.View(m.sourceListKeyMap))
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

func (m *model) buildWindowBasedOnFocus(primaryPanels, configContentPanel, stateDetailPanel, reviewCombiPanels, contentPanels, reviewPanels, sourceDetailPanel, bottomLine string) string {
	switch m.focusState {
	case ConfigSummaryPanelFocus:
		return m.buildWindow(primaryPanels, configContentPanel, bottomLine)
	case StatePanelFocus:
		return m.buildWindow(primaryPanels, stateDetailPanel, bottomLine)
	case ListPanelFocus, ReviewStackPanelFocus, ContextPanelFocus:
		return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
	case ContentPanelFocus:
		return m.buildWindowBasedOnZoom(primaryPanels, reviewCombiPanels, contentPanels, bottomLine)
	case ReviewPanelFocus, InstantPromptPanelFocus:
		return m.buildWindowBasedOnZoom(primaryPanels, reviewCombiPanels, reviewPanels, bottomLine)
	case SourceListPanelFocus:
		return m.buildWindow(primaryPanels, sourceDetailPanel, bottomLine)
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

func (m *model) setPanelSize() {
	m.setPrimaryPanelSizes()
	m.setSecondaryPanelSizes()
}

func (m *model) setPrimaryPanelSizes() {
	const (
		stateSummaryPanelOuterHeight  = stateSummaryPanelHeight + borderHeight*2
		reviewStackPanelOuterHeight   = reviewStackPanelHeight + borderHeight*2
		contextPanelOuterHeight       = contextPanelHeight + borderHeight*2
		configSummaryPanelOuterHeight = configSummaryPanelHeight + borderHeight*2
	)

	listPanelHeight := m.winSize.height - stateSummaryPanelOuterHeight - reviewStackPanelOuterHeight - contextPanelOuterHeight - configSummaryPanelOuterHeight - footerHeight - borderHeight*2 - sourceListPanelHeight - borderHeight*2
	primaryAreaWidth, _ := m.calcAreaSize()

	m.list.SetSize(primaryAreaWidth-borderWidth*2, listPanelHeight)
	m.configSummaryPanel.Width = primaryAreaWidth - borderWidth*2
	m.configSummaryPanel.Height = configSummaryPanelHeight
	m.reviewStackPanel.Width = primaryAreaWidth - borderWidth*2
	m.reviewStackPanel.Height = reviewStackPanelHeight
	m.statePanel.Width = primaryAreaWidth - borderWidth*2
	m.statePanel.Height = stateSummaryPanelHeight
	m.contextPanel.Width = primaryAreaWidth - borderWidth*2
	m.contextPanel.Height = contextPanelHeight
	m.sourceListPanel.SetSize(primaryAreaWidth-borderWidth*2, sourceListPanelHeight)
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

	m.contentPanel.Width = m.panelSize.itemPreviewPanelWidth - borderWidth*2
	m.contentPanel.Height = m.winSize.height - instantPromptPanelOuterHeight - borderHeight*2 - footerHeight

	m.configContentPanel.Width = secondlyAreaWidth - borderWidth*2
	m.configContentPanel.Height = m.winSize.height - borderHeight*2 - footerHeight

	m.reviewPanel.Width = m.panelSize.itemReviewPanelWidth - borderWidth*2
	m.reviewPanel.Height = m.winSize.height - instantPromptPanelOuterHeight - borderHeight*2 - footerHeight

	m.instantPromptPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.instantPromptPanel.SetHeight(instantPromptPanelHeight)

	m.stateDetailPanel.Width = secondlyAreaWidth - borderWidth*2
	m.stateDetailPanel.Height = m.winSize.height - borderHeight*2 - footerHeight

	m.sourceDetailPanel.Width = secondlyAreaWidth - borderWidth*2
	m.sourceDetailPanel.Height = m.winSize.height - borderHeight*2 - footerHeight
}

func (m *model) buildPrimaryPanels(statePanel, listPanel, reviewStackPanel, contextPanel, sourceListPanel, configPanel string) string {
	return lipgloss.JoinVertical(lipgloss.Top, statePanel, listPanel, reviewStackPanel, contextPanel, sourceListPanel, configPanel)
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
	if state == ListPanelFocus || state == ConfigSummaryPanelFocus || state == StatePanelFocus || state == ContextPanelFocus || state == ReviewStackPanelFocus || state == SourceListPanelFocus {
		return true
	}
	return false
}

func isFocusItemPreviewPanel(state FocusState) bool {
	return state == ContentPanelFocus
}
