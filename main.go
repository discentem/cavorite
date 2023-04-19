package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/logger"
	"github.com/spf13/afero"

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

func getPrereqs(c *cli.Context) {
	// backend - type == local/s3

	// sourceRepo - string

	// fsys afero.FS
}

var (
	globalStoreOpts stores.Options

)

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
			Aliases:  []string{"verbose"},
			Required: false,
			Value:    false,
			Usage:    "displays logger.V(2).* messages",
		},
	}
	app := &cli.App{
		Before: func(cCtx *cli.Context) error {
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
			sourceRepo := c.String("source_repo")
			pantriAddress := c.String("pantri_address")
			globalStoreOpts = stores.Options{
				RemoveFromSourceRepo:  &remove,
				MetaDataFileExtension: fileExt,
				PantriAddress: pantriAddress,
			}
	
			// store agnostic initialization, specific initialization determined by backend
			// return pantri.Initialize(context.Background(), afero.NewOsFs(), sourceRepo, backend, pantriAddress, opts)

			fmt.Fprintf(cCtx.App.Writer, "BEFORE WAS RUN\n")
			return nil
		},
		Flags: flags,
		Usage: `pantri: but in go!`,
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{},
				Usage:   "Initalize pantri.",
				Flags: []cli.Flag{
					// &cli.StringFlag{
					// 	Name:     "pantri_address",
					// 	Aliases:  []string{"p", "pa"},
					// 	Required: true,
					// 	Usage:    "path to pantri storage",
					// },
					// &cli.StringFlag{
					// 	Name:     "metadata_file_extension",
					// 	Required: false,
					// 	Aliases:  []string{"ext"},
					// 	Usage:    "extension for object metadata files",
					// 	Value:    "pfile",
					// },
					&cli.StringFlag{
						Name:     "region",
						Required: false,
						Aliases:  []string{},
						Usage:    "region is used for Storage providers that have a geographical concept. Mostly for cloud providers.",
						Value:    "us-east-1",
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

					fmt.Println(globalStoreOpts)

					fs := afero.NewOsFs()
					ctx := context.Background()

					// switch backend {
					// case "s3":
					// 	s, err := stores.NewS3StoreClient(ctx, fs, awsRegion, sourceRepo, opts)
					// 	if err != nil {
					// 		return err
					// 	}
					// 	err = s.WriteConfig(ctx, sourceRepo)
					// 	if err != nil {
					// 		return err
					// 	}
					// default:
					// 	return nil, fmt.Errorf("%s: %w", b, ErrUnsupportedStore)
					// }

					return nil
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
					&cli.StringFlag{
						Name:     "destination",
						Required: true,
						Usage:    "The directory in source_repo to 'upload' the object to",
					},
				},
				Action: func(c *cli.Context) error {
					setLoggerOpts(c)

					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to upload")
					}
					sourceRepo := c.String("source_repo")
					destination := c.String("destination")
					fsys := afero.NewOsFs()
					s, err := pantri.Load(fsys, sourceRepo)
					if err != nil {
						return err
					}
					// TODO(discentem) improve log message to include pantriAddress
					logger.Infof("Uploading %s to %s", c.Args().Slice(), destination)
					if err := s.Upload(context.Background(), fsys, sourceRepo, destination, c.Args().Slice()...); err != nil {
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
					fsys := afero.NewOsFs()

					//example start - lets init a standardized --> pantri client
					// s3Client := s3.New()

					s, err := pantri.InitWithS3(s3.New(ctx, "us-east-1", fsys, sourceRepo, address, opts))

					s, err := pantri.InitWithLocalStore(localstore.New(fsys, sourceRepo, address, opts))

					s.Retrieve()
					s.Upload()

					// example end

					pantri.InitializeS3()
					s, err := pantri.Load(fsys, sourceRepo)
					if err != nil {
						return err
					}
					// TODO(discentem) improve log message to include pantriAddress
					log.Printf("Retrieving %s", c.Args().Slice())
					if err := s.Retrieve(context.Background(), fsys, sourceRepo, c.Args().Slice()...); err != nil {
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
