package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	"kctx/internal/kube"
	"kctx/internal/model"
	"kctx/internal/provider"
	"kctx/internal/provider/gcloud"
	"kctx/pkg/ui"
)

type Switcher struct {
	Providers []provider.Provider
}

func (s *Switcher) Run(ctx context.Context) error {
	var p provider.Provider
	var err error

	if len(s.Providers) == 1 {
		p = s.Providers[0]
		fmt.Println("✔", p.Name())
	} else {
		p, err = ui.SelectProvider(s.Providers)
		if err != nil {
			return err
		}
	}

	if ac, ok := p.(provider.AuthChecker); ok {
		if err := ac.CheckAuth(ctx); err != nil {
			return err
		}
	}

	accs, err := p.ListAccounts(ctx)
	if err != nil {
		return err
	}
	if len(accs) == 0 {
		return fmt.Errorf("no accounts found for provider %s", p.Name())
	}

	acc, err := ui.SelectAccount(accs)
	if err != nil {
		return err
	}

	if err := p.UseAccount(ctx, acc); err != nil {
		return err
	}

	if p.Name() == "gcp" {
		if gc, ok := p.(*gcloud.GCloud); ok {
			projects, err := gc.ListProjects(ctx)
			if err != nil {
				return err
			}

			projDisplay, err := ui.SelectString("Project", projects)
			if err != nil {
				return err
			}

			proj := projDisplay
			if i := strings.Index(projDisplay, " ("); i != -1 {
				proj = projDisplay[:i]
			}

			os.Setenv("KCTX_GCP_PROJECT", proj)
		}
	}

	clusters, err := p.ListClusters(ctx)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		fmt.Println("⚠️ no clusters found")
		return nil
	}

	cluster, err := ui.SelectCluster(clusters)
	if err != nil {
		return err
	}

	if err := p.GetCredentials(ctx, cluster); err != nil {
		return err
	}

	realName := acc.Name
	if acc.Meta != nil {
		if v, ok := acc.Meta["profile"]; ok {
			realName = v
		}
		if v, ok := acc.Meta["project"]; ok {
			realName = v
		}
		if v, ok := acc.Meta["account"]; ok {
			realName = v
		}
		if v, ok := acc.Meta["context"]; ok {
			realName = v
		}
	}

	ctxName := kube.NormalizeContext(p.Name(), model.Account{
		Name: realName,
	}, cluster)

	if err := kube.RenameCurrentContext(ctxName); err != nil {
		return err
	}

	fmt.Println("Switched to:", ctxName)
	return nil
}
