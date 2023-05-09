package cli

import (
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func loadConfig(fs afero.Fs) error {
	// Defaults set here will be used if they do not exist in the config file
	viper.SetFs(fs)
	viper.SetDefault("store_type", stores.StoreTypeUndefined)
	viper.SetDefault("metadata_file_extension", "pfile")
	// Set up the config file details
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".pantri")

	// Retrieve from EnvVars if they exist...
	viper.AutomaticEnv()

	// Read from the config file path
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&cfg)
}
