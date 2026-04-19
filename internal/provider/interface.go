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
	GetCredentials(ctx context.Context, cluster model.Cluster) error
}
