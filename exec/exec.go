package exec

import (
	"bufio"
	"io"
	"os/exec"
)

type Executor interface {
	Start() error
	Wait() error
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
}

type RealExecutor struct {
	*exec.Cmd
}

func WriteOutput(in io.ReadCloser, post io.WriteCloser) error {
	r := bufio.NewScanner(in)
	for r.Scan() {
		m := r.Text()
		_, err := post.Write([]byte(m + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
