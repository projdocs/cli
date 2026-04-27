package services

import (
	"github.com/projdocs/cli/pkg/services/caddy"
	"github.com/projdocs/cli/pkg/services/kong"
	"github.com/projdocs/cli/pkg/services/minio"
	"github.com/projdocs/cli/pkg/services/postgres"
	"github.com/projdocs/cli/pkg/services/postgrest"
	"github.com/projdocs/cli/pkg/services/storage"
	"github.com/projdocs/cli/pkg/types"
)

func GetAll() []types.ServiceConstructor {
	return []types.ServiceConstructor{
		postgres.ServiceConstructor,
		kong.ServiceConstructor,
		caddy.ServiceConstructor,
		postgrest.ServiceConstructor,
		minio.ServiceConstructor,
		storage.ServiceConstructor,
	}
}
