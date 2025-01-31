package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/config"
	"github.com/shutils/lazyreview/pkg/openai"
	state "github.com/shutils/lazyreview/pkg/state"
)

type panelSize struct {
	secondlyPanelWidth,
	itemPreviewPanelWidth, itemReviewPanelWidth int
}

type winSize struct {
	height, width int
}

type listItem struct {
	title, param, sourceName, id string
}

type sourceItem struct {
	name      string
	collector []string
	previewer []string
	enabled   bool
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.param }
func (i listItem) FilterValue() string { return i.title }

func (i sourceItem) Title() string {
	if i.enabled {
		return "☑ " + i.name
	}
	return "☐ " + i.name
}
func (i sourceItem) Description() string {
	if i.enabled {
		return "☑ collector: " + strings.Join(i.collector, ", ") + " previewer: " + strings.Join(i.previewer, ", ")
	}
	return "☐ collector: " + strings.Join(i.collector, ", ") + " previewer: " + strings.Join(i.previewer, ", ")
}
func (i sourceItem) FilterValue() string { return i.name }

type updateSourceListMsg struct {
}

type model struct {
	panels                 panels
	keyMaps                keyMaps
	panelSize              panelSize
	winSize                winSize
	reviewList             []reviewInfo
	targetDir              string
	outputFile             string
	stateFile              string
	conf                   config.Config
	client                 openai.Client
	zoomState              ZoomState
	focusState             FocusState
	reviewState            ReviewState
	reviewStack            []int
	reviewStackDenominator int // reviewStackが0になるまでにたまったreviewの数
	instantPrompt          string
	uiState                state.State
	currentHistoryIndex    int
	state                  state.State
	message                string
	initialized            bool
}

func NewUi(conf config.Config, client openai.Client) model {
	m := model{
		panels:              NewPanels(),
		keyMaps:             DefaultKeyMap(),
		reviewList:          []reviewInfo{},
		targetDir:           conf.Target,
		outputFile:          conf.Output,
		stateFile:           conf.State,
		conf:                conf,
		client:              client,
		focusState:          ItemListPanelFocus,
		reviewState:         NoAction,
		reviewStack:         []int{},
		instantPrompt:       "",
		uiState:             state.LoadState(conf.State),
		currentHistoryIndex: 0,
		state:               state.State{},
		message:             "",
		initialized:         true,
	}
	m.panels.configDetailPanel.SetContent(strings.Join(conf.ToStringArray(), "\n"))
	m.panels.configSummaryPanel.SetContent("Config path: " + conf.ConfigPath)

	m.UpdateState()
	m.currentHistoryIndex = len(m.uiState.PromptHistory)
	m.loadReviews()
	m.panels.itemListPanel.SetItems(getItems(m.conf, m.reviewList))
	m.panels.sourceListPanel.SetItems(getSourceItems(m.conf.Sources))
	m.onChangeListSelectedItem()
	return m
}

