package cli

import (
	"context"
	"fmt"

	"github.com/discentem/cavorite/internal/config"
	"github.com/discentem/cavorite/internal/metadata"
	"github.com/discentem/cavorite/internal/program"
	"github.com/discentem/cavorite/internal/stores"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// These vars are available to every sub command
	debug         bool
	vv            bool
	cfg           config.Config
	pluginAddress string

	// TODO (@radsec) Update this to be dynamic with GH action on new release and tagging....
	version string = "development"
)

func ExecuteWithContext(ctx context.Context) error {
	return rootCmd().ExecuteContext(ctx)
}

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf(program.Name),
		Short: "A source control friendly binary storage system",
		Long:  "A source control friendly binary storage system",
		// PersistentPreRun -- all downstream cmds will inherit this fn()
		// Set the global logger opts
		/*
			// The *Run functions are executed in the following order:
			//   * PersistentPreRun() [X]
			//   * PreRun()
			//   * Run()
			//   * PostRun()
			//   * PersistentPostRun()
			// All functions get the same args, the arguments after the command name.
		*/
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLoggerOpts()
		},
		// RunE
		// Return the help page if an error occurs
		/*
			// The *Run functions are executed in the following order:
			//   * PersistentPreRunE()
			//   * PreRunE()
			//   * RunE() [X]
			//   * PostRunE()
			//   * PersistentPostRunE()
			// All functions get the same args, the arguments after the command name.
		*/
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
		Args:    cobra.NoArgs,
		Version: version,
	}

	// At the rootCmd level, set these global flags that will be available to downstream cmds
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")

	// Defaults set here will be used if they do not exist in the config file
	viper.SetDefault("store_type", stores.StoreTypeUndefined)
	viper.SetDefault("metadata_file_extension", metadata.MetadataFileExtension)

	if vv {
		logger.SetLevel(2)
	}

	// Import subCmds into the rootCmd
	rootCmd.AddCommand(
		initCmd(),
		retrieveCmd(),
		uploadCmd(),
	)

	return rootCmd
}
