package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

const NetworkName = "projdocs-net"

func MakeNetworkConfig(aliases ...string) *network.NetworkingConfig {
	cfg := make(map[string]*network.EndpointSettings, 1)
	cfg[NetworkName] = &network.EndpointSettings{
		Aliases: aliases,
	}
	return &network.NetworkingConfig{
		EndpointsConfig: cfg,
	}
}

func (docker *Client) EnsureNetwork(ctx context.Context) error {
	_, err := docker.api.NetworkInspect(ctx, NetworkName, client.NetworkInspectOptions{Verbose: true})
	if err != nil {
		if isNetworkNotFoundErr(err) {
			if _, err := docker.api.NetworkCreate(ctx, NetworkName, client.NetworkCreateOptions{
				Driver:     "bridge",
				Scope:      "local",
				EnableIPv4: new(true),
				EnableIPv6: new(true),
				Internal:   false, // true = no external connectivity (usually keep false)
				Attachable: true,  // allow standalone containers to attach/detach
			}); err != nil {
				return fmt.Errorf("failed to create network: %v", err)
			}
		} else {
			return fmt.Errorf("failed to inspect network: %v", err)
		}
	}
	return nil
}
