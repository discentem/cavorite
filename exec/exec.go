package exec

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

type Executor interface {
	Command(path string, args ...string)
	Stream(posters ...io.WriteCloser) error
}

var (
	_ = Executor(&RealExecutor{})
)

type RealExecutor struct {
	*exec.Cmd
	posters []io.WriteCloser
}

type RealExecutorOption func(e *RealExecutor)

func WithPosters(posters ...io.WriteCloser) RealExecutorOption {
	return func(e *RealExecutor) {
		e.posters = posters
	}
}

func NewRealExecutor(opts ...RealExecutorOption) *RealExecutor {
	re := &RealExecutor{}
	for _, opt := range opts {
		opt(re)
	}
	return re
}

func (e *RealExecutor) Command(bin string, args ...string) {
	if e.Cmd == nil {
		e.Cmd = exec.Command(bin, args...)
		return
	}
	e.Cmd.Path = bin
	e.Cmd.Args = args
}

func (e *RealExecutor) Stream(posters ...io.WriteCloser) error {
	stdout, err := e.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := e.StderrPipe()
	if err != nil {
		return err
	}
	inputPipes := []io.ReadCloser{stdout, stderr}

	if err := e.Start(); err != nil {
		return err
	}

	if posters == nil {
		if e.posters == nil {
			posters = []io.WriteCloser{os.Stdout}
		} else {
			posters = e.posters
		}

	}

	for _, pipe := range inputPipes {
		for _, post := range posters {
			//nolint:errcheck
			go WriteOutput(pipe, post)
		}
	}

	if err := e.Wait(); err != nil {
		return err
	}
	return nil
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
