package do

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"kctx/internal/cache"
	"kctx/internal/model"
)

type DO struct{}

func (d *DO) Name() string {
	return "do"
}

func (d *DO) ListAccounts(ctx context.Context) ([]model.Account, error) {
	out, err := exec.CommandContext(ctx, "doctl", "auth", "list").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")

	type result struct {
		context string
		team    string
		count   int
	}

	var contexts []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		name := strings.ReplaceAll(l, " (current)", "")
		contexts = append(contexts, name)
	}

	ch := make(chan result, len(contexts))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, ctxName := range contexts {
		wg.Add(1)

		go func(c string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			team := getTeamName(ctx, c)
			count := getClusterCount(ctx, c)

			ch <- result{
				context: c,
				team:    team,
				count:   count,
			}
		}(ctxName)
	}

	wg.Wait()
	close(ch)

	var res []model.Account

	for r := range ch {
		label := r.context

		if r.team != "" {
			label = fmt.Sprintf("%s [%s]", r.context, r.team)
		}

		if r.count > 0 {
			label = fmt.Sprintf("%s (%d)", label, r.count)
		}

		res = append(res, model.Account{
			Name: label,
			Meta: map[string]string{
				"context": r.context,
			},
		})
	}

	return res, nil
}

func (d *DO) UseAccount(ctx context.Context, acc model.Account) error {
	contextName := acc.Name

	if acc.Meta != nil {
		if v, ok := acc.Meta["context"]; ok {
			contextName = v
		}
	}

	if err := ensureContextExists(ctx, contextName); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx,
		"doctl", "auth", "switch",
		"--context", contextName,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (d *DO) ListClusters(ctx context.Context) ([]model.Cluster, error) {
	cmd := exec.CommandContext(ctx,
		"doctl", "kubernetes", "cluster", "list",
		"--format", "Name", "--no-header",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	var res []model.Cluster
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		res = append(res, model.Cluster{Name: l})
	}

	return res, nil
}

func (d *DO) GetCredentials(ctx context.Context, c model.Cluster) error {
	cmd := exec.CommandContext(ctx,
		"doctl", "kubernetes", "cluster", "kubeconfig", "save", c.Name,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ensureContextExists(ctx context.Context, name string) error {
	out, _ := exec.CommandContext(ctx, "doctl", "auth", "list").Output()

	if !strings.Contains(string(out), name) {
		return fmt.Errorf("DO context '%s' not found. Run: kctx do add", name)
	}

	return nil
}

func getTeamName(ctx context.Context, context string) string {
	_ = exec.CommandContext(ctx, "doctl", "auth", "switch", "--context", context).Run()

	out, err := exec.CommandContext(ctx,
		"doctl", "account", "get",
		"--format", "Team",
		"--no-header",
	).Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func getClusterCount(ctx context.Context, context string) int {
	key := "do_clusters_" + context

	if data, ok := cache.Get(key); ok {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			return 0
		}
		return len(lines)
	}

	_ = exec.CommandContext(ctx, "doctl", "auth", "switch", "--context", context).Run()

	cmd := exec.CommandContext(ctx,
		"doctl", "kubernetes", "cluster", "list",
		"--format", "Name", "--no-header",
	)

	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	cache.Set(key, out, 2*time.Minute)

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0
	}
	return len(lines)
}
