package stores

import (
	"context"
	"errors"
	"fmt"
	"net/rpc"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

type TransferArgs struct {
	Objects []string
}

var (
	PluginSet = plugin.PluginSet{
		"store": &StorePlugin{},
	}
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "blah",
	}
)

type ClientStore struct{ *rpc.Client }

type ErrResp struct {
	Err string
}

type OptionsResp struct {
	Opts Options
	ErrResp
}

func (p *ClientStore) Upload(ctx context.Context, objects ...string) error {
	// context is not actually transferred
	resp := new(ErrResp)
	err := p.Call("Plugin.Upload", TransferArgs{Objects: objects}, resp)
	if err != nil {
		return fmt.Errorf("could not call rpc for Upload: %w", err)
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func (p *ClientStore) Retrieve(ctx context.Context, objects ...string) error {
	// context is not actually transferred
	resp := new(ErrResp)
	err := p.Call("Plugin.Retrieve", TransferArgs{Objects: objects}, resp)
	if err != nil {
		return fmt.Errorf("could not call rpc for Retrieve: %w", err)
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func (p *ClientStore) GetOptions() (Options, error) {
	resp := new(OptionsResp)
	err := p.Call("Plugin.GetOptions", new(interface{}), resp)
	if err != nil {
		return Options{}, err
	}
	if resp.Err != "" {
		return Options{}, errors.New(resp.Err)
	}
	return resp.Opts, err
}

type ServerStore struct {
	Store
}

func (p *ServerStore) Upload(resp *ErrResp, objects ...string) error {
	err := p.Store.Upload(context.Background(), objects...)
	if err != nil {
		resp.Err = err.Error()
	}
	return nil
}

func (p *ServerStore) Retrieve(resp *ErrResp, objects ...string) error {
	err := p.Store.Retrieve(context.Background(), objects...)
	if err != nil {
		resp.Err = err.Error()
	}
	return nil
}

func (p *ServerStore) GetOptions(resp *OptionsResp) error {
	opts, err := p.Store.GetOptions()
	if err != nil {
		resp.Err = err.Error()
		resp.Opts = Options{}
	}
	resp.Opts = opts
	return nil
}

// StorePlugin implements plugin.Plugin
type StorePlugin struct {
	Store
}

func (p *StorePlugin) Server(_ *plugin.MuxBroker) (interface{}, error) {
	return &ServerStore{Store: p.Store}, nil
}

func (p *StorePlugin) Client(_ *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &ClientStore{Client: client}, nil
}

func StartPlugin(address string) (Store, error) {
	// start plugin
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         PluginSet,
		Cmd:             exec.Command(address),
		// Logger: hclog.New(&hclog.LoggerOptions{
		// 	Name:   "plugin",
		// 	Output: os.Stdout,
		// 	Level:  hclog.Debug,
		// }),
	})

	// connect to RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not create rpc client: %v", err)
	}

	// get plugin as <any>
	raw, err := rpcClient.Dispense("store")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not get plugin: %v", err)
	}
	s := raw.(Store)
	return s, nil

}

type PluggableStore struct{}

func (s *PluggableStore) Upload(ctx context.Context, objects ...string) error {
	opts, err := s.GetOptions()
	if err != nil {
		return err
	}
	ps, err := StartPlugin(opts.BackendAddress)
	if err != nil {
		return err
	}
	return ps.Upload(context.Background(), objects...)
}

func (s *PluggableStore) Retrieve(ctx context.Context, objects ...string) error {
	opts, err := s.GetOptions()
	if err != nil {
		return err
	}
	ps, err := StartPlugin(opts.BackendAddress)
	if err != nil {
		return err
	}
	return ps.Retrieve(ctx, objects...)
}

func (s *PluggableStore) GetOptions() (Options, error) {
	opts, err := s.GetOptions()
	if err != nil {
		return Options{}, nil
	}
	ps, err := StartPlugin(opts.BackendAddress)
	if err != nil {
		return Options{}, err
	}
	return ps.GetOptions()
}
