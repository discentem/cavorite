package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/discentem/pantri_but_go/stores"
	localstore "github.com/discentem/pantri_but_go/stores/local"
	"github.com/urfave/cli/v2"
)

type store interface {
	Upload(objects []string) error
	Retrieve(objects []string) error
}

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
		&cli.BoolFlag{
			Name:  "create_pantri",
			Value: true,
			Usage: "Set debug to true for enhanced logging",
		},
	}
	app := &cli.App{
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{},
				Usage:   "Initialize pantri repo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "storage",
						Aliases: []string{"s"},
						Value:   "local",
						Usage:   "Specify the storage type when initalizing a pantri repo.",
					},
					&cli.BoolFlag{
						Name:  "remove",
						Value: true,
						Usage: "Remove the file from local repo if present",
					},
				},
				Action: func(c *cli.Context) error {
					storage := c.String("storage")
					if storage == "" {
						return errors.New("--storage must be specified")
					}
					create := c.Bool("create_pantri")
					remove := c.Bool("remove")
					if storage == "local" {
						_, err := localstore.New("repo", "pantri", &create, &remove)
						if err != nil {
							return err
						}
					} else {
						return fmt.Errorf("init is not implemented yet for non %q storage", "local")
					}
					return nil
				},
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload the specified file",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to upload")
					}
					storage, err := stores.StorageType()
					if err != nil {
						return err
					}
					if *storage == "" {
						return errors.New("storage not detected")
					}
					if *storage == "local" {
						ls, err := localstore.Load()
						if err != nil {
							return err
						}
						s := store(ls)
						log.Printf("Uploading %s", c.Args().Slice())
						if err := s.Upload(c.Args().Slice()); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:    "retrieve",
				Aliases: []string{"r", "ret"},
				Usage:   "Retrieve the specified file",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return errors.New("you must pass the path of an object to retrieve")
					}

					storage, err := stores.StorageType()
					if err != nil {
						return err
					}
					if *storage == "" {
						return errors.New("storage not detected")
					}
					if *storage == "local" {
						ls, err := localstore.Load()
						if err != nil {
							return err
						}
						s := store(ls)
						log.Printf("Retrieving %s", c.Args().Slice())
						if err := s.Retrieve(c.Args().Slice()); err != nil {
							return err
						}
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
