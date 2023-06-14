package fileutils

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/discentem/cavorite/internal/bindetector"
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

func GetBinariesWalkPath(fsys afero.Afero, path string) ([]string, error) {
	var paths []string

	err := fsys.Walk(path, func(p string, info fs.FileInfo, err error) error {
		isDir, err := fsys.IsDir(path)
		if err != nil {
			return err
		}

		if isDir {
			return nil
		}

		if bindetector.IsBinary(p) {
			paths = append(paths, p)
		}

		return nil
	})

	return paths, err
}
