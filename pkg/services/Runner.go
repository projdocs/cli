package services

import (
	"context"
	"fmt"

	"github.com/projdocs/cli/config"
	"github.com/projdocs/cli/pkg/dkr"
	"github.com/projdocs/cli/pkg/types"
)

type Builder struct {
	built        bool
	constructors []types.ServiceConstructor
	docker       *dkr.DockerClient
}

func NewRunner(docker *dkr.DockerClient, constructors ...types.ServiceConstructor) *Builder {
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

func (r *Builder) Build(cfg config.File) *Runner {

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
	docker   *dkr.DockerClient
}

func (r *Runner) Start(ctx context.Context) error {

	if r.started {
		return fmt.Errorf("already started")
	}

	r.started = true
	for _, service := range r.services {
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
				// attempt to stop/cleanup
				r.docker.Stop(ctx, *containerID)
			}
		}
	}
	return nil
}
