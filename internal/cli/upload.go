package cli

import (
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func uploadCommand() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a file to pantri",
		Long:  "Upload a file to pantri",
		Args:  cobra.MinimumNArgs(1),
		// Load the config with OsFs
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return loadConfig(afero.NewOsFs())
		},
		RunE: uploadRunE,
	}

	return uploadCmd
}

func uploadRunE(cmd *cobra.Command, objects []string) error {
	fsys := afero.NewOsFs()

	s, err := buildStoresFromConfig(
		cmd.Context(),
		cfg,
		fsys,
		cfg.Options,
	)
	if err != nil {
		return err
	}

	logger.Infof("Uploading to: %s", s.GetOptions().PantriAddress)
	logger.Infof("Uploading file: %s", objects)
	if err := s.Upload(cmd.Context(), objects...); err != nil {
		return err
	}
	return nil

}
