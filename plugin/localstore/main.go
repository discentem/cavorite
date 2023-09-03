package main

import (
	"context"
	"os"

	"github.com/discentem/cavorite/internal/stores"
	"github.com/hashicorp/go-hclog"
)

type LocalStore struct {
	logger hclog.Logger
}

func (s *LocalStore) Upload(ctx context.Context, objects ...string) error {
	s.logger.Info("logging for localStore plugin during Upload")
	err := os.WriteFile("/tmp/dat1", []byte("hello\nplugin\n"), 0644)
	return err
}
func (s *LocalStore) Retrieve(ctx context.Context, objects ...string) error {
	s.logger.Info("logging for localStore plugin during Retrieve")
	return nil
}
func (s *LocalStore) GetOptions() (stores.Options, error) {
	s.logger.Info("logging for localStore plugin during GetOptions")
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
