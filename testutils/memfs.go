package testutils

import (
	"bytes"
	"io"
	iofs "io/fs"
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

	return &memfsys, nil
}

func WalkFs(fs afero.Fs, w io.Writer) error {
	return afero.Walk(fs, "", func(path string, info iofs.FileInfo, _ error) error {
		b, err := afero.ReadFile(fs, path)
		if bytes.Equal(b, []byte(``)) {
			return nil
		}
		if err != nil {
			return err
		}
		p := path + "/" + info.Name() + "\n" + string(b) + "\n"
		_, err = w.Write([]byte(p))
		if err != nil {
			return err
		}
		return nil
	})
}
