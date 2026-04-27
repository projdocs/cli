package caddy

import (
	_ "embed"
	"net/netip"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/config"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/types"
	"github.com/projdocs/cli/pkg/types/embeds"
)

//go:embed Caddyfile
var Caddyfile []byte

var ServiceConstructor types.ServiceConstructor = func(cfg config.Config) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Embeds: []*embeds.EmbeddedFile{
			{
				Path: embeds.MustParsePath("/etc/caddy/Caddyfile"),
				Data: Caddyfile,
			},
		},
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-caddy",
			Image: "caddy/caddy:2.11.2-alpine",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Env: []string{
					"PROXY_DOMAIN=" + cfg.File.URLs.Supabase,
					"WEB_DOMAIN=" + cfg.File.URLs.Web,
				},
				Cmd: []string{
					"/bin/sh",
					"-c",
					`caddy run --config /etc/caddy/Caddyfile --adapter caddyfile`,
				},
			},
			HostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
				PortBindings: network.PortMap{
					network.MustParsePort("80/tcp"):  []network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "80"}},
					network.MustParsePort("443/tcp"): []network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "443"}},
					network.MustParsePort("443/udp"): []network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "443"}},
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeVolume,
						Source: "caddy_data",
						Target: "/data",
					},
					{
						Type:   mount.TypeVolume,
						Source: "caddy_config",
						Target: "/config",
					},
				},
			},
			NetworkingConfig: &network.NetworkingConfig{},
		},
	}
}
