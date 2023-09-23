package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/program"
	"github.com/discentem/cavorite/stores"
	"github.com/google/logger"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func shouldRetrieve(fsys afero.Fs, m *metadata.ObjectMetaData, cfile string) (bool, error) {
	expectedHash := m.Checksum
	f, err := fsys.Open(strings.TrimSuffix(cfile, filepath.Ext(cfile)))
	if err != nil {
		return true, err
	}
	actualHash, err := metadata.SHA256FromReader(f)
	if err != nil {
		return false, err
	}
	if actualHash != expectedHash {
		return true, nil
	}
	return false, nil
}

func Retrieve(ctx context.Context, fsys afero.Fs, s stores.Store, cfiles ...string) error {
	var result *multierr.Error
	var objects []string
	var cmap metadata.CfileMetadataMap
	for _, cfile := range cfiles {
		m, err := metadata.ParseCfile(fsys, cfile)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		doRetrieve, err := shouldRetrieve(fsys, m, cfile)
		if !doRetrieve {
			if err != nil {
				result = multierr.Append(result, err)
			}
			// we don't need to retrieve because we already have it
			continue
		}
		objects = append(objects, m.Name)
		cmap[m.Name] = *m
	}
	retrieveErr := s.Retrieve(ctx, cmap, objects...)
	return multierr.Append(result, retrieveErr).ErrorOrNil()

}

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
			return config.Load(afero.NewOsFs())
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
		config.Cfg,
		fsys,
		config.Cfg.Options,
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
	return Retrieve(cmd.Context(), fsys, s, objects...)
}
