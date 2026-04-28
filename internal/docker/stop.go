package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func (docker *Client) Stop(ctx context.Context, containerID string) error {

	if _, err := docker.api.ContainerStop(ctx, containerID, client.ContainerStopOptions{}); err != nil {
		return fmt.Errorf("could not shutdown %s: %w", containerID, err)
	}

	if _, err := docker.api.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("could not remove %s: %w", containerID, err)
	}

	return nil
}
