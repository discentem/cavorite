package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"

	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/stores"
)

type LocalStore struct {
	logger hclog.Logger
	fsys   afero.Fs
	opts   *stores.Options
}

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

func (s *LocalStore) Retrieve(ctx context.Context, objects ...string) error {
	s.logger.Info(fmt.Sprintf("Retrieving %v via localstore plugin", objects))
	return nil
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
