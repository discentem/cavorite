package bindetector

import (
	shell "github.com/discentem/cavorite/exec"
)

func IsBinary(filepath string) (bool, error) {
	// for now, we must shell out to /usr/bin/file or respective windows exe called file.exe provided by git
	// someday, binary determination could be implemented in-house without the need for external tools but until then
	// shelling out is a necessary evil
	e := shell.NewRealExecutor()
	binary, err := fileIsABinary(e, filepath)
	if err != nil {
		return false, err
	}
	return binary, nil
}
