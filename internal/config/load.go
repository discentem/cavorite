package config

import (
	"github.com/discentem/cavorite/internal/stores"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func Load(fs afero.Fs, cfgDest ...*Config) error {
	// Defaults set here will be used if they do not exist in the config file
	viper.SetFs(fs)
	viper.SetDefault("store_type", stores.StoreTypeUndefined)
	viper.SetDefault("metadata_file_extension", "cfile")
	// Set up the config file details
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".cavorite")

	// Retrieve from EnvVars if they exist...
	viper.AutomaticEnv()

	// Read from the config file path
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// if cfgDest is passed, unmarshal into cfgDest[0]. cfgDest[1-len(cfgDest)] is always ignored.
	if cfgDest != nil {
		viper.Unmarshal(cfgDest[0])
	}
	return viper.Unmarshal(&Cfg)
}
