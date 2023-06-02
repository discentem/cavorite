package testutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

type MapFile struct {
	Content []byte
	ModTime *time.Time
}

// memMapFsWith creates a afero.MemMapFs from a map[string]mapFile
func MemMapFsWith(files map[string]MapFile) (*afero.Fs, error) {
	memfsys := afero.NewMemMapFs()

	mfs := memfsys.(*afero.MemMapFs)

	for fname, mfile := range files {
		err := memfsys.MkdirAll(filepath.Dir(fname), os.ModeDir)
		if err != nil {
			return nil, err
		}
		afile, err := memfsys.Create(fname)
		if err != nil {
			return nil, err
		}
		_, err = afile.Write(mfile.Content)
		if err != nil {
			return nil, err
		}
		if mfile.ModTime != nil {
			err := memfsys.Chtimes(fname, time.Time{}, *mfile.ModTime)
			if err != nil {
				return nil, err
			}
		}
	}
	fmt.Println("listing mfs")
	err := afero.Walk(mfs, "/", func(path string, info fs.FileInfo, err error) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &memfsys, nil
}
