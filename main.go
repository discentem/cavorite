package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/stores"

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
		&cli.StringFlag{
			Name:     "aws_region",
			Required: false,
			Aliases:  []string{"region"},
			Usage:    "optional: AWS region: defaults to us-east-1",
			Value:    "us-east-1",
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
					awsRegion := c.String("aws_region")

					ctx := context.Background()
					fsys := afero.NewOsFs()

					// store agnostic initialization, specific initialization determined by backend
					var storeType stores.StoreType
					var cfg config.Config
					switch storeType.FromString(backend) {
					case stores.StoreTypeS3:
						cfg = config.InitializeStoreTypeS3Config(
							ctx,
							fsys,
							sourceRepo,
							pantriAddress,
							awsRegion,
							opts,
						)
					default:
						return config.ErrUnsupportedStore
					}
					return cfg.Write(fsys, sourceRepo)
				},
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload the specified file",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "remove",
						Value: false,
						Usage: "Remove the file from local sourceRepo if present",
					},
				},
				Action: func(c *cli.Context) error {
					setLoggerOpts(c)
					var store stores.Store
					var cfg config.Config

					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to upload")
					}
					// Retrieve required and optional flag values
					awsRegion := c.String("region")
					sourceRepo := c.String("source_repo")

					ctx := context.Background()
					fsys := afero.NewOsFs()
					if sourceRepo == "" {
						return errors.New("sourceRepo must be specified")
					}

					b, err := config.ReadConfig(fsys, sourceRepo)
					if err != nil {
						return errors.New("")
					}
					err = json.Unmarshal(b, &cfg)
					if err != nil {
						return err
					}

					opts := cfg.Options

					switch cfg.StoreType {
					case stores.StoreTypeS3:
						s3, err := stores.NewS3StoreClient(ctx, fsys, awsRegion, sourceRepo, opts)
						if err != nil {
							return fmt.Errorf("improper stores.S3Client init: %v", err)
						}
						store = stores.Store(s3)
					default:
						return fmt.Errorf("type %s is not supported", cfg.StoreType.String())
					}

					logger.Infof("Uploading to: %s", store.GetOptions().PantriAddress)
					logger.Infof("Uploading file: %s", c.Args().Slice())
					if err := store.Upload(ctx, sourceRepo, c.Args().Slice()...); err != nil {
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
					var store stores.Store
					var cfg config.Config

					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to retrieve")
					}
					// Retrieve required and optional flag values
					awsRegion := c.String("region")
					sourceRepo := c.String("source_repo")

					ctx := context.Background()
					fsys := afero.NewOsFs()
					if sourceRepo == "" {
						return errors.New("sourceRepo must be specified")
					}

					b, err := config.ReadConfig(fsys, sourceRepo)
					if err != nil {
						return errors.New("")
					}
					err = json.Unmarshal(b, &cfg)
					if err != nil {
						return err
					}

					opts := cfg.Options

					switch cfg.StoreType {
					case stores.StoreTypeS3:
						s3, err := stores.NewS3StoreClient(ctx, fsys, awsRegion, sourceRepo, opts)
						if err != nil {
							return fmt.Errorf("improper stores.S3Client init: %v", err)
						}
						store = stores.Store(s3)
					default:
						return fmt.Errorf("type %s is not supported", cfg.StoreType.String())
					}

					logger.Infof("Downloading files from: %s", store.GetOptions().PantriAddress)
					logger.Infof("Downloading file: %s", c.Args().Slice())
					if err := store.Retrieve(ctx, sourceRepo, c.Args().Slice()...); err != nil {
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
