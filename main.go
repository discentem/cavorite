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
			Name:     "backend",
			Value:    "local",
			Aliases:  []string{"b"},
			Usage:    "Specify the storage backend to use",
			FilePath: ".config",
			EnvVars:  []string{"BACKEND"},
		},
		&cli.BoolFlag{
			Name:  "remove",
			Value: false,
			Usage: "Remove the file from local repo if present",
		},
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
					if !c.Args().Present() {
						return errors.New("You must pass the path of an object to upload")
					}
					remove := c.Bool("remove")
					ls, err := localstore.New("repo", "pantri", &remove)
					if err != nil {
						return err
					}
					s := store(ls)
					log.Printf("Uploading %s", c.Args().Slice())
					s.Upload(c.Args().Slice())
					return err
				},
			},
			{
				Name:    "retrieve",
				Aliases: []string{"r"},
				Usage:   "Retrieve the specified file",
				Action: func(c *cli.Context) error {
					if !c.Args().Present() {
						return errors.New("You must pass the path of an object to retrieve")
					}
					remove := c.Bool("remove")
					ls, err := localstore.New("repo", "pantri", &remove)
					if err != nil {
						return err
					}
					s := store(ls)
					s.Retrieve(c.Args().Slice())
					return err
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete the specified file",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
