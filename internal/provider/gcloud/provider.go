package gcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"kctx/internal/cache"
	"kctx/internal/model"
)

type GCloud struct{}

func (g *GCloud) Name() string {
	return "gcp"
}

func (g *GCloud) ListAccounts(ctx context.Context) ([]model.Account, error) {
	out, err := exec.CommandContext(ctx, "gcloud", "auth", "list", "--format=json").Output()
	if err != nil {
		return nil, err
	}

	var data []struct {
		Account string `json:"account"`
	}

	if err := json.Unmarshal(out, &data); err != nil {
		return nil, err
	}

	var res []model.Account
	for _, a := range data {
		if a.Account == "" {
			continue
		}
		res = append(res, model.Account{
			Name: a.Account,
			Meta: map[string]string{
				"account": a.Account,
			},
		})
	}

	return res, nil
}

func (g *GCloud) ListProjects(ctx context.Context) ([]string, error) {
	account := os.Getenv("CLOUDSDK_CORE_ACCOUNT")

	out, err := exec.CommandContext(
		ctx,
		"gcloud",
		"projects", "list",
		"--account="+account,
		"--format=value(projectId)",
	).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")

	var res []string
	var fallback []string

	for _, p := range lines {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		fallback = append(fallback, p)
		count := getClusterCount(ctx, p)

		if count > 0 {
			res = append(res, fmt.Sprintf("%s (%d)", p, count))
		}
	}

	if len(res) == 0 {
		return fallback, nil
	}

	return res, nil
}

func (g *GCloud) UseAccount(ctx context.Context, acc model.Account) error {
	account := acc.Name
	if acc.Meta != nil {
		if v, ok := acc.Meta["account"]; ok {
			account = v
		}
	}

	os.Setenv("CLOUDSDK_CORE_ACCOUNT", account)
	return nil
}

func (g *GCloud) ListClusters(ctx context.Context) ([]model.Cluster, error) {
	project := os.Getenv("KCTX_GCP_PROJECT")
	account := os.Getenv("CLOUDSDK_CORE_ACCOUNT")

	if project == "" || account == "" {
		return []model.Cluster{}, nil
	}

	key := "gcp_clusters_" + project

	if data, ok := cache.Get(key); ok {
		return parseClustersFromValue(data), nil
	}

	cmd := exec.CommandContext(
		ctx,
		"gcloud",
		"container", "clusters", "list",
		"--account="+account,
		"--project="+project,
		"--format=value(name,location)",
	)

	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return []model.Cluster{}, nil
	}

	cache.Set(key, out, 90*time.Second)
	return parseClustersFromValue(out), nil
}

func (g *GCloud) GetCredentials(ctx context.Context, c model.Cluster) error {
	account := os.Getenv("CLOUDSDK_CORE_ACCOUNT")

	cmd := exec.CommandContext(
		ctx,
		"gcloud",
		"container", "clusters", "get-credentials",
		c.Name,
		"--account="+account,
		"--zone", c.Location,
	)
	if err := cmd.Run(); err == nil {
		return nil
	}

	cmd = exec.CommandContext(
		ctx,
		"gcloud",
		"container", "clusters", "get-credentials",
		c.Name,
		"--account="+account,
		"--region", c.Location,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getClusterCount(ctx context.Context, project string) int {
	account := os.Getenv("CLOUDSDK_CORE_ACCOUNT")
	key := "gcp_clusters_" + project

	if data, ok := cache.Get(key); ok {
		return len(strings.Fields(string(data)))
	}

	cmd := exec.CommandContext(
		ctx,
		"gcloud",
		"container", "clusters", "list",
		"--account="+account,
		"--project="+project,
		"--format=value(name)",
	)

	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	cache.Set(key, out, 2*time.Minute)
	return len(strings.Fields(string(out)))
}

func parseClustersFromValue(data []byte) []model.Cluster {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	var res []model.Cluster
	for _, l := range lines {
		parts := strings.Fields(l)
		if len(parts) < 2 {
			continue
		}
		res = append(res, model.Cluster{
			Name:     parts[0],
			Location: parts[1],
		})
	}

	return res
}
