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
// 			Name:    "retrieve",
// 			Aliases: []string{"r"},
// 			Usage:   "Retrieve the specified file",
// 			Action: func(c *cli.Context) error {
// 				setLoggerOpts(c)
// 				var store stores.Store
// 				var cfg config.Config

// 				if c.NArg() == 0 {
// 					return errors.New("you must pass the path of an object to retrieve")
// 				}
// 				// Retrieve required and optional flag values
// 				awsRegion := c.String("region")
// 				sourceRepo := c.String("source_repo")

// 				ctx := context.Background()
// 				fsys := afero.NewOsFs()
// 				if sourceRepo == "" {
// 					return errors.New("sourceRepo must be specified")
// 				}

// 				b, err := config.ReadConfig(fsys, sourceRepo)
// 				if err != nil {
// 					return errors.New("")
// 				}
// 				err = json.Unmarshal(b, &cfg)
// 				if err != nil {
// 					return err
// 				}

// 				opts := cfg.Options

// 				switch cfg.StoreType {
// 				case stores.StoreTypeS3:
// 					s3, err := stores.NewS3StoreClient(ctx, fsys, awsRegion, sourceRepo, opts)
// 					if err != nil {
// 						return fmt.Errorf("improper stores.S3Client init: %v", err)
// 					}
// 					store = stores.Store(s3)
// 				default:
// 					return fmt.Errorf("type %s is not supported", cfg.StoreType.String())
// 				}

// 				logger.Infof("Downloading files from: %s", store.GetOptions().PantriAddress)
// 				logger.Infof("Downloading file: %s", c.Args().Slice())
// 				if err := store.Retrieve(ctx, sourceRepo, c.Args().Slice()...); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		},
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
