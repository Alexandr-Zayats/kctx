package core

import (
	"context"
	"fmt"

	"kctx/internal/kube"
	"kctx/pkg/ui"
)

func SwitchNamespace(ctx context.Context) error {
	ns, err := kube.ListNamespaces(ctx)
	if err != nil {
		return err
	}

	selected, err := ui.SelectString("Namespace", ns)
	if err != nil {
		return err
	}

	if err := kube.SetNamespace(ctx, selected); err != nil {
		return err
	}

	fmt.Println("Switched namespace to:", selected)
	return nil
}
