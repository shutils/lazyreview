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
	itemPreviewPanel    simpleViewPortPanel
	itemReviewPanel     simpleViewPortPanel
	stateSummaryPanel   simpleViewPortPanel
	stateDetailPanel    simpleViewPortPanel
	configSummaryPanel  simpleViewPortPanel
	configDetailPanel   simpleViewPortPanel
	reviewProgressPanel reviewProgressPanel
	reviewStackPanel    simpleViewPortPanel
	sourceListPanel     compactListPanel
	sourceDetailPanel   simpleViewPortPanel
	contextListPanel    compactListPanel
	contextDetailPanel  simpleViewPortPanel
	promptPanel         promptPanel
	spinner             spinner.Model
	messagePanel        simpleViewPortPanel
}

func NewPanels() panels {
	p := panels{
		itemListPanel:       NewItemListPanel("Items"),
		itemPreviewPanel:    NewSimpleViewPort("Content"),
		itemReviewPanel:     NewSimpleViewPort("Review"),
		stateSummaryPanel:   NewSimpleViewPort("State"),
		stateDetailPanel:    NewSimpleViewPort("State detail"),
		configSummaryPanel:  NewSimpleViewPort("Config"),
		configDetailPanel:   NewSimpleViewPort("Config content"),
		reviewProgressPanel: NewReviewProgress("Review progress"),
		reviewStackPanel:    NewSimpleViewPort("Review stack"),
		sourceListPanel:     NewCompactListPanel("Source list"),
		sourceDetailPanel:   NewSimpleViewPort("Source detail"),
		contextListPanel:    NewCompactListPanel("Context"),
		contextDetailPanel:  NewSimpleViewPort("Context detail"),
		promptPanel:         NewPromptPanel("Instant prompt"),
		spinner:             spinner.New(),
		messagePanel:        NewSimpleViewPort("Message"),
	}

	p.setInitSetting()
	return p
}

