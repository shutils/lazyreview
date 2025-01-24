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
	listPanelStyle := baseStyle
	contentPanelStyle := baseStyle
	reviewPanelStyle := baseStyle
	reviewStackPanelStyle := baseStyle
	instantPromptPanelStyle := baseStyle
	configPanelStyle := baseStyle
	configContentPanelStyle := baseStyle
	statePanelStyle := baseStyle
	stateDetailPanelStyle := baseStyle
	contextPanelStyle := baseStyle
	state := "/"

	m.setPanelSize()

	if m.reviewState == Reviewing {
		state = m.spinner.View()
	}

	var helpString string

	helpModel := help.New()
	globalHelp := helpModel.View(m.globalKeyMap)

	// フォーカスされているパネルにスタイルを適用
	switch m.focusState {
	case ListPanelFocus:
		listPanelStyle = listPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.listKeyMap))
	case ContentPanelFocus:
		contentPanelStyle = contentPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.contentKeyMap))
	case ReviewPanelFocus:
		reviewPanelStyle = reviewPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.reviewKeyMap))
	case ReviewStackPanelFocus:
		reviewStackPanelStyle = reviewStackPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.reviewStackKeyMap))
	case InstantPromptPanelFocus:
		instantPromptPanelStyle = instantPromptPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.promptKeyMap))
	case ConfigSummaryPanelFocus:
		configPanelStyle = configPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.configSummaryKeyMap))
	case StatePanelFocus:
		statePanelStyle = statePanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.stateKeyMap))
	case ContextPanelFocus:
		contextPanelStyle = contextPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.contextKeyMap))
	}

	listPanel := m.buildPanel(m.list.View(), listPanelStyle, m.list.Width(), "List")
	contentPanel := m.buildPanel(m.contentPanel.View(), contentPanelStyle, m.contentPanel.Width, "Content")
	reviewPanel := m.buildPanel(m.reviewPanel.View(), reviewPanelStyle, m.reviewPanel.Width, "Review")
	reviewStackPanel := m.buildPanel(m.reviewStackPanel.View(), reviewStackPanelStyle, m.reviewStackPanel.Width, "Review stack")
	configPanel := m.buildPanel(m.configSummaryPanel.View(), configPanelStyle, m.configSummaryPanel.Width, "Config")
	configContentPanel := m.buildPanel(m.configContentPanel.View(), configContentPanelStyle, m.configContentPanel.Width, "Config content")
	statePanel := m.buildPanel(m.statePanel.View(), statePanelStyle, m.statePanel.Width, "State")
	stateDetailPanel := m.buildPanel(m.stateDetailPanel.View(), stateDetailPanelStyle, m.stateDetailPanel.Width, "State detail")
	instantPromptPanel := m.buildPanel(m.instantPromptPanel.View(), instantPromptPanelStyle, m.panelSize.secondlyPanelWidth, "Instant prompt")
	contextPanel := m.buildPanel(m.contextPanel.View(), contextPanelStyle, m.contextPanel.Width, "Context")

	primaryPanels := m.buildPrimaryPanels(statePanel, listPanel, reviewStackPanel, contextPanel, configPanel)
	reviewCombiPanels := m.buildReviewCombiPanels(contentPanel, reviewPanel, instantPromptPanel)
	contentPanels := m.buildContentPanels(contentPanel, instantPromptPanel)
	reviewPanels := m.buildReviewPanels(reviewPanel, instantPromptPanel)

	bottomLine := lipgloss.JoinHorizontal(lipgloss.Left, state, " ", helpString)
	switch m.focusState {
	case ConfigSummaryPanelFocus:
		return m.buildWindow(primaryPanels, configContentPanel, bottomLine)
	case StatePanelFocus:
		return m.buildWindow(primaryPanels, stateDetailPanel, bottomLine)
	case ListPanelFocus:
		return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
	case ReviewStackPanelFocus:
		return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
	case ContextPanelFocus:
		return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
	case ContentPanelFocus:
		switch m.zoomState {
		case Normal:
			return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
		case Middle:
			return m.buildWindowZoom(reviewCombiPanels, bottomLine)
		case Max:
			return m.buildWindowZoom(contentPanels, bottomLine)
		}
	case ReviewPanelFocus:
		switch m.zoomState {
		case Normal:
			return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
		case Middle:
			return m.buildWindowZoom(reviewCombiPanels, bottomLine)
		case Max:
			return m.buildWindowZoom(reviewPanels, bottomLine)
		}
	case InstantPromptPanelFocus:
		switch m.zoomState {
		case Normal:
			return m.buildWindow(primaryPanels, reviewCombiPanels, bottomLine)
		case Middle:
			return m.buildWindowZoom(reviewCombiPanels, bottomLine)
		case Max:
			return m.buildWindowZoom(reviewPanels, bottomLine)
		}
	}
	return ""
}

