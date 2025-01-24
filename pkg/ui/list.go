package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shutils/lazyreview/pkg/config"
)

func (m *model) onChangeListSelectedItem() (tea.Model, tea.Cmd) {
	selectedItem, ok := m.list.SelectedItem().(listItem)
	reviewContent := "No review"
	itemContent := ""
	if ok && m.getReviewIndex(selectedItem.param) != -1 {
		reviewContent = getRendered(m.reviewList[m.getReviewIndex(selectedItem.param)].Review, m.conf.Glamour, m.reviewPanel.Width)
	}
	if m.conf.Previewer != "" {
		itemContent = customPreviewer(m.conf.Previewer, selectedItem.param)
	} else {
		itemContent = defaultPreviewer(selectedItem.param)
	}
	m.reviewPanel.SetContent(reviewContent)
	m.contentPanel.SetContent(itemContent)
	return m, nil
}

// Returns a index of the item with the given param
func findIndex(items []list.Item, param string) int {
	for i, item := range items {
		if item.(listItem).param == param {
			return i
		}
	}
	return -1
}

func getItems(conf config.Config, reviewList []reviewInfo) []list.Item {
	var items []list.Item

	if conf.Collector != "" {
		items = customCollector(conf)
	} else {
		items = defaultItemCollector(conf)
	}

	reviewStateMap := make(map[string]string)
	for _, review := range reviewList {
		reviewStateMap[review.Param] = review.State
	}

	for i, item := range items {
		_item, ok := item.(listItem)
		if !ok {
			continue
		}

		title := _item.Title()
		if state, exists := reviewStateMap[_item.Description()]; exists {
			if state == "finish" {
				title = "☑ " + title
			} else {
				title = "☐ " + title
			}
		} else {
			title = "☐ " + title
		}

		items[i] = listItem{
			title:     title,
			param:     _item.Description(),
			aiContext: false,
		}
	}

	return items
}

func getContextItems(items []list.Item) []list.Item {
	var _items []list.Item

	for _, item := range items {
		_item, ok := item.(listItem)
		if !ok {
			continue
		}
		if _item.aiContext {
			_items = append(_items, _item)
		}
	}

	return _items
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
