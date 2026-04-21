package provider

import (
	"context"

	"kctx/internal/model"
)

type Provider interface {
	Name() string
	ListAccounts(ctx context.Context) ([]model.Account, error)
	UseAccount(ctx context.Context, acc model.Account) error
	ListClusters(ctx context.Context) ([]model.Cluster, error)
	GetCredentials(ctx context.Context, c model.Cluster) error
}

// Optional capability.
// Реализуется только провайдерами, которым нужен pre-flight auth check.
type AuthChecker interface {
	CheckAuth(ctx context.Context) error
}
