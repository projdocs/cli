package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/moby/moby/client"
	"github.com/projdocs/cli/pkg/types"
	"github.com/projdocs/cli/pkg/types/embeds"
)

func (docker *Client) Create(ctx context.Context, constructor *types.ServiceConstructorResult) (*string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

		// ensure image exists
		if _, err := docker.api.ImageInspect(ctx, constructor.Container.Image); err != nil {
			if isImageNotFoundErr(err) {
				if rc, err := docker.api.ImagePull(ctx, constructor.Container.Image, client.ImagePullOptions{}); err != nil {
					return nil, fmt.Errorf("unable to pull image %s: %w", constructor.Container.Image, err)
				} else {
					defer rc.Close()
					_, err = io.Copy(io.Discard, rc)
					if err != nil {
						return nil, fmt.Errorf("unable to pull image %s: %w", constructor.Container.Image, err)
					}
				}
			} else {
				return nil, fmt.Errorf("unable to inspect image: %w", err)
			}
		}

		// create container
		// TODO: consider removing and moving this logic to Runner
		ctr, err := docker.api.ContainerCreate(ctx, *constructor.Container)
		if err != nil {
			if isNameConflictErr(err) {
				// remove dangling container
				if _, e := docker.api.ContainerRemove(ctx, constructor.Container.Name, client.ContainerRemoveOptions{
					RemoveVolumes: false,
					Force:         true,
				}); e != nil {
					return nil, fmt.Errorf("container name ('%v') conflicts and could not be removed: %v", constructor.Container.Name, e)
				}

				// retry to create container
				ctr, err = docker.api.ContainerCreate(ctx, *constructor.Container)
				if err != nil {
					return nil, fmt.Errorf("could not create container after cleaning-up dangling contained: %v", err)
				}
			} else {
				return nil, err
			}
		}

		// write any embedded files
		if constructor.Embeds != nil && len(constructor.Embeds) > 0 {
			for f, file := range constructor.Embeds {
				if _, err := docker.copyToContainer(ctx, ctr.ID, file); err != nil {
					return &ctr.ID, fmt.Errorf("failed to create file %d in container '%s': %v", f, constructor.Container.Name, err)
				}
			}
		}

		return &ctr.ID, nil
	}
}

func (docker *Client) copyToContainer(ctx context.Context, containerID string, file *embeds.EmbeddedFile) (*client.CopyToContainerResult, error) {

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	now := time.Now()

	// Emit mkdir -p style directory headers (0755)
	if file.Path.Parent() != "." && file.Path.Parent() != "/" {
		parts := strings.Split(file.Path.Parent(), "/")
		cur := ""
		for _, seg := range parts {
			if seg == "" {
				continue
			}
			if cur == "" {
				cur = seg
			} else {
				cur = cur + "/" + seg
			}
			hdr := &tar.Header{
				Name:     cur + "/", // relative dir entry
				Typeflag: tar.TypeDir,
				Mode:     0o755, // drwxr-xr-x
				ModTime:  now,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				_ = tw.Close()
				return nil, fmt.Errorf("write dir header %q: %w", cur, err)
			}
		}
	}

	// File header at full relative path (parent/base) with 0644
	fhdr := &tar.Header{
		Name:     path.Join(file.Path.Parent(), file.Path.Base()),
		Typeflag: tar.TypeReg,
		Mode:     0o555, // xr-xr-xr-
		Size:     int64(len(file.Data)),
		ModTime:  now,
	}
	if err := tw.WriteHeader(fhdr); err != nil {
		_ = tw.Close()
		return nil, fmt.Errorf("write file header %q: %w", file.Path, err)
	}
	if _, err := tw.Write(file.Data); err != nil {
		_ = tw.Close()
		return nil, fmt.Errorf("write file data %q: %w", file.Path, err)
	}
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("close tar: %w", err)
	}

	// Extract under "/" so the tar's relative paths land at absolute locations
	cpy, cpyErr := docker.api.CopyToContainer(ctx, containerID, client.CopyToContainerOptions{
		DestinationPath:           "/",
		Content:                   bytes.NewReader(buf.Bytes()),
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                false, // do NOT preserve uid/gid from headers; use container defaults
	})
	return &cpy, cpyErr

}
