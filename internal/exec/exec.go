package exec

import "os/exec"

type execCommand struct {
	exec.Cmd
}

func Command(name string, arg ...string) *execCommand {
	cmd := exec.Command(
		name,
		arg...,
	)

	return &execCommand{
		*cmd,
	}
}

func (c *execCommand) Run() error {
	return c.Run()
}
