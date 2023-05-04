package main

import (
	"io"
	"log"

	"github.com/discentem/pantri_but_go/internal/cmd/pantri"
	"github.com/google/logger"
)

func main() {
	defer logger.Init("pantri_but_go", true, false, io.Discard).Close()
	err := pantri.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
