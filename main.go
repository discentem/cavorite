package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
		&cli.StringFlag{
			Name:    "backend",
			Value:   "local",
			Aliases: []string{"b"},
			Usage:   "Specify the storage backend to use",
			EnvVars: []string{"BACKEND"},
		},
		&cli.BoolFlag{
			Name:  "remove",
			Value: false,
			Usage: "Remove the file from local repo if present",
		},
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
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload the specified file",
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
