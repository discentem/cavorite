package fileutils

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/discentem/cavorite/internal/bindetector"
	"github.com/google/logger"
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

func GetBinariesWalkPath(fsys afero.Fs, objects []string) ([]string, error) {
	afs := afero.Afero{Fs: fsys}

	var paths []string
	logger.V(2).Infof("walking objects: %s", objects)
	for _, object := range objects {
		err := afs.Walk(object, func(p string, info fs.FileInfo, err error) error {
			logger.V(2).Infof("walkfunc() p: %s", p)
			isDir, err := afs.IsDir(p)
			if err != nil {
				logger.Errorf("walkfunc() fsys.IsDir error: %v", err)
				return err
			}
			logger.V(2).Infof("walkfunc() isDir: %t", isDir)

			if isDir {
				return nil
			}

			if bindetector.IsBinary(p) {
				logger.V(2).Infof("walkfunc() - bindetector().IsBinary appending path: %s", p)
				paths = append(paths, p)
			}

			return nil
		})
		if err != nil {
			return paths, err
		}
	}

	return paths, nil
}
