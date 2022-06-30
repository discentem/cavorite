package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/discentem/pantri_but_go/stores"
	storesinit "github.com/discentem/pantri_but_go/stores/initialize"
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
						Name:    "source_repo",
						Value:   "repo",
						Aliases: []string{"sr"},
						Usage:   "path to source repo",
					},
				},
				Action: func(c *cli.Context) error {
					remove := c.Bool("remove")
					backend := c.String("backend")
					opts := stores.Options{
						RemoveFromSourceRepo: &remove,
					}
					sourceRepo := c.String("source_repo")
					// store agnostic initialization, specific initialization determined by backend
					err := storesinit.Initalize(sourceRepo, backend, "pantri", opts)
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
					remove := c.Bool("remove")
					ls, err := localstore.New("repo", "pantri", &remove)
					if err != nil {
						return err
					}
					s := store(ls)
					log.Printf("Uploading %s", c.Args().Slice())
					if err := s.Upload(c.Args().Slice()); err != nil {
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
					remove := c.Bool("remove")
					ls, err := localstore.New("repo", "pantri", &remove)
					if err != nil {
						return err
					}
					s := store(ls)
					log.Printf("Retrieving %s", c.Args().Slice())
					if err := s.Retrieve(c.Args().Slice()); err != nil {
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
