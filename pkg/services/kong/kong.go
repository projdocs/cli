package kong

import (
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/config"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/types"
)

var ServiceConstructor types.ServiceConstructor = func(cfg config.File) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Embeds: nil,
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-kong",
			Image: "kong/kong:3.9.1",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Env: []string{
					"KONG_PORT_MAPS=443:8000,443:8443",
				},
			},
			HostConfig: &container.HostConfig{
				PortBindings: network.PortMap{},
			},
			NetworkingConfig: &network.NetworkingConfig{},
		},
	}
}
