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
	for _, pipe := range inputPipes {
		for _, post := range posters {
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
