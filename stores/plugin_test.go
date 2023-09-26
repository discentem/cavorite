package stores

import (
	"testing"
	"time"

	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores/pluginproto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMetadataMapToPluginProtoMap(t *testing.T) {
	expected := map[string]*pluginproto.ObjectMetadata{
		"thing.cfile": {
			Name:         "thing",
			Checksum:     "whatever",
			DateModified: timestamppb.New(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)),
		},
	}
	actual := MetadataMapToPluginProtoMap(metadata.CfileMetadataMap{
		"thing.cfile": metadata.ObjectMetaData{
			Name:     "thing",
			Checksum: "whatever",
		},
	})
	require.Equal(t, expected, actual)
}

func TestPluginProtoMapToMetadataMap(t *testing.T) {
	expected := metadata.CfileMetadataMap{
		"thing.cfile": metadata.ObjectMetaData{
			Name:     "thing",
			Checksum: "whatever",
		},
	}
	actual := PluginProtoMapToMetadataMap(&pluginproto.ObjectsAndMetadataMap{
		Map: map[string]*pluginproto.ObjectMetadata{
			"thing.cfile": {
				Name:         "thing",
				Checksum:     "whatever",
				DateModified: timestamppb.New(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	})
	require.Equal(t, expected, actual)
}
