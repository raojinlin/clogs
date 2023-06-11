package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"io"
)

func Logs(container string, logFile string, options *types.ContainerLogsOptions) (io.ReadCloser, error) {
	ctx := context.Background()
	cli, err := NewCli()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	if options == nil {
		options = &types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Tail: "50"}
	}

	if logFile == "stdout" {
		return cli.ContainerLogs(ctx, container, *options)
	}

	logFileCmd := []string{"tail"}
	if options.Follow {
		logFileCmd = append(logFileCmd, "-f")
	}

	if options.Tail != "" {
		logFileCmd = append(logFileCmd, "-n", options.Tail)
	}

	logFileCmd = append(logFileCmd, logFile)
	return Exec(container, logFileCmd)
}
