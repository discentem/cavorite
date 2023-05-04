package pantri

import (
	"context"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Pantri repo or directory",
	Long:  "Initialize a new Pantri repo or directory",
	Args:  cobra.ExactArgs(1),
	RunE:  Init,
}

func Init(cmd *cobra.Command, args []string) error {
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
}

func init() {
	initCmd.PersistentFlags().String("backend_address", "", "Address for the storage backend")
	initCmd.PersistentFlags().String("region", "us-east-1", "Default region for the storage backend")
	initCmd.PersistentFlags().String("store_type", "", "Storage backend to use")
	// Bind all the flags to a viper setting so we can use viper everywhere without thinking about it
	viper.BindPFlag("backend_address", initCmd.PersistentFlags().Lookup("backend_address"))
	viper.BindPFlag("region", initCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("store_type", initCmd.PersistentFlags().Lookup("store_type"))
}
