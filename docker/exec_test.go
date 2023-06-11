package docker

import (
	"io"
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	r, err := Exec("fc74dae7159c", []string{"bash", "-c", "for i in {0..10}; do sleep 1;echo $i;done"})
	if err != nil {
		t.Fatal(err)
	}

	io.Copy(os.Stdout, r)
	r.Close()
}
