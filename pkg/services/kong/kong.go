package kong

import (
	_ "embed"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/internal/docker"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/types"
	"github.com/projdocs/cli/pkg/types/embeds"
)

//go:embed kong.yml
var Yaml []byte

//go:embed kong-entrypoint.sh
var Entrypoint []byte

var ServiceConstructor types.ServiceConstructor = func(cfg config.Config) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Embeds: []*embeds.EmbeddedFile{
			{
				Path: embeds.MustParsePath("/home/kong/temp.yml"),
				Data: Yaml,
			},
			{
				Path: embeds.MustParsePath("/home/kong/kong-entrypoint.sh"),
				Data: Entrypoint,
			},
		},
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-kong",
			Image: "kong/kong:3.9.1",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Entrypoint: []string{"/home/kong/kong-entrypoint.sh"},
				Healthcheck: &container.HealthConfig{
					Test:     []string{"CMD", "kong", "health"},
					Interval: 5 * time.Second,
					Timeout:  5 * time.Second,
					Retries:  5,
				},
				Env: []string{
					`KONG_DATABASE=off`,
					`KONG_DECLARATIVE_CONFIG=/usr/local/kong/kong.yml`,
					`KONG_DNS_ORDER=LAST,A,CNAME`,
					`KONG_DNS_NOT_FOUND_TTL=1`,
					`KONG_PLUGINS=request-transformer,cors,key-auth,acl,basic-auth,request-termination,ip-restriction,post-function`,
					`KONG_NGINX_PROXY_PROXY_BUFFER_SIZE=160k`,
					`KONG_NGINX_PROXY_PROXY_BUFFERS=64 160k`,
					`KONG_PROXY_ACCESS_LOG=/dev/stdout combined`,
					"SUPABASE_ANON_KEY=" + cfg.Supabase.Keys.Anon.Symmetric,
					"SUPABASE_SERVICE_KEY=" + cfg.Supabase.Keys.Service.Symmetric,
					"SUPABASE_PUBLISHABLE_KEY=" + cfg.Supabase.Keys.Publishable,
					"SUPABASE_SECRET_KEY=" + cfg.Supabase.Keys.Secret,
					"ANON_KEY_ASYMMETRIC=" + cfg.Supabase.Keys.Anon.Asymmetric,
					"SERVICE_ROLE_KEY_ASYMMETRIC=" + cfg.Supabase.Keys.Service.Asymmetric,
					"DASHBOARD_USERNAME=supabase",
					"DASHBOARD_PASSWORD=" + cfg.Supabase.Studio.Password,
				},
			},
			HostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
				//PortBindings: network.PortMap{
				//	network.MustParsePort("8000/tcp"): []network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "8000"}},
				//	network.MustParsePort("8443/tcp"): []network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "8443"}},
				//},
			},
			NetworkingConfig: docker.MakeNetworkConfig("kong", "api-gw"),
		},
	}
}
