//go:build windows
// +build windows

package bindetector

import "errors"

func fileIsABinary(filepath string) (*bool, error) {
	return nil, errors.New("fileIsABinary not implemented on Windows")

}
