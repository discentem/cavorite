package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/google/logger"
	"github.com/spf13/afero"

	pantri "github.com/discentem/pantri_but_go/pantri/loader"

	"github.com/discentem/pantri_but_go/stores"

	"github.com/urfave/cli/v2"
)

var (
	defaultStore = "local"
)

func setLoggerOpts(c *cli.Context) {
	if c.Bool("vv") {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}

func main() {
	defer logger.Init("pantri_but_go", true, false, io.Discard).Close()

	flags := []cli.Flag{
		// Debug is not currently being used. Remove this line once we add logging
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "Set debug to true for enhanced logging",
		},
		&cli.StringFlag{
			Name:     "source_repo",
			Required: true,
			Aliases:  []string{"sr"},
			Usage:    "path to source repo",
		},
		&cli.BoolFlag{
			Name:     "vv",
			Required: false,
			Value:    false,
			Usage:    "displays logger.V(2).* messages",
		},
	}
	app := &cli.App{
		Flags: flags,
		Usage: `pantri: but in go!`,
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{},
				Usage:   "Initalize pantri.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "pantri_address",
						Aliases:  []string{"p", "pa"},
						Required: true,
						Usage:    "path to pantri storage",
					},
					&cli.StringFlag{
						Name:     "metadata_file_extension",
						Required: false,
						Aliases:  []string{"ext"},
						Usage:    "extension for object metadata files",
						Value:    "pfile",
					},
				},
				Action: func(c *cli.Context) error {
					setLoggerOpts(c)

					remove := c.Bool("remove")
					var backend string
					if c.NArg() == 0 {
						log.Printf("defaulting to %s for pantri storage", defaultStore)
						backend = defaultStore
					} else if c.NArg() == 1 {
						backend = c.Args().First()
					} else {
						return errors.New("specifying multiple backends not allowed, try again")
					}
					fileExt := c.String("metadata_file_extension")
					opts := stores.Options{
						RemoveFromSourceRepo:  &remove,
						MetaDataFileExtension: fileExt,
					}
					sourceRepo := c.String("source_repo")
					pantriAddress := c.String("pantri_address")
					// store agnostic initialization, specific initialization determined by backend
					err := pantri.Initialize(context.Background(), afero.NewOsFs(), sourceRepo, backend, pantriAddress, opts)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload the specified file",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "remove",
						Value: true,
						Usage: "Remove the file from local sourceRepo if present",
					},
				},
				Action: func(c *cli.Context) error {
					setLoggerOpts(c)

					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to upload")
					}
					sourceRepo := c.String("source_repo")
					s, err := pantri.Load(afero.NewOsFs(), sourceRepo)
					if err != nil {
						return err
					}
					// TODO(discentem) improve log message to include pantriAddress
					logger.Infof("Uploading %s", c.Args().Slice())
					if err := s.Upload(context.Background(), sourceRepo, c.Args().Slice()...); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "retrieve",
				Aliases: []string{"r"},
				Usage:   "Retrieve the specified file",
				Action: func(c *cli.Context) error {
					setLoggerOpts(c)

					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to retrieve")
					}
					sourceRepo := c.String("source_repo")
					s, err := pantri.Load(afero.NewOsFs(), sourceRepo)
					if err != nil {
						return err
					}
					// TODO(discentem) improve log message to include pantriAddress
					log.Printf("Retrieving %s", c.Args().Slice())
					if err := s.Retrieve(context.Background(), sourceRepo, c.Args().Slice()...); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete the specified file",
				Action: func(c *cli.Context) error {
					return errors.New("delete is not implemented yet")
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
