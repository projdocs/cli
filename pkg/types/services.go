package types

import (
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg/types/embeds"
)

type ServiceConstructorResult struct {
	Embeds         []*embeds.EmbeddedFile
	Container      *client.ContainerCreateOptions
	AfterStartExec []string
}

type ServiceConstructor = func(cfg config.Config) *ServiceConstructorResult
