package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/ioutils"
	"io"
)

func Exec(container string, command []string) (io.ReadCloser, error) {
	cli, err := NewCli()
	if err != nil {
		return nil, err
	}

	execCreate, err := cli.ContainerExecCreate(context.Background(), container, types.ExecConfig{
		User:         "",
		Privileged:   false,
		Tty:          false,
		ConsoleSize:  nil,
		AttachStdin:  false,
		AttachStderr: true,
		AttachStdout: true,
		Detach:       false,
		DetachKeys:   "",
		Env:          nil,
		WorkingDir:   "",
		Cmd:          command,
	})

	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerExecAttach(context.Background(), execCreate.ID, types.ExecStartCheck{
		Detach:      false,
		Tty:         false,
		ConsoleSize: nil,
	})
	if err != nil {
		return nil, err
	}

	return ioutils.NewReadCloserWrapper(resp.Reader, func() error {
		return resp.Conn.Close()
	}), nil
}
