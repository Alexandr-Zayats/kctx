package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "kctx",
	Short: "Multi-cloud Kubernetes context switcher",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
