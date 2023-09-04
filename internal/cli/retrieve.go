package cli

import (
	"errors"
	"fmt"

	"github.com/discentem/cavorite/internal/program"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func retrieveCmd() *cobra.Command {
	retrieveCmd := &cobra.Command{
		Use:   "retrieve",
		Short: fmt.Sprintf("retrieve a file from %s", program.Name),
		Long:  fmt.Sprintf("retrieve a file from %s", program.Name),
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
	defer s.Close()

	sourceRepoRoot, err := rootOfSourceRepo()
	if err != nil {
		return err
	}
	if sourceRepoRoot == nil {
		return errors.New("sourceRepoRoot cannot be nil")
	}

	// We need to remove the prefix from the path so it is relative
	objects, err = removePathPrefix(objects, *sourceRepoRoot)
	if err != nil {
		return fmt.Errorf("retrieve error: %w", err)
	}

	opts, err := s.GetOptions()
	if err != nil {
		return err
	}

	logger.Infof("Downloading files from: %s", opts.BackendAddress)
	logger.Infof("Downloading file: %s", objects)
	if err := s.Retrieve(cmd.Context(), objects...); err != nil {
		return err
	}
	return nil
}
