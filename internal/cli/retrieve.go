package cli

import (
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func retrieveCmd() *cobra.Command {
	retrieveCmd := &cobra.Command{
		Use:   "retrieve",
		Short: "retrieve a file from pantri",
		Long:  "retrieve a file from pantri",
		Args:  cobra.MinimumNArgs(1),
		// PersistentPreRunE
		// Loads the config with OsFs
		/*
			// The *Run functions are executed in the following order:
			//   * PersistentPreRunE()
			//   * PreRunE() [X]
			//   * RunE()
			//   * PostRunE()
			//   * PersistentPostRunE()
			// All functions get the same args, the arguments after the command name.
		*/
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return loadConfig(afero.NewOsFs())
		},
		RunE: retrieveFn,
	}

	return retrieveCmd
}

// retrieveFn is the execution runtime for the retrieveCmd functionality
// in Cobra, this is the RunE phase
/*
	// The *Run functions are executed in the following order:
	//   * PersistentPreRunE()
	//   * PreRunE()
	//   * RunE() [X]
	//   * PostRunE()
	//   * PersistentPostRunE()
	// All functions get the same args, the arguments after the command name.
*/
func retrieveFn(cmd *cobra.Command, objects []string) error {
	fsys := afero.NewOsFs()

	s, err := initStoreFromConfig(
		cmd.Context(),
		cfg,
		fsys,
		cfg.Options,
	)
	if err != nil {
		return err
	}

	logger.Infof("Downloading files from: %s", s.GetOptions().PantriAddress)
	logger.Infof("Downloading file: %s", objects)
	if err := s.Retrieve(cmd.Context(), objects...); err != nil {
		return err
	}
	return nil
}
