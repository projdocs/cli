package docker

import (
	"context"

	"github.com/moby/moby/client"
)

func (docker *Client) Start(ctx context.Context, containerID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, err := docker.api.ContainerStart(ctx, containerID, client.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
}
