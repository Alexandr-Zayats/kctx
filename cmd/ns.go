package cmd

import (
	"context"

	"kctx/internal/core"

	"github.com/spf13/cobra"
)

var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "Switch namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.SwitchNamespace(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(nsCmd)
}
