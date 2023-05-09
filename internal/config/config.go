package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/google/logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

type Config struct {
	StoreType stores.StoreType `json:"store_type" mapstructure:"store_type"`
	Options   stores.Options   `json:"options" mapstructure:"options"`
	Validate  func() error     `json:"-"`
}

type dirExpander func(string) (string, error)

var (
	// dirExpander can be overwritten for tests
	// expander can be overwritten with fakes for tests
	expander             dirExpander = homedir.Expand
	ErrValidateNil                   = errors.New("pantri config must have a Validate() function")
	ErrValidate                      = errors.New("validate() failed")
	ErrDirExpander                   = errors.New("dirExpander failed")
	ErrUnsupportedStore              = errors.New("not a supported store type")
	ErrConfigNotExist                = errors.New("config file does not exist")
	ErrConfigDirNotExist             = errors.New("config directory does not exist")
)

func InitializeStoreTypeS3(
	ctx context.Context,
	fsys afero.Fs,
	sourceRepo, pantriAddress, awsRegion string,
	opts stores.Options,
) Config {
	return Config{
		StoreType: stores.StoreTypeS3,
		Options:   opts,
		Validate: func() error {
			return nil
		},
	}
}

func (c *Config) Write(fsys afero.Fs, sourceRepo string) error {
	if c.Validate == nil {
		return ErrValidateNil
	}
	if err := c.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidate, err)
	}
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	esr, err := expander(sourceRepo)
	if err != nil {
		return fmt.Errorf(
			"%w: %v",
			ErrDirExpander,
			err,
		)
	}
	cfile := filepath.Join(esr, ".pantri/config")
	if _, err := fsys.Stat(esr); err != nil {
		return fmt.Errorf("%s does not exist, so we can't make it a pantri repo", esr)
	}

	if err := fsys.MkdirAll(filepath.Dir(cfile), os.ModePerm); err != nil {
		return err
	}
	logger.Infof("initializing pantri config at %s", cfile)
	return afero.WriteFile(fsys, cfile, b, os.ModePerm)
}

func ReadConfig(fsys afero.Fs, sourceRepo string) ([]byte, error) {
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
