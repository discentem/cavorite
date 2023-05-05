package pantri

import (
	"context"
	"fmt"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "retrieve a file from pantri",
	Long:  "retrieve a file from pantri",
	Args:  cobra.MinimumNArgs(1),
	RunE:  Retrieve,
}

func Retrieve(cmd *cobra.Command, objects []string) error {
	setLoggerOpts()
	var store stores.Store
	var cfg config.Config

	ctx := context.Background()
	fsys := afero.NewOsFs()

	// Unmarshal entire config from viper
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}
	// Because viper can't go from s3 => iota, let's cheese it
	cfg.StoreType = cfg.StoreType.FromString(viper.GetString("store_type"))

	opts := cfg.Options

	switch cfg.StoreType {
	case stores.StoreTypeS3:
		s3, err := stores.NewS3StoreClient(ctx, fsys, opts)
		if err != nil {
			return fmt.Errorf("improper stores.S3Client init: %v", err)
		}
		store = stores.Store(s3)
	default:
		return fmt.Errorf("type %s is not supported", cfg.StoreType.String())
	}

	// We need to remove the prefix from the path so it is relative
	objects, err = removePathPrefix(objects)
	if err != nil {
		return err
	}

	logger.V(2).Infof("Downloading file list: %v", objects)
	logger.Infof("Downloading files from: %s", store.GetOptions().PantriAddress)
	logger.Infof("Downloading file: %s", objects)
	if err := store.Retrieve(ctx, objects...); err != nil {
		return err
	}
	return nil
}
