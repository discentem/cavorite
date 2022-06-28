package initialize

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/discentem/pantri_but_go/stores"
	localstore "github.com/discentem/pantri_but_go/stores/local"
	"github.com/mitchellh/mapstructure"
)

func Initalize(sourceRepo, backend, address string, opts stores.Options) error {
	if backend != "local" {
		return fmt.Errorf("%s is not supported. only %q is currently supported", backend, "local")
	}
	_, err := localstore.NewWithOptions(sourceRepo, address, opts)
	if err != nil {
		return err
	}
	return nil

}

func Load(sourceRepo string) (stores.Store, error) {
	f, err := os.Open(fmt.Sprintf("%s/.pantri/config", sourceRepo))
	if err != nil {
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
	if m["type"] != "local" {
		return nil, errors.New("load only supports local type")
	}
	log.Printf("type %q detected in ./pantri/config", "local")
	var s *localstore.Store
	if err := mapstructure.Decode(m, &s); err != nil {
		return nil, err
	}
	return stores.Store(s), nil

}
