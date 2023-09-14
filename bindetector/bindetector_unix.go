//go:build darwin || freebsd || linux
// +build darwin freebsd linux

package bindetector

import (
	"bytes"
	"os/exec"

	"github.com/google/logger"
)

func execFile(filepath string) []byte {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(
		"/usr/bin/file",
		"--mime-encoding",
		"-b",
		filepath,
	)

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		logger.Errorf("isBinary cmd run error: %v", err)
	}

	return stdoutBuf.Bytes()

}
