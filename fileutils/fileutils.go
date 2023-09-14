package fileutils

import (
	"fmt"
	"io"
	"os"

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

func OpenOrCreateFile(fsys afero.Fs, filename string) (afero.File, error) {
	file, err := fsys.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
