package main

import (
	"io"
	"log"

	"github.com/discentem/cavorite/internal/cmd/root"
	"github.com/google/logger"
)

func main() {
	defer logger.Init("cavorite", true, false, io.Discard).Close()
	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
