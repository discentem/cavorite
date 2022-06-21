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
	var (
		obj      = flag.String("object", "", "path to object")
		remove   = flag.Bool("remove", false, "remove objects from repo dir")
		upload   = flag.Bool("upload", false, "upload objects to pantri")
		retrieve = flag.Bool("retrieve", false, "retrieve objects from pantri")
	)
	flag.Parse()
	if !*upload && !*retrieve {
		log.Fatal("one of {upload, retrieve} must be passed")
	}

	if *obj == "" {
		log.Fatal(ErrObjectEmpty)
	}

	ls, err := localstore.New("repo", "pantri", remove)
	if err != nil {
		log.Fatal(err)
	}
	s := store(ls)
	objs := []string{*obj}

	if *retrieve {
		if err := s.Retrieve(objs); err != nil {
			log.Fatal(err)
		}
	} else if *upload {
		if err := s.Upload(objs); err != nil {
			log.Fatal(err)
		}
	}

}
