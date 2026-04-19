package cmd

import "github.com/spf13/cobra"

var previewCmd = &cobra.Command{
	Use:    "internal-preview",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
