package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/program"
	"github.com/discentem/cavorite/stores"
)

func initCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: fmt.Sprintf("Initialize a new %s repo", program.Name),
		Long:  fmt.Sprintf("Initialize a new %s repo", program.Name),
		Args: func(cmd *cobra.Command, args []string) error {
			fn := cobra.ExactArgs(1)
			err := fn(cmd, args)
			if err != nil {
				return fmt.Errorf("you must specify a path to a repo you want %s to track", program.Name)
			}
			return nil
		},
		PreRunE: initPreExecFn,
		RunE:    initFn,
	}

	initCmd.PersistentFlags().String("backend_address", "", "Address for the storage backend")
	initCmd.PersistentFlags().String("region", "us-east-1", "Default region for the storage backend")
	initCmd.PersistentFlags().String("store_type", "", "Storage backend to use")
	initCmd.PersistentFlags().String("plugin_address", "", "Address for go-plugin that provides implementation for Store")

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
	if err := viper.BindPFlag("plugin_address", cmd.PersistentFlags().Lookup("plugin_address")); err != nil {
		return errors.New("Failed to bind plugin_address to viper")
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

	pluginAddress := viper.GetString("plugin_address")
	fmt.Printf("pluginAddress: %s\n", pluginAddress)
	if pluginAddress != "" {
		opts.PluginAddress = pluginAddress
	}

	fsys := afero.NewOsFs()

	sb := stores.StoreType(backend)

	switch sb {
	case stores.StoreTypeS3:
		fallthrough
	case stores.StoreTypeGCS:
		fallthrough
	case stores.StoreTypeAzureBlob:
		fallthrough
	case stores.StoreTypeGoPlugin:
		config.Cfg = config.InitalizeStoreTypeOf(
			cmd.Context(),
			sb,
			fsys,
			repoToInit,
			backendAddress,
			opts,
		)
	default:
		return config.ErrUnsupportedStore
	}
	return config.Cfg.Write(fsys, repoToInit)
}
