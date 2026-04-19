package provider

import (
	"kctx/internal/provider/aws"
	"kctx/internal/provider/do"
	"kctx/internal/provider/gcloud"
)

func All() []Provider {
	return []Provider{
		&aws.AWS{},
		&gcloud.GCloud{},
		&do.DO{},
	}
}
