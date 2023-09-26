package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/google/logger"
	"github.com/hashicorp/go-hclog"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"

	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/fileutils"
	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores"
)

type LocalStore struct {
	logger hclog.Logger
	fsys   afero.Fs
	opts   *stores.Options
}

var (
	_                   = stores.Store(&LocalStore{})
	ErrCfilesLengthZero = errors.New("at least one cfile must be specified")
)

func (s *LocalStore) Upload(ctx context.Context, objects ...string) error {
	s.logger.Info(fmt.Sprintf("Uploading %v via localstore plugin", objects))

	if !filepath.IsAbs(s.opts.BackendAddress) {
		return fmt.Errorf("s.Opts.BackendAddress %q is not absolute", s.opts.BackendAddress)
	}

	var result *multierr.Error
	for _, o := range objects {
		objp := path.Join(s.opts.BackendAddress, o)
		if err := s.fsys.MkdirAll(path.Dir(objp), os.ModePerm); err != nil {
			return err
		}
		var srcf afero.File
		srcf, err := s.fsys.Open(o)
		if err != nil {
			result = multierr.Append(err)
			continue
		}
		defer srcf.Close()
		dst, err := s.fsys.Create(objp)
		if err != nil {
			result = multierr.Append(err)
			continue
		}
		defer dst.Close()
		_, err = io.Copy(dst, srcf)
		result = multierr.Append(err)
		// fmt.Printf("final error: %v\n", err)
	}

	return result.ErrorOrNil()
}

func (s *LocalStore) Retrieve(ctx context.Context, mmap metadata.CfileMetadataMap, cfiles ...string) error {
	var result *multierr.Error
	if len(cfiles) == 0 {
		return ErrCfilesLengthZero
	}
	s.logger.Info("", mmap)
	for _, cfile := range cfiles {
		srcFilePath := filepath.Join(s.opts.BackendAddress, mmap[cfile].Name)
		s.logger.Info(fmt.Sprintf("srcFilePath: %s", srcFilePath))
		sf, err := s.fsys.Open(srcFilePath)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		defer sf.Close()
		_, err = sf.Seek(0, io.SeekStart)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		// create the destination file
		df, err := fileutils.OpenOrCreateFile(s.fsys, mmap[cfile].Name)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		defer df.Close()
		fmt.Printf("mmap: %s\n", mmap)
		m, ok := mmap[cfile]
		if !ok {
			result = multierr.Append(result, fmt.Errorf("%q not found in mmap", cfile))
			continue
		}
		// copy from backend
		logger.Infof("copying %q to %q", srcFilePath, mmap[cfile].Name)
		_, err = io.Copy(df, sf)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}

		matches, err := metadata.HashFromCfileMatches(s.fsys, cfile, m.Checksum)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		if !matches {
			fmt.Printf("hash for %s did not match expected hash (%q) in %q\n", mmap[cfile].Name, mmap[cfile].Checksum, cfile)
			if err := s.fsys.Remove(m.Name); err != nil {
				result = multierr.Append(result, err)
				continue
			}
			result = multierr.Append(result, metadata.ErrRetrieveFailureHashMismatch)
			continue
		}
	}
	return result.ErrorOrNil()
}

func (s *LocalStore) GetOptions() (stores.Options, error) {
	if s.opts != nil {
		return *s.opts, nil
	}
	opts, err := config.LoadOptions(s.fsys)
	if err != nil {
		return stores.Options{}, nil
	}
	s.opts = &opts
	return *s.opts, nil

}

func (s *LocalStore) Close() error {
	return nil
}

func main() {

	hlog := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	ls := &LocalStore{
		logger: hlog,
		fsys:   afero.NewOsFs(),
	}

	stores.ListenAndServePlugin(ls, hlog)
}
