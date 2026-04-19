package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Manage DigitalOcean contexts",
}

var doAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new DigitalOcean context (team)",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Context name (e.g. team-prod): ")
		ctxName, _ := reader.ReadString('\n')
		ctxName = strings.TrimSpace(ctxName)

		fmt.Print("DigitalOcean API Token: ")
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		c := exec.Command("doctl", "auth", "init", "--context", ctxName)
		c.Stdin = strings.NewReader(token + "\n")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return err
		}

		fmt.Println("✔ Context added:", ctxName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
	doCmd.AddCommand(doAddCmd)
}
