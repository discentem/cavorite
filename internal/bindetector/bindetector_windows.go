//go:build windows
// +build windows

package bindetector

func execFile(filepath string) []byte {
	var stdoutBuf bytes.Buffer
	// cmd := exec.Command(
	// 	"/usr/bin/file",
	// 	"--mime-encoding",
	// 	"-b",
	// 	filepath,
	// )

	// cmd.Stdout = &stdoutBuf
	// cmd.Stderr = &stderrBuf

	// err := cmd.Run()
	// if err != nil {
	// 	logger.Errorf("isBinary cmd run error: %v", err)
	// }

	// return stdoutBuf.Bytes()

	return stdoutBuf.Bytes()

}
