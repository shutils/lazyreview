package ui

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/shutils/lazyreview/pkg/config"
)

func defaultPreviewer(param string) string {
	if param == "" {
		return "Error: No param"
	}
	var fallbackText = "This item is not text"
	content, err := os.ReadFile(param)
	if err != nil {
		return string(fallbackText)
	}
	ty := http.DetectContentType(content)
	switch ty {
	case "text/plain; charset=utf-8":
		return string(content)
	case "text/xml; charset=utf-8":
		return string(content)
	case "text/html; charset=utf-8":
		return string(content)
	default:
		return string(fallbackText)
	}
}

func customPreviewer(cmds []string, param string) string {
	if param == "" {
		return "Error: No param"
	}
	if len(cmds) == 0 {
		return "Error: No previewer"
	}
	args := cmds[1:]
	args = append(args, param)
	cmd := exec.Command(cmds[0], args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}
	return output
}

func previewContent(item listItem, sources []config.Source) string {
	if item.sourceName != "" {
		source, _ := getSource(item.sourceName, sources)
		if len(source.Previewer) != 0 {
			return customPreviewer(source.Previewer, item.param)
		}
	}
	return defaultPreviewer(item.param)
}

func previewContext(item contextItem, sources []config.Source) string {
	if item.isEdited {
		return item.content
	}
	if item.sourceName != "" {
		source, _ := getSource(item.sourceName, sources)
		if len(source.Previewer) != 0 {
			return customPreviewer(source.Previewer, item.param)
		}
	}
	return defaultPreviewer(item.param)
}
