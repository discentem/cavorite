package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/discentem/cavorite/internal/stores"
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
	ErrValidateNil                   = errors.New("cavorite config must have a Validate() function")
	ErrValidate                      = errors.New("validate() failed")
	ErrDirExpander                   = errors.New("dirExpander failed")
	ErrUnsupportedStore              = errors.New("not a supported store type")
	ErrConfigNotExist                = errors.New("config file does not exist")
	ErrConfigDirNotExist             = errors.New("config directory does not exist")
)

func InitializeStoreTypeS3(
	ctx context.Context,
	fsys afero.Fs,
	sourceRepo, backendAddress, awsRegion string,
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

func InitializeStoreTypeGCS(
	ctx context.Context,
	fsys afero.Fs,
	sourceRepo, backendAddress string,
	opts stores.Options,
) Config {
	return Config{
		StoreType: stores.StoreTypeGCS,
		Options:   opts,
		Validate: func() error {
			return nil
		},
	}
}

func InitializeStoreTypeAzureBlob(
	ctx context.Context,
	fsys afero.Fs,
	sourceRepo, backendAddress string,
	opts stores.Options,
) Config {
	return Config{
		StoreType: stores.StoreTypeAzureBlob,
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
