package stores

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores/pluginproto"
)

var (
	PluginSet = plugin.PluginSet{
		"store": &storePlugin{},
	}

	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "cavorite",
	}

	// FIXME: make configurable?
	HLog = hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
)

// clientStore is the client cavorite uses to communicate with the plugin
type clientStore struct {
	pluginproto.PluginClient
}

func (p *clientStore) Upload(ctx context.Context, objects ...string) error {
	_, err := p.PluginClient.Upload(ctx, &pluginproto.Objects{Objects: objects})
	return err
}

func MetadataMapToPluginProtoMap(mmap metadata.CfileMetadataMap) map[string]*pluginproto.ObjectMetadata {
	ppmm := make(map[string]*pluginproto.ObjectMetadata)
	for k, v := range mmap {
		ppmm[k] = &pluginproto.ObjectMetadata{
			Name:         v.Name,
			Checksum:     v.Checksum,
			DateModified: timestamppb.New(v.DateModified),
		}
	}
	return ppmm
}

func PluginProtoMapToMetadataMap(ppm *pluginproto.ObjectsAndMetadataMap) metadata.CfileMetadataMap {
	mm := make(metadata.CfileMetadataMap)
	for k, v := range ppm.Map {
		mm[k] = metadata.ObjectMetaData{
			Name:         v.Name,
			Checksum:     v.Checksum,
			DateModified: v.DateModified.AsTime(),
		}
	}
	return mm
}

func (p *clientStore) Retrieve(ctx context.Context, mmap metadata.CfileMetadataMap, objects ...string) error {
	_, err := p.PluginClient.Retrieve(
		ctx,
		&pluginproto.ObjectsAndMetadataMap{
			Objects: &pluginproto.Objects{
				Objects: objects,
			},
			Map: MetadataMapToPluginProtoMap(mmap),
		},
	)
	return err
}

func (p *clientStore) GetOptions() (Options, error) {
	opts, err := p.PluginClient.GetOptions(context.Background(), &emptypb.Empty{})
	if err != nil {
		return Options{}, err
	}

	return Options{
		BackendAddress:        opts.BackendAddress,
		MetadataFileExtension: opts.MetadataFileExtension,
		Region:                opts.Region,
	}, nil
}

func (p *clientStore) Close() error {
	return nil
}

// serverStore is the server the plugin uses to communicate with cavorite
type serverStore struct {
	Store
	pluginproto.UnimplementedPluginServer
}

func (p *serverStore) Upload(ctx context.Context, objects *pluginproto.Objects) (*emptypb.Empty, error) {
	err := p.Store.Upload(ctx, objects.Objects...)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (p *serverStore) Retrieve(ctx context.Context, mmap *pluginproto.ObjectsAndMetadataMap) (*emptypb.Empty, error) {
	err := p.Store.Retrieve(ctx, PluginProtoMapToMetadataMap(mmap), mmap.Objects.Objects...)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (p *serverStore) GetOptions(_ context.Context, _ *emptypb.Empty) (*pluginproto.Options, error) {
	opts, err := p.Store.GetOptions()
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pluginproto.Options{
		BackendAddress:        opts.BackendAddress,
		MetadataFileExtension: opts.MetadataFileExtension,
		Region:                opts.Region,
	}, nil
}

// storePlugin implements plugin.GRPCPlugin
type storePlugin struct {
	plugin.Plugin
	Store
}

func (p *storePlugin) GRPCServer(_ *plugin.GRPCBroker, server *grpc.Server) error {
	pluginproto.RegisterPluginServer(server, &serverStore{Store: p.Store})
	return nil
}

func (p *storePlugin) GRPCClient(ctx context.Context, _ *plugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return &clientStore{PluginClient: pluginproto.NewPluginClient(client)}, nil
}

// PluggableStore is the Store used by cavorite that wraps go-plugin
type PluggableStore struct {
	client *plugin.Client
	Store
}

func NewPluggableStore(_ context.Context, opts Options) (*PluggableStore, error) {
	cmd := exec.Command(opts.PluginAddress)
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  HandshakeConfig,
		Plugins:          PluginSet,
		Cmd:              cmd,
		Logger:           HLog,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// connect to RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not create rpc client: %w", err)
	}

	// get plugin as <any>
	raw, err := rpcClient.Dispense("store")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not dispense plugin: %w", err)
	}
	// assert to StoreWithOptions
	s := raw.(Store)

	// assert to Store
	// store := raw.(Store)

	return &PluggableStore{
		client: client,
		Store:  s,
	}, nil
}

func (p *PluggableStore) Close() error {
	p.client.Kill()
	return nil
}

// ListenAndServePlugin is used by plugins to start listening to requests
func ListenAndServePlugin(store Store, logger hclog.Logger) {
	PluginSet["store"] = &storePlugin{Store: store}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         PluginSet,
		Logger:          logger,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
