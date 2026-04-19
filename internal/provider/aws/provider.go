package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"kctx/internal/cache"
	"kctx/internal/model"
	"kctx/internal/project"
)

type AWS struct{}

func (a *AWS) Name() string {
	return "aws"
}

// =====================
// ACCOUNTS
// =====================
func (a *AWS) ListAccounts(ctx context.Context) ([]model.Account, error) {
	root, path := project.DetectAWSConfig()

	if path != "" {
		cfg, err := LoadConfig(path)
		if err != nil {
			return nil, err
		}

		projectName := cfg.Project
		if projectName == "" {
			projectName = filepath.Base(root)
		}

		if err := syncAWSConfig(cfg, projectName); err != nil {
			return nil, err
		}

		var res []model.Account
		for name, acc := range cfg.Accounts {
			profile := fmt.Sprintf("%s-%s", projectName, name)

			res = append(res, model.Account{
				Name: profile,
				Meta: map[string]string{
					"profile": profile,
					"region":  acc.Region,
				},
			})
		}

		return res, nil
	}

	return listProfiles(ctx)
}

func listProfiles(ctx context.Context) ([]model.Account, error) {
	out, err := exec.CommandContext(ctx, "aws", "configure", "list-profiles").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")

	var res []model.Account

	for _, p := range lines {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		res = append(res, model.Account{
			Name: p,
			Meta: map[string]string{
				"profile": p,
			},
		})
	}

	return res, nil
}

// =====================
// CONFIG SYNC
// =====================
func syncAWSConfig(cfg *Config, project string) error {
	role := cfg.SSO.RoleName
	if role == "" {
		role = "AdministratorAccess"
	}

	for name, acc := range cfg.Accounts {
		profile := fmt.Sprintf("%s-%s", project, name)

		cmds := [][]string{
			{"configure", "set", "profile." + profile + ".sso_start_url", cfg.SSO.StartURL},
			{"configure", "set", "profile." + profile + ".sso_region", cfg.SSO.Region},
			{"configure", "set", "profile." + profile + ".sso_account_id", acc.AccountID},
			{"configure", "set", "profile." + profile + ".sso_role_name", role},
			{"configure", "set", "profile." + profile + ".region", acc.Region},
			{"configure", "set", "profile." + profile + ".output", "json"},
		}

		for _, args := range cmds {
			cmd := exec.Command("aws", args...)
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

// =====================
// USE ACCOUNT
// =====================
func (a *AWS) UseAccount(ctx context.Context, acc model.Account) error {
	profile := acc.Name

	if acc.Meta != nil {
		if v, ok := acc.Meta["profile"]; ok {
			profile = v
		}
	}

	os.Setenv("AWS_PROFILE", profile)
	os.Setenv("AWS_DEFAULT_PROFILE", profile)
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	return ensureSSO(profile)
}

// =====================
// CLUSTERS
// =====================
func (a *AWS) ListClusters(ctx context.Context) ([]model.Cluster, error) {
	profile := os.Getenv("AWS_PROFILE")

	if profile == "" {
		return []model.Cluster{}, nil
	}

	if err := ensureSSO(profile); err != nil {
		return nil, err
	}

	key := "aws_clusters_" + profile

	if data, ok := cache.Get(key); ok {
		return parseClusters(data), nil
	}

	cmd := exec.CommandContext(ctx,
		"aws", "eks", "list-clusters",
		"--output", "json",
		"--profile", profile,
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	cache.Set(key, out, 60*time.Second)

	return parseClusters(out), nil
}

// =====================
// 🔥 
// =====================
func (a *AWS) GetCredentials(ctx context.Context, c model.Cluster) error {
	profile := os.Getenv("AWS_PROFILE")

	cmd := exec.CommandContext(ctx,
		"aws", "eks", "update-kubeconfig",
		"--name", c.Name,
		"--profile", profile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// =====================
// SSO
// =====================
func ensureSSO(profile string) error {
	key := "aws_valid_" + profile

	if _, ok := cache.Get(key); ok {
		return nil
	}

	cmd := exec.Command("aws", "sts", "get-caller-identity", "--profile", profile)
	if err := cmd.Run(); err == nil {
		cache.Set(key, []byte("ok"), 10*time.Minute)
		return nil
	}

	fmt.Println("🔐 AWS SSO login:", profile)

	login := exec.Command("aws", "sso", "login", "--profile", profile)
	login.Stdout = os.Stdout
	login.Stderr = os.Stderr

	if err := login.Run(); err != nil {
		return err
	}

	cache.Set(key, []byte("ok"), 10*time.Minute)

	return nil
}

// =====================
// HELPERS
// =====================
func parseClusters(data []byte) []model.Cluster {
	var parsed struct {
		Clusters []string `json:"clusters"`
	}

	_ = json.Unmarshal(data, &parsed)

	var res []model.Cluster
	for _, c := range parsed.Clusters {
		res = append(res, model.Cluster{Name: c})
	}

	return res
}
