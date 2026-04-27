package dkr

import (
	"bytes"
	"context"
	"fmt"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func (docker *DockerClient) ExecInContainer(ctx context.Context, containerID string, cmd []string) (string, error) {
	execResp, err := docker.api.ExecCreate(ctx, containerID, client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          false,
	})
	if err != nil {
		return "", fmt.Errorf("ExecInContainer create failed: %w", err)
	}

	// attach
	att, err := docker.api.ExecAttach(ctx, execResp.ID, client.ExecAttachOptions{TTY: false})
	if err != nil {
		return "", fmt.Errorf("ExecInContainer attach failed: %w", err)
	}
	defer att.Close()

	// copy output to local stdout/stderr
	var buf bytes.Buffer
	_, _ = stdcopy.StdCopy(&buf, &buf, att.Reader)
	output := buf.String()

	// check exit code
	if inspect, err := docker.api.ExecInspect(ctx, execResp.ID, client.ExecInspectOptions{}); err != nil {
		return output, fmt.Errorf("ExecInContainer inspect failed: %w", err)
	} else {
		if inspect.ExitCode != 0 {
			return output, fmt.Errorf("ExecInContainer command exited with code %d", inspect.ExitCode)
		}
		return output, nil
	}
}
