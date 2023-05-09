package cli

import (
	"context"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/metadata"
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// These vars are available to every sub command
	debug bool
	vv    bool
	cfg   config.Config

	// TODO (@radsec) Update this to be dynamic with GH action on new release and tagging....
	version string = "development"
)

func ExecuteWithContext(ctx context.Context) error {
	return rootCommand().ExecuteContext(ctx)
}

func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pantri",
		Short: "A source control friendly binary storage system",
		Long:  "A source control friendly binary storage system",
		PreRun: func(cmd *cobra.Command, args []string) {
			// set global logger opts
			setLoggerOpts()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
		Args:    cobra.NoArgs,
		Version: version,
	}

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")

	// Defaults set here will be used if they do not exist in the config file
	viper.SetDefault("store_type", stores.StoreTypeUndefined)
	viper.SetDefault("metadata_file_extension", metadata.MetaDataFileExtension)

	rootCmd.AddCommand(
		initCommand(),
		retrieveCommand(),
		uploadCommand(),
	)

	return rootCmd
}
