package kube

import (
	"context"
	"os/exec"
	"strings"
)

func ListNamespaces(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx,
		"kubectl", "get", "ns",
		"-o", "jsonpath={.items[*].metadata.name}",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(out)), nil
}

func SetNamespace(ctx context.Context, ns string) error {
	cmd := exec.CommandContext(ctx,
		"kubectl", "config", "set-context", "--current", "--namespace", ns,
	)

	return cmd.Run()
}
