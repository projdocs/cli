package storage

import (
	"strconv"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/docker"
	"github.com/projdocs/cli/pkg/services/postgres"
	"github.com/projdocs/cli/pkg/types"
)

var ServiceConstructor types.ServiceConstructor = func(cfg config.Config) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-storage",
			Image: "supabase/storage-api:v1.48.26",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Env: []string{
					"STORAGE_BACKEND=s3",
					"GLOBAL_S3_ENDPOINT=http://minio:9000",
					"GLOBAL_S3_PROTOCOL=http",
					"GLOBAL_S3_FORCE_PATH_STYLE=true",
					"AWS_ACCESS_KEY_ID=supa-storage",
					"AWS_SECRET_ACCESS_KEY=" + cfg.Supabase.Storage.MinioPassword,
					"ANON_KEY=" + cfg.Supabase.Keys.Anon.Symmetric,
					"SERVICE_KEY=" + cfg.Supabase.Keys.Service.Symmetric,
					"POSTGREST_URL=http://rest:3000",
					"AUTH_JWT_SECRET=" + cfg.Supabase.Keys.JWTSecret,
					"DATABASE_URL=postgres://supabase_storage_admin:" + cfg.Supabase.Postgres.Password + "@" + postgres.ContainerName + ":5432/postgres",
					"STORAGE_PUBLIC_URL=" + cfg.File.URLs.Web + ":8000",
					"REQUEST_ALLOW_X_FORWARDED_PATH=true",
					"FILE_SIZE_LIMIT=" + strconv.Itoa(cfg.File.Storage.SizeLimitInMB*1024*1024), // convert to bytes
					"STORAGE_BACKEND=file",
					"GLOBAL_S3_BUCKET=supabase",
					"FILE_STORAGE_BACKEND_PATH=/var/lib/storage",
					"TENANT_ID=supabase",
					"REGION=local",
					"ENABLE_IMAGE_TRANSFORMATION=true",
					"IMGPROXY_URL=http://imgproxy:5001",
					"S3_PROTOCOL_ACCESS_KEY_ID=" + cfg.Supabase.Storage.S3AccessKeyID,
					"S3_PROTOCOL_ACCESS_KEY_SECRET=" + cfg.Supabase.Storage.S3SecretKey,
				},
				Healthcheck: &container.HealthConfig{
					Test: []string{"CMD",
						"wget",
						"--no-verbose",
						"--tries=1",
						"--spider",
						"http://storage:5000/status"},
					Interval: 5 * time.Second,
					Timeout:  5 * time.Second,
					Retries:  3,
				},
			},
			HostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
			},
			NetworkingConfig: docker.MakeNetworkConfig("storage"),
		},
	}
}