func (p *panels) setInitSetting() {
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
	stateSummaryPanelHeight     = 1
	configSummaryPanelHeight    = 1
	reviewProgressPanelHeight   = 1
	contextListPanelMaxHeight   = 5
	instantPromptPanelMinHeight = 1
	sourceListPanelMaxHeight    = 5
	footerHeight                = 1

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
	m.resetHighLightPanel()
	m.setHighLightPanel()
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

	listPanel := m.panels.itemListPanel.View()
	contentPanel := m.panels.itemPreviewPanel.View()
	reviewPanel := m.panels.itemReviewPanel.View()
	reviewStackPanel := m.panels.reviewStackPanel.View()
	configPanel := m.panels.configSummaryPanel.View()
	configContentPanel := m.panels.configDetailPanel.View()
	statePanel := m.panels.stateSummaryPanel.View()
	stateDetailPanel := m.panels.stateDetailPanel.View()
	instantPromptPanel := m.panels.promptPanel.View()
	contextPanel := m.panels.contextListPanel.View()
	sourceListPanel := m.panels.sourceListPanel.View()
	sourceDetailPanel := m.panels.sourceDetailPanel.View()
	contextDetailPanel := m.panels.contextDetailPanel.View()
	reviewProgressPanel := m.panels.reviewProgressPanel.View()

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

func (m *model) calcInstantPromptPanelHeight() {
	height := math.Min(float64(m.winSize.height/2), float64(m.panels.promptPanel.LineCount()))
	m.panels.promptPanel.SetHeight(int(height))
}

func (m *model) setPanelSize() (tea.Model, tea.Cmd) {
	m.setPrimaryPanelSizes()
	m.setSecondaryPanelSizes()

	m.panels.messagePanel.SetWidth(m.winSize.width)
	m.panels.messagePanel.SetHeight(m.winSize.height)
	return m, nil
}

func (m *model) setHighLightPanel() {
	switch m.focusState {
	case ItemListPanelFocus:
		m.panels.itemListPanel.SetHighlight(true)
	case ContentPanelFocus:
		m.panels.itemPreviewPanel.SetHighlight(true)
	case ReviewPanelFocus:
		m.panels.itemReviewPanel.SetHighlight(true)
	case ReviewStackProgressPanelFocus:
		m.panels.reviewProgressPanel.SetHighlight(true)
	case ConfigSummaryPanelFocus:
		m.panels.configSummaryPanel.SetHighlight(true)
	case StatePanelFocus:
		m.panels.stateSummaryPanel.SetHighlight(true)
	case InstantPromptPanelFocus:
		m.panels.promptPanel.SetHighlight(true)
	case SourceListPanelFocus:
		m.panels.sourceListPanel.SetHighlight(true)
	case ContextPanelFocus:
		m.panels.contextListPanel.SetHighlight(true)
	}
}

func (m *model) resetHighLightPanel() {
	m.panels.reviewStackPanel.SetHighlight(false)
	m.panels.configSummaryPanel.SetHighlight(false)
	m.panels.stateSummaryPanel.SetHighlight(false)
	m.panels.itemPreviewPanel.SetHighlight(false)
	m.panels.itemReviewPanel.SetHighlight(false)
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

	m.panels.itemListPanel.SetWidth(primaryAreaWidth - borderWidth*2)
	m.panels.itemListPanel.SetHeight(listPanelHeight)
	m.panels.configSummaryPanel.SetWidth(primaryAreaWidth - borderWidth*2)
	m.panels.configSummaryPanel.SetHeight(configSummaryPanelHeight)
	m.panels.reviewProgressPanel.SetWidth(primaryAreaWidth - borderWidth*2)
	m.panels.stateSummaryPanel.SetWidth(primaryAreaWidth - borderWidth*2)
	m.panels.stateSummaryPanel.SetHeight(stateSummaryPanelHeight)
	m.panels.contextListPanel.SetWidth(primaryAreaWidth - borderWidth*2)
	m.panels.contextListPanel.SetHeight(contextListPanelHeight)
	m.panels.sourceListPanel.SetWidth(primaryAreaWidth - borderWidth*2)
	m.panels.sourceListPanel.SetHeight(sourceListPanelHeight)
}

func (m *model) setSecondaryPanelSizes() {
	m.calcInstantPromptPanelHeight()
	var (
		instantPromptPanelOuterHeight = m.panels.promptPanel.Height() + borderHeight*2
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

	m.panels.itemPreviewPanel.SetWidth(m.panelSize.itemPreviewPanelWidth - borderWidth*2)
	m.panels.itemPreviewPanel.SetHeight(m.winSize.height - instantPromptPanelOuterHeight - borderHeight*2 - footerHeight)

	m.panels.itemReviewPanel.SetWidth(m.panelSize.itemReviewPanelWidth - borderWidth*2)
	m.panels.itemReviewPanel.SetHeight(m.winSize.height - instantPromptPanelOuterHeight - borderHeight*2 - footerHeight)

	m.panels.configDetailPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.configDetailPanel.SetHeight(m.winSize.height - borderHeight*2 - footerHeight)

	m.panels.stateDetailPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.stateDetailPanel.SetHeight(m.winSize.height - borderHeight*2 - footerHeight)

	m.panels.promptPanel.SetWidth(secondlyAreaWidth - borderWidth*2)

	m.panels.stateDetailPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.stateDetailPanel.SetHeight(m.winSize.height - borderHeight*2 - footerHeight)

	m.panels.sourceDetailPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.sourceDetailPanel.SetHeight(m.winSize.height - borderHeight*2 - footerHeight)

	m.panels.contextDetailPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.contextDetailPanel.SetHeight(m.winSize.height - borderHeight*2 - footerHeight)

	m.panels.reviewStackPanel.SetWidth(secondlyAreaWidth - borderWidth*2)
	m.panels.reviewStackPanel.SetHeight(m.winSize.height - borderHeight*2 - footerHeight)
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
	width, height   int
	title           string
	model           list.Model
	showDescription bool
	normalDelegate  list.DefaultDelegate
	narrowDelegate  list.DefaultDelegate
	isHighlight     bool
}

func NewItemListPanel(title string) itemListPanel {
	l := itemListPanel{}
	l.model = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.showDescription = true
	l.normalDelegate = l.NewDefaultNormalDelegate()
	l.narrowDelegate = l.NewDefaultNarrowDelegate()
	l.title = title
	l.model.SetShowHelp(false)
	l.model.SetShowTitle(false)
	l.model.KeyMap.Quit.Unbind()
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

func (l *itemListPanel) SetWidth(w int) {
	l.width = w
	l.model.SetWidth(w)
}

func (l *itemListPanel) SetHeight(h int) {
	l.height = h
	l.model.SetHeight(h)
}

func (l *itemListPanel) SetHighlight(highlight bool) {
	l.isHighlight = highlight
}

func (l *itemListPanel) Width() int {
	return l.width
}

func (l *itemListPanel) Height() int {
	return l.height
}

func (l *itemListPanel) Init() tea.Cmd {
	return nil
}

func (l *itemListPanel) Update(msg tea.Msg) (itemListPanel, tea.Cmd) {

	var cmd tea.Cmd
	l.model, cmd = l.model.Update(msg)
	return *l, cmd
}

func (l *itemListPanel) View() string {
	style := baseStyle
	if l.isHighlight {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return InsertTitleWithOffset(style.Width(l.width).Height(l.height).Render(l.model.View()), l.title)
}

type promptPanel struct {
	width, height int
	MaxHeight     int
	model         textarea.Model
	isHighlight   bool
	title         string
}

func NewPromptPanel(title string) promptPanel {
	p := promptPanel{}
	p.model = textarea.New()
	p.model.SetHeight(instantPromptPanelMinHeight)
	p.height = instantPromptPanelMinHeight
	return p
}

func (p *promptPanel) UpdateSize() (promptPanel, tea.Cmd) {
	p.SetWidth(p.width)
	p.SetHeight(p.height)
	return *p, nil
}

func (p *promptPanel) SetWidth(w int) {
	p.width = w
	p.model.SetWidth(w)
}

func (p *promptPanel) SetHeight(h int) {
	p.height = h
	p.model.SetHeight(p.height)
}

func (p *promptPanel) SetHighlight(highlight bool) {
	p.isHighlight = highlight
}

func (p *promptPanel) Width() int {
	return p.width
}

func (p *promptPanel) Height() int {
	return p.height
}

func (p *promptPanel) SetValue(content string) {
	p.model.SetValue(content)
	p.SetHeight(p.model.LineCount())
}

func (p *promptPanel) LineCount() int {
	return p.model.LineCount()
}

func (p *promptPanel) Focus() {
	p.model.Focus()
}

func (p *promptPanel) Blur() {
	p.model.Blur()
}

func (p *promptPanel) Value() string {
	return p.model.Value()
}

func (p *promptPanel) Init() tea.Cmd {
	return nil
}

func (p *promptPanel) Update(msg tea.Msg) (promptPanel, tea.Cmd) {
	var cmd tea.Cmd
	p.model, cmd = p.model.Update(msg)
	p.SetHeight(p.model.LineCount())
	p.model.SetHeight(p.height)
	return *p, cmd
}

func (p *promptPanel) View() string {
	style := baseStyle
	if p.isHighlight {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return InsertTitleWithOffset(style.Width(p.width).Height(p.height).Render(p.model.View()), p.title)
}

type simpleViewPortPanel struct {
	width, height int
	title         string
	model         viewport.Model

	isHighlight bool
}

func NewSimpleViewPort(title string) simpleViewPortPanel {
	v := simpleViewPortPanel{}
	v.model = viewport.New(0, 0)
	v.title = title
	return v
}

func (v *simpleViewPortPanel) SetWidth(w int) {
	v.width = w
	v.model.Width = w
}

func (v *simpleViewPortPanel) SetHeight(h int) {
	v.height = h
	v.model.Height = h
}

func (v *simpleViewPortPanel) SetHighlight(highlight bool) {
	v.isHighlight = highlight
}

func (v *simpleViewPortPanel) Width() int {
	return v.width
}

func (v *simpleViewPortPanel) Height() int {
	return v.height
}

func (v *simpleViewPortPanel) SetContent(content string) {
	v.model.SetContent(content)
}

func (v *simpleViewPortPanel) LineDown(n int) {
	v.model.LineDown(n)
}

func (v *simpleViewPortPanel) LineUp(n int) {
	v.model.LineUp(n)
}

func (v *simpleViewPortPanel) HalfViewDown() {
	v.model.HalfViewDown()
}

func (v *simpleViewPortPanel) HalfViewUp() {
	v.model.HalfViewUp()
}

func (v *simpleViewPortPanel) GotoTop() {
	v.model.GotoTop()
}

func (v *simpleViewPortPanel) GotoBottom() {
	v.model.GotoBottom()
}

func (v *simpleViewPortPanel) Init() tea.Cmd {
	return nil
}

func (v *simpleViewPortPanel) Update(msg tea.Msg) (simpleViewPortPanel, tea.Cmd) {
	var cmd tea.Cmd
	v.model, cmd = v.model.Update(msg)
	return *v, cmd
}

func (v *simpleViewPortPanel) View() string {
	style := baseStyle
	if v.isHighlight {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return InsertTitleWithOffset(style.Width(v.width).Height(v.height).Render(v.model.View()), v.title)
}

type reviewProgressPanel struct {
	title       string
	width       int
	model       progress.Model
	isHighlight bool
}

func NewReviewProgress(title string) reviewProgressPanel {
	r := reviewProgressPanel{}
	r.model = progress.New()
	r.title = title
	return r
}

func (r *reviewProgressPanel) SetWidth(w int) {
	r.width = w
	r.model.Width = w
}

func (r *reviewProgressPanel) SetHighlight(highlight bool) {
	r.isHighlight = highlight
}

func (r *reviewProgressPanel) SetPercent(percent float64) tea.Cmd {
	return r.model.SetPercent(percent)
}

func (r *reviewProgressPanel) Init() tea.Cmd {
	return nil
}

func (r *reviewProgressPanel) Update(msg tea.Msg) (reviewProgressPanel, tea.Cmd) {
	var cmd tea.Cmd
	progressModel, cmd := r.model.Update(msg)
	m, ok := progressModel.(progress.Model)
	if ok {
		r.model = m
		return *r, cmd
	}
	return *r, cmd
}

func (r *reviewProgressPanel) View() string {
	style := baseStyle
	if r.isHighlight {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return InsertTitleWithOffset(style.Width(r.width).Height(1).Render(r.model.View()), r.title)
}

type compactListPanel struct {
	width, height int
	title         string
	model         list.Model
	isHighlight   bool
}

func NewCompactListPanel(title string) compactListPanel {
	l := compactListPanel{}
	l.model = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.title = title
	l.model.SetDelegate(list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	})
	l.model.SetShowTitle(false)
	l.model.SetShowHelp(false)
	l.model.SetShowStatusBar(false)
	l.model.SetShowFilter(false)
	l.model.KeyMap.Quit.Unbind()
	l.model.KeyMap.Filter.Unbind()
	return l
}

func (l *compactListPanel) SetWidth(w int) {
	l.width = w
	l.model.SetWidth(w)
}

func (l *compactListPanel) SetHeight(h int) {
	l.height = h
	l.model.SetHeight(h)
}

func (l *compactListPanel) SetHighlight(highlight bool) {
	l.isHighlight = highlight
}

func (l *compactListPanel) Width() int {
	return l.width
}

func (l *compactListPanel) Height() int {
	return l.height
}

func (l *compactListPanel) Items() []list.Item {
	return l.model.Items()
}

func (l *compactListPanel) SelectedItem() list.Item {
	return l.model.SelectedItem()
}

func (l *compactListPanel) SetShowPagination(show bool) {
	l.model.SetShowPagination(show)
}

func (l *compactListPanel) SetItems(items []list.Item) {
	l.model.SetItems(items)
}

func (l *compactListPanel) Init() tea.Cmd {
	return nil
}

func (l *compactListPanel) Update(msg tea.Msg) (compactListPanel, tea.Cmd) {
	var cmd tea.Cmd
	l.model, cmd = l.model.Update(msg)
	return *l, cmd
}

func (l *compactListPanel) View() string {
	style := baseStyle
	if l.isHighlight {
		style = style.BorderForeground(lipgloss.Color("62"))
	}
	return InsertTitleWithOffset(style.Width(l.width).Height(l.height).Render(l.model.View()), l.title)
}
