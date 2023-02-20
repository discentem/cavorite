package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/stores"
	"github.com/google/logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

type Config struct {
	Type          string         `json:"type"`
	PantriAddress string         `json:"pantri_address"`
	Opts          stores.Options `json:"options"`
	Validate      func() error   `json:"-"`
}

type dirExpanderer func(string) (string, error)

// dirExpander can be overwritten for tests
var dirExpander dirExpanderer = homedir.Expand

var ErrValidateNil = errors.New("pantri config must have a Validate() function")
var ErrValidate = errors.New("validate() failed")
var ErrDirExpander = errors.New("dirExpander failed")

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
	esr, err := dirExpander(sourceRepo)
	if err != nil {
		return fmt.Errorf(
			"%w: %v",
			ErrDirExpander,
			err,
		)
	}
	cfile := filepath.Join(esr, ".pantri/config")
	if _, err := os.Stat(esr); err != nil {
		return fmt.Errorf("%s does not exist, so we can't make it a pantri repo", esr)
	}

	if err := os.MkdirAll(filepath.Dir(cfile), os.ModePerm); err != nil {
		return err
	}
	logger.Infof("initializing pantri config at %s", cfile)
	return os.WriteFile(cfile, b, os.ModePerm)
}
