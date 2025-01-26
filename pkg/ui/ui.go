package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/config"
	"github.com/shutils/lazyreview/pkg/openai"
	state "github.com/shutils/lazyreview/pkg/state"
)

type panelSize struct {
	primaryPanelWidth, secondlyPanelWidth,
	primaryPanelHeight, secondlyPanelHeight,
	itemPreviewPanelWidth, itemReviewPanelWidth,
	listPanelHeight, configPanelHeight, itemPreviewPanelHeight, itemReviewPanelHeight, instantPromptPanelHeight,
	statePanelHeight, contextPanelHeight, reviewStackPanelHeight int
}

type winSize struct {
	height, width int
}

type listItem struct {
	title, param, sourceName string
	aiContext                bool
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.param }
func (i listItem) FilterValue() string { return i.title }

type model struct {
	list                list.Model
	contentPanel        viewport.Model
	reviewPanel         viewport.Model
	reviewStackPanel    viewport.Model
	instantPromptPanel  textarea.Model
	configSummaryPanel  viewport.Model
	configContentPanel  viewport.Model
	statePanel          viewport.Model
	stateDetailPanel    viewport.Model
	contextPanel        viewport.Model
	panelSize           panelSize
	winSize             winSize
	reviewList          []reviewInfo
	targetDir           string
	outputFile          string
	stateFile           string
	conf                config.Config
	client              openai.Client
	zoomState           ZoomState
	focusState          FocusState
	reviewState         ReviewState
	reviewStack         []int
	spinner             spinner.Model
	instantPrompt       string
	globalKeyMap        globalKeyMap
	listKeyMap          listKeyMap
	contentKeyMap       contentKeyMap
	reviewKeyMap        reviewKeyMap
	reviewStackKeyMap   reviewStackKeyMap
	promptKeyMap        promptKeyMap
	configSummaryKeyMap configSummaryKeyMap
	stateKeyMap         stateKeyMap
	contextKeyMap       contextKeyMap
	uiState             state.State
	currentHistoryIndex int
	state               state.State
}

func NewUi(conf config.Config, client openai.Client) model {
	itemList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0) // listの高さを設定
	itemList.SetShowHelp(false)
	itemList.KeyMap.Quit.Unbind()
	itemList.SetShowTitle(false)
	contentPanel := viewport.New(0, 20)
	reviewPanel := viewport.New(0, 20)
	instantPromptPanel := textarea.New()
	configPanel := viewport.New(0, 20)
	configPanel.SetContent("Config path: " + conf.ConfigPath)
	configContentPanel := viewport.New(0, 20)
	configContentPanel.SetContent(strings.Join(conf.ToStringArray(), "\n"))
	m := model{
		list:                itemList,
		contentPanel:        contentPanel,
		reviewPanel:         reviewPanel,
		instantPromptPanel:  instantPromptPanel,
		configSummaryPanel:  configPanel,
		configContentPanel:  configContentPanel,
		statePanel:          viewport.New(0, 0),
		stateDetailPanel:    viewport.New(0, 0),
		contextPanel:        viewport.New(0, 0),
		reviewList:          []reviewInfo{},
		targetDir:           conf.Target,
		outputFile:          conf.Output,
		stateFile:           conf.State,
		conf:                conf,
		client:              client,
		focusState:          ListPanelFocus,
		reviewState:         NoAction,
		reviewStack:         []int{},
		spinner:             spinner.New(),
		instantPrompt:       "",
		globalKeyMap:        GetGlobalKeymap(),
		listKeyMap:          GetListKeymap(),
		contentKeyMap:       GetContentKeymap(),
		reviewKeyMap:        GetReviewKeymap(),
		reviewStackKeyMap:   GetReviewStackKeymap(),
		promptKeyMap:        GetPromptKeymap(),
		configSummaryKeyMap: GetConfigSummaryKeymap(),
		stateKeyMap:         GetStateKeymap(),
		contextKeyMap:       GetContextKeymap(),
		uiState:             state.LoadState(conf.State),
		currentHistoryIndex: 0,
		state:               state.State{},
	}
	m.UpdateState()
	m.currentHistoryIndex = len(m.uiState.PromptHistory)
	m.loadReviews()
	m.list.SetItems(getItems(m.conf, m.reviewList))
	m.onChangeListSelectedItem()
	return m
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
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
		selectedItem := m.list.SelectedItem().(listItem)
		index := findIndex(m.list.Items(), msg.param)
		item := m.list.Items()[index].(listItem)
		review := reviewInfo{
			Param:  item.param,
			Review: msg.content,
			State:  "finish",
		}
		if m.isReviewExist(msg.param) {
			m.reviewList[m.getReviewIndex(msg.param)] = review
		} else {
			m.reviewList = append(m.reviewList, review)
		}
		m.saveReviews()
		if selectedItem.param == msg.param {
			m.reviewPanel.SetContent(msg.content)
			m.onChangeListSelectedItem()
		}
		cmd = func() tea.Msg {
			return reviewStackMsg{
				param:     msg.param,
				operation: Remove,
			}
		}
		m.UpdateState()
		return m, cmd
	case reviewStateMsg:
		m.reviewState = msg.state
	case reviewStackMsg:
		index := findIndex(m.list.Items(), msg.param)
		if msg.operation == Add {
			m.addReviewStack(index)
		} else {
			m.removeReviewStack(index)
			m.changeItemTitlePrefix(index, "☑ ")
		}
		m.updateReviewStackPanel()
		return m, nil
	case aiContextMsg:
		index := findIndex(m.list.Items(), msg.itemParam)
		if msg.method == AddContext {
			m.addContextStack(index)
		} else {
			m.removeContextStack(index)
		}
		m.updateContextPanel()
		return m, nil
	default:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == ListPanelFocus {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == InstantPromptPanelFocus {
		m.instantPromptPanel, cmd = m.instantPromptPanel.Update(msg)
		cmds = append(cmds, cmd)
		m.instantPrompt = m.instantPromptPanel.Value()
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.makeView()
}

func (m *model) UpdateState() (tea.Model, tea.Cmd) {
	m.state = state.LoadState(m.stateFile)
	m.statePanel.SetContent(m.state.ShowUsage(m.conf.ModelCost))
	m.stateDetailPanel.SetContent(m.state.ShowUsedToken())
	return m, nil
}

func (m *model) isReviewExist(itemParam string) bool {
	return m.getReviewIndex(itemParam) != -1
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
	itemTitleList := getItemListString(getReviewStackItems(m.list.Items(), m.reviewStack))
	m.reviewStackPanel.SetContent(itemTitleList)
}

func (m *model) updateContextPanel() {
	itemTitleList := getItemListString(getContextItems(m.list.Items()))
	m.contextPanel.SetContent(itemTitleList)
}

func (m *model) addContextStack(index int) {
	item := m.list.Items()[index].(listItem)
	item.aiContext = true
	m.list.SetItem(index, item)
}

func (m *model) removeContextStack(index int) {
	item := m.list.Items()[index].(listItem)
	item.aiContext = false
	m.list.SetItem(index, item)
}

func (m *model) changeItemTitlePrefix(index int, prefix string) {
	item := m.list.Items()[index].(listItem)
	item.title = replacePrefix(item.title, prefix)
	m.list.SetItem(index, item)
}
