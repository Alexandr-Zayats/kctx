package ui

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Select(label string, items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select")
	}

	cmd := exec.Command("fzf", "--prompt", label+": ")

	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func SelectWithPreview(label string, items []string, previewCmd string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select")
	}

	cmd := exec.Command("fzf",
		"--prompt", label+": ",
		"--preview", previewCmd,
		"--preview-window=right:60%",
	)

	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
