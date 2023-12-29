package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

func ListContainers(opts *types.ContainerListOptions) ([]types.Container, error) {
	ctx := context.Background()
	cli, err := NewCli()
	if err != nil {
		return nil, err
	}

	return cli.ContainerList(ctx, *opts)
}