package kube

import (
	"fmt"

	"kctx/internal/model"

	"k8s.io/client-go/tools/clientcmd"
)

func NormalizeContext(p string, acc model.Account, c model.Cluster) string {
	if c.Location != "" {
		return fmt.Sprintf("%s-%s-%s-%s", p, acc.Name, c.Location, c.Name)
	}
	return fmt.Sprintf("%s-%s-%s", p, acc.Name, c.Name)
}

func RenameCurrentContext(newName string) error {
	cfg, err := clientcmd.LoadFromFile(clientcmd.RecommendedHomeFile)
	if err != nil {
		return err
	}

	old := cfg.CurrentContext

	ctxObj, ok := cfg.Contexts[old]
	if !ok {
		return fmt.Errorf("current context not found: %s", old)
	}

	finalName := newName
	i := 1
	for {
		if _, exists := cfg.Contexts[finalName]; !exists {
			break
		}
		finalName = fmt.Sprintf("%s-%d", newName, i)
		i++
	}

	cfg.Contexts[finalName] = ctxObj
	cfg.CurrentContext = finalName

	return clientcmd.WriteToFile(*cfg, clientcmd.RecommendedHomeFile)
}
