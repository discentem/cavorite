package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/discentem/pantri_but_go/env"
	localstore "github.com/discentem/pantri_but_go/stores/local"
	"github.com/spf13/cobra"
)

type store interface {
	Upload(objects []string) error
	Retrieve(objects []string) error
}

var (
	ErrObjectEmpty = fmt.Errorf("object can't be %q", "")
)

var (
	obj    string
	remove bool
)

func main() {
	// var (
	// 	obj = flag.String("object", "", "path to object")
	// )

	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "",
		Long:  `root command`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Stat(".config.yaml")
			if os.IsNotExist(err) {
				return errors.New("cannot find .config.yaml, aborting")
			}
			if obj == "" {
				return ErrObjectEmpty
			}
			return nil
		},
		// Make sure we are in a repo with a config

	}
	rootCmd.PersistentFlags().StringVarP(&obj, "object", "o", "", "path to object")

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "",
		Long:  `upload objects to pantri`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ls, err := localstore.New("repo", "pantri", &remove)
			if err != nil {
				return err
			}
			s := store(ls)
			objs := []string{obj}
			if err := s.Upload(objs); err != nil {
				return err
			}
			return nil
		},
	}

	uploadCmd.Flags().BoolVarP(&remove, "remove", "r", env.Bool("REMOVE", false), "remove object locally after upload")

	retrieveCmd := &cobra.Command{
		Use:   "retrieve",
		Short: "",
		Long:  `retrieve objects from pantri`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ls, err := localstore.New("repo", "pantri", &remove)
			if err != nil {
				return err
			}
			s := store(ls)
			objs := []string{obj}
			if err := s.Retrieve(objs); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(retrieveCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
