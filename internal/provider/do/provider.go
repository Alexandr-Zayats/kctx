package do

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"kctx/internal/cache"
	"kctx/internal/model"
)

type DO struct{}

func (d *DO) Name() string {
	return "do"
}

func (d *DO) CheckAuth(ctx context.Context) error {
	if _, err := exec.LookPath("doctl"); err != nil {
		return fmt.Errorf("DigitalOcean CLI not found. Install 'doctl' and try again")
	}

	checkCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(checkCtx, "doctl", "account", "get")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	msg := strings.TrimSpace(string(out))

	if errors.Is(checkCtx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("DigitalOcean auth check timed out. Re-run 'doctl auth init'")
	}

	if msg == "" {
		return fmt.Errorf("DigitalOcean not authenticated.\n\nRun:\n  doctl auth init")
	}

	return fmt.Errorf("DigitalOcean not authenticated.\n\nRun:\n  doctl auth init\n\nDetails:\n  %s", msg)
}

func (d *DO) ListAccounts(ctx context.Context) ([]model.Account, error) {
	key := "do_contexts"

	var out []byte
	if data, ok := cache.Get(key); ok {
		out = data
	} else {
		var err error
		out, err = exec.CommandContext(ctx, "doctl", "auth", "list").Output()
		if err != nil {
			return nil, err
		}
		cache.Set(key, out, 15*time.Minute)
	}

	lines := strings.Split(string(out), "\n")

	var res []model.Account
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		contextName := strings.ReplaceAll(l, " (current)", "")
		label := contextName

		team := getTeamName(ctx, contextName)
		if team != "" {
			label = fmt.Sprintf("%s [%s]", contextName, team)
		}

		res = append(res, model.Account{
			Name: label,
			Meta: map[string]string{
				"context": contextName,
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
	currentContext, err := currentDOContext(ctx)
	if err != nil {
		return nil, err
	}

	key := "do_clusters_" + currentContext
	if data, ok := cache.Get(key); ok {
		return parseDOClusters(data), nil
	}

	cmd := exec.CommandContext(ctx,
		"doctl", "kubernetes", "cluster", "list",
		"--format", "Name", "--no-header",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	cache.Set(key, out, 10*time.Minute)
	return parseDOClusters(out), nil
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
	key := "do_contexts"

	var out []byte
	if data, ok := cache.Get(key); ok {
		out = data
	} else {
		var err error
		out, err = exec.CommandContext(ctx, "doctl", "auth", "list").Output()
		if err != nil {
			return err
		}
		cache.Set(key, out, 15*time.Minute)
	}

	if !strings.Contains(string(out), name) {
		return fmt.Errorf("DO context '%s' not found. Run: kctx do add", name)
	}

	return nil
}

func currentDOContext(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "doctl", "auth", "list").Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasSuffix(l, " (current)") {
			return strings.TrimSpace(strings.TrimSuffix(l, " (current)")), nil
		}
	}

	return "", fmt.Errorf("current DO context not found")
}

func getTeamName(ctx context.Context, contextName string) string {
	key := "do_team_" + contextName
	if data, ok := cache.Get(key); ok {
		return strings.TrimSpace(string(data))
	}

	out, err := exec.CommandContext(ctx,
		"doctl",
		"--context", contextName,
		"account", "get",
		"--format", "Team",
		"--no-header",
	).Output()

	if err != nil {
		return ""
	}

	cache.Set(key, out, 30*time.Minute)
	return strings.TrimSpace(string(out))
}

func parseDOClusters(data []byte) []model.Cluster {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	var res []model.Cluster
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		res = append(res, model.Cluster{Name: l})
	}

	return res
}
