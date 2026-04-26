package types

import (
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/config"
	"github.com/projdocs/cli/pkg/types/embeds"
)

type ServiceConstructorResult struct {
	Embeds    []*embeds.EmbeddedFile
	Container *client.ContainerCreateOptions
}

type ServiceConstructor = func(cfg config.File) *ServiceConstructorResult
