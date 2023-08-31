package stores

import (
	"context"
	"net/rpc"

	"github.com/spf13/afero"
)

type TransferArgs struct {
	Ctx     context.Context
	Objects []string
}

type PluginStoreRPCServer struct {
	Impl GoPluginStore
}

func (p *PluginStoreRPCServer) Upload(ctx context.Context, objects ...string) error {
	return p.Impl.Upload(ctx, objects...)
}

func (p *PluginStoreRPCServer) Retrieve(ctx context.Context, objects ...string) error {
	return p.Impl.Retrieve(ctx, objects...)
}

func (p *PluginStoreRPCServer) GetOptions() (Options, error) {
	return p.Impl.GetOptions()
}

func (p *PluginStoreRPCServer) GetFsys() (afero.Fs, error) {
	return p.Impl.GetFsys()
}

type PluginStoreClient struct{ client *rpc.Client }

func (p *PluginStoreClient) Upload(ctx context.Context, objects ...string) error {

	return p.client.Call("Plugin.Upload", TransferArgs{Ctx: ctx, Objects: objects}, nil)
}

func (p *PluginStoreClient) Retrieve(ctx context.Context, objects ...string) error {
	return p.client.Call("Plugin.Retrieve", TransferArgs{Ctx: ctx, Objects: objects}, nil)
}

func (p *PluginStoreClient) GetOptions() (Options, error) {
	var opts Options
	err := p.client.Call("Plugin.GetOptions", nil, opts)
	if err != nil {
		return Options{}, err
	}
	return opts, err
}

func (p *PluginStoreClient) GetFsys() (afero.Fs, error) {
	var fs afero.Fs
	err := p.client.Call("Plugin.GetFsys", nil, fs)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

type GoPluginStore struct{}

func (s *GoPluginStore) Upload(ctx context.Context, objects ...string) error {
	return nil
}

func (s *GoPluginStore) Retrieve(ctx context.Context, objects ...string) error {
	return nil
}

func (s *GoPluginStore) GetOptions() (Options, error) {
	return Options{}, nil
}

func (s *GoPluginStore) GetFsys() (afero.Fs, error) {
	return afero.NewMemMapFs(), nil
}
