package cmd

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current context",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
		if err != nil {
			return err
		}

		ctx := cfg.CurrentContext

		ns := cfg.Contexts[ctx].Namespace
		if ns == "" {
			ns = "default"
		}

		fmt.Println("Context:", ctx)
		fmt.Println("Namespace:", ns)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
