package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/discentem/pantri_but_go/stores"
	localstore "github.com/discentem/pantri_but_go/stores/local"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
)

func Initalize(sourceRepo, backend, address string, opts stores.Options) error {
	if backend != "local" {
		return fmt.Errorf("%s is not supported. only %q is currently supported", backend, "local")
	}
	_, err := localstore.New(sourceRepo, address, opts)
	if err != nil {
		return err
	}
	return nil

}

func LoadLocalStore(m map[string]interface{}) (stores.Store, error) {
	log.Printf("type %q detected in pantri %q", m["type"], m["pantri_address"])
	var s *localstore.Store
	if err := mapstructure.Decode(m, &s); err != nil {
		return nil, err
	}
	return stores.Store(s), nil
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
	if m["type"] == "local" {
		return LoadLocalStore(m)
	}
	return nil, errors.New("only local stores are supported for now")
}
