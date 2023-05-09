package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/discentem/pantri_but_go/internal/cli"
	"github.com/google/logger"
)

func main() {
	defer logger.Init("pantri_but_go", true, false, io.Discard).Close()
	ctx := context.Background()
	err := cli.ExecuteWithContext(ctx)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
