package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/logger"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/metadata"
	cavoriteObjLib "github.com/discentem/cavorite/objects"
	"github.com/discentem/cavorite/program"
	"github.com/discentem/cavorite/stores"
)

func uploadCmd() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: fmt.Sprintf("Upload a file to %s", program.Name),
		Long:  fmt.Sprintf("Upload a file to %s", program.Name),
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
		RunE: uploadFn,
	}

	return uploadCmd
}

// uploadFn is the execution runtime for the uploadCmd functionality
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

var (
	ErrWriteMetadataToFsys = errors.New("failed to write metadata to fsys")
	ErrOpen                = errors.New("failed to open")
	ErrUpload              = errors.New("failed to upload")
)

func upload(ctx context.Context, fsys afero.Fs, s stores.Store, objects ...string) error {
	opts, err := s.GetOptions()
	if err != nil {
		return err
	}
	logger.Info("Options:", opts)

	logger.Infof("Uploading to: %s", opts.BackendAddress)
	logger.Infof("Uploading file: %s", objects)

	var derivedKeys []string
	prefixOp := cavoriteObjLib.AddPrefixToKey{Prefix: opts.ObjectKeyPrefix}
	derivedKeys = objects
	if opts.ObjectKeyPrefix != "" {
		derivedKeys = cavoriteObjLib.ModifyMultipleKeys(
			prefixOp,
			objects...,
		)
	}

	if err := s.Upload(ctx, derivedKeys...); err != nil {
		logger.Error(err)
		return fmt.Errorf("%w for %v", ErrUpload, objects)
	}

	var errResult error
	for _, obj := range objects {
		f, err := fsys.Open(obj)
		if err != nil {
			errResult = multierr.Append(fmt.Errorf("%w for %s", ErrOpen, obj))
			continue
		}
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			errResult = multierr.Append(err)
			continue
		}
		mon := prefixOp.Modify(obj)
		err = metadata.WriteToFsys(metadata.FsysWriteRequest{
			Object:       mon,
			Fsys:         fsys,
			Fi:           f,
			MetadataPath: obj,
			Extension:    opts.MetadataFileExtension,
		})
		if err != nil {
			errResult = multierr.Append(fmt.Errorf("%w for %s", ErrWriteMetadataToFsys, obj))
		}
	}
	return errResult
}
func uploadFn(cmd *cobra.Command, objects []string) error {
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
		return fmt.Errorf("upload error: %w", err)
	}
	return upload(cmd.Context(), fsys, s, objects...)
}
