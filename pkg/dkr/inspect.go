package dkr

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
	"github.com/moby/moby/client"
)

func (docker *DockerClient) InspectContainer(ctx context.Context, containerID string) error {
	for retries := 0; retries < 5; retries++ {
		if inspect, inspectErr := docker.api.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{}); inspectErr != nil {
			color.Yellow("failed to inspect container: %s", inspectErr)
			break
		} else {
			if inspect.Container.State.Health.Status == "healthy" {
				return nil
			} else if inspect.Container.State.Health.Status == "unhealthy" {
				return errors.New("container is not healthy")
			} else if inspect.Container.State.Health.Status != "starting" {
				color.Yellow("container %s is not healthy (status=%s;retry=%d)", inspect.Container.Name, inspect.Container.State.Health.Status, retries)
			}
		}
		time.Sleep(time.Duration(math.Pow(2, float64(retries))) * time.Second)
	}
	return fmt.Errorf("container not healthy (retries exceeded)")
}
