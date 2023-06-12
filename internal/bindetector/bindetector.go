package bindetector

import (
	"regexp"

	"github.com/google/logger"
)

func IsBinary(filepath string) bool {
	// for now, we must shell out to /usr/bin/file or respective windows exe called file.exe provided by git
	// someday, binary determination could be implemented in-house without the need for external tools but until then
	// shelling out is a necessary evil
	return isBinary(execFile(filepath))
}

func isBinary(bytes []byte) bool {
	matched, err := regexp.Match(`binary`, bytes)
	if err != nil {
		logger.Errorf("isBinary regex matching error: %v", err)
	}

	return matched
}
