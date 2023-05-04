package main

import (
	"io"
	"log"

	"github.com/discentem/pantri_but_go/internal/cmd/pantri"
	"github.com/google/logger"
)

var (
	defaultStore = "local"
)

func main() {
	defer logger.Init("pantri_but_go", true, false, io.Discard).Close()
	err := pantri.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

// app := &cli.App{
// 	Flags: flags,
// 	Usage: `pantri: but in go!`,
// 	Commands: []*cli.Command{
// 		{
// 		{
// 			Name:    "delete",
// 			Aliases: []string{"d"},
// 			Usage:   "Delete the specified file",
// 			Action: func(c *cli.Context) error {
// 				return errors.New("delete is not implemented yet")
// 			},
// 		},
// 	},
// }
