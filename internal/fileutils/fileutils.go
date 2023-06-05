package fileutils

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
)

func BytesFromAferoFile(f afero.File) ([]byte, error) {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	objInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	b := make([]byte, objInfo.Size())
	_, err = f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from objectHandle: %w", err)
	}
	return b, nil
}
