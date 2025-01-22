package ui

import (
	"bytes"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/shutils/lazyreview/pkg/config"
)

func defaultItemCollector(conf config.Config) []list.Item {
	items := []list.Item{}
	compiledPatterns := make([]*regexp.Regexp, len(conf.Ignores))
	for i, p := range conf.Ignores {
		compiledPatterns[i] = regexp.MustCompile(p)
	}
	err := filepath.WalkDir(conf.Target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// 絞り込み処理
			for _, re := range compiledPatterns {
				if re.MatchString(path) {
					return nil
				}
			}
			items = append(items, listItem{title: d.Name(), param: path})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return items
}

func customCollector(conf config.Config) []list.Item {
	items := []list.Item{}
	cmdArray := strings.Split(conf.Collector, " ")
	args := cmdArray[1:]
	cmd := exec.Command(cmdArray[0], args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return []list.Item{}
	}

	output := stdout.String()
	paramStrings := strings.Split(strings.TrimSpace(output), "\n")
	for _, param := range paramStrings {
		items = append(items, listItem{title: filepath.Base(param), param: param})
	}
	return items
}
