package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

type Client struct {
	api client.APIClient
}

func NewClient(docker client.APIClient) *Client {
	return &Client{
		api: docker,
	}
}

func (docker *Client) Ping(ctx context.Context) error {
	if _, err := docker.api.Ping(ctx, client.PingOptions{}); err != nil {
		return fmt.Errorf("could not ping docker client: %w", err)
	}
	return nil
}
