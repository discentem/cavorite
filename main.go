package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	pantri "github.com/discentem/pantri_but_go/config/loader"
	"github.com/discentem/pantri_but_go/stores"
	"github.com/urfave/cli/v2"
)

var (
	ErrObjectEmpty = fmt.Errorf("object can't be %q", "")
)

func main() {
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
			Aliases:  []string{},
			Usage:    "path to source repo",
		},
	}
	app := &cli.App{
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{},
				Usage:   "Initalize pantri.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "backend",
						Value:   "local",
						Aliases: []string{"b"},
						Usage:   "Specify the storage backend to use",
						EnvVars: []string{"BACKEND"},
					},
					&cli.StringFlag{
						Name:     "pantri_address",
						Required: true,
						Usage:    "path to pantri storage",
					},
				},
				Action: func(c *cli.Context) error {
					remove := c.Bool("remove")
					backend := c.String("backend")
					opts := stores.Options{
						RemoveFromSourceRepo: &remove,
					}
					sourceRepo := c.String("source_repo")
					pantriAddress := c.String("pantri_address")
					// store agnostic initialization, specific initialization determined by backend
					err := pantri.Initalize(sourceRepo, backend, pantriAddress, opts)
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
					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to upload")
					}
					sourceRepo := c.String("source_repo")
					s, err := pantri.Load(sourceRepo)
					if err != nil {
						return err
					}
					// TODO(discentem) improve log message to include pantriAddress
					log.Printf("Uploading %s", c.Args().Slice())
					if err := s.Upload(sourceRepo, c.Args().Slice()...); err != nil {
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
					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to retrieve")
					}
					sourceRepo := c.String("source_repo")
					s, err := pantri.Load(sourceRepo)
					if err != nil {
						return err
					}
					log.Printf("Retrieving %s", c.Args().Slice())
					if err := s.Retrieve(sourceRepo, c.Args().Slice()...); err != nil {
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
