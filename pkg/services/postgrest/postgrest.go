package postgrest

import (
	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	config2 "github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/internal/docker"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/services/postgres"
	"github.com/projdocs/cli/pkg/types"
)

var ServiceConstructor types.ServiceConstructor = func(cfg config2.Config) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-postgrest",
			Image: "postgrest/postgrest:v14.10",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Env: []string{
					"PGRST_ADMIN_SERVER_PORT=3001",

					fmt.Sprintf("PGRST_DB_URI=postgres://authenticator:%s@%s:5432/postgres", cfg.Supabase.Postgres.Password, postgres.ContainerName),
					"PGRST_DB_SCHEMAS=public",
					"PGRST_DB_MAX_ROWS=1000",
					"PGRST_DB_EXTRA_SEARCH_PATH=public",
					"PGRST_DB_ANON_ROLE=anon",

					fmt.Sprintf("PGRST_JWT_SECRET=%s", cfg.Supabase.Keys.JWTJWKS),
					"PGRST_DB_USE_LEGACY_GUCS=false",
					fmt.Sprintf("PGRST_APP_SETTINGS_JWT_SECRET=%s", cfg.Supabase.Keys.JWTSecret),
					"PGRST_APP_SETTINGS_JWT_EXP=3600",
				},
				Cmd: []string{
					"postgrest",
				},
				//Healthcheck: &container.HealthConfig{
				//	Test:     []string{"CMD", "bash", "-c", "exec 3<>/dev/tcp/localhost/3001 && echo -e 'GET /live HTTP/1.0\\r\\nHost: localhost\\r\\n\\r\\n' >&3 && head -1 <&3 | grep -q '200'"},
				//	Interval: 5 * time.Second,
				//	Timeout:  5 * time.Second,
				//	Retries:  10,
				//},
			},
			HostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
			},
			NetworkingConfig: docker.MakeNetworkConfig("rest"),
		},
	}
}
