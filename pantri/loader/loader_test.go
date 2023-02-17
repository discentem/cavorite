package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func fsysMissingConfigDir(sourceRepo string) (afero.Fs, error) {
	return afero.NewMemMapFs(), nil
}

func fsysMissingConfigFile(sourceRepo string) (afero.Fs, error) {
	confDir := filepath.Join(sourceRepo, ".pantri")
	fsys := afero.NewMemMapFs()
	if err := fsys.MkdirAll(confDir, os.ModeAppend); err != nil {
		return nil, err
	}
	return fsys, nil
}

func fsysWithConfig(content, sourceRepo string) (afero.Fs, error) {
	confDir := filepath.Join(sourceRepo, ".pantri")
	fsys := afero.NewMemMapFs()
	if err := fsys.MkdirAll(confDir, os.ModeAppend); err != nil {
		return nil, err
	}
	if err := afero.WriteFile(
		fsys,
		filepath.Join(confDir, "config"),
		[]byte(content),
		os.ModeAppend,
	); err != nil {
		return nil, err
	}
	return fsys, nil
}

func TestConfigDirMissing(t *testing.T) {
	sourceRepo := "blah"
	// fake filesystem with repo sourceRepo, where `/blah/.pantri` does not exist
	fsys, err := fsysMissingConfigDir(sourceRepo)
	if err != nil {
		t.Error(err)
	}
	b, err := readConfig(fsys, sourceRepo)
	assert.ErrorIs(t, err, ErrConfigDirNotExist)
	assert.Empty(t, b)
}

func TestConfigMissing(t *testing.T) {
	sourceRepo := "blah"
	// fake filesystem with repo sourceRepo, where `/blah/.pantri` exists but `/blah/.pantri/config` does not
	fsys, err := fsysMissingConfigFile(sourceRepo)
	if err != nil {
		t.Error(err)
	}
	b, err := readConfig(fsys, sourceRepo)
	assert.ErrorIs(t, err, ErrConfigNotExist)
	assert.Empty(t, b)
}

func TestConfigExists(t *testing.T) {
	sourceRepo := "blah"
	// fake filesystem with repo sourceRepo with config file
	fsys, err := fsysWithConfig("blah", sourceRepo)
	if err != nil {
		t.Error(err)
	}
	b, err := readConfig(fsys, sourceRepo)
	assert.NoError(t, err)
	assert.Equal(t, []byte("blah"), b)
}

func TestReadConfig(t *testing.T) {
	// override expander so we don't make queries against the real filesystem
	// on real filesystems, default expander will resolve homedir. Here
	expander = func(path string) (string, error) {
		t.Log("using overriden expander for tests")
		return path, nil
	}

	// Ensure ErrConfigDirNotExist is returned if .pantri directory missing
	t.Run(".pantri directory missing", TestConfigDirMissing)
	// Ensure ErrConfigNotExist is returned if .pantri/config file is missing
	t.Run(".pantri/config file missing", TestConfigMissing)
	// Ensure ReadConfig() works when config exists in .pantri/config
	t.Run(".pantri/config exists", TestConfigExists)
}
