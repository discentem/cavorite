//go:build darwin || freebsd || linux
// +build darwin freebsd linux

package bindetector

import (
	"bytes"
	"errors"

	shell "github.com/discentem/cavorite/exec"
)

func fileIsABinary(executor shell.Executor, filepath string) (bool, error) {
	buf := new(bytes.Buffer)

	executor.Command(
		"/usr/bin/file",
		"--mime-encoding",
		"-b",
		filepath,
	)

	if err := executor.Stream(&nopWriteCloser{buf}); err != nil {
		return false, err
	}

	output := bytes.TrimSpace(buf.Bytes())
	if len(output) == 0 {
		return false, errors.New("fileIsABinary: no output from file command")
	}

	return bytes.Contains(output, []byte("binary")), nil
}

type nopWriteCloser struct {
	*bytes.Buffer
}

func (n *nopWriteCloser) Close() error { return nil }
