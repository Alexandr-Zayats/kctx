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

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [provider|alias]",
	Short: "Switch Kubernetes context",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		providers := provider.All()

		if len(args) == 1 {
			name := strings.ToLower(args[0])

			for _, p := range providers {
				if p.Name() == name {
					s := core.Switcher{
						Providers: []provider.Provider{p},
					}
					return s.Run(context.Background())
				}
			}

			return fmt.Errorf("unknown provider: %s", args[0])
		}

		s := core.Switcher{
			Providers: providers,
		}

		return s.Run(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
