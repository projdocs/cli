package minio

import (
	"path/filepath"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/internal/docker"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/types"
)

var ServiceConstructor types.ServiceConstructor = func(cfg config.Config) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-minio",
			Image: "cgr.dev/chainguard/minio",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Env: []string{
					"MINIO_ROOT_USER=supa-storage",
					"MINIO_ROOT_PASSWORD=" + cfg.Supabase.Storage.MinioPassword,
				},
				Cmd: []string{
					"server",
					"--console-address",
					":9001",
					"/data",
				},
				Healthcheck: &container.HealthConfig{
					Test:     []string{"CMD", "mc", "ready", "local"},
					Interval: 2 * time.Second,
					Timeout:  10 * time.Second,
					Retries:  5,
				},
			},
			HostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: filepath.Join(config.GetConfigDirPath(), "storage"),
						Target: "/data",
					},
				},
			},
			NetworkingConfig: docker.MakeNetworkConfig("minio"),
		},
		AfterStartExec: []string{
			"sh", "-c",
			"mc alias set supa-minio http://minio:9000 supa-storage " + cfg.Supabase.Storage.MinioPassword + " && " +
				"mc mb --ignore-existing supa-minio/supabase",
		},
	}
}
