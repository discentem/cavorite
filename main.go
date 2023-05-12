package main

import (
	"context"
	"io"
	"log"

	"github.com/discentem/cavorite/internal/cli"
	"github.com/discentem/cavorite/internal/program"
	"github.com/google/logger"
)

func main() {
	defer logger.Init(program.Name, true, false, io.Discard).Close()
	ctx := context.Background()
	err := cli.ExecuteWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
