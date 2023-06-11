package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"os"
	"testing"
)

func TestLogs(t *testing.T) {
	out, err := Logs("88eadee7d25f5f9419b16d48265a4367f195fa204583b9520cc1573d413d47c9", "stdout", &types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      "",
		Until:      "",
		Timestamps: false,
		Follow:     false,
		//Details:    true,
	})

	if err != nil {
		t.Fatal(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
