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

var (
	winWidth, winHeight,
	primaryPanelWidth, secondlyPanelWidth,
	primaryPanelHeight, secondlyPanelHeight,
	itemPreviewPanelWidth, itemReviewPanelWidth,
	listPanelHeight, configPanelHeight, itemPreviewPanelHeight, itemReviewPanelHeight, instantPromptPanelHeight,
	statePanelHeight, contextPanelHeight, reviewStackPanelHeight int
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
	winWidth = msg.Width
	winHeight = msg.Height
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

	primaryPanelHeight = winHeight - 6
	secondlyPanelHeight = winHeight - 3
	listPanelHeight = winHeight - 23
	configPanelHeight = 1
	statePanelHeight = 1
	contextPanelHeight = 5
	reviewStackPanelHeight = 5
	instantPromptPanelHeight = 5
	itemPreviewPanelHeight = winHeight - 10
	itemReviewPanelHeight = winHeight - 10

	switch m.zoomState {
	case Normal:
		if isFocusPrimary(m.focusState) {
			primaryPanelWidth = winWidth/3 - 2
			secondlyPanelWidth = winWidth/3*2 - 2
			itemPreviewPanelWidth = winWidth/3 - 2
			itemReviewPanelWidth = winWidth/3 - 2
		} else {
			primaryPanelWidth = winWidth / 10 * 2
			secondlyPanelWidth = winWidth/10*8 + 2
			itemPreviewPanelWidth = winWidth / 10 * 4
			itemReviewPanelWidth = winWidth / 10 * 4
		}
	case Middle:
		if isFocusPrimary(m.focusState) {
			primaryPanelWidth = winWidth / 2
			secondlyPanelWidth = winWidth/2 - 4
			itemPreviewPanelWidth = winWidth/4 - 3
			itemReviewPanelWidth = winWidth/4 - 3
		} else {
			secondlyPanelWidth = winWidth - 2
			itemPreviewPanelWidth = winWidth/2 - 2
			itemReviewPanelWidth = winWidth/2 - 2
		}
	case Max:
		if isFocusPrimary(m.focusState) {
			primaryPanelWidth = winWidth - 2
		} else {
			secondlyPanelWidth = winWidth - 2
			if isFocusItemPreviewPanel(m.focusState) {
				itemPreviewPanelWidth = winWidth - 2
			} else {
				itemReviewPanelWidth = winWidth - 2
			}
		}
	}

	m.list.SetSize(primaryPanelWidth, listPanelHeight)

	m.configSummaryPanel.Width = primaryPanelWidth
	m.configSummaryPanel.Height = configPanelHeight
	m.contentPanel.Width = itemPreviewPanelWidth
	m.contentPanel.Height = itemPreviewPanelHeight
	m.configContentPanel.Width = secondlyPanelWidth
	m.configContentPanel.Height = secondlyPanelHeight
	m.reviewPanel.Width = itemReviewPanelWidth
	m.reviewPanel.Height = itemReviewPanelHeight
	m.reviewStackPanel.Width = primaryPanelWidth
	m.reviewStackPanel.Height = reviewStackPanelHeight
	m.instantPromptPanel.SetWidth(secondlyPanelWidth)
	m.instantPromptPanel.SetHeight(instantPromptPanelHeight)
	m.statePanel.Width = primaryPanelWidth
	m.statePanel.Height = statePanelHeight
	m.stateDetailPanel.Width = secondlyPanelWidth
	m.stateDetailPanel.Height = secondlyPanelHeight
	m.contextPanel.Width = primaryPanelWidth
	m.contextPanel.Height = contextPanelHeight

	state := "/"

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
	instantPromptPanel := m.buildPanel(m.instantPromptPanel.View(), instantPromptPanelStyle, secondlyPanelWidth, "Instant prompt")
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
