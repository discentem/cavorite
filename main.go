package main

import (
	"flag"
	"fmt"
	"log"

	localstore "github.com/discentem/pantri_but_go/stores/local"
)

type store interface {
	Upload(objects []string) error
	Retrieve(objects []string) error
}

var (
	ErrObjectEmpty = fmt.Errorf("object can't be %q", "")
)

func main() {
	obj := flag.String("object", "", "path to object")
	remove := flag.Bool("remove", false, "remove objects from repo dir")
	flag.Parse()
	if *obj == "" {
		log.Fatal(ErrObjectEmpty)
	}

	ls, err := localstore.New("repo", "pantri", remove)
	if err != nil {
		log.Fatal(err)
	}
	s := store(ls)
	if err := s.Upload([]string{*obj}); err != nil {
		log.Fatal(err)
	}

}
