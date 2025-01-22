package ui

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func defaultPreviewer(param string) string {
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

func customPreviewer(previewer string, param string) string {
	cmdArray := strings.Split(previewer, " ")
	args := cmdArray[1:]
	args = append(args, param)
	cmd := exec.Command(cmdArray[0], args...)

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
