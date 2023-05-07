package pantri

import (
	"errors"
	"log"
	"os"

	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var (
	// These vars are available to every sub command
	debug bool
	vv    bool
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "A source control friendly binary storage system",
		Long:  "A source control friendly binary storage system",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	return rootCmd
}

func setLoggerOpts() {
	if vv {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}

func Execute() error {
	rootCmd := getRootCmd()
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")
	// Defaults set here will be used if they do not exist in the config file
	viper.SetDefault("store_type", nil)
	viper.SetDefault("metadata_file_extension", "pfile")
	// Set up the config file details
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".pantri")

	initCmd, err := getInitCmd()
	if err != nil {
		return err
	}
	rootCmd.AddCommand(initCmd, getUploadCmd(), retrieveCmd)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ToDo - Figure out if there is a smarter way to do this. This is pretty janky.
			bypassFor := []string{"init", "help", "-h"}
			skipConfigCheck := false
			for _, cmd := range bypassFor {
				if slices.Contains(os.Args, cmd) {
					skipConfigCheck = true
				}
			}
			if len(os.Args) < 2 {
				logger.Info("no subcommand was issued, displaying help page...")
				if err := rootCmd.Help(); err != nil {
					return err
				}
				os.Exit(0)

			} else if !skipConfigCheck {
				return errors.New("No config file found, please run init in the base of the repo.")
			}
			// Config file not found; ignore error if desired
		} else {
			return errors.New("An error occured loading the configuration. Please confirm it is in the correct format.")
			// Config file was found but another error was produced
		}
	} else {
		// We may want to iterate a slice of settings that must be set before continuing
		if viper.Get("store_type") == "" {
			log.Fatal("No store_type has been specified, exiting...")
		}
	}
	// run root command
	return rootCmd.Execute()
}
