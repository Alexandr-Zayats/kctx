package project

import (
	"os"
	"path/filepath"
)

func DetectAWSConfig() (string, string) {
	dir, _ := os.Getwd()

	for {
		// 🔥 yaml
		yaml := filepath.Join(dir, ".kctx", "aws.yaml")
		if _, err := os.Stat(yaml); err == nil {
			return dir, yaml
		}

		// 🔥 conf
		conf := filepath.Join(dir, ".kctx", "aws.conf")
		if _, err := os.Stat(conf); err == nil {
			return dir, conf
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", ""
}
