package dkr

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

type DockerClient struct {
	api client.APIClient
}

func NewClient(docker client.APIClient) *DockerClient {
	return &DockerClient{
		api: docker,
	}
}

func (docker *DockerClient) Ping(ctx context.Context) error {
	if _, err := docker.api.Ping(ctx, client.PingOptions{}); err != nil {
		return fmt.Errorf("could not ping docker client: %w", err)
	}
	return nil
}
