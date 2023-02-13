package loader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/stores"
	localstore "github.com/discentem/pantri_but_go/stores/local"
	"github.com/discentem/pantri_but_go/stores/s3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

func Initialize(ctx context.Context, sourceRepo, backend, address string, opts stores.Options) error {
	switch b := (backend); b {
	case "local":
		_, err := localstore.New(sourceRepo, address, opts)
		if err != nil {
			return err
		}
	case "s3":
		_, err := s3.New(ctx, sourceRepo, address, opts)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a supported store type", b)
	}
	return nil

}

var (
	ErrConfigNotExist    = errors.New("because the config file does not exist")
	ErrConfigDirNotExist = errors.New("the config directory does not exist")
)

type dirExpander func(string) (string, error)

// expander can be overwritten with fakes for tests
var expander dirExpander = homedir.Expand

func readConfig(fsys afero.Fs, sourceRepo string) ([]byte, error) {
	cfile, err := expander(
		fmt.Sprintf("%s/.pantri/config", sourceRepo),
	)
	if err != nil {
		return nil, err
	}

	if _, err := fsys.Stat(filepath.Dir(cfile)); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"%s has not be initialized as a pantri repo: %w: %s",
				sourceRepo,
				ErrConfigDirNotExist,
				filepath.Dir(cfile),
			)
		}
		return nil, err
	}
	f, err := fsys.Open(cfile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"%s has not be initialized as a pantri repo: %w: %s",
				sourceRepo,
				ErrConfigNotExist,
				cfile,
			)
		}
		return nil, err
	}
	return afero.ReadAll(f)
}

func Load(fsys afero.Fs, sourceRepo string) (stores.Store, error) {
	b, err := readConfig(fsys, sourceRepo)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	switch t := (m["type"]); t {
	case "local":
		return localstore.Load(m)
	case "s3":
		return s3.Load(m)
	default:
		return nil, fmt.Errorf("%s is not a support store type", t)
	}
}
