package exec

import (
	"bufio"
	"io"
	"os/exec"
)

type Executor interface {
	Command(string, ...string)
	Start() error
	Wait() error
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
}

var (
	_ = Executor(&RealExecutor{})
)

type RealExecutor struct {
	*exec.Cmd
}

func (e *RealExecutor) Command(bin string, args ...string) {
	e.Cmd.Path = bin
	e.Cmd.Args = args
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
