package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/stores"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Type          string         `json:"type"`
	PantriAddress string         `json:"pantri"`
	Opts          stores.Options `json:"options"`
}

func (c *Config) WriteToDisk(sourceRepo string) error {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	cfile, err := homedir.Expand(
		fmt.Sprintf("%s/.pantri/config", sourceRepo),
	)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfile), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(cfile, b, os.ModePerm)
}
