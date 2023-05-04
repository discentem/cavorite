package pantri

import (
	"fmt"
	"log"
	"os"

	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var (
	// These vars are available to every sub command
	debug   bool
	vv      bool
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "testing 1 2 3",
		Long:  "A test of the Cobra broadcasting system",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func setLoggerOpts() {
	if vv {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Run in debug mode")
	rootCmd.PersistentFlags().BoolVar(&vv, "vv", false, "Run in verbose logging mode")
	// Defaults set here will be used if they do not exist in the config file
	viper.SetDefault("store_type", nil)
	viper.SetDefault("metadata_file_extension", "pfile")
	// Set up the config file details
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".pantri")
	// err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if slices.Contains(os.Args, "init") {
				fmt.Println("Doing an init, ignore missing config")
			}
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	} else {
		// We may want to iterate a slice of settings that must be set before continuing
		if viper.Get("store_type") == "" {
			log.Fatal("No store_type has been specified, exiting...")
		}
	}

	rootCmd.AddCommand(initCmd, uploadCmd)
}
