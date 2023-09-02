package stores

import (
	"context"
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
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
	HLog = hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
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
	err := p.Call("Plugin.Upload", TransferArgs{Objects: objects}, &resp)
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
	err := p.Call("Plugin.Retrieve", TransferArgs{Objects: objects}, &resp)
	if err != nil {
		return fmt.Errorf("could not call rpc for Retrieve: %w", err)
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func (p *ClientStore) GetOptions() (Options, error) {
	fmt.Println("ClientStore GetOptions")
	resp := new(OptionsResp)
	fmt.Println("before Plugin.GetOptions")
	err := p.Call("Plugin.GetOptions", new(interface{}), &resp)
	if err != nil {
		return Options{}, err
	}
	fmt.Println("after Plugin.GetOptions")
	if resp.Err != "" {
		return Options{}, errors.New(resp.Err)
	}
	return resp.Opts, err
}

type ServerStore struct {
	Store
}

func (p *ServerStore) Upload(objects []string, resp *ErrResp) error {
	err := p.Store.Upload(context.Background(), objects...)
	if err != nil {
		resp.Err = err.Error()
	}
	return nil
}

func (p *ServerStore) Retrieve(objects []string, resp *ErrResp) error {
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

func GetClient(address string) *plugin.Client {
	// start plugin
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         PluginSet,
		Cmd:             exec.Command(address),
		Logger:          HLog,
	})

	return client
}

func DispensePlugin(rpcClient *plugin.Client) (Store, error) {
	client, err := rpcClient.Client()
	if err != nil {
		return nil, err
	}
	// get plugin as <any>
	raw, err := client.Dispense("store")
	if err != nil {
		rpcClient.Kill()
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
	c := GetClient(opts.BackendAddress)
	defer c.Kill()
	ps, err := DispensePlugin(c)
	if err != nil {
		return err
	}
	fmt.Println("uploading from PluggableStore")
	err = ps.Upload(context.Background(), objects...)
	c.Kill()
	return err
}

func (s *PluggableStore) Retrieve(ctx context.Context, objects ...string) error {
	opts, err := s.GetOptions()
	if err != nil {
		return err
	}
	c := GetClient(opts.BackendAddress)
	ps, err := DispensePlugin(c)
	if err != nil {
		return err
	}
	err = ps.Retrieve(ctx, objects...)
	c.Kill()
	return err
}

func (s *PluggableStore) GetOptions() (Options, error) {
	opts, err := s.GetOptions()
	if err != nil {
		return Options{}, nil
	}
	c := GetClient(opts.BackendAddress)
	defer c.Kill()
	ps, err := DispensePlugin(c)
	if err != nil {
		return Options{}, err
	}
	return ps.GetOptions()
}
