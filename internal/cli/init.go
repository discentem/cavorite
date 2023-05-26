package cli

import (
	"errors"
	"fmt"

	"github.com/discentem/cavorite/internal/config"
	"github.com/discentem/cavorite/internal/program"
	"github.com/discentem/cavorite/internal/stores"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "init",
		Short:   fmt.Sprintf("Initialize a new %s repo", program.Name),
		Long:    fmt.Sprintf("Initialize a new %s repo", program.Name),
		Args:    cobra.ExactArgs(1),
		PreRunE: initPreExecFn,
		RunE:    initFn,
	}

	initCmd.PersistentFlags().String("backend_address", "", "Address for the storage backend")
	initCmd.PersistentFlags().String("region", "us-east-1", "Default region for the storage backend")
	initCmd.PersistentFlags().String("store_type", "", "Storage backend to use")

	return initCmd
}

// initPreExecFn is the pre-execution runtime for the initCmd functionality
// in Cobra, PreRunE is a concept of Cobra. It runs before RunE, similar to init() running before
/*
	// The *Run functions are executed in the following order:
	//   * PersistentPreRunE()
	//   * PreRunE() [X]
	//   * RunE()
	//   * PostRunE()
	//   * PersistentPostRunE()
	// All functions get the same args, the arguments after the command name.
*/
func initPreExecFn(cmd *cobra.Command, args []string) error {
	// Bind all the flags to a viper setting so we can use viper everywhere without thinking about it
	if err := viper.BindPFlag("backend_address", cmd.PersistentFlags().Lookup("backend_address")); err != nil {
		return errors.New("Failed to bind backend_address to viper")
	}
	if err := viper.BindPFlag("region", cmd.PersistentFlags().Lookup("region")); err != nil {
		return errors.New("Failed to bind region to viper")
	}
	if err := viper.BindPFlag("store_type", cmd.PersistentFlags().Lookup("store_type")); err != nil {
		return errors.New("Failed to bind store_type to viper")
	}

	return nil
}

// initFn is the execution runtime for the initCmd functionality
// in Cobra, this is the RunE phase
/*
	// The *Run functions are executed in the following order:
	//   * PersistentPreRunE()
	//   * PreRunE() []
	//   * RunE() [X]
	//   * PostRunE()
	//   * PersistentPostRunE()
	// All functions get the same args, the arguments after the command name.
*/
func initFn(cmd *cobra.Command, args []string) error {
	repoToInit := args[0]
	backend := viper.GetString("store_type")
	fileExt := viper.GetString("metadata_file_extension")
	backendAddress := viper.GetString("backend_address")
	region := viper.GetString("region")

	opts := stores.Options{
		BackendAddress:        backendAddress,
		MetadataFileExtension: fileExt,
		Region:                region,
	}

	fsys := afero.NewOsFs()

	switch stores.StoreType(backend) {
	case stores.StoreTypeS3:
		cfg = config.InitializeStoreTypeS3(
			cmd.Context(),
			fsys,
			repoToInit,
			backendAddress,
			region,
			opts,
		)
	case stores.StoreTypeGCS:
		cfg = config.InitializeStoreTypeGCS(
			cmd.Context(),
			fsys,
			repoToInit,
			backendAddress,
			opts,
		)
	default:
		return config.ErrUnsupportedStore
	}
	return cfg.Write(fsys, repoToInit)
}
