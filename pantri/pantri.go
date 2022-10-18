package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/stores"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Type          string         `json:"type"`
	PantriAddress string         `json:"pantri_address"`
	Opts          stores.Options `json:"options"`
	Validate      func() error   `json:"-"`
}

func (c *Config) WriteToDisk(sourceRepo string) error {
	if c.Validate == nil {
		return errors.New("pantri config must have a Validate() function")
	}
	if err := c.Validate(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	esr, err := homedir.Expand(sourceRepo)
	if err != nil {
		return err
	}
	cfile := filepath.Join(esr, "/.pantri/config")
	if _, err := os.Stat(esr); err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s does not exist, so we can't make it a pantri repo", esr)
	}

	if err := os.MkdirAll(filepath.Dir(cfile), os.ModePerm); err != nil {
		return err
	}
	log.Printf("initializing pantri config at %s", cfile)
	return os.WriteFile(cfile, b, os.ModePerm)
}