func (m *model) setPanelSize() {
	m.panelSize.primaryPanelHeight = m.winSize.height - 6
	m.panelSize.secondlyPanelHeight = m.winSize.height - 3
	m.panelSize.listPanelHeight = m.winSize.height - 23
	m.panelSize.configPanelHeight = 1
	m.panelSize.statePanelHeight = 1
	m.panelSize.contextPanelHeight = 5
	m.panelSize.reviewStackPanelHeight = 5
	m.panelSize.instantPromptPanelHeight = 5
	m.panelSize.itemPreviewPanelHeight = m.winSize.height - 10
	m.panelSize.itemReviewPanelHeight = m.winSize.height - 10

	switch m.zoomState {
	case Normal:
		if isFocusPrimary(m.focusState) {
			m.panelSize.primaryPanelWidth = m.winSize.width/3 - 2
			m.panelSize.secondlyPanelWidth = m.winSize.width/3*2 - 2
			m.panelSize.itemPreviewPanelWidth = m.winSize.width/3 - 2
			m.panelSize.itemReviewPanelWidth = m.winSize.width/3 - 2
		} else {
			m.panelSize.primaryPanelWidth = m.winSize.width / 10 * 2
			m.panelSize.secondlyPanelWidth = m.winSize.width/10*8 + 2
			m.panelSize.itemPreviewPanelWidth = m.winSize.width / 10 * 4
			m.panelSize.itemReviewPanelWidth = m.winSize.width / 10 * 4
		}
	case Middle:
		if isFocusPrimary(m.focusState) {
			m.panelSize.primaryPanelWidth = m.winSize.width / 2
			m.panelSize.secondlyPanelWidth = m.winSize.width/2 - 4
			m.panelSize.itemPreviewPanelWidth = m.winSize.width/4 - 3
			m.panelSize.itemReviewPanelWidth = m.winSize.width/4 - 3
		} else {
			m.panelSize.secondlyPanelWidth = m.winSize.width - 2
			m.panelSize.itemPreviewPanelWidth = m.winSize.width/2 - 2
			m.panelSize.itemReviewPanelWidth = m.winSize.width/2 - 2
		}
	case Max:
		if isFocusPrimary(m.focusState) {
			m.panelSize.primaryPanelWidth = m.winSize.width - 2
		} else {
			m.panelSize.secondlyPanelWidth = m.winSize.width - 2
			if isFocusItemPreviewPanel(m.focusState) {
				m.panelSize.itemPreviewPanelWidth = m.winSize.width - 2
			} else {
				m.panelSize.itemReviewPanelWidth = m.winSize.width - 2
			}
		}
	}

	m.list.SetSize(m.panelSize.primaryPanelWidth, m.panelSize.listPanelHeight)

	m.configSummaryPanel.Width = m.panelSize.primaryPanelWidth
	m.configSummaryPanel.Height = m.panelSize.configPanelHeight
	m.contentPanel.Width = m.panelSize.itemPreviewPanelWidth
	m.contentPanel.Height = m.panelSize.itemPreviewPanelHeight
	m.configContentPanel.Width = m.panelSize.secondlyPanelWidth
	m.configContentPanel.Height = m.panelSize.secondlyPanelHeight
	m.reviewPanel.Width = m.panelSize.itemReviewPanelWidth
	m.reviewPanel.Height = m.panelSize.itemReviewPanelHeight
	m.reviewStackPanel.Width = m.panelSize.primaryPanelWidth
	m.reviewStackPanel.Height = m.panelSize.reviewStackPanelHeight
	m.instantPromptPanel.SetWidth(m.panelSize.secondlyPanelWidth)
	m.instantPromptPanel.SetHeight(m.panelSize.instantPromptPanelHeight)
	m.statePanel.Width = m.panelSize.primaryPanelWidth
	m.statePanel.Height = m.panelSize.statePanelHeight
	m.stateDetailPanel.Width = m.panelSize.secondlyPanelWidth
	m.stateDetailPanel.Height = m.panelSize.secondlyPanelHeight
	m.contextPanel.Width = m.panelSize.primaryPanelWidth
	m.contextPanel.Height = m.panelSize.contextPanelHeight
}

func (m *model) buildPrimaryPanels(statePanel, listPanel, reviewStackPanel, contextPanel, configPanel string) string {
	return lipgloss.JoinVertical(lipgloss.Top, statePanel, listPanel, reviewStackPanel, contextPanel, configPanel)
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

func (m *model) buildPanel(p string, s lipgloss.Style, w int, title string) string {
	return InsertTitleWithOffset(s.Width(w).Render(p), title)
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
	if state == ListPanelFocus || state == ConfigSummaryPanelFocus || state == StatePanelFocus || state == ContextPanelFocus || state == ReviewStackPanelFocus {
		return true
	}
	return false
}

func isFocusItemPreviewPanel(state FocusState) bool {
	return state == ContentPanelFocus
}
