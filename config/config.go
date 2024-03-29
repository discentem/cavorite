package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/discentem/cavorite/stores"
	"github.com/google/logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

type Config struct {
	StoreType stores.StoreType             `json:"store_type" mapstructure:"store_type"`
	Options   stores.Options               `json:"options" mapstructure:"options"`
	Validate  func() error                 `json:"-"`
	Expander  func(string) (string, error) `json:"-"`
	Marshal   func(v any) ([]byte, error)  `json:"-"`
}

var (
	Cfg Config
)

var (
	ErrValidateNil       = errors.New("cavorite config must have a Validate() function")
	ErrValidate          = errors.New("validate() failed")
	ErrDirExpander       = errors.New("dirExpander failed")
	ErrDirExpanderNil    = errors.New("dirExpander cannot be nil")
	ErrMarshalNil        = errors.New("marshal cannot be nil")
	ErrUnsupportedStore  = errors.New("not a supported store type")
	ErrConfigNotExist    = errors.New("config file does not exist")
	ErrConfigDirNotExist = errors.New("config directory does not exist")
)

func InitalizeStoreTypeOf(
	ctx context.Context,
	storeType stores.StoreType,
	fsys afero.Fs,
	sourceRepo string,
	opts stores.Options,
) Config {
	return Config{
		StoreType: storeType,
		Options:   opts,
		Validate: func() error {
			return nil
		},
		Marshal: func(v any) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		},
		Expander: homedir.Expand,
	}
}

func (c *Config) Write(fsys afero.Fs, sourceRepo string) error {
	if c.Validate == nil {
		return ErrValidateNil
	}
	if err := c.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidate, err)
	}

	if c.Marshal == nil {
		return ErrMarshalNil
	}

	b, err := c.Marshal(c)
	if err != nil {
		return err
	}
	if c.Expander == nil {
		return ErrDirExpanderNil
	}

	esr, err := c.Expander(sourceRepo)
	if err != nil {
		return fmt.Errorf(
			"%w: %v",
			ErrDirExpander,
			err,
		)
	}
	cfile := filepath.Join(esr, ".cavorite/config")
	if _, err := fsys.Stat(esr); err != nil {
		return fmt.Errorf("%s does not exist, so we can't make it a cavorite repo", esr)
	}

	if err := fsys.MkdirAll(filepath.Dir(cfile), os.ModePerm); err != nil {
		return err
	}
	logger.Infof("initializing cavorite config at %s", cfile)
	return afero.WriteFile(fsys, cfile, b, os.ModePerm)
}