func (m model) Init() tea.Cmd {
	return m.panels.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if action := m.handleKey(msg); action != nil {
			return action()
		}
	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)
	case reviewMsg:
		selectedItem := m.panels.itemListPanel.SelectedItem().(listItem)
		index := findIndex(m.panels.itemListPanel.Items(), msg.id)
		if index == -1 {
			return m, nil
		}
		item := m.panels.itemListPanel.Items()[index].(listItem)
		review := reviewInfo{
			ID:     item.id,
			Param:  item.param,
			Review: msg.content,
			State:  "finish",
		}
		if m.isReviewExist(msg.id) {
			m.reviewList[m.getReviewIndex(msg.id)] = review
		} else {
			m.reviewList = append(m.reviewList, review)
		}
		m.saveReviews()
		if selectedItem.id == msg.id {
			m.panels.itemReviewPanel.SetContent(msg.content)
			m.onChangeListSelectedItem()
		}
		cmd = func() tea.Msg {
			return reviewStackMsg{
				id:        msg.id,
				operation: Remove,
			}
		}
		m.UpdateState()
		return m, cmd
	case reviewStateMsg:
		m.reviewState = msg.state
	case reviewStackMsg:
		index := findIndex(m.panels.itemListPanel.Items(), msg.id)
		if index == -1 {
			return m, nil
		}
		if msg.operation == Add {
			m.addReviewStack(index)
			m.reviewStackDenominator++
		} else {
			m.removeReviewStack(index)
			m.changeItemTitlePrefix(index, "☑ ")
		}
		m.updateReviewStackPanel()
		if len(m.reviewStack) == 0 {
			cmd := m.panels.reviewProgressPanel.SetPercent(1)
			m.reviewStackDenominator = 0
			return m, cmd
		}
		cmd := m.updateReviewProgressPanel()
		return m, cmd
	case updateSourceListMsg:
		m.panels.sourceListPanel.SetItems(getSourceItems(m.conf.Sources))
		m.panels.contextListPanel.Update(msg)
		m.panels.itemListPanel.SetItems(getItems(m.conf, m.reviewList))
	case progress.FrameMsg:
		progressModel, cmd := m.panels.reviewProgressPanel.Update(msg)
		m.panels.reviewProgressPanel = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	case showMessageMsg:
		if msg.message != "" {
			m.message = msg.message + "\n\n" + "Press enter to return..."
			m.focusState = MessagePanelFocus
		}
		return m, nil
	case setPromptMsg:
		m.panels.promptPanel.SetValue(msg.text)
		return m, nil
	case closedEditorMsg:
		if msg.err != nil {
			return m, func() tea.Msg {
				return showMessageMsg{
					message: fmt.Sprintf("err: \n\n%s", msg.err.Error()),
				}
			}
		}
	default:
		switch m.focusState {
		case SourceListPanelFocus:
			selectedSourceName := m.panels.sourceListPanel.SelectedItem().(sourceItem)
			selectedSource := m.conf.GetSourceFromName(selectedSourceName.name)
			m.panels.contextDetailPanel.SetContent(selectedSource.String())
		case ContextPanelFocus:
			m.panels.contextDetailPanel.SetContent(m.getContextString())
		}
		m.panels.spinner, cmd = m.panels.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == ItemListPanelFocus {
		m.panels.itemListPanel, cmd = m.panels.itemListPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == SourceListPanelFocus {
		m.panels.sourceListPanel, cmd = m.panels.sourceListPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == ContextPanelFocus {
		m.panels.contextListPanel, cmd = m.panels.contextListPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == InstantPromptPanelFocus {
		m.panels.promptPanel, cmd = m.panels.promptPanel.Update(msg)
		cmds = append(cmds, cmd)
		m.instantPrompt = m.panels.promptPanel.Value()
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.makeView()
}

func (m *model) UpdateState() (tea.Model, tea.Cmd) {
	m.state = state.LoadState(m.stateFile)
	m.panels.stateSummaryPanel.SetContent(m.state.ShowUsage(m.conf.ModelCost))
	m.panels.stateDetailPanel.SetContent(m.state.ShowUsedToken())
	return m, nil
}

func (m *model) isReviewExist(id string) bool {
	return m.getReviewIndex(id) != -1
}

func (m *model) addReviewStack(index int) {
	m.reviewStack = append(m.reviewStack, index)
	m.reviewState = Reviewing
}

func (m *model) removeReviewStack(index int) {
	for i, v := range m.reviewStack {
		if v == index {
			m.reviewStack = append(m.reviewStack[:i], m.reviewStack[i+1:]...)
			break
		}
	}
	if len(m.reviewStack) == 0 {
		m.reviewState = NoAction
	}
}

func (m *model) updateReviewStackPanel() {
	itemTitleList := getItemListString(getReviewStackItems(m.panels.itemListPanel.Items(), m.reviewStack))
	m.panels.reviewStackPanel.SetContent(itemTitleList)
}

func (m *model) updateReviewProgressPanel() tea.Cmd {
	percent := float64(1)
	if m.reviewStackDenominator != 0 {
		percent = 1 - float64(len(m.reviewStack))/float64(m.reviewStackDenominator)
	}
	return m.panels.reviewProgressPanel.SetPercent(percent)
}

func (m *model) addContextStack(id string) (tea.Model, tea.Cmd) {
	index := findIndex(m.panels.itemListPanel.Items(), id)
	if index == -1 {
		return *m, nil
	}
	item := m.panels.itemListPanel.Items()[index]

	contextList := m.panels.contextListPanel.Items()
	contextList = append(contextList, item)
	m.panels.contextListPanel.SetItems(contextList)
	return *m, nil
}

func (m *model) removeContextStack(id string) (tea.Model, tea.Cmd) {
	index := findIndex(m.panels.contextListPanel.Items(), id)
	contextList := m.panels.contextListPanel.Items()

	if index < 0 || index >= len(contextList) {
		return *m, nil
	}

	var newContextList []list.Item
	if index == len(contextList)-1 {
		newContextList = contextList[:index]
	} else {
		newContextList = append(contextList[:index], contextList[index+1:]...)
	}

	m.panels.contextListPanel.SetItems(newContextList)
	return *m, nil
}

func (m *model) changeItemTitlePrefix(index int, prefix string) {
	item := m.panels.itemListPanel.Items()[index].(listItem)
	item.title = replacePrefix(item.title, prefix)
	m.panels.itemListPanel.SetItem(index, item)
}
