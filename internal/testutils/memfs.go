package testutils

import (
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
	for fname, mfile := range files {
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
	return &memfsys, nil
}
