package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/stores"
	"github.com/discentem/pantri_but_go/stores/local"
	localstore "github.com/discentem/pantri_but_go/stores/local"
	"github.com/discentem/pantri_but_go/stores/s3"
	"github.com/mitchellh/go-homedir"
)

func Initialize(ctx context.Context, sourceRepo, backend, address string, opts stores.Options) error {
	fmt.Println("loader:backend", backend)
	switch b := (backend); b {
	case "local":
		_, err := localstore.New(sourceRepo, address, opts)
		if err != nil {
			return err
		}
	case "s3":
		_, err := s3.New(ctx, sourceRepo, address, opts)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a supported store type", b)
	}
	return nil

}

func Load(sourceRepo string) (stores.Store, error) {
	cfile, err := homedir.Expand(
		fmt.Sprintf("%s/.pantri/config", sourceRepo),
	)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filepath.Dir(cfile)); err != nil {
		if os.IsNotExist(err) {
			log.Printf("%s has not be initialized as a pantri repo", sourceRepo)
			return nil, err
		}
	}
	f, err := os.Open(cfile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("%s has not be initialized as a pantri repo", sourceRepo)
			return nil, err
		}
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	switch t := (m["type"]); t {
	case "local":
		return local.Load(m)
	case "s3":
		return s3.Load(m)
	default:
		return nil, fmt.Errorf("%s is not a support store type", t)
	}
}
