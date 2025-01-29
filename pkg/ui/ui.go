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
	title, param, sourceName, id string
}

type sourceItem struct {
	name, collector, previewer string
	enabled                    bool
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
		return "☑ collector: " + i.collector + " previewer: " + i.previewer
	}
	return "☐ collector: " + i.collector + " previewer: " + i.previewer
}
func (i sourceItem) FilterValue() string { return i.name }

type updateSourceListMsg struct {
}

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
	contextPanel        list.Model
	contextDetailPanel  viewport.Model
	sourceDetailPanel   viewport.Model
	sourceListPanel     list.Model
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
	sourceListKeyMap    sourceListKeyMap
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
		list:               itemList,
		contentPanel:       contentPanel,
		reviewPanel:        reviewPanel,
		instantPromptPanel: instantPromptPanel,
		configSummaryPanel: configPanel,
		configContentPanel: configContentPanel,
		statePanel:         viewport.New(0, 0),
		stateDetailPanel:   viewport.New(0, 0),
		contextPanel: list.New([]list.Item{}, list.DefaultDelegate{
			ShowDescription: false,
			Styles:          list.NewDefaultItemStyles(),
		}, 0, 0),
		contextDetailPanel: viewport.New(0, 0),
		sourceListPanel: list.New([]list.Item{}, list.DefaultDelegate{
			ShowDescription: false,
			Styles:          list.NewDefaultItemStyles(),
		}, 0, 0),
		sourceDetailPanel:   viewport.New(0, 0),
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
		sourceListKeyMap:    GetSourceListKeymap(),
		uiState:             state.LoadState(conf.State),
		currentHistoryIndex: 0,
		state:               state.State{},
	}
	m.sourceListPanel.SetShowTitle(false)
	m.sourceListPanel.SetShowHelp(false)
	m.sourceListPanel.SetShowStatusBar(false)
	m.sourceListPanel.SetShowFilter(false)
	m.sourceListPanel.KeyMap.Filter.Unbind()

	m.contextPanel.SetShowTitle(false)
	m.contextPanel.SetShowHelp(false)
	m.contextPanel.SetShowStatusBar(false)
	m.contextPanel.SetShowFilter(false)
	m.contextPanel.KeyMap.Filter.Unbind()

	m.UpdateState()
	m.currentHistoryIndex = len(m.uiState.PromptHistory)
	m.loadReviews()
	m.list.SetItems(getItems(m.conf, m.reviewList))
	m.sourceListPanel.SetItems(getSourceItems(m.conf.Sources))
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
		index := findIndex(m.list.Items(), msg.id)
		if index == -1 {
			return m, nil
		}
		item := m.list.Items()[index].(listItem)
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
			m.reviewPanel.SetContent(msg.content)
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
		index := findIndex(m.list.Items(), msg.id)
		if index == -1 {
			return m, nil
		}
		if msg.operation == Add {
			m.addReviewStack(index)
		} else {
			m.removeReviewStack(index)
			m.changeItemTitlePrefix(index, "☑ ")
		}
		m.updateReviewStackPanel()
		return m, nil
	case updateSourceListMsg:
		m.sourceListPanel.SetItems(getSourceItems(m.conf.Sources))
		m.contextPanel.Update(msg)
		m.list.SetItems(getItems(m.conf, m.reviewList))
	default:
		switch m.focusState {
		case SourceListPanelFocus:
			selectedSourceName := m.sourceListPanel.SelectedItem().(sourceItem)
			selectedSource := m.conf.GetSourceFromName(selectedSourceName.name)
			m.sourceDetailPanel.SetContent(selectedSource.String())
		case ContextPanelFocus:
			m.contextDetailPanel.SetContent(m.getContextString())
		}
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == ListPanelFocus {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == SourceListPanelFocus {
		m.sourceListPanel, cmd = m.sourceListPanel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == ContextPanelFocus {
		m.contextPanel, cmd = m.contextPanel.Update(msg)
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
	itemTitleList := getItemListString(getReviewStackItems(m.list.Items(), m.reviewStack))
	m.reviewStackPanel.SetContent(itemTitleList)
}

func (m *model) addContextStack(id string) (tea.Model, tea.Cmd) {
	index := findIndex(m.list.Items(), id)
	if index == -1 {
		return *m, nil
	}
	item := m.list.Items()[index]

	contextList := m.contextPanel.Items()
	contextList = append(contextList, item)
	m.contextPanel.SetItems(contextList)
	return *m, nil
}

func (m *model) removeContextStack(id string) (tea.Model, tea.Cmd) {
	index := findIndex(m.contextPanel.Items(), id)
	contextList := m.contextPanel.Items()

	if index < 0 || index >= len(contextList) {
		return *m, nil
	}

	var newContextList []list.Item
	if index == len(contextList)-1 {
		newContextList = contextList[:index]
	} else {
		newContextList = append(contextList[:index], contextList[index+1:]...)
	}

	m.contextPanel.SetItems(newContextList)
	return *m, nil
}

func (m *model) changeItemTitlePrefix(index int, prefix string) {
	item := m.list.Items()[index].(listItem)
	item.title = replacePrefix(item.title, prefix)
	m.list.SetItem(index, item)
}
