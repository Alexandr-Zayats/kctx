package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"kctx/internal/config"
	"kctx/internal/core"
	"kctx/internal/provider"
	"kctx/internal/provider/aws"
	"kctx/internal/provider/do"
	"kctx/internal/provider/gcloud"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [provider|alias]",
	Short: "Switch Kubernetes context",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Alias mode: kctx switch <alias>
		if len(args) == 1 {
			cfg, err := config.Load()
			if err == nil {
				if target, ok := cfg.Aliases[args[0]]; ok {
					c := exec.Command("kubectl", "config", "use-context", target)
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					return c.Run()
				}
			}
		}

		// Explicit provider mode: kctx switch aws|gcp|do
		if len(args) == 1 {
			p, err := resolveProvider(args[0])
			if err == nil {
				s := core.Switcher{
					Providers: []provider.Provider{p},
				}
				return s.Run(context.Background())
			}
		}

		// Generic mode: kctx switch
		s := core.Switcher{
			Providers: provider.All(),
		}

		return s.Run(context.Background())
	},
}

func resolveProvider(name string) (provider.Provider, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "aws":
		return &aws.AWS{}, nil
	case "do", "digitalocean":
		return &do.DO{}, nil
	case "gcp", "gcloud", "google":
		return &gcloud.GCloud{}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
