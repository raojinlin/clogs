package docker

import "github.com/docker/docker/client"

func NewCli() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
