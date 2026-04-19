package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var awsCleanCmd = &cobra.Command{
	Use:   "aws-clean",
	Short: "Remove invalid AWS profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		out, err := exec.CommandContext(ctx, "aws", "configure", "list-profiles").Output()
		if err != nil {
			return err
		}

		lines := strings.Split(string(out), "\n")

		var bad []string

		for _, p := range lines {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			check := exec.CommandContext(ctx,
				"aws", "sts", "get-caller-identity",
				"--profile", p,
			)

			if err := check.Run(); err != nil {
				fmt.Println("❌ invalid:", p)
				bad = append(bad, p)
			} else {
				fmt.Println("✅ valid:", p)
			}
		}

		if len(bad) == 0 {
			fmt.Println("No invalid profiles 🎉")
			return nil
		}

		fmt.Println("\nProfiles to delete:")
		for _, b := range bad {
			fmt.Println(" -", b)
		}

		fmt.Print("\nDelete them? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')

		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("Aborted")
			return nil
		}

		for _, p := range bad {
			removeProfile(p)
		}

		fmt.Println("Done 🧹")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(awsCleanCmd)
}

func removeProfile(profile string) {
	files := []string{
		os.Getenv("HOME") + "/.aws/config",
		os.Getenv("HOME") + "/.aws/credentials",
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")

		var out []string
		skip := false

		for _, l := range lines {
			// config формат
			if strings.Contains(l, "[profile "+profile+"]") ||
				strings.TrimSpace(l) == "["+profile+"]" {
				skip = true
				continue
			}

			if skip {
				if strings.HasPrefix(l, "[") {
					skip = false
				} else {
					continue
				}
			}

			out = append(out, l)
		}

		_ = os.WriteFile(file, []byte(strings.Join(out, "\n")), 0644)
	}
}
