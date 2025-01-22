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
	listPanelHeight, configPanelHeight, itemPreviewPanelHeight, itemReviewPanelHeight, instantPromptPanelHeight int
)

const (
	ListPanelFocus FocusState = iota
	ContentPanelFocus
	ReviewPanelFocus
	InstantPromptPanelFocus
	ConfigSummaryPanelFocus
)

var (
	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
)

func replacePrefix(text string, prefix string) string {
	runes := []rune(text) // 文字列をruneスライスに変換
	var updatedTitle string
	if len(runes) >= 2 {
		updatedTitle = prefix + string(runes[2:]) // 先頭2文字を置換して文字列に戻す
	} else {
		updatedTitle = prefix // 文字数が2未満の場合も安全に処理
	}
	return updatedTitle
}

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
	instantPromptPanelStyle := baseStyle
	configPanelStyle := baseStyle
	configContentPanelStyle := baseStyle

	primaryPanelHeight = winHeight - 6
	secondlyPanelHeight = winHeight - 3
	listPanelHeight = winHeight - 6
	configPanelHeight = 1
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
	m.instantPromptPanel.SetWidth(secondlyPanelWidth)
	m.instantPromptPanel.SetHeight(instantPromptPanelHeight)

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
	case InstantPromptPanelFocus:
		instantPromptPanelStyle = instantPromptPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.promptKeyMap))
	case ConfigSummaryPanelFocus:
		configPanelStyle = configPanelStyle.BorderForeground(lipgloss.Color("62"))
		helpString = MakeBottomLine(globalHelp, helpModel.View(m.configSummaryKeyMap))
	}

	listPanel := listPanelStyle.Width(primaryPanelWidth).Render(m.list.View())
	listPanel = InsertTitleWithOffset(listPanel, "List")

	contentPanel := contentPanelStyle.Render(m.contentPanel.View())
	contentPanel = InsertTitleWithOffset(contentPanel, "Content")

	reviewPanel := reviewPanelStyle.Render(m.reviewPanel.View())
	reviewPanel = InsertTitleWithOffset(reviewPanel, "Review")

	configPanel := configPanelStyle.Render(m.configSummaryPanel.View())
	configPanel = InsertTitleWithOffset(configPanel, "Config")

	configContentPanel := configContentPanelStyle.Render(m.configContentPanel.View())
	configContentPanel = InsertTitleWithOffset(configContentPanel, "Config content")

	instantPromptPanel := instantPromptPanelStyle.Render(m.instantPromptPanel.View())
	instantPromptPanel = InsertTitleWithOffset(instantPromptPanel, "Instant prompt")

	bottomLine := lipgloss.JoinHorizontal(lipgloss.Left, state, " ", helpString)
	switch m.focusState {
	case ConfigSummaryPanelFocus:
		return lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, listPanel, configPanel),
				lipgloss.JoinVertical(lipgloss.Top, configContentPanel),
			),
			bottomLine,
		)
	case ListPanelFocus:
		return lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, listPanel, configPanel),
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					contentPanel,
					reviewPanel,
				),
					instantPromptPanel,
				),
			),
			bottomLine,
		)
	case ContentPanelFocus:
		switch m.zoomState {
		case Normal:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					lipgloss.JoinVertical(lipgloss.Top, listPanel, configPanel),
					lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
						contentPanel,
						reviewPanel,
					),
						instantPromptPanel,
					),
				),
				bottomLine,
			)
		case Middle:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					contentPanel,
					reviewPanel,
				),
					instantPromptPanel,
				),
				bottomLine,
			)
		case Max:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					contentPanel,
				),
					instantPromptPanel,
				),
				bottomLine,
			)
		}
	case ReviewPanelFocus:
		switch m.zoomState {
		case Normal:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					lipgloss.JoinVertical(lipgloss.Top, listPanel, configPanel),
					lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
						contentPanel,
						reviewPanel,
					),
						instantPromptPanel,
					),
				),
				bottomLine,
			)
		case Middle:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					contentPanel,
					reviewPanel,
				),
					instantPromptPanel,
				),
				bottomLine,
			)
		case Max:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					reviewPanel,
				),
					instantPromptPanel,
				),
				bottomLine,
			)
		}
	case InstantPromptPanelFocus:
		switch m.zoomState {
		case Normal:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					lipgloss.JoinVertical(lipgloss.Top, listPanel, configPanel),
					lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
						contentPanel,
						reviewPanel,
					),
						instantPromptPanel,
					),
				),
				bottomLine,
			)
		case Middle:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					contentPanel,
					reviewPanel,
				),
					instantPromptPanel,
				),
				bottomLine,
			)
		case Max:
			return lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left,
					reviewPanel,
				),
					instantPromptPanel,
				),
				bottomLine,
			)
		}
	}
	return ""
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
	if state == ListPanelFocus || state == ConfigSummaryPanelFocus {
		return true
	}
	return false
}

func isFocusItemPreviewPanel(state FocusState) bool {
	return state == ContentPanelFocus
}
