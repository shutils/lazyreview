package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/config"
)

func (m *model) onChangeListSelectedItem() (tea.Model, tea.Cmd) {
	selectedItem, ok := m.panels.itemListPanel.SelectedItem().(listItem)
	reviewContent := "No review"
	itemContent := previewContent(selectedItem, m.conf.Sources)
	if ok && m.getReviewIndex(selectedItem.id) != -1 {
		reviewContent = getRendered(m.reviewList[m.getReviewIndex(selectedItem.id)].Review, m.conf.Glamour, m.panels.itemReviewPanel.Width)
	}
	m.loadReviewPanel(reviewContent)
	m.loadContentPanel(itemContent)
	return m, nil
}

func (m *model) loadReviewPanel(itemContent string) {
	m.panels.itemReviewPanel.SetContent(itemContent)
	m.panels.itemReviewPanel.GotoTop()
}

func (m *model) loadContentPanel(itemContent string) {
	m.panels.itemPreviewPanel.SetContent(itemContent)
	m.panels.itemPreviewPanel.GotoTop()
}

// Returns a index of the item with the given param
func findIndex(items []list.Item, id string) int {
	for i, item := range items {
		if item.(listItem).id == id {
			return i
		}
	}
	return -1
}

func getItems(conf config.Config, reviewList []reviewInfo) []list.Item {
	var items []list.Item

	if len(conf.Sources) > 0 && !isDisabledAllSource(conf.Sources) {
		items = collectItemsFromSources(conf.Sources)
	} else if conf.Collector != "" {
		items = customCollector(conf.Collector, "")
	} else {
		items = defaultItemCollector(conf)
	}

	reviewStateMap := make(map[string]string)
	for _, review := range reviewList {
		reviewStateMap[review.ID] = review.State
	}

	for i, item := range items {
		_item, ok := item.(listItem)
		if !ok {
			continue
		}

		title := _item.Title()
		id := makeHash(_item)
		if state, exists := reviewStateMap[id]; exists {
			if state == "finish" {
				title = "☑ " + title
			} else {
				title = "☐ " + title
			}
		} else {
			title = "☐ " + title
		}

		items[i] = listItem{
			title:      title,
			param:      _item.Description(),
			sourceName: _item.sourceName,
			id:         id,
		}
	}

	return items
}

func getReviewStackItems(items []list.Item, indexes []int) []list.Item {
	var filteredItems []list.Item

	// Create a set of valid indexes for quick lookup
	indexSet := make(map[int]struct{})
	for _, index := range indexes {
		indexSet[index] = struct{}{}
	}

	// Iterate through the original items and check if the index is in the set
	for i, item := range items {
		if _, exists := indexSet[i]; exists {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func getItemListString(items []list.Item) string {
	var params []string

	for _, item := range items {
		_item, ok := item.(listItem)
		if ok {
			params = append(params, _item.Description())
		}
	}

	return strings.Join(params, "\n")
}

func isDisabledAllSource(sources []config.Source) bool {
	for _, source := range sources {
		if source.Enabled {
			return false
		}
	}
	return true
}

func collectItemsFromSources(sources []config.Source) []list.Item {
	var items []list.Item

	for _, source := range sources {
		if source.Enabled {
			if source.Collector == "" {
				continue
			}
			collectedItems := customCollector(source.Collector, source.Name)
			if len(collectedItems) != 0 {
				items = append(items, collectedItems...)
			}
		}
	}

	return items
}

func getSource(name string, sources []config.Source) (config.Source, error) {
	for _, source := range sources {
		if source.Name == name {
			return source, nil
		}
	}
	return config.Source{}, fmt.Errorf("source with name '%s' not found", name)
}

func makeHash(item listItem) string {
	seed := item.param + item.sourceName
	return fmt.Sprintf("%x", seed)
}

func getSourceItems(sources []config.Source) []list.Item {
	var items []list.Item

	if len(sources) == 0 {
		return []list.Item{}
	}

	for _, source := range sources {
		items = append(items, sourceItem{
			name:      source.Name,
			collector: source.Collector,
			previewer: source.Previewer,
			enabled:   source.Enabled,
		})
	}

	return items
}
