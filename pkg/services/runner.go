package services

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/projdocs/cli/internal"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg/docker"
	"github.com/projdocs/cli/pkg/types"
)

type Builder struct {
	built        bool
	constructors []types.ServiceConstructor
	docker       *docker.Client
}

func NewRunner(docker *docker.Client, constructors ...types.ServiceConstructor) *Builder {
	return &Builder{
		built:        false,
		constructors: constructors,
		docker:       docker,
	}
}

func (r *Builder) Register(constructor types.ServiceConstructor) error {
	if r.built {
		return fmt.Errorf("already built")
	}
	r.constructors = append(r.constructors, constructor)
	return nil
}

func (r *Builder) MustRegister(constructor types.ServiceConstructor) *Builder {
	if err := r.Register(constructor); err != nil {
		panic(err)
	}
	return r
}

func (r *Builder) Build(cfg config.Config) *Runner {

	services := make([]*types.ServiceConstructorResult, len(r.constructors))
	for i, constructor := range r.constructors {
		services[i] = constructor(cfg)
	}

	return &Runner{
		docker:   r.docker,
		started:  false,
		services: services,
	}
}

type Runner struct {
	started  bool
	services []*types.ServiceConstructorResult
	docker   *docker.Client
}

func (r *Runner) Stop() {

	spin := internal.NewSpinner("Shutting down...")
	var errors []error
	for i := len(r.services) - 1; i >= 0; i-- {
		service := r.services[i]
		spin.Update(fmt.Sprintf("Starting %s...", service.Container.Name))
		if err := r.docker.Stop(context.Background(), service.Container.Name); err != nil {
			errors = append(errors, fmt.Errorf("could not stop container %s: %w", service.Container.Name, err))
		}
	}

	if len(errors) == 0 {
		spin.Success("Stopped successfully!")
	} else {
		spin.Fail("Stopped with errors: ")
		for _, err := range errors {
			color.Red(err.Error())
		}
	}
}

func (r *Runner) Start(ctx context.Context) error {

	if r.started {
		return fmt.Errorf("already started")
	}
	r.started = true

	spin := internal.NewSpinner("Ensuring docker network...")
	err := r.docker.EnsureNetwork(ctx)
	if err != nil {
		spin.Fail("could not ensure docker network exists!")
		return err
	}

	spin.Update("Starting services...")
	for _, service := range r.services {
		spin.Update(fmt.Sprintf("Starting %s...", service.Container.Name))
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// create container
			if containerID, err := r.docker.Create(ctx, service); err != nil {
				if containerID != nil {
					// created, but with errors
					// attempt to stop/cleanup
					r.docker.Stop(ctx, *containerID)
				}
				return err
			} else
			// start container
			if err := r.docker.Start(ctx, *containerID); err != nil {
				color.Yellow("failed to start %s docker container: %s", service.Container.Name, err)
			} else {
				healthy := true
				if service.Container.Config != nil && service.Container.Config.Healthcheck != nil {
					if hcErr := r.docker.InspectContainer(ctx, *containerID); hcErr != nil {
						color.Yellow("failed to inspect %s docker container: %s", service.Container.Name, hcErr)
						healthy = false
					}
				}
				if healthy && service.AfterStartExec != nil {
					if output, err := r.docker.ExecInContainer(ctx, *containerID, service.AfterStartExec); err != nil {
						color.Yellow("%s after-start hook failed: %s", service.Container.Name, err)
						color.Red(output)
					}
				}
			}
		}
	}
	spin.Success("All services started!")
	return nil
}
