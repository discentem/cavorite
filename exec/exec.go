package exec

import (
	"os/exec"
)

type ExecCommand struct {
	exec.Cmd
}

func Command(name string, arg ...string) *ExecCommand {
	cmd := exec.Command(
		name,
		arg...,
	)

	return &ExecCommand{
		*cmd,
	}
}

func (c *ExecCommand) Run() error {
	return c.Cmd.Run()
}
