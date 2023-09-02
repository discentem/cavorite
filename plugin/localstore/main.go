package main

import (
	"context"

	"github.com/discentem/cavorite/internal/stores"
	"github.com/google/logger"
	"github.com/hashicorp/go-plugin"
)

type LocalStore struct{}

func (s *LocalStore) Upload(ctx context.Context, objects ...string) error {
	logger.Info("logging for localStore plugin during Upload")
	return nil
}
func (s *LocalStore) Retrieve(ctx context.Context, objects ...string) error {
	logger.Info("logging for localStore plugin during Retrieve")
	return nil
}
func (s *LocalStore) GetOptions() (stores.Options, error) {
	logger.Info("logging for localStore plugin during GetOptions")
	return stores.Options{
		Region: "plugin region",
	}, nil
}

func main() {
	ls := LocalStore{}
	stores.PluginSet["store"] = &stores.StorePlugin{Store: &ls}

	// logger.Info("", stores.PluginSet)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: stores.HandshakeConfig,
		Plugins:         stores.PluginSet,
		// Logger:          logger,
	})
}
