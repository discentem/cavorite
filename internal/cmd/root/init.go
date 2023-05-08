package root

import (
	"context"
	"errors"

	"github.com/discentem/cavorite/internal/config"
	"github.com/discentem/cavorite/internal/stores"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getInitCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Pantri repo",
		Long:  "Initialize a new Pantri repo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setLoggerOpts()

			repoToInit := args[0]
			backend := viper.GetString("store_type")
			fileExt := viper.GetString("metadata_file_extension")
			pantriAddress := viper.GetString("backend_address")
			region := viper.GetString("region")

			opts := stores.Options{
				PantriAddress:         pantriAddress,
				MetaDataFileExtension: fileExt,
				Region:                region,
			}

			ctx := context.Background()
			fsys := afero.NewOsFs()

			var storeType stores.StoreType
			var cfg config.Config
			switch storeType.FromString(backend) {
			case stores.StoreTypeS3:
				cfg = config.InitializeStoreTypeS3Config(
					ctx,
					fsys,
					repoToInit,
					pantriAddress,
					region,
					opts,
				)
			default:
				return config.ErrUnsupportedStore
			}
			return cfg.Write(fsys, repoToInit)
		},
	}
	cmd.PersistentFlags().String("backend_address", "", "Address for the storage backend")
	cmd.PersistentFlags().String("region", "us-east-1", "Default region for the storage backend")
	cmd.PersistentFlags().String("store_type", "", "Storage backend to use")
	// Bind all the flags to a viper setting so we can use viper everywhere without thinking about it
	if err := viper.BindPFlag("backend_address", cmd.PersistentFlags().Lookup("backend_address")); err != nil {
		return nil, errors.New("Failed to bind backend_address to viper")
	}
	if err := viper.BindPFlag("region", cmd.PersistentFlags().Lookup("region")); err != nil {
		return nil, errors.New("Failed to bind region to viper")
	}
	if err := viper.BindPFlag("store_type", cmd.PersistentFlags().Lookup("store_type")); err != nil {
		return nil, errors.New("Failed to bind store_type to viper")
	}
	return cmd, nil
}
