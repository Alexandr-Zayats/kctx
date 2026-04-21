package ui

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func FzfSelect(title string, items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select")
	}

	input := strings.Join(items, "\n")

	args := []string{
		"--height=20",
		"--layout=reverse",
		"--border",
		"--ansi",
		"--prompt=" + title + "> ",
	}

	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(input)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("fzf failed: %s", msg)
	}

	return strings.TrimSpace(stdout.String()), nil
}
