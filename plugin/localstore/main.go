package main

import (
	"context"
	"fmt"
	"os"

	"github.com/discentem/cavorite/internal/stores"
	"github.com/hashicorp/go-hclog"
)

type LocalStore struct {
	logger hclog.Logger
}

func (s *LocalStore) Upload(ctx context.Context, objects ...string) error {
	s.logger.Info(fmt.Sprintf("Uploading %v via localstore plugin", objects))
	return nil
}
func (s *LocalStore) Retrieve(ctx context.Context, objects ...string) error {
	s.logger.Info(fmt.Sprintf("Retrieving %v via localstore plugin", objects))
	return nil
}
func (s *LocalStore) GetOptions() (stores.Options, error) {
	s.logger.Info("GetOptions() called on localstore plugin")
	return stores.Options{
		Region: "plugin region",
	}, nil
}

func (s *LocalStore) Close() error {
	return nil
}

func main() {
	hlog := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	ls := &LocalStore{
		logger: hlog,
	}

	stores.ListenAndServePlugin(ls, hlog)
}
