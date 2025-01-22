package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	config "github.com/shutils/lazyreview/pkg/config"
	openai "github.com/shutils/lazyreview/pkg/openai"
	ui "github.com/shutils/lazyreview/pkg/ui"
)

func main() {
	conf := config.NewConfig()
	client := openai.NewClient(conf)

	m := ui.NewUi(conf, client)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
